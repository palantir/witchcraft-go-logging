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

package zeroimpl

import (
	"io"

	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/rs/zerolog"
)

func LoggerProvider() wlog.LoggerProvider {
	return &loggerProvider{}
}

type loggerProvider struct{}

func (lp *loggerProvider) NewLogger(w io.Writer) wlog.Logger {
	return &zeroLogger{
		logger: zerolog.New(w),
	}
}

func (lp *loggerProvider) NewLeveledLogger(w io.Writer, level wlog.LogLevel) wlog.LeveledLogger {
	return &zeroLogger{
		logger:         zerolog.New(w),
		AtomicLogLevel: wlog.NewAtomicLogLevel(level),
	}
}
