package database

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/internal/database/tag"
	"git.maxset.io/web/knaxim/pkg/srverror"
)

func InjestFile(ctx context.Context, file FileI, contenttype string, stream io.Reader, db Database) (fs *FileStore, err error) {
	defer func() {
		if r := recover(); r != nil {
			fs = nil
			switch v := r.(type) {
			case srverror.Error:
				err = v
			case error:
				err = srverror.New(v, 500, "Server Error", "unable to process input")
			default:
				err = srverror.New(fmt.Errorf("Error Injecting File: %+#v", v), 500, "Server Error")
			}
		}
	}()
	fs, err = NewFileStore(stream)
	if err != nil {
		panic(err)
	}
	fs.ContentType = contenttype

	{
		ownerbase := db.Owner(ctx)
		currentspace, err := ownerbase.GetSpace(file.GetOwner().GetID())
		if err != nil {
			panic(err)
		}
		totalspace, err := ownerbase.GetTotalSpace(file.GetOwner().GetID())
		if err != nil {
			panic(err)
		}
		if currentspace+fs.FileSize > totalspace {
			panic(srverror.Basic(460, "No Space"))
		}
	}
	{
		sb := db.Store(ctx)
		matches, err := sb.MatchHash(fs.ID.Hash)
		if err != nil {
			panic(err)
		}
		var matched bool
		for _, m := range matches {
			if bytes.Equal(fs.Content, m.Content) {
				fs = m
				matched = true
				break
			}
		}
		if !matched {
			fs.ID, err = sb.Reserve(fs.ID)
			if err != nil {
				panic(err)
			}
			err = sb.Insert(fs)
			if err != nil {
				panic(err)
			}
		}
	}
	{
		fb := db.File(ctx)
		tempID, err := fb.Reserve(filehash.NewFileID(fs.ID))
		if err != nil {
			panic(err)
		}
		file.setID(tempID)
		err = fb.Insert(file)
		if err != nil {
			panic(err)
		}
	}
	return fs, nil
}

func generateContentTags(ctx context.Context, fs *FileStore, db Database) {
	go func() {
		rcontent, err := fs.Reader()
		if err != nil {
			return
		}
		tags, err := tag.ExtractContentTags(rcontent)
		if err != nil {
			return
		}
		tb := db.Tag(ctx)
		tb.UpsertStore(fs.ID, tags...)
	}()
}
