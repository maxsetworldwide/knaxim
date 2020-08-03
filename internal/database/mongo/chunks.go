// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongo

import (
	"context"
	"errors"

	"git.maxset.io/web/knaxim/internal/database/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initChunkIndex(ctx context.Context, d *Database, client *mongo.Client) error {
	I := client.Database(d.DBName).Collection(d.CollNames["chunk"]).Indexes()
	_, err := I.CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{bson.E{Key: "id", Value: 1}, bson.E{Key: "idx", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

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
