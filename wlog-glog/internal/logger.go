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

package glogimpl

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/palantir/witchcraft-go-logging/wlog"
)

type gLogger struct{}

func (*gLogger) Log(params ...wlog.Param) {
	glog.Info(createGLogMsg("", params))
}

func (*gLogger) Debug(msg string, params ...wlog.Param) {
	glog.Info(createGLogMsg(msg, params))
}

func (*gLogger) Info(msg string, params ...wlog.Param) {
	glog.Info(createGLogMsg(msg, params))
}

func (*gLogger) Warn(msg string, params ...wlog.Param) {
	glog.Warning(createGLogMsg(msg, params))
}

func (*gLogger) Error(msg string, params ...wlog.Param) {
	glog.Error(createGLogMsg(msg, params))
}

func (*gLogger) SetLevel(level wlog.LogLevel) {
	// intentionally not implemented as glog uses the globally defined level
}

func (*gLogger) Enabled(_ wlog.LogLevel) bool {
	return true
}

func createGLogMsg(msg string, params []wlog.Param) string {
	entry := wlog.NewMapLogEntry()
	wlog.ApplyParams(entry, wlog.ParamsWithMessage(msg, params))

	// TODO: ignore/omit unsafe params?
	return strings.Join(paramsToLog(entry), ", ")
}

// paramsToLog returns the parameters to log as strings of the form "<key>: <value>".
func paramsToLog(entry wlog.MapLogEntry) []string {
	var params []string
	for k, v := range entry.StringValues() {
		params = append(params, fmt.Sprintf("%s: %s", k, v))
	}
	for k, v := range entry.SafeLongValues() {
		params = append(params, fmt.Sprintf("%s: %v", k, v))
	}
	for k, v := range entry.IntValues() {
		params = append(params, fmt.Sprintf("%s: %v", k, v))
	}
	for k, v := range entry.StringListValues() {
		params = append(params, fmt.Sprintf("%s: %v", k, v))
	}
	for k, v := range entry.StringMapValues() {
		params = append(params, fmt.Sprintf("%s: %v", k, v))
	}
	for k, v := range entry.AnyMapValues() {
		params = append(params, fmt.Sprintf("%s: %v", k, v))
	}
	for k, v := range entry.ObjectValues() {
		params = append(params, fmt.Sprintf("%s: %v", k, v))
	}
	return params
}
