// Copyright (c) 2018 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conjuretype_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/palantir/witchcraft-go-logging/internal/conjuretype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var dateTimeJSONs = []struct {
	sec        int64
	zoneOffset int
	str        string
	json       string
}{
	{
		sec:  1483326245,
		str:  `2017-01-02T03:04:05Z`,
		json: `"2017-01-02T03:04:05Z"`,
	},
	{
		sec:  1483326245,
		str:  `2017-01-02T03:04:05Z`,
		json: `"2017-01-02T03:04:05.000Z"`,
	},
	{
		sec:  1483326245,
		str:  `2017-01-02T03:04:05Z`,
		json: `"2017-01-02T03:04:05.000000000Z"`,
	},
	{
		sec:        1483326245,
		zoneOffset: 3600,
		str:        `2017-01-02T04:04:05+01:00`,
		json:       `"2017-01-02T04:04:05.000000000+01:00"`,
	},
	{
		sec:        1483326245,
		zoneOffset: 7200,
		str:        `2017-01-02T05:04:05+02:00`,
		json:       `"2017-01-02T05:04:05.000000000+02:00"`,
	},
	{
		sec:        1483326245,
		zoneOffset: 3600,
		str:        `2017-01-02T04:04:05+01:00`,
		json:       `"2017-01-02T04:04:05.000000000+01:00[Europe/Berlin]"`,
	},
}

func TestDateTimeString(t *testing.T) {
	for i, currCase := range dateTimeJSONs {
		currDateTime := conjuretype.DateTime(time.Unix(currCase.sec, 0).In(time.FixedZone("", currCase.zoneOffset)))
		assert.Equal(t, currCase.str, currDateTime.String(), "Case %d", i)
	}
}

func TestDateTimeMarshal(t *testing.T) {
	for i, currCase := range dateTimeJSONs {
		currDateTime := conjuretype.DateTime(time.Unix(currCase.sec, 0).In(time.FixedZone("", currCase.zoneOffset)))
		bytes, err := json.Marshal(currDateTime)
		require.NoError(t, err, "Case %d", i)

		var unmarshaledFromMarshal conjuretype.DateTime
		err = json.Unmarshal(bytes, &unmarshaledFromMarshal)
		require.NoError(t, err, "Case %d", i)

		var unmarshaledFromCase conjuretype.DateTime
		err = json.Unmarshal([]byte(currCase.json), &unmarshaledFromCase)
		require.NoError(t, err, "Case %d", i)

		assert.Equal(t, unmarshaledFromCase, unmarshaledFromMarshal, "Case %d", i)
	}
}

func TestDateTimeUnmarshal(t *testing.T) {
	for i, currCase := range dateTimeJSONs {
		wantDateTime := time.Unix(currCase.sec, 0).UTC()
		if currCase.zoneOffset != 0 {
			wantDateTime = wantDateTime.In(time.FixedZone("", currCase.zoneOffset))
		}

		var gotDateTime conjuretype.DateTime
		err := json.Unmarshal([]byte(currCase.json), &gotDateTime)
		require.NoError(t, err, "Case %d", i)

		assert.Equal(t, wantDateTime, time.Time(gotDateTime), "Case %d", i)
	}
}

func TestDateTimeUnmarshalInvalid(t *testing.T) {
	for i, currCase := range []struct {
		input   string
		wantErr string
	}{
		{
			input:   `"foo"`,
			wantErr: "parsing time \"foo\" as \"2006-01-02T15:04:05.999999999Z07:00\": cannot parse \"foo\" as \"2006\"",
		},
		{
			input:   `"2017-01-02T04:04:05.000000000+01:00[Europe/Berlin"`,
			wantErr: "parsing time \"2017-01-02T04:04:05.000000000+01:00[Europe/Berlin\": extra text: [Europe/Berlin",
		},
		{
			input:   `"2017-01-02T04:04:05.000000000+01:00[[Europe/Berlin]]"`,
			wantErr: "parsing time \"2017-01-02T04:04:05.000000000+01:00[\": extra text: [",
		},
	} {
		var gotDateTime *conjuretype.DateTime
		err := json.Unmarshal([]byte(currCase.input), &gotDateTime)
		assert.EqualError(t, err, currCase.wantErr, "Case %d", i)
	}
}
