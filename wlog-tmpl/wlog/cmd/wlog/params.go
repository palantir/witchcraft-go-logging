// Copyright (c) 2020 Palantir Technologies. All rights reserved.
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

package wlog

import (
	"github.com/palantir/witchcraft-go-logging/wlog-tmpl/logentryformatter"
)

type Param interface {
	apply(*convertParams)
}

type paramFunc func(*convertParams)

func (f paramFunc) apply(p *convertParams) {
	f(p)
}

func getConvertParams(params ...Param) convertParams {
	rv := convertParams{
		formatters: make(map[logentryformatter.LogType]logentryformatter.Formatter),
		only:       make(map[logentryformatter.LogType]struct{}),
		exclude:    make(map[logentryformatter.LogType]struct{}),
	}
	for _, p := range params {
		if p == nil {
			continue
		}
		p.apply(&rv)
	}
	return rv
}

type convertParams struct {
	// if true, any error encountered while trying to parse a line as a witchcraft log entry will be written to the output.
	// Default behavior is to write such lines to the output unmodified.
	strict bool
	// map from log type to the formatter that will be used to format the type
	formatters map[logentryformatter.LogType]logentryformatter.Formatter
	// if the map is non-empty, then only log entries that match the types in this map will be outputted as log
	// entries (log entries with a type that is not in this map will be ignored).
	only map[logentryformatter.LogType]struct{}
	// any log entry type specified in this map will be excluded from output
	exclude map[logentryformatter.LogType]struct{}
}

func Strict(strict bool) Param {
	return paramFunc(func(p *convertParams) {
		p.strict = strict
	})
}

func Formatters(formatters map[logentryformatter.LogType]logentryformatter.Formatter) Param {
	return paramFunc(func(p *convertParams) {
		for k, v := range formatters {
			p.formatters[k] = v
		}
	})
}

func OnlyString(types ...string) Param {
	return paramFunc(func(p *convertParams) {
		for _, curr := range types {
			if curr == "" {
				continue
			}
			p.only[logentryformatter.LogType(curr)] = struct{}{}
		}
	})
}

func ExcludeString(types ...string) Param {
	return paramFunc(func(p *convertParams) {
		for _, curr := range types {
			if curr == "" {
				continue
			}
			p.exclude[logentryformatter.LogType(curr)] = struct{}{}
		}
	})
}
