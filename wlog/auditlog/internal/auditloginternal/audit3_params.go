// Copyright (c) 2025 Palantir Technologies. All rights reserved.
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
	"slices"

	"github.com/palantir/pkg/datetime"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logentry"
	"github.com/palantir/witchcraft-go-logging/wlog/internal/logtypes"
)

const (
	Audit3TypeValue = "audit.3"

	Audit3DeploymentKey     = "deployment"
	Audit3HostKey           = "host"
	Audit3ProductKey        = "product"
	Audit3ProductVersionKey = "productVersion"
	Audit3StackKey          = "stack"
	Audit3ServiceKey        = "service"
	Audit3EnvironmentKey    = "environment"
	Audit3ProducerTypeKey   = "producerType"
	Audit3OrganizationsKey  = "organizations"
	Audit3EventIDKey        = "eventId"
	Audit3LogEntryIDKey     = "logEntryId"
	Audit3UserAgentKey      = "userAgent"
	Audit3CategoriesKey     = "categories"
	Audit3EntitiesKey       = "entities"
	Audit3UsersKey          = "users"
	Audit3OriginsKey        = "origins"
	Audit3SourceOriginKey   = "sourceOrigin"
	Audit3RequestFieldsKey  = "requestFields"
	Audit3ResultFieldsKey   = "resultFields"
	Audit3OriginKey         = "origin"
	Audit3NameKey           = "name"
	Audit3ResultKey         = "result"
)

type Audit3AuditProducerType string

const (
	Audit3AuditProducerServer Audit3AuditProducerType = "SERVER"
	Audit3AuditProducerClient Audit3AuditProducerType = "CLIENT"
)

type Audit3Organization struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}

type Audit3ContextualizedUser struct {
	UID       string   `json:"uid"`
	UserName  *string  `json:"userName,omitempty"`
	FirstName *string  `json:"firstName,omitempty"`
	LastName  *string  `json:"lastName,omitempty"`
	Groups    []string `json:"groups,omitempty"`
	Realm     *string  `json:"realm,omitempty"`
}

// Audit3ToParams builds parameters backed by a fresh audit.3 log that can outlive the log call.
func Audit3ToParams[P any](name string, result AuditResultType, params []P, get func(P) AuditParam) []wlog.Param {
	log := new(logtypes.AuditLogV3)
	populateAuditLogV3(log, &atLogTimeValues{}, name, result, params, get)
	return []wlog.Param{logentry.Param(log)}
}

// populateAuditLogV3 populates a caller-provided AuditLogV3 with the required fields, applies the provided params, and
// records the resolved logEntryId in values for a sibling v2 emission.
func populateAuditLogV3[P any](log *logtypes.AuditLogV3, values *atLogTimeValues, name string, result AuditResultType, params []P, get func(P) AuditParam) {
	log.Type = Audit3TypeValue
	log.Time = datetime.DateTime(values.getOrComputeTimeValue())
	log.Name = name
	log.Result = logtypes.AuditResult(result)
	log.LogEntryId = values.getOrGenerateLogEntryID()
	for _, p := range params {
		if fn := get(p).Audit3ParamFn; fn != nil {
			fn(log)
		}
	}
	values.logEntryID = log.LogEntryId
}

func Audit3Deployment(deployment string) AuditParam {
	return audit3OnlyParam(func(log *logtypes.AuditLogV3) {
		log.Deployment = deployment
	})
}

func Audit3Host(host string) AuditParam {
	return audit3OnlyParam(func(log *logtypes.AuditLogV3) {
		log.Host = host
	})
}

func Audit3Product(product string) AuditParam {
	return audit3OnlyParam(func(log *logtypes.AuditLogV3) {
		log.Product = product
	})
}

func Audit3ProductVersion(productVersion string) AuditParam {
	return audit3OnlyParam(func(log *logtypes.AuditLogV3) {
		log.ProductVersion = productVersion
	})
}

func Audit3Stack(stack string) AuditParam {
	return audit3OnlyParam(func(log *logtypes.AuditLogV3) {
		log.Stack = stack
	})
}

func Audit3Service(service string) AuditParam {
	return audit3OnlyParam(func(log *logtypes.AuditLogV3) {
		log.Service = service
	})
}

func Audit3Environment(environment string) AuditParam {
	return audit3OnlyParam(func(log *logtypes.AuditLogV3) {
		log.Environment = environment
	})
}

func Audit3ProducerType(producerType Audit3AuditProducerType) AuditParam {
	return audit3OnlyParam(func(log *logtypes.AuditLogV3) {
		log.ProducerType = string(producerType)
	})
}

