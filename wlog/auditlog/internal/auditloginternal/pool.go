// Copyright (c) 2026 Palantir Technologies. All rights reserved.
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

package auditloginternal

import (
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
)

var (
	poolAuditLogV2 = wloginternal.NewPool[logging.AuditLogV2](func(log *logging.AuditLogV2) {
		wloginternal.ReleaseMap(&log.RequestParams)
		wloginternal.ReleaseMap(&log.ResultParams)
		*log = logging.AuditLogV2{}
	})

	poolAuditLogV3 = wloginternal.NewPool[logging.AuditLogV3](func(log *logging.AuditLogV3) {
		wloginternal.ReleaseMap(&log.RequestParams)
		wloginternal.ReleaseMap(&log.RequestFields)
		wloginternal.ReleaseMap(&log.ResultParams)
		wloginternal.ReleaseMap(&log.ResultFields)
		*log = logging.AuditLogV3{}
	})
)
