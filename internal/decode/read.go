package decode

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"strings"
	"sync"
	"time"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/process"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"git.maxset.io/web/knaxim/pkg/asyncreader"
	"git.maxset.io/web/knaxim/pkg/skyset"
)

//Read generates meta data from the content of a filestore
func Read(ctx context.Context, cncl context.CancelFunc, name string, fs *types.FileStore, dbconfig database.Database, tika string, gotenburg string) {
	errlock := new(sync.Mutex)
	var errs []error
	pusherr := func(e error) {
		if e != nil {
			errlock.Lock()
			defer errlock.Unlock()
			errs = append(errs, e)
		}
	}
	db, err := dbconfig.Connect(ctx)
	pusherr(err)
	if err == nil {
		wg := new(sync.WaitGroup)
		wg.Add(5)
		tagch := make(chan []tag.Tag, 2)
		tagfinished := new(sync.WaitGroup)
		tagfinished.Add(1)
		go func() {
			defer tagfinished.Done()
			tb := db.Tag()
			sudofileid := types.FileID{
				StoreID: fs.ID,
			}
			for tags := range tagch {
				filetags := []tag.FileTag{}
				for _, t := range tags {
					filetags = append(filetags, tag.FileTag{
						File: sudofileid,
						Tag:  t,
					})
				}
				err := tb.Upsert(filetags...)
				if err != nil {
					pusherr(err)
					return
				}
			}
		}()

		// Tika
		writetext, readtexts := asyncreader.New(2)
		go func(w io.WriteCloser) {
			defer wg.Done()
			defer w.Close()
			// feed result of tika to writer
			fileRead, err := fs.Reader()
			if err != nil {
				pusherr(err)
				return
			}
			tread, err := tikaTextExtract(ctx, fileRead, tika)
			if err != nil {
				pusherr(err)
				return
			}
			defer tread.Close()
			_, err = io.Copy(w, tread)
			if err != nil {
				pusherr(err)
				return
			}
		}(writetext)
		go func(r io.Reader) {
			defer wg.Done()
			scanner := bufio.NewScanner(r)
			scanner.Split(SentenceSplitter)
			nlpch := make(chan string, 5)
			go func() {
				defer wg.Done()
				var nlp nlpaggregate
				for sent := range nlpch {
					nlp.add(skyset.BuildPhrases(sent))
				}
				var nlptags []tag.Tag
				for syn, data := range nlp.report() {
					var typ tag.Type
					switch syn {
					case skyset.TOPIC:
						typ = tag.TOPIC
					case skyset.ACTION:
						typ = tag.ACTION
					case skyset.PROCESS:
						typ = tag.PROCESS
					case skyset.RESOURCE:
						typ = tag.RESOURCE
					}
					nlptags = append(nlptags, data.tags(typ)...)
				}
				select {
				case tagch <- nlptags:
				case <-ctx.Done():
				}
			}()
			var ContentLines []types.ContentLine
			for i := 0; scanner.Scan(); i++ {
				sentence := scanner.Text()
				nlpch <- sentence
				ContentLines = append(ContentLines, types.ContentLine{
					ID:       fs.ID,
					Position: i,
					Content:  []string{sentence},
				})
			}
			close(nlpch)
			if err := scanner.Err(); err != nil {
				pusherr(err)
				return
			}
			cb := db.Content()
			if err := cb.Insert(ContentLines...); err != nil {
				pusherr(err)
			}
		}(readtexts[0])
		go func(r io.Reader) {
			defer wg.Done()
			tags, err := tag.ExtractContentTags(strings.NewReader(name))
			if err != nil {
				pusherr(err)
				return
			}
			//   Split Words > ContentTags
			ftags, err := tag.ExtractContentTags(r)
			if err != nil {
				pusherr(err)
				return
			}
			tags = append(tags, ftags...)
			select {
			case tagch <- tags:
			case <-ctx.Done():
				pusherr(ctx.Err())
				return
			}
		}(readtexts[1])

		// Gotenburg > View
		go func() {
			defer wg.Done()
			var result []byte
			extConst, ok := process.ExtMap[fs.ContentType]
			if !ok || extConst == process.PDF {
				// no conversions available. do not put a view in the db. retrieval of this
				// view should return 404 or 302 or 303 to indicate that sentences should be used
				// OR is PDF
				// do not store a copy in the viewbase
				// have the /view api just return the store by checking the content type
				return
			}
			buf := &bytes.Buffer{}
			r, err := fs.Reader()
			if err != nil {
				pusherr(err)
				return
			}
			if _, err = io.Copy(buf, r); err != nil {
				pusherr(err)
				return
			}
			converter := process.NewFileConverter(gotenburg)
			gotenFinished := make(chan error)
			go func() {
				var err error
				switch extConst {
				case process.OFFICE:
					result, err = converter.ConvertOffice(name, buf.Bytes())
				}
				gotenFinished <- err
			}()
			select {
			case err := <-gotenFinished:
				if err != nil {
					pusherr(err)
					return
				}
			case <-ctx.Done():
				pusherr(ctx.Err())
				return
			}
			vs, err := types.NewViewStore(fs.ID, bytes.NewReader(result))
			if err != nil {
				pusherr(err)
				return
			}
			vb := db.View()
			if err = vb.Insert(vs); err != nil {
				pusherr(err)
				return
			}
		}()

		//After all jobs are done
		wg.Wait()
		close(tagch)
		tagfinished.Wait()
		if cncl != nil {
			cncl()
		}
		db.Close(ctx)
	}
	if len(errs) == 0 {
		fs.Perr = nil
	} else {
		sb := new(strings.Builder)
		sb.WriteString("Processing Errors:")
		for _, e := range errs {
			sb.WriteByte(' ')
			sb.WriteString(e.Error())
			sb.WriteByte('.')
		}
		fs.Perr = &errors.Processing{
			Status:  242,
			Message: sb.String(),
		}
	}
	errctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	db, err = dbconfig.Connect(errctx)
	if err != nil {
		return
	}
	db.Store().UpdateMeta(fs)
}