func Audit3Organizations(organizations []Audit3Organization) AuditParam {
	return audit3OnlyParam(func(log *logtypes.AuditLogV3) {
		log.Organizations = slices.Grow(log.Organizations, len(organizations))
		for _, o := range organizations {
			log.Organizations = append(log.Organizations, logtypes.AuditLogV3Organization{Id: o.ID, Reason: o.Reason})
		}
	})
}

func Audit3EventID(eventID string) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			wloginternal.SetMapParam(&log.RequestParams, "_auditEventId", any(eventID))
		},
		Audit3ParamFn: func(log *logtypes.AuditLogV3) {
			log.EventId = eventID
		},
	}
}

func Audit3LogEntryID(logEntryID string) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			wloginternal.SetMapParam(&log.RequestParams, "_auditLogEntryId", any(logEntryID))
		},
		Audit3ParamFn: func(log *logtypes.AuditLogV3) {
			log.LogEntryId = logEntryID
		},
	}
}

func Audit3UserAgent(userAgent string) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			wloginternal.SetMapParam(&log.RequestParams, "_userAgent", any(userAgent))
		},
		Audit3ParamFn: func(log *logtypes.AuditLogV3) {
			log.UserAgent = userAgent
		},
	}
}

func Audit3Categories(categories []string) AuditParam {
	return audit3OnlyParam(func(log *logtypes.AuditLogV3) {
		log.Categories = append(log.Categories, categories...)
	})
}

func Audit3Entities(entities []any) AuditParam {
	return audit3OnlyParam(func(log *logtypes.AuditLogV3) {
		log.Entities = append(log.Entities, entities...)
	})
}

func Audit3Users(users []Audit3ContextualizedUser) AuditParam {
	return audit3OnlyParam(func(log *logtypes.AuditLogV3) {
		log.Users = slices.Grow(log.Users, len(users))
		for _, u := range users {
			log.Users = append(log.Users, logtypes.AuditLogV3ContextualizedUser{
				Uid:       u.UID,
				UserName:  u.UserName,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Groups:    u.Groups,
				Realm:     u.Realm,
			})
		}
	})
}

func Audit3Origins(origins []string) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			wloginternal.SetMapParam(&log.RequestParams, "_forwardedOrigins", any(origins))
		},
		Audit3ParamFn: func(log *logtypes.AuditLogV3) {
			log.Origins = append(log.Origins, origins...)
		},
	}
}

func Audit3SourceOrigin(sourceOrigin string) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			wloginternal.SetMapParam(&log.RequestParams, "_sourceOrigin", any(sourceOrigin))
		},
		Audit3ParamFn: func(log *logtypes.AuditLogV3) {
			log.SourceOrigin = sourceOrigin
		},
	}
}

func Audit3RequestField(key string, value any) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			wloginternal.SetMapParam(&log.RequestParams, key, value)
		},
		Audit3ParamFn: func(log *logtypes.AuditLogV3) {
			wloginternal.SetMapParam(&log.RequestFields, key, value)
		},
	}
}

func Audit3RequestFields(requestFields map[string]any) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			wloginternal.SetMapParams(&log.RequestParams, requestFields)
		},
		Audit3ParamFn: func(log *logtypes.AuditLogV3) {
			wloginternal.SetMapParams(&log.RequestFields, requestFields)
		},
	}
}

func Audit3ResultField(key string, value any) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			wloginternal.SetMapParam(&log.ResultParams, key, value)
		},
		Audit3ParamFn: func(log *logtypes.AuditLogV3) {
			wloginternal.SetMapParam(&log.ResultFields, key, value)
		},
	}
}

func Audit3ResultFields(resultFields map[string]any) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			wloginternal.SetMapParams(&log.ResultParams, resultFields)
		},
		Audit3ParamFn: func(log *logtypes.AuditLogV3) {
			wloginternal.SetMapParams(&log.ResultFields, resultFields)
		},
	}
}

func Audit3TokenID(tokenID string) AuditParam {
	return AuditParam{
		Audit2ParamFn: func(log *logtypes.AuditLogV2) {
			log.TokenId = tokenID
			wloginternal.SetMapParam(&log.RequestParams, "_tokenId", any(tokenID))
		},
		Audit3ParamFn: func(log *logtypes.AuditLogV3) {
			log.TokenId = tokenID
		},
	}
}

func audit3OnlyParam(fn func(*logtypes.AuditLogV3)) AuditParam {
	return AuditParam{Audit3ParamFn: fn}
}
