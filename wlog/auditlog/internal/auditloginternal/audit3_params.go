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
	"github.com/palantir/pkg/datetime"
	"github.com/palantir/pkg/uuid"
	"github.com/palantir/witchcraft-go-logging/conjure/witchcraft-logging-api/witchcraft/api/logging"
	"github.com/palantir/witchcraft-go-logging/wlog"
	wloginternal "github.com/palantir/witchcraft-go-logging/wlog/internal"
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

// Audit3Param marks an AuditParam as native to audit.3.
type Audit3Param struct {
	Param AuditParam
}

func Audit3ToParams(atLogTimeValuesProvider *atLogTimeValues, name string, result AuditResultType, inParams []Audit3Param) []wlog.Param {
	if atLogTimeValuesProvider == nil {
		atLogTimeValuesProvider = &atLogTimeValues{}
	}
	log := BuildAuditLogV3(atLogTimeValuesProvider, name, result, inParams)
	return []wlog.Param{audit3ObjectToParam(log, atLogTimeValuesProvider.getLogEntryIDValue())}
}

func Audit3ToDirectObject(name string, result AuditResultType, inParams []Audit3Param) (logging.AuditLogV3, wlog.Param) {
	values := &atLogTimeValues{}
	log := BuildAuditLogV3(values, name, result, inParams)
	if CanEmitAudit3Direct(log, values.logEntryIDValue) {
		return *log, nil
	}
	return *log, audit3ObjectToParam(log, values.getLogEntryIDValue())
}

// BuildAuditLogV3 allocates a fresh logging.AuditLogV3, populates required fields, applies the provided params, and
// returns the populated struct. Check CanEmitAudit3Direct before emitting the object directly.
//
// Hot-path callers may use populateAuditLogV3 with a pre-allocated (e.g., pooled) struct instead.
func BuildAuditLogV3(atLogTimeValuesProvider *atLogTimeValues, name string, result AuditResultType, inParams []Audit3Param) *logging.AuditLogV3 {
	log := new(logging.AuditLogV3)
	populateAuditLogV3(log, atLogTimeValuesProvider, name, result, inParams)
	return log
}

// populateAuditLogV3 populates a caller-provided logging.AuditLogV3 with the required fields, applies the provided
// params, and propagates the resolved logEntryId back into atLogTimeValuesProvider for a sibling v2 emission.
func populateAuditLogV3(log *logging.AuditLogV3, atLogTimeValuesProvider *atLogTimeValues, name string, result AuditResultType, inParams []Audit3Param) {
	if atLogTimeValuesProvider == nil {
		atLogTimeValuesProvider = &atLogTimeValues{}
	}
	log.Type = Audit3TypeValue
	log.Time = datetime.DateTime(atLogTimeValuesProvider.getOrComputeTimeValue())
	log.Name = name
	log.Result = logging.New_AuditResult(logging.AuditResult_Value(result))
	for _, p := range inParams {
		if p.Param.Audit3ParamFn != nil {
			p.Param.Audit3ParamFn(log)
		}
		if p.Param.LogEntryID != nil {
			atLogTimeValuesProvider.logEntryIDValue = p.Param.LogEntryID
		}
	}
	wloginternal.OmitEmpty(&log.Stack)
	wloginternal.OmitEmpty(&log.Service)
	wloginternal.OmitEmpty(&log.Environment)
	wloginternal.OmitEmpty(&log.UserAgent)
	wloginternal.OmitEmpty(&log.SourceOrigin)
	wloginternal.OmitEmpty(&log.Uid)
	wloginternal.OmitEmpty(&log.Sid)
	wloginternal.OmitEmpty(&log.TokenId)
	wloginternal.OmitEmpty(&log.OrgId)
	wloginternal.OmitEmpty(&log.TraceId)
	wloginternal.OmitEmpty(&log.Origin)
	// Resolve the default after parameters so supplied IDs do not generate discarded UUIDs.
	if log.LogEntryId == nil && atLogTimeValuesProvider.logEntryIDValue == nil {
		if atLogTimeValuesProvider.generatedLogEntryID == nil {
			atLogTimeValuesProvider.generatedLogEntryID = new(uuid.NewUUID())
		}
		log.LogEntryId = atLogTimeValuesProvider.generatedLogEntryID
	}
}

