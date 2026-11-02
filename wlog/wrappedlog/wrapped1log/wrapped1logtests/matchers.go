// Copyright (c) 2026 Palantir Technologies. All rights reserved.
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

package wrapped1logtests

import (
	"maps"
	"time"
)

// Entry providers may add a timestamp to the envelope as well as the payload.
func withoutProviderTime(got map[string]any) map[string]any {
	if value, exists := got["time"]; exists {
		timestamp, ok := value.(string)
		if !ok {
			return got
		}
		if _, err := time.Parse(time.RFC3339Nano, timestamp); err != nil {
			return got
		}
	}
	got = maps.Clone(got)
	delete(got, "time")
	return got
}
