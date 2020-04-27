package query

import "errors"

type CType uint8

const (
	OWNER CType = iota
	FILE
)

func decodeCType(s string) (CType, error) {
	switch s {
	case "o":
		fallthrough
	case "owner":
		return OWNER, nil
	case "f":
		fallthrough
	case "file":
		return FILE, nil
	default:
		return 0, errors.New("unrecognized Context Type")
	}
}

type CRestriction uint8

const (
	ALL CRestriction = iota
	OWNED
	VIEW
)

func decodeCRestriction(s string) (CRestriction, error) {
	switch s {
	case "":
		fallthrough
	case "all":
		return ALL, nil
	case "o":
		fallthrough
	case "owned":
		return OWNED, nil
	case "viewable":
		fallthrough
	case "v":
		return VIEW, nil
	default:
		return 0, errors.New("Unrecognized Context Restriction")
	}
}

type C struct {
	Type  CType        `json:"type"`
	ID    string       `json:"id"`
	Limit CRestriction `json:"only,omitempty"`
}

func decodeC(i interface{}) (contexts []C, err error) {
	switch v := i.(type) {
	case []interface{}:
		for _, ele := range v {
			var temp []C
			temp, err = decodeC(ele)
			if err != nil {
				return
			}
			contexts = append(contexts, temp...)
		}
	case map[string]interface{}:
		tstr, ok := v["type"].(string)
		if !ok {
			return nil, errors.New("Missing Context Type")
		}
		var t CType
		t, err = decodeCType(tstr)
		if err != nil {
			return
		}
		id, ok := v["id"].(string)
		if !ok {
			return nil, errors.New("Missing ID of context")
		}
		restriction := ALL
		if i, assigned := v["only"]; assigned {
			if r, ok := i.(string); ok {
				restriction, err = decodeCRestriction(r)
			}
		}
		contexts = append(contexts, C{
			Type:  t,
			ID:    id,
			Limit: restriction,
		})
	case string:
		contexts = append(contexts, C{
			Type:  OWNER,
			ID:    v,
			Limit: ALL,
		})
	default:
		return nil, errors.New("unrecognized Context Value")
	}
	return
}
