package mongo

import (
	"errors"

	"git.maxset.io/web/knaxim/internal/database/types"
)

type contentchunk struct {
	ID    types.StoreID `bson:"id"`
	Index uint32        `bson:"idx"`
	Data  []byte        `bson:"data"`
}

const chunksize = 15 << 20

func chunkify(ID types.StoreID, content []byte) []interface{} {
	var chunks []interface{}
	var i uint32
	for start := 0; start < len(content); start += chunksize {
		end := start + chunksize
		if end > len(content) {
			end = len(content)
		}
		chunks = append(chunks, &contentchunk{
			ID:    ID,
			Index: i,
			Data:  content[start:end],
		})
		i++
	}
	return chunks
}

func chunksort(list []*contentchunk) []*contentchunk {
	pos := 0
	for pos < len(list) {
		target := int(list[pos].Index)
		if target == pos {
			pos++
		} else {
			if target == int(list[target].Index) {
				panic(errors.New("Improper chunk list"))
			}
			list[pos], list[target] = list[target], list[pos]
		}
	}
	return list
}

func appendchunks(list []*contentchunk) []byte {
	out := make([]byte, 0, (len(list)-1)*chunksize+len(list[len(list)-1].Data))
	for _, chunk := range list {
		out = append(out, chunk.Data...)
	}
	return out
}

func filterchunks(list []*contentchunk) [][]*contentchunk {
	out := make([][]*contentchunk, 0)
	for _, ch := range list {
		added := false
		for i, outlist := range out {
			if ch.ID.Equal(outlist[0].ID) {
				out[i] = append(outlist, ch)
				added = true
				break
			}
		}
		if !added {
			out = append(out, []*contentchunk{ch})
		}
	}
	return out
}
