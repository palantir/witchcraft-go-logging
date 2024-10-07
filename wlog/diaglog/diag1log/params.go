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

package diag1log

import (
	"time"

	"github.com/palantir/pkg/datetime"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wtypes"
)

const (
	TypeValue = "diagnostic.1"
)

type Param = wloginternal.Param[wtypes.DiagnosticLogV1]

type paramFunc = wloginternal.ParamFunc[wtypes.DiagnosticLogV1]

func defaultParam(diagnostic wtypes.Diagnostic) Param {
	return paramFunc(func(l *wtypes.DiagnosticLogV1) {
		l.Type = TypeValue
		l.Time = datetime.DateTime(time.Now())
		l.Diagnostic = diagnostic
	})
}

func UnsafeParam(key string, value interface{}) Param {
	return paramFunc(func(l *wtypes.DiagnosticLogV1) {
		wloginternal.SetMapParam(&l.UnsafeParams, key, value)
	})
}

func UnsafeParams(unsafe map[string]interface{}) Param {
	return paramFunc(func(l *wtypes.DiagnosticLogV1) {
		wloginternal.SetMapParams(&l.UnsafeParams, unsafe)
	})
}
