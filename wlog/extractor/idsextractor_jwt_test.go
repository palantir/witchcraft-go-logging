// Copyright (c) 2022 Palantir Technologies. All rights reserved.
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

package extractor

import (
	"encoding/base64"
	puuid "github.com/palantir/pkg/uuid"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_uuidFromBase64StdEncodedString(t *testing.T) {
	in := "vp9kXVLgSem6MdsyknYV2w=="
	inBytes, err := base64.StdEncoding.DecodeString(in)
	require.NoError(t, err)
	expected := "be9f645d-52e0-49e9-ba31-db32927615db"
	got, err := uuid.FromBytes(inBytes)
	assert.NoError(t, err)
	assert.Equal(t, expected, got.String())
	assert.Equal(t, expected, uuidFromBase64StdEncodedString(in))

	_, err = puuid.ParseUUID(string(inBytes))
	// doesn't pass
	assert.NoError(t, err)

	_, err = puuid.ParseUUID(in)
	// doesn't pass
	assert.NoError(t, err)
}
