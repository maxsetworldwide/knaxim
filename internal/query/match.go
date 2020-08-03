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

package query

import (
	"errors"

	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/internal/database/types/tag"
)

// M is the matching condition to filter file ids by
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
			owner, err = types.DecodeOwnerIDString(o)
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

// SearchTag builds a tag.FileTag the represents the same data as the M so that it can be used to search tags
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