func audit3ObjectToParam(log *logging.AuditLogV3, rawLogEntryID *string) wlog.Param {
	return wlog.NewParam(func(entry wlog.LogEntry) {
		entry.StringValue(wlog.TypeKey, log.Type)
		entry.StringValue(wlog.TimeKey, log.Time.String())
		entry.OptionalStringValue(Audit3DeploymentKey, log.Deployment)
		entry.OptionalStringValue(Audit3HostKey, log.Host)
		entry.OptionalStringValue(Audit3ProductKey, log.Product)
		entry.OptionalStringValue(Audit3ProductVersionKey, log.ProductVersion)
		if log.Stack != nil {
			entry.OptionalStringValue(Audit3StackKey, *log.Stack)
		}
		if log.Service != nil {
			entry.OptionalStringValue(Audit3ServiceKey, *log.Service)
		}
		if log.Environment != nil {
			entry.OptionalStringValue(Audit3EnvironmentKey, *log.Environment)
		}
		if pt := log.ProducerType.String(); pt != "" {
			entry.StringValue(Audit3ProducerTypeKey, pt)
		}
		if len(log.Organizations) > 0 {
			orgs := make([]any, len(log.Organizations))
			for i, o := range log.Organizations {
				orgs[i] = Audit3Organization{ID: o.Id, Reason: o.Reason}
			}
			entry.ObjectListAppendValue(Audit3OrganizationsKey, orgs)
		}
		if log.EventId != (uuid.UUID{}) {
			entry.StringValue(Audit3EventIDKey, log.EventId.String())
		}
		if rawLogEntryID != nil {
			entry.StringValue(Audit3LogEntryIDKey, *rawLogEntryID)
		} else if log.LogEntryId != nil {
			entry.StringValue(Audit3LogEntryIDKey, log.LogEntryId.String())
		}
		if log.UserAgent != nil {
			entry.OptionalStringValue(Audit3UserAgentKey, *log.UserAgent)
		}
		if len(log.Categories) > 0 {
			entry.StringListAppendValue(Audit3CategoriesKey, log.Categories)
		}
		if len(log.Entities) > 0 {
			entry.ObjectListAppendValue(Audit3EntitiesKey, log.Entities)
		}
		if len(log.Users) > 0 {
			users := make([]any, len(log.Users))
			for i, u := range log.Users {
				users[i] = Audit3ContextualizedUser{
					UID:       string(u.Uid),
					UserName:  u.UserName,
					FirstName: u.FirstName,
					LastName:  u.LastName,
					Groups:    u.Groups,
					Realm:     u.Realm,
				}
			}
			entry.ObjectListAppendValue(Audit3UsersKey, users)
		}
		if len(log.Origins) > 0 {
			entry.StringListAppendValue(Audit3OriginsKey, log.Origins)
		}
		if log.SourceOrigin != nil {
			entry.OptionalStringValue(Audit3SourceOriginKey, *log.SourceOrigin)
		}
		if len(log.RequestFields) > 0 {
			entry.AnyMapValue(Audit3RequestFieldsKey, log.RequestFields)
		}
		if len(log.ResultFields) > 0 {
			entry.AnyMapValue(Audit3ResultFieldsKey, log.ResultFields)
		}
		if log.Uid != nil {
			entry.OptionalStringValue(wlog.UIDKey, string(*log.Uid))
		}
		if log.Sid != nil {
			entry.OptionalStringValue(wlog.SIDKey, string(*log.Sid))
		}
		if log.TokenId != nil {
			entry.OptionalStringValue(wlog.TokenIDKey, string(*log.TokenId))
		}
		if log.OrgId != nil {
			entry.OptionalStringValue(wlog.OrgIDKey, string(*log.OrgId))
		}
		if log.TraceId != nil {
			entry.OptionalStringValue(wlog.TraceIDKey, string(*log.TraceId))
		}
		if log.Origin != nil {
			entry.OptionalStringValue(Audit3OriginKey, *log.Origin)
		}
		entry.StringValue(Audit3NameKey, log.Name)
		entry.StringValue(Audit3ResultKey, string(log.Result.Value()))
	})
}

