package custom_type

import (
	"encoding/json"
	"strconv"
)

type JsonInt64 int64

func (i JsonInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(i), 10))
}

func (i *JsonInt64) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		value, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*i = JsonInt64(value)
		return nil
	}
	return json.Unmarshal(b, (*int64)(i))
}
