package conjuretype

import (
	"encoding/json"
)

type Bearertoken string

func (b Bearertoken) String() string {
	return string(b)
}

func (b *Bearertoken) UnmarshalJSON(inputBytes []byte) error {
	var val string
	if err := json.Unmarshal(inputBytes, &val); err != nil {
		return err
	}

	*b = Bearertoken(val)
	return nil
}

func (b Bearertoken) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}
