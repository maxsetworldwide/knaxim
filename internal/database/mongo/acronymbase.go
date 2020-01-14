package mongo

import (
	"bytes"
	"strings"

	"git.maxset.io/server/knaxim/srverror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type acronym struct {
	Acronym  string `bson:"acronym"`
	Complete string `bson:"complete"`
}

type Acronymbase struct {
	Database
}

func stripAcronym(in string) string {
	abytes := []byte(in)
	abytes = bytes.ToUpper(abytes)
	keybuilder := new(strings.Builder)
	for _, c := range abytes {
		if 'A' <= c && 'Z' >= c {
			keybuilder.WriteByte(c)
		}
	}
	return keybuilder.String()
}

func (ab *Acronymbase) Put(a, c string) error {
	val := acronym{
		Acronym:  stripAcronym(a),
		Complete: c,
	}
	_, err := ab.client.Database(ab.DBName).Collection(ab.CollNames["acronym"]).UpdateOne(ab.ctx, val, bson.M{
		"$setOnInsert": val,
	}, options.Update().SetUpsert(true))
	return srverror.New(err, 500, "Database Error A1", "Failed to insert acronym")
}

func (ab *Acronymbase) Get(a string) ([]string, error) {
	cursor, err := ab.client.Database(ab.DBName).Collection(ab.CollNames["acronym"]).Find(ab.ctx, bson.M{
		"acronym": stripAcronym(a),
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, srverror.New(err, 500, "Database Error A2", "Failed to find acronym")
	}
	var result []acronym
	if err := cursor.All(ab.ctx, &result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, srverror.New(err, 500, "Database Error A2.1", "Failed to decode acronym")
	}
	var out []string
	for _, r := range result {
		out = append(out, r.Complete)
	}
	return out, nil
}
