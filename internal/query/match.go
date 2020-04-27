package query

import (
	"errors"

	"git.maxset.io/web/knaxim/internal/database/types/tag"
)

type M struct {
	Tag   tag.Type    `json:"tagtype"`
	Word  string      `json:"word"`
	Regex interface{} `json:"regex,omitempty"`
}

func decodeM(i interface{}) (matches []M, err error) {
	switch v := i.(type) {
	case []interface{}:
		for _, ele := range v {
			var temp []M
			temp, err = decodeM(ele)
			if err != nil {
				return
			}
			matches = append(matches, temp...)
		}
	case map[string]interface{}:
		tstr, ok := v["tagtype"].(string)
		if !ok {
			return nil, errors.New("Missing tagtype")
		}
		var t tag.Type
		t, err = tag.DecodeType(tstr)
		if err != nil {
			return
		}
		w, ok := v["word"].(string)
		if !ok {
			return nil, errors.New("Missing word")
		}
		matches = append(matches, M{
			Tag:   t,
			Word:  w,
			Regex: v["regex"],
		})
	case string:
		matches = append(matches, M{
			Tag:   tag.CONTENT,
			Word:  v,
			Regex: true,
		})
	default:
		return nil, errors.New("Unrecognized Match Value")
	}
	return
}
