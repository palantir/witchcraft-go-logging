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

package wlogtmpl

import (
	"io"

	"github.com/palantir/pkg/bytesbuffers"
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog-tmpl/logentryformatter"
	"github.com/palantir/witchcraft-go-logging/wtypes"
)

type tmplPrinter struct {
	w   io.Writer
	cfg *Config

	delegate   wlog.LoggerProvider
	bufferPool bytesbuffers.Pool
}

func (l *tmplPrinter) Print(log wtypes.LogType) error {
	buf := l.bufferPool.Get()
	defer l.bufferPool.Put(buf)

	if err := l.delegate(buf).Print(log); err != nil {
		return err
	}

	out, err := logentryformatter.FormatLogLine(buf.String(), l.cfg.UnwrapperMap, l.cfg.FormatterMap, l.cfg.Only, l.cfg.Exclude)
	if err != nil {
		if !l.cfg.Strict {
			out = buf.String()
		}
		out = err.Error()
	}
	if len(out) > 0 {
		if out[len(out)-1] != '\n' {
			out += "\n"
		}
		_, err = l.w.Write([]byte(out))
		return err
	}
	return nil
}
