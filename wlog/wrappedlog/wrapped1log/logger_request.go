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

package wrapped1log

import (
	"github.com/palantir/witchcraft-go-logging/wlog"
	"github.com/palantir/witchcraft-go-logging/wlog/extractor"
	"github.com/palantir/witchcraft-go-logging/wlog/reqlog/req2log"
)

type wrappedReq2Logger struct {
	name             string
	version          string
	idsExtractor     extractor.IDsFromRequest
	pathParamPerms   req2log.ParamPerms
	queryParamPerms  req2log.ParamPerms
	headerParamPerms req2log.ParamPerms

	logger wlog.Logger
}

func (l *wrappedReq2Logger) Request(r req2log.Request) {
	l.logger.Log(l.toRequestParams(r)...)
}

func (l *wrappedReq2Logger) PathParamPerms() req2log.ParamPerms {
	return l.pathParamPerms
}

func (l *wrappedReq2Logger) QueryParamPerms() req2log.ParamPerms {
	return l.queryParamPerms
}

func (l *wrappedReq2Logger) HeaderParamPerms() req2log.ParamPerms {
	return l.headerParamPerms
}

func (l *wrappedReq2Logger) toRequestParams(r req2log.Request) []wlog.Param {
	outParams := make([]wlog.Param, len(defaultTypeParam)+2)
	copy(outParams, defaultTypeParam)
	outParams[len(defaultTypeParam)] = wlog.NewParam(wrappedTypeParams(l.name, l.version).apply)
	outParams[len(defaultTypeParam)+1] = wlog.NewParam(req2PayloadParams(r, l.idsExtractor, l.pathParamPerms, l.queryParamPerms, l.headerParamPerms).apply)
	return outParams
}
