package conjuretype

import (
	"encoding/json"
)

type Rid string

func (r Rid) String() string {
	return string(r)
}

func (r *Rid) UnmarshalJSON(b []byte) error {
	var val string
	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}

	*r = Rid(val)
	return nil
}

func (r Rid) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}
