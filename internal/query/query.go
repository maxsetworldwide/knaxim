package query

import (
	"encoding/json"
)

type Q struct {
	Context []C `json:"context"`
	Match   []M `json:"match"`
}

func (q *Q) UnmarshalJSON(b []byte) error {
	var target struct {
		C interface{} `json:"context"`
		M interface{} `json:"match"`
	}
	err := json.Unmarshal(b, &target)
	if err != nil {
		return err
	}
	if q.Context, err = decodeC(target.C); err != nil {
		return err
	}
	if q.Match, err = decodeM(target.M); err != nil {
		return err
	}
	return nil
}
