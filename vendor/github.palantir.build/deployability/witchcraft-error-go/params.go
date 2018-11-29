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

package werror

import (
	"github.palantir.build/deployability/witchcraft-params-go"
)

type Param interface {
	apply(*werror)
}

type param func(*werror)

func (p param) apply(e *werror) {
	p(e)
}

func SafeParam(key string, val interface{}) Param {
	return SafeParams(map[string]interface{}{key: val})
}

func SafeParams(vals map[string]interface{}) Param {
	return paramsHelper(vals, true)
}

func UnsafeParam(key string, val interface{}) Param {
	return UnsafeParams(map[string]interface{}{key: val})
}

func UnsafeParams(vals map[string]interface{}) Param {
	return paramsHelper(vals, false)
}

func paramsHelper(vals map[string]interface{}, safe bool) Param {
	return param(func(z *werror) {
		for k, v := range vals {
			z.params[k] = paramValue{
				safe:  safe,
				value: v,
			}
		}
	})
}

func SafeAndUnsafeParams(safe, unsafe map[string]interface{}) Param {
	return param(func(z *werror) {
		SafeParams(safe).apply(z)
		UnsafeParams(unsafe).apply(z)
	})
}

func Params(object wparams.ParamStorer) Param {
	return param(func(z *werror) {
		SafeParams(object.SafeParams()).apply(z)
		UnsafeParams(object.UnsafeParams()).apply(z)
	})
}
