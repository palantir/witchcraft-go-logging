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

package audit3log

import (
	"os"
	"reflect"

	"github.com/palantir/witchcraft-go-logging/wlog"
)

const (
	TypeValue = "audit.3"

	OtherUIDsKey     = "otherUids"
	OriginKey        = "origin"
	NameKey          = "name"
	ResultKey        = "result"
	RequestParamsKey = "requestParams"
	ResultParamsKey  = "resultParams"
	//
	DeploymentKey     = "deployment"
	HostKey           = "host"
	ProductKey        = "product"
	ProductVersionKey = "productVersion"
	StackKey          = "stack"
	ServiceKey        = "service"
	EnvironmentKey    = "environment"
	ProducerTypeKey   = "producerType"
	OrganizationsKey  = "organizations"
	EventIdKey        = "eventId"
	UserAgentKey      = "userAgent"
	CategoriesKey     = "categories"
	// EntitiesKey       = "entities"
	UsersKey        = "users"
	OriginsKey      = "origins"
	SourceOriginKey = "sourceOrigin"
	// requestParams
	// resultParams
	// time
	// uid
	// sid
	// tokenId
	// traceId
	// origin
	// name
	// result
)

type Param interface {
	apply(entry wlog.LogEntry)
}

func ApplyParam(p Param, entry wlog.LogEntry) {
	if p == nil {
		return
	}
	p.apply(entry)
}

type paramFunc func(entry wlog.LogEntry)

func (f paramFunc) apply(entry wlog.LogEntry) {
	f(entry)
}

func auditRequiredParams(name string, resultType AuditResultType, deployment string, product string, productVersion string) Param {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}
	return paramFunc(func(logger wlog.LogEntry) {
		logger.StringValue(NameKey, name)
		logger.StringValue(ResultKey, string(resultType))
		logger.StringValue(DeploymentKey, deployment)
		logger.StringValue(HostKey, hostname)
		logger.StringValue(ProductKey, product)
		logger.StringValue(ProductVersionKey, productVersion)
	})
}

func Stack(stack string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(StackKey, stack)
	})
}

func Service(service string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(ServiceKey, service)
	})
}

func Environment(environment string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(EnvironmentKey, environment)
	})
}

func ProducerType(producerType AuditProducerType) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(ServiceKey, string(producerType))
	})
}

func Organizations(organizations ...AuditOrganizationType) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.ObjectValue(OrganizationsKey, organizations, reflect.TypeOf(organizations))
	})
}

func EventID(eventID string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(EventIdKey, eventID)
	})
}

func UserAgent(userAgent string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(UserAgentKey, userAgent)
	})
}

func Categories(categories ...string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.StringListValue(CategoriesKey, categories)
	})
}

// TODO entities

func Users(users ...AuditContextualizedUserType) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.ObjectValue(UsersKey, users, reflect.TypeOf(users))
	})
}

func Origins(origins ...string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.StringListValue(OriginsKey, origins)
	})
}

func SourceOrigin(sourceOrigin string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(SourceOriginKey, sourceOrigin)
	})
}

func UID(uid string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(wlog.UIDKey, uid)
	})
}

func SID(sid string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(wlog.SIDKey, sid)
	})
}

func TokenID(tokenID string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(wlog.TokenIDKey, tokenID)
	})
}

func TraceID(traceID string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(wlog.TraceIDKey, traceID)
	})
}

func OtherUIDs(otherUIDs ...string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.StringListValue(OtherUIDsKey, otherUIDs)
	})
}

func Origin(origin string) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.OptionalStringValue(OriginKey, origin)
	})
}

func RequestParam(key string, value interface{}) Param {
	return RequestParams(map[string]interface{}{
		key: value,
	})
}

func RequestParams(requestParams map[string]interface{}) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.AnyMapValue(RequestParamsKey, requestParams)
	})
}

// key="key1", value={"level": "Data", "payload": "val2" }
func ResultParam(key string, value interface{}) Param {
	return ResultParams(map[string]interface{}{
		key: value,
	})
}

// resultParams={ "key1": { "level": "UserInput", "payload": "val1" }, "key2": { "level": "Data", "payload": "val2" } }
func ResultParams(resultParams map[string]interface{}) Param {
	return paramFunc(func(entry wlog.LogEntry) {
		entry.AnyMapValue(ResultParamsKey, resultParams)
	})
}