func Audit3Deployment(deployment string) Audit3Param {
	return audit3OnlyParamFn(func(log *logging.AuditLogV3) {
		log.Deployment = deployment
	})
}

func Audit3Host(host string) Audit3Param {
	return audit3OnlyParamFn(func(log *logging.AuditLogV3) {
		log.Host = host
	})
}

func Audit3Product(product string) Audit3Param {
	return audit3OnlyParamFn(func(log *logging.AuditLogV3) {
		log.Product = product
	})
}

func Audit3ProductVersion(productVersion string) Audit3Param {
	return audit3OnlyParamFn(func(log *logging.AuditLogV3) {
		log.ProductVersion = productVersion
	})
}

func Audit3Stack(stack string) Audit3Param {
	return audit3OnlyParamFn(func(log *logging.AuditLogV3) {
		log.Stack = &stack
	})
}

func Audit3Service(service string) Audit3Param {
	return audit3OnlyParamFn(func(log *logging.AuditLogV3) {
		log.Service = &service
	})
}

func Audit3Environment(environment string) Audit3Param {
	return audit3OnlyParamFn(func(log *logging.AuditLogV3) {
		log.Environment = &environment
	})
}

func Audit3ProducerType(producerType Audit3AuditProducerType) Audit3Param {
	return audit3OnlyParamFn(func(log *logging.AuditLogV3) {
		log.ProducerType = logging.New_AuditProducer(logging.AuditProducer_Value(producerType))
	})
}

func Audit3Organizations(organizations []Audit3Organization) Audit3Param {
	return audit3OnlyParamFn(func(log *logging.AuditLogV3) {
		for _, o := range organizations {
			log.Organizations = append(log.Organizations, logging.Organization{Id: o.ID, Reason: o.Reason})
		}
	})
}

func Audit3EventID(eventID string) Audit3Param {
	id, err := uuid.ParseUUID(eventID)
	if err != nil {
		id = uuid.UUID{}
	}
	return Audit3Param{
		Param: AuditParam{
			Audit2ParamFn: func(log *logging.AuditLogV2) {
				wloginternal.SetMapParam(&log.RequestParams, "_auditEventId", any(eventID))
			},
			Audit3ParamFn: func(log *logging.AuditLogV3) {
				log.EventId = id
			},
		},
	}
}

func Audit3LogEntryID(logEntryID string) Audit3Param {
	var id *uuid.UUID
	if parsed, err := uuid.ParseUUID(logEntryID); err == nil {
		id = &parsed
	}
	return Audit3Param{
		Param: AuditParam{
			LogEntryID: &logEntryID,
			Audit2ParamFn: func(log *logging.AuditLogV2) {
				wloginternal.SetMapParam(&log.RequestParams, "_auditLogEntryId", any(logEntryID))
			},
			Audit3ParamFn: func(log *logging.AuditLogV3) {
				log.LogEntryId = id
			},
		},
	}
}

func Audit3UserAgent(userAgent string) Audit3Param {
	return Audit3Param{
		Param: AuditParam{
			Audit2ParamFn: func(log *logging.AuditLogV2) {
				wloginternal.SetMapParam(&log.RequestParams, "_userAgent", any(userAgent))
			},
			Audit3ParamFn: func(log *logging.AuditLogV3) {
				log.UserAgent = &userAgent
			},
		},
	}
}

func Audit3Categories(categories []string) Audit3Param {
	return audit3OnlyParamFn(func(log *logging.AuditLogV3) {
		log.Categories = append(log.Categories, categories...)
	})
}

func Audit3Entities(entities []any) Audit3Param {
	return audit3OnlyParamFn(func(log *logging.AuditLogV3) {
		log.Entities = append(log.Entities, entities...)
	})
}

func Audit3Users(users []Audit3ContextualizedUser) Audit3Param {
	return audit3OnlyParamFn(func(log *logging.AuditLogV3) {
		for _, u := range users {
			var groups []string
			if len(u.Groups) > 0 {
				groups = u.Groups
			}
			log.Users = append(log.Users, logging.ContextualizedUser{
				Uid:       logging.UserId(u.UID),
				UserName:  u.UserName,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Groups:    groups,
				Realm:     u.Realm,
			})
		}
	})
}

