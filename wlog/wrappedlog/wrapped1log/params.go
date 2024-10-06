// Copyright (c) 2021 Palantir Technologies. All rights reserved.
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

package wrapped1log

import (
	"github.com/palantir/witchcraft-go-logging/wapi/logging"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

const (
	TypeValue = "wrapped.1"
)

type Param = wloginternal.Param[logging.WrappedLogV1]

type paramFunc = wloginternal.ParamFunc[logging.WrappedLogV1]

func defaultParam(payload logging.WrappedLogV1Payload) Param {
	return paramFunc(func(l *logging.WrappedLogV1) {
		l.Type = TypeValue
		l.Payload = payload
	})
}

func EntityName(name string) Param {
	return paramFunc(func(l *logging.WrappedLogV1) {
		l.EntityName = name
	})
}

func EntityVersion(version string) Param {
	return paramFunc(func(l *logging.WrappedLogV1) {
		l.EntityVersion = version
	})
}
