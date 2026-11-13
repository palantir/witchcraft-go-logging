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

package svc1log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncPackage(t *testing.T) {
	for fn, want := range map[string]string{
		"main.main":                       "main",
		"a.b/c.F":                         "a.b/c",
		"a.b/c.(*T).M":                    "a.b/c",
		"a.b/c.F.func1":                   "a.b/c",
		"a.b/c.F[...]":                    "a.b/c",
		"a.b/c.F[go.shape.*a.b/d.T]":      "a.b/c",
		"a.b/c_test.TestF.func1":          "a.b/c_test",
		"github.com/x/wlog-zap.(*l).Info": "github.com/x/wlog-zap",
	} {
		assert.Equal(t, want, funcPackage(fn), fn)
	}
	assert.True(t, inPackage(wlogPackage+"/svclog/svc1log", wlogPackage))
	assert.True(t, inPackage(wlogPackage, wlogPackage))
	assert.False(t, inPackage(wlogPackage+"-zap", wlogPackage))
}
