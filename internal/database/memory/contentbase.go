package memory

import (
	"regexp"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/pkg/srverror"
)

type Contentbase struct {
	Database
}

func (cb *Contentbase) Insert(lines ...database.ContentLine) error {
	lock.Lock()
	defer lock.Unlock()
	for _, line := range lines {
		if len(cb.Lines[line.ID.String()]) <= line.Position {
			if cap(cb.Lines[line.ID.String()]) <= line.Position {
				newarr := make([]database.ContentLine, line.Position, line.Position*2+2)
				copy(newarr, cb.Lines[line.ID.String()])
				newarr = append(newarr, line)
				cb.Lines[line.ID.String()] = newarr
			} else {
				for len(cb.Lines[line.ID.String()]) < line.Position {
					cb.Lines[line.ID.String()] = append(cb.Lines[line.ID.String()], database.ContentLine{})
				}
				cb.Lines[line.ID.String()] = append(cb.Lines[line.ID.String()], line)
			}
		} else {
			cb.Lines[line.ID.String()][line.Position] = line
		}
	}
	return nil
}

func (cb *Contentbase) Len(id filehash.StoreID) (int64, error) {
	lock.RLock()
	defer lock.RUnlock()
	return int64(len(cb.Lines[id.String()])), nil
}

func (cb *Contentbase) Slice(id filehash.StoreID, start int, end int) ([]database.ContentLine, error) {
	lock.RLock()
	defer lock.RUnlock()
	return cb.slice(id, start, end)
}

func (cb *Contentbase) slice(id filehash.StoreID, start int, end int) ([]database.ContentLine, error) {
	var perr error
	{
		sb := cb.store(nil).(*Storebase)
		fs, err := sb.get(id)
		if err != nil {
			sb.close()
			return nil, err
		}
		if fs.Perr != nil {
			perr = fs.Perr
		}
		sb.close()
	}
	if len(cb.Lines[id.String()]) < end {
		end = len(cb.Lines[id.String()])
	}
	if start >= end {
		return nil, nil
	}
	return cb.Lines[id.String()][start:end], perr
}

func (cb *Contentbase) RegexSearchFile(regex string, file filehash.StoreID, start int, end int) ([]database.ContentLine, error) {
=======
	lock.RLock()
	defer lock.RUnlock()
	return cb.regexSearchFile(regex, file, start, end)
}

func (cb *Contentbase) regexSearchFile(regex string, file filehash.StoreID, start int, end int) ([]database.ContentLine, error) {
	var perr error
	{
		sb := cb.store(nil).(*Storebase)
		fs, err := sb.get(file)
		if err != nil {
			sb.close()
			return nil, err
		}
		if fs.Perr != nil {
			perr = fs.Perr
		}
		sb.close()
	}
	rgx, err := regexp.Compile(regex)
	if err != nil {
		return nil, srverror.New(err, 400, "Bad Search", "search string failed to compile to regex")
	}
	slice, _ := cb.slice(file, start, end)
	var out []database.ContentLine
	for _, line := range slice {
		for _, content := range line.Content {
			if rgx.MatchString(content) {
				out = append(out, line)
				break
			}
		}
	}
	return out, perr
}
