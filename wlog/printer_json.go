// Copyright (c) 2024 Palantir Technologies. All rights reserved.
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
	"encoding/json"
	"io"

	"github.com/palantir/witchcraft-go-logging/wtypes"
)

// jsonPrinter is a Printer that writes JSON-marshaled log objects to an output.
// It minimizes allocations by using a shared buffer pool.
type jsonPrinter struct {
	out io.Writer
}

func JSONPrinter(out io.Writer) Printer {
	return jsonPrinter{out: out}
}

func (p jsonPrinter) Print(log wtypes.LogType) error {
	enc := json.NewEncoder(p.out)
	enc.SetEscapeHTML(false)
	// json.Encoder appends a newline after each JSON object.
	return enc.Encode(log)
}
