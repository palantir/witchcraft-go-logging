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
	"bytes"
	"io"

	"github.com/palantir/witchcraft-go-logging/wlog-tmpl/logentryformatter"
)

// formatWriter formats each JSON log line that the delegate logger writes.
type formatWriter struct {
	w   io.Writer
	cfg *Config
}

func (f *formatWriter) Write(p []byte) (int, error) {
	line := string(bytes.TrimRight(p, "\n"))
	out, err := logentryformatter.FormatLogLine(line, f.cfg.UnwrapperMap, f.cfg.FormatterMap, f.cfg.Only, f.cfg.Exclude)
	switch {
	case err != nil && f.cfg.Strict:
		out = err.Error()
	case err != nil:
		out = line
	case out == "":
		// Only and Exclude filter the line out.
		return len(p), nil
	}
	if _, err := io.WriteString(f.w, out+"\n"); err != nil {
		return 0, err
	}
	return len(p), nil
}
