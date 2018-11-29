package conjuretype

import (
	"encoding/json"
	"strings"
	"time"
)

type DateTime time.Time

func (d DateTime) String() string {
	return time.Time(d).Format(time.RFC3339Nano)
}

func (d *DateTime) UnmarshalJSON(b []byte) error {
	var val string
	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}

	// Conjure supports DateTime inputs that end with an optional zone identifier enclosed in square brackets
	// (for example, "2017-01-02T04:04:05.000000000+01:00[Europe/Berlin]"). If the input string ends in a ']' and
	// contains a '[', parse the string up to '['.
	if strings.HasSuffix(val, "]") {
		if openBracketIdx := strings.LastIndex(val, "["); openBracketIdx != -1 {
			val = val[:openBracketIdx]
		}
	}

	timeVal, err := time.Parse(time.RFC3339Nano, val)
	if err != nil {
		return err
	}
	*d = DateTime(timeVal)

	return nil
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
