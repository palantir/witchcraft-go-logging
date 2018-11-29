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
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	// only modification is changing "3" to "4" here. Because the stack trace is always taken by the werror package,
	// omit one extra frame (caller should not see werror package as part of the output stack).
	n := runtime.Callers(4, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

// stack represents a stack of program counters.
type stack []uintptr

func (s *stack) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case state.Flag('+'):
			for _, pc := range *s {
				f := errors.Frame(pc)
				fmt.Fprintf(state, "\n%+v", f)
			}
		}
	}
}

func (s *stack) StackTrace() errors.StackTrace {
	f := make([]errors.Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = errors.Frame((*s)[i])
	}
	return f
}