func Audit3Origins(origins []string) Audit3Param {
	return Audit3Param{
		Param: AuditParam{
			Audit2ParamFn: func(log *logging.AuditLogV2) {
				wloginternal.SetMapParam(&log.RequestParams, "_forwardedOrigins", any(origins))
			},
			Audit3ParamFn: func(log *logging.AuditLogV3) {
				log.Origins = append(log.Origins, origins...)
			},
		},
	}
}

func Audit3SourceOrigin(sourceOrigin string) Audit3Param {
	return Audit3Param{
		Param: AuditParam{
			Audit2ParamFn: func(log *logging.AuditLogV2) {
				wloginternal.SetMapParam(&log.RequestParams, "_sourceOrigin", any(sourceOrigin))
			},
			Audit3ParamFn: func(log *logging.AuditLogV3) {
				log.SourceOrigin = &sourceOrigin
			},
		},
	}
}

func Audit3RequestField(key string, value any) Audit3Param {
	return Audit3Param{Param: AuditParam{
		Audit2ParamFn: func(log *logging.AuditLogV2) {
			wloginternal.SetMapParam(&log.RequestParams, key, value)
		},
		Audit3ParamFn: func(log *logging.AuditLogV3) {
			wloginternal.SetMapParam(&log.RequestFields, key, value)
		},
	}}
}

func Audit3RequestFields(requestFields map[string]any) Audit3Param {
	return Audit3Param{
		Param: AuditParam{
			Audit2ParamFn: func(log *logging.AuditLogV2) {
				wloginternal.SetMapParams(&log.RequestParams, requestFields)
			},
			Audit3ParamFn: func(log *logging.AuditLogV3) {
				wloginternal.SetMapParams(&log.RequestFields, requestFields)
			},
		},
	}
}

func Audit3ResultField(key string, value any) Audit3Param {
	return Audit3Param{Param: AuditParam{
		Audit2ParamFn: func(log *logging.AuditLogV2) {
			wloginternal.SetMapParam(&log.ResultParams, key, value)
		},
		Audit3ParamFn: func(log *logging.AuditLogV3) {
			wloginternal.SetMapParam(&log.ResultFields, key, value)
		},
	}}
}

func Audit3ResultFields(resultFields map[string]any) Audit3Param {
	return Audit3Param{
		Param: AuditParam{
			Audit2ParamFn: func(log *logging.AuditLogV2) {
				wloginternal.SetMapParams(&log.ResultParams, resultFields)
			},
			Audit3ParamFn: func(log *logging.AuditLogV3) {
				wloginternal.SetMapParams(&log.ResultFields, resultFields)
			},
		},
	}
}

func Audit3UID(uid string) Audit3Param {
	return Audit3Param{Param: uidParam(uid)}
}

func Audit3SID(sid string) Audit3Param {
	return Audit3Param{Param: sidParam(sid)}
}

func Audit3TokenID(tokenID string) Audit3Param {
	return Audit3Param{
		Param: AuditParam{
			Audit2ParamFn: func(log *logging.AuditLogV2) {
				log.TokenId = (*logging.TokenId)(&tokenID)
				wloginternal.SetMapParam(&log.RequestParams, "_tokenId", any(tokenID))
			},
			Audit3ParamFn: func(log *logging.AuditLogV3) {
				log.TokenId = (*logging.TokenId)(&tokenID)
			},
		},
	}
}

func Audit3OrgID(orgID string) Audit3Param {
	return Audit3Param{Param: orgIDParam(orgID)}
}

func Audit3TraceID(traceID string) Audit3Param {
	return Audit3Param{Param: traceIDParam(traceID)}
}

func Audit3Origin(origin string) Audit3Param {
	return Audit3Param{Param: originParam(origin)}
}

func audit3OnlyParamFn(fn func(*logging.AuditLogV3)) Audit3Param {
	return Audit3Param{
		Param: AuditParam{
			Audit3ParamFn: fn,
		},
	}
}
