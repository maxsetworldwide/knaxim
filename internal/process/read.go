package process

import (
	"context"
	"io"
	"strings"
	"sync"
	"time"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/errors"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
	"git.maxset.io/web/knaxim/pkg/asyncreader"
)

//Read generates meta data from the content of a filestore
func Read(ctx context.Context, fs *types.FileStore, db database.Database, tika string, gotenburg string) {
	errlock := new(sync.Mutex)
	var errs []error
	pusherr := func(e error) {
		if e != nil {
			errlock.Lock()
			defer errlock.Unlock()
			errs = append(errs, e)
		}
	}
	tagch := make(chan []tag.Tag)
	tb := db.Tag(ctx)
	go func() {
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
	go func(r io.Reader) {
		//   Split Sentences
		//     > ContentLines
		//     skyset
		//       aggregate data > NLP Tags
	}(readtexts[0])
	go func(r io.Reader) {
		//   Split Words > ContentTags
		tags, err := tag.ExtractContentTags(r)
		if err != nil {
			pusherr(err)
			return
		}
		select {
		case tagch <- tags:
		case <-ctx.Done():
		}
	}(readtexts[1])
	go func(w io.WriteCloser) {
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

	// Gotenburg > View

	//After all jobs are done
	reportErrors(errs, fs, db)
}

func reportErrors(errs []error, fs *types.FileStore, db database.Database) {
	if len(errs) == 0 {
		return
	}
	sb := new(strings.Builder)
	sb.WriteString("Processing Errors:")
	for _, e := range errs {
		sb.WriteByte(' ')
		sb.WriteString(e.Error())
		sb.WriteByte('.')
	}
	errctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	fs.Perr = &errors.Processing{
		Status:  500,
		Message: sb.String(),
	}
	db.Store(errctx).UpdateMeta(fs)
}
