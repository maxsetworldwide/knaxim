package query

import (
	"errors"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
)

type M struct {
	Tag   tag.Type      `json:"tagtype"`
	Word  string        `json:"word"`
	Owner types.OwnerID `json:"owner,omitempty"`
	Regex interface{}   `json:"regex,omitempty"`
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
		var owner types.OwnerID
		if v["owner"] != nil {
			o, ok := v["owner"].(string)
			if !ok {
				return nil, errors.New("owner must be a string in match condition")
			}
			owner, err = types.DecodeObjectIDString(o)
		}
		matches = append(matches, M{
			Tag:   t,
			Word:  w,
			Regex: v["regex"],
			Owner: owner,
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

func (m M) SearchTag() tag.FileTag {
	ft := tag.FileTag{
		Owner: m.Owner,
		Tag: tag.Tag{
			Word: m.Word,
			Type: m.Tag,
		},
	}
	if m.Regex != nil {
		ft.Type = ft.Type | tag.SEARCH
		ft.Data = tag.Data{
			tag.SEARCH: map[string]interface{}{
				"regex":        true,
				"regexoptions": "i",
			},
		}
	}
	return ft
}
