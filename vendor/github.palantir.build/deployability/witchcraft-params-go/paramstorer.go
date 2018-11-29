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

package wparams

// ParamStorer is a type that stores safe and unsafe parameters.
type ParamStorer interface {
	SafeParams() map[string]interface{}
	UnsafeParams() map[string]interface{}
}

type mapParamStorer struct {
	safeParams   map[string]interface{}
	unsafeParams map[string]interface{}
}

func NewParamStorer(paramStorers ...ParamStorer) ParamStorer {
	safeParams := make(map[string]interface{})
	unsafeParams := make(map[string]interface{})
	for _, storer := range paramStorers {
		if storer == nil {
			continue
		}
		for k, v := range storer.SafeParams() {
			safeParams[k] = v
			delete(unsafeParams, k)
		}
		for k, v := range storer.UnsafeParams() {
			unsafeParams[k] = v
			delete(safeParams, k)
		}
	}
	return NewSafeAndUnsafeParamStorer(safeParams, unsafeParams)
}

func NewSafeParamStorer(safeParams map[string]interface{}) ParamStorer {
	return NewSafeAndUnsafeParamStorer(safeParams, nil)
}

func NewUnsafeParamStorer(unsafeParams map[string]interface{}) ParamStorer {
	return NewSafeAndUnsafeParamStorer(nil, unsafeParams)
}

func NewSafeAndUnsafeParamStorer(safeParams, unsafeParams map[string]interface{}) ParamStorer {
	storer := &mapParamStorer{
		safeParams:   make(map[string]interface{}),
		unsafeParams: make(map[string]interface{}),
	}
	for k, v := range safeParams {
		storer.safeParams[k] = v
	}
	for k, v := range unsafeParams {
		storer.unsafeParams[k] = v
		delete(storer.safeParams, k)
	}
	return storer
}

func (m *mapParamStorer) SafeParams() map[string]interface{} {
	return m.safeParams
}

func (m *mapParamStorer) UnsafeParams() map[string]interface{} {
	return m.unsafeParams
}
