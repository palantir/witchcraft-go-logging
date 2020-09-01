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

package integration_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/palantir/godel/pkg/products/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cli string

func TestMain(m *testing.M) {
	var err error
	cli, err = products.Bin("wlog")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build wlog: %v", err)
	}
	os.Exit(m.Run())
}

func TestStrictFlag(t *testing.T) {

	for i, currCase := range []struct {
		input string
		want  string
	}{
		{
			input: `{"type":"service.1","time":"2017-04-12T17:41:07.744Z","level":"ERROR","message":"Error handling request {}, safe: {}","origin":"com.palantir.remoting2.servers.jersey.JsonExceptionMapper","thread":"qtp1360518503-16","params":{"request": "/foo"},"uid":null,"sid":null,"tokenId":null,"traceId":"fa4f6a37ac662fbd","stacktrace":"java.lang.NullPointerException: {throwableMessage}\n\tcom.palantir.edu.profiles.resource.ProfileResource.getUserId(ProfileResource.java:36)\n\tcom.palantir.edu.profiles.resource.ProfileResource.getProfile(ProfileResource.java:32)\n\tsun.reflect.NativeMethodAccessorImpl.invoke0(Native Method)\n\tsun.reflect.NativeMethodAccessorImpl.invoke(NativeMethodAccessorImpl.java:62)\n\tsun.reflect.DelegatingMethodAccessorImpl.invoke(DelegatingMethodAccessorImpl.java:43)\n\tjava.lang.reflect.Method.invoke(Method.java:498)\n\torg.glassfish.jersey.server.model.internal.ResourceMethodInvocationHandlerFactory$1.invoke(ResourceMethodInvocationHandlerFactory.java:81)\n\torg.glassfish.jersey.server.model.internal.AbstractJavaResourceMethodDispatcher$1.run(AbstractJavaResourceMethodDispatcher.java:144)\n\torg.glassfish.jersey.server.model.internal.AbstractJavaResourceMethodDispatcher.invoke(AbstractJavaResourceMethodDispatcher.java:161)\n\torg.glassfish.jersey.server.model.internal.JavaResourceMethodDispatcherProvider$TypeOutInvoker.doDispatch(JavaResourceMethodDispatcherProvider.java:205)\n\torg.glassfish.jersey.server.model.internal.AbstractJavaResourceMethodDispatcher.dispatch(AbstractJavaResourceMethodDispatcher.java:99)\n\torg.glassfish.jersey.server.model.ResourceMethodInvoker.invoke(ResourceMethodInvoker.java:389)\n\torg.glassfish.jersey.server.model.ResourceMethodInvoker.apply(ResourceMethodInvoker.java:347)\n\torg.glassfish.jersey.server.model.ResourceMethodInvoker.apply(ResourceMethodInvoker.java:102)\n\torg.glassfish.jersey.server.ServerRuntime$2.run(ServerRuntime.java:326)\n\torg.glassfish.jersey.internal.Errors$1.call(Errors.java:271)\n\torg.glassfish.jersey.internal.Errors$1.call(Errors.java:267)\n\torg.glassfish.jersey.internal.Errors.process(Errors.java:315)\n\torg.glassfish.jersey.internal.Errors.process(Errors.java:297)\n\torg.glassfish.jersey.internal.Errors.process(Errors.java:267)\n\torg.glassfish.jersey.process.internal.RequestScope.runInScope(RequestScope.java:317)\n\torg.glassfish.jersey.server.ServerRuntime.process(ServerRuntime.java:305)\n\torg.glassfish.jersey.server.ApplicationHandler.handle(ApplicationHandler.java:1154)\n\torg.glassfish.jersey.servlet.WebComponent.serviceImpl(WebComponent.java:473)\n\torg.glassfish.jersey.servlet.WebComponent.service(WebComponent.java:427)\n\torg.glassfish.jersey.servlet.ServletContainer.service(ServletContainer.java:388)\n\torg.glassfish.jersey.servlet.ServletContainer.service(ServletContainer.java:341)\n\torg.glassfish.jersey.servlet.ServletContainer.service(ServletContainer.java:228)\n\torg.eclipse.jetty.servlet.ServletHolder.handle(ServletHolder.java:838)\n\torg.eclipse.jetty.servlet.ServletHandler.doHandle(ServletHandler.java:543)\n\torg.eclipse.jetty.server.handler.ScopedHandler.nextHandle(ScopedHandler.java:188)\n\torg.eclipse.jetty.server.session.SessionHandler.doHandle(SessionHandler.java:1584)\n\torg.eclipse.jetty.server.handler.ScopedHandler.nextHandle(ScopedHandler.java:188)\n\torg.eclipse.jetty.server.handler.ContextHandler.doHandle(ContextHandler.java:1228)\n\torg.eclipse.jetty.server.handler.ScopedHandler.nextScope(ScopedHandler.java:168)\n\torg.eclipse.jetty.servlet.ServletHandler.doScope(ServletHandler.java:481)\n\torg.eclipse.jetty.server.session.SessionHandler.doScope(SessionHandler.java:1553)\n\torg.eclipse.jetty.server.handler.ScopedHandler.nextScope(ScopedHandler.java:166)\n\torg.eclipse.jetty.server.handler.ContextHandler.doScope(ContextHandler.java:1130)\n\torg.eclipse.jetty.server.handler.ScopedHandler.handle(ScopedHandler.java:141)\n\torg.eclipse.jetty.server.handler.RequestLogHandler.handle(RequestLogHandler.java:56)\n\torg.eclipse.jetty.server.handler.gzip.GzipHandler.handle(GzipHandler.java:530)\n\torg.eclipse.jetty.server.handler.HandlerWrapper.handle(HandlerWrapper.java:132)\n\torg.eclipse.jetty.server.Server.handle(Server.java:564)\n\torg.eclipse.jetty.server.HttpChannel.handle(HttpChannel.java:318)\n\torg.eclipse.jetty.server.HttpConnection.onFillable(HttpConnection.java:251)\n\torg.eclipse.jetty.io.AbstractConnection$ReadCallback.succeeded(AbstractConnection.java:279)\n\torg.eclipse.jetty.io.FillInterest.fillable(FillInterest.java:112)\n\torg.eclipse.jetty.io.ssl.SslConnection.onFillable(SslConnection.java:261)\n\torg.eclipse.jetty.io.ssl.SslConnection$3.succeeded(SslConnection.java:150)\n\torg.eclipse.jetty.io.FillInterest.fillable(FillInterest.java:112)\n\torg.eclipse.jetty.io.ChannelEndPoint$2.run(ChannelEndPoint.java:124)\n\torg.eclipse.jetty.util.thread.QueuedThreadPool.runJob(QueuedThreadPool.java:672)\n\torg.eclipse.jetty.util.thread.QueuedThreadPool$2.run(QueuedThreadPool.java:590)\n\tjava.lang.Thread.run(Thread.java:745)\n","unsafeParams":{"0":"8df8ace6-a068-4094-a7ff-0273469302f5","throwableMessage":"Message"}}`,
			want: `ERROR [2017-04-12T17:41:07.744Z] com.palantir.remoting2.servers.jersey.JsonExceptionMapper: Error handling request 8df8ace6-a068-4094-a7ff-0273469302f5, safe: {} (request: /foo) (0: 8df8ace6-a068-4094-a7ff-0273469302f5, throwableMessage: Message)
java.lang.NullPointerException: Message
	com.palantir.edu.profiles.resource.ProfileResource.getUserId(ProfileResource.java:36)
	com.palantir.edu.profiles.resource.ProfileResource.getProfile(ProfileResource.java:32)
	sun.reflect.NativeMethodAccessorImpl.invoke0(Native Method)
	sun.reflect.NativeMethodAccessorImpl.invoke(NativeMethodAccessorImpl.java:62)
	sun.reflect.DelegatingMethodAccessorImpl.invoke(DelegatingMethodAccessorImpl.java:43)
	java.lang.reflect.Method.invoke(Method.java:498)
	org.glassfish.jersey.server.model.internal.ResourceMethodInvocationHandlerFactory$1.invoke(ResourceMethodInvocationHandlerFactory.java:81)
	org.glassfish.jersey.server.model.internal.AbstractJavaResourceMethodDispatcher$1.run(AbstractJavaResourceMethodDispatcher.java:144)
	org.glassfish.jersey.server.model.internal.AbstractJavaResourceMethodDispatcher.invoke(AbstractJavaResourceMethodDispatcher.java:161)
	org.glassfish.jersey.server.model.internal.JavaResourceMethodDispatcherProvider$TypeOutInvoker.doDispatch(JavaResourceMethodDispatcherProvider.java:205)
	org.glassfish.jersey.server.model.internal.AbstractJavaResourceMethodDispatcher.dispatch(AbstractJavaResourceMethodDispatcher.java:99)
	org.glassfish.jersey.server.model.ResourceMethodInvoker.invoke(ResourceMethodInvoker.java:389)
	org.glassfish.jersey.server.model.ResourceMethodInvoker.apply(ResourceMethodInvoker.java:347)
	org.glassfish.jersey.server.model.ResourceMethodInvoker.apply(ResourceMethodInvoker.java:102)
	org.glassfish.jersey.server.ServerRuntime$2.run(ServerRuntime.java:326)
	org.glassfish.jersey.internal.Errors$1.call(Errors.java:271)
	org.glassfish.jersey.internal.Errors$1.call(Errors.java:267)
	org.glassfish.jersey.internal.Errors.process(Errors.java:315)
	org.glassfish.jersey.internal.Errors.process(Errors.java:297)
	org.glassfish.jersey.internal.Errors.process(Errors.java:267)
	org.glassfish.jersey.process.internal.RequestScope.runInScope(RequestScope.java:317)
	org.glassfish.jersey.server.ServerRuntime.process(ServerRuntime.java:305)
	org.glassfish.jersey.server.ApplicationHandler.handle(ApplicationHandler.java:1154)
	org.glassfish.jersey.servlet.WebComponent.serviceImpl(WebComponent.java:473)
	org.glassfish.jersey.servlet.WebComponent.service(WebComponent.java:427)
	org.glassfish.jersey.servlet.ServletContainer.service(ServletContainer.java:388)
	org.glassfish.jersey.servlet.ServletContainer.service(ServletContainer.java:341)
	org.glassfish.jersey.servlet.ServletContainer.service(ServletContainer.java:228)
	org.eclipse.jetty.servlet.ServletHolder.handle(ServletHolder.java:838)
	org.eclipse.jetty.servlet.ServletHandler.doHandle(ServletHandler.java:543)
	org.eclipse.jetty.server.handler.ScopedHandler.nextHandle(ScopedHandler.java:188)
	org.eclipse.jetty.server.session.SessionHandler.doHandle(SessionHandler.java:1584)
	org.eclipse.jetty.server.handler.ScopedHandler.nextHandle(ScopedHandler.java:188)
	org.eclipse.jetty.server.handler.ContextHandler.doHandle(ContextHandler.java:1228)
	org.eclipse.jetty.server.handler.ScopedHandler.nextScope(ScopedHandler.java:168)
	org.eclipse.jetty.servlet.ServletHandler.doScope(ServletHandler.java:481)
	org.eclipse.jetty.server.session.SessionHandler.doScope(SessionHandler.java:1553)
	org.eclipse.jetty.server.handler.ScopedHandler.nextScope(ScopedHandler.java:166)
	org.eclipse.jetty.server.handler.ContextHandler.doScope(ContextHandler.java:1130)
	org.eclipse.jetty.server.handler.ScopedHandler.handle(ScopedHandler.java:141)
	org.eclipse.jetty.server.handler.RequestLogHandler.handle(RequestLogHandler.java:56)
	org.eclipse.jetty.server.handler.gzip.GzipHandler.handle(GzipHandler.java:530)
	org.eclipse.jetty.server.handler.HandlerWrapper.handle(HandlerWrapper.java:132)
	org.eclipse.jetty.server.Server.handle(Server.java:564)
	org.eclipse.jetty.server.HttpChannel.handle(HttpChannel.java:318)
	org.eclipse.jetty.server.HttpConnection.onFillable(HttpConnection.java:251)
	org.eclipse.jetty.io.AbstractConnection$ReadCallback.succeeded(AbstractConnection.java:279)
	org.eclipse.jetty.io.FillInterest.fillable(FillInterest.java:112)
	org.eclipse.jetty.io.ssl.SslConnection.onFillable(SslConnection.java:261)
	org.eclipse.jetty.io.ssl.SslConnection$3.succeeded(SslConnection.java:150)
	org.eclipse.jetty.io.FillInterest.fillable(FillInterest.java:112)
	org.eclipse.jetty.io.ChannelEndPoint$2.run(ChannelEndPoint.java:124)
	org.eclipse.jetty.util.thread.QueuedThreadPool.runJob(QueuedThreadPool.java:672)
	org.eclipse.jetty.util.thread.QueuedThreadPool$2.run(QueuedThreadPool.java:590)
	java.lang.Thread.run(Thread.java:745)

`,
		},
		// output always has a trailing newline
		{
			input: `{}`,
			want:  "Log line JSON \"{}\" does not have a \"type\" key so its log type cannot be determined\n",
		},
		{
			input: "{}\n",
			want:  "Log line JSON \"{}\" does not have a \"type\" key so its log type cannot be determined\n",
		},
		// strict spec doesn't allow blank lines
		{
			input: "\n\n\n",
			want:  "Log line \"\" is not valid JSON\nLog line \"\" is not valid JSON\nLog line \"\" is not valid JSON\n",
		},
		{
			input: `{"type":"service.1","message":"foo\nbar"}`,
			want:  "      [0001-01-01T00:00:00Z]     foo\nbar\n",
		},
		{
			input: `{`,
			want:  "Failed to parse log line \"{\" as JSON: unexpected end of JSON input\n",
		},
		{
			input: `{"type":"service.FOOBARBAZ","message":"hi"}`,
			want:  "Skipping unknown log line type: service.FOOBARBAZ\n",
		},
		{
			input: `{"type":"service.1","level":"ERROR"}
			        {"type":"service.1","level":"INFO"}`,
			want: "ERROR [0001-01-01T00:00:00Z]     \nINFO  [0001-01-01T00:00:00Z]     \n",
		},
		{
			input: `{"type":"request.1","time":"2017-05-25T01:20:57.126Z","method":"GET","protocol":"HTTP/1.1","path":"/readiness","pathParams":{},"queryParams":{},"headerParams":{"Accept":"*/*","Host":"myhost.mydomain.com:8901","User-Agent":"curl/7.40.0"},"status":200,"requestSize":"0","responseSize":"2","duration":1000,"uid":null,"sid":null,"tokenId":null,"traceId":"64181265c6a8daf6","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}
{"type":"request.1","time":"2017-05-25T01:21:12.125Z","method":"GET","protocol":"HTTP/1.1","path":"/readiness","pathParams":{},"queryParams":{},"headerParams":{"Accept":"*/*","Host":"myhost.mydomain.com:8901","User-Agent":"curl/7.40.0"},"status":200,"requestSize":"0","responseSize":"2","duration":2000,"uid":null,"sid":null,"tokenId":null,"traceId":"6b272ef9831f17a5","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}`,
			want: `[2017-05-25T01:20:57.126Z] "GET /readiness HTTP/1.1" 200 2 1000
[2017-05-25T01:21:12.125Z] "GET /readiness HTTP/1.1" 200 2 2000
`,
		},
	} {
		cmd := exec.Command(cli, "--strict")
		output := runWithStdin(t, cmd, currCase.input)
		assert.Equal(t, currCase.want, output, "Case %d: %s %s", i, currCase.input, output)
	}
}

func TestWLogStdin(t *testing.T) {

	for i, currCase := range []struct {
		input string
		want  string
	}{
		{
			input: `{"type":"service.1","message":"foo\nbar"}
`,
			want: "      [0001-01-01T00:00:00Z]     foo\nbar\n",
		},
		{
			input: "foo\n",
			want:  "foo\n",
		},
		{
			input: "{}\n",
			want:  "{}\n",
		},
		{
			input: "{\n",
			want:  "{\n",
		},
		// preserve empty lines
		{
			input: "a\n\nb\n\n\n",
			want:  "a\n\nb\n\n\n",
		},
		// for e.g. tail -F var/log/*log scenarios
		{
			input: `==> var/log/request.log <==
{"type":"request.1","time":"2017-05-25T07:27:27.123Z","method":"GET","protocol":"HTTP/1.1","path":"/readiness","pathParams":{},"queryParams":{},"headerParams":{"Accept":"*/*","Host":"myhost.mydomain.com:8901","User-Agent":"curl/7.40.0"},"status":200,"requestSize":"0","responseSize":"2","duration":1000,"uid":null,"sid":null,"tokenId":null,"traceId":"58074cb9b5934d07","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}

==> var/log/event.log <==
{"type":"event.nomatch","time":"2017-05-25T07:27:31.81Z","eventName":"com.palantir.tritium.metrics.MetricRegistries.reservoir.type","eventType":"gauge","values":{},"uid":null,"sid":null,"tokenId":null,"unsafeParams":{"value":"org.mpierce.metrics.reservoir.hdrhistogram.HdrHistogramReservoir"}}`,
			want: `==> var/log/request.log <==
[2017-05-25T07:27:27.123Z] "GET /readiness HTTP/1.1" 200 2 1000

==> var/log/event.log <==
{"type":"event.nomatch","time":"2017-05-25T07:27:31.81Z","eventName":"com.palantir.tritium.metrics.MetricRegistries.reservoir.type","eventType":"gauge","values":{},"uid":null,"sid":null,"tokenId":null,"unsafeParams":{"value":"org.mpierce.metrics.reservoir.hdrhistogram.HdrHistogramReservoir"}}
`,
		},
	} {
		cmd := exec.Command(cli)
		output := runWithStdin(t, cmd, currCase.input)
		assert.Equal(t, currCase.want, output, "Case %d: %s", i, currCase.input)
	}
}

func TestInputFlag(t *testing.T) {

	for i, currCase := range []struct {
		input string
		want  string
	}{
		{
			input: `{"type":"service.1","message":"foo\nbar"}`,
			want:  "      [0001-01-01T00:00:00Z]     foo\nbar\n",
		},
		{
			input: `{"type":"request.1","time":"2017-05-25T01:20:57.126Z","method":"GET","protocol":"HTTP/1.1","path":"/readiness","pathParams":{},"queryParams":{},"headerParams":{"Accept":"*/*","Host":"myhost.mydomain.com:8901","User-Agent":"curl/7.40.0"},"status":200,"requestSize":"0","responseSize":"2","duration":1000,"uid":null,"sid":null,"tokenId":null,"traceId":"64181265c6a8daf6","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}
{"type":"request.1","time":"2017-05-25T01:21:12.125Z","method":"GET","protocol":"HTTP/1.1","path":"/readiness","pathParams":{},"queryParams":{},"headerParams":{"Accept":"*/*","Host":"myhost.mydomain.com:8901","User-Agent":"curl/7.40.0"},"status":200,"requestSize":"0","responseSize":"2","duration":2000,"uid":null,"sid":null,"tokenId":null,"traceId":"6b272ef9831f17a5","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}`,
			want: `[2017-05-25T01:20:57.126Z] "GET /readiness HTTP/1.1" 200 2 1000
[2017-05-25T01:21:12.125Z] "GET /readiness HTTP/1.1" 200 2 2000
`,
		},
		{
			input: `foo`,
			want:  "foo\n",
		},
	} {
		cmd := exec.Command(cli, "--input", currCase.input)
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Case %d\nOutput: %s", i, string(output))
		assert.Equal(t, currCase.want, string(output), "Case %d: %s", i, currCase.input)
	}
}

func TestOnlyExcludeFlags(t *testing.T) {

	for i, currCase := range []struct {
		name         string
		input        string
		onlyFlags    []string
		excludeFlags []string
		want         string
	}{
		{
			name: "only includes only specified type",
			input: `{"type":"service.1","message":"foo\nbar"}
{"type":"request.1","time":"2017-05-25T01:20:57.126Z","method":"GET","protocol":"HTTP/1.1","path":"/readiness","pathParams":{},"queryParams":{},"headerParams":{"Accept":"*/*","Host":"myhost.mydomain.com:8901","User-Agent":"curl/7.40.0"},"status":200,"requestSize":"0","responseSize":"2","duration":1000,"uid":null,"sid":null,"tokenId":null,"traceId":"64181265c6a8daf6","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}
{"type":"custom-type"}`,
			onlyFlags: []string{
				"service.1",
			},
			want: "      [0001-01-01T00:00:00Z]     foo\nbar\n",
		},
		{
			name: "only can specify multiple types",
			input: `{"type":"service.1","message":"foo\nbar"}
{"type":"request.1","time":"2017-05-25T01:20:57.126Z","method":"GET","protocol":"HTTP/1.1","path":"/readiness","pathParams":{},"queryParams":{},"headerParams":{"Accept":"*/*","Host":"myhost.mydomain.com:8901","User-Agent":"curl/7.40.0"},"status":200,"requestSize":"0","responseSize":"2","duration":1000,"uid":null,"sid":null,"tokenId":null,"traceId":"64181265c6a8daf6","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}
{"type":"custom-type"}`,
			onlyFlags: []string{
				"service.1",
				"request.1",
			},
			want: "      [0001-01-01T00:00:00Z]     foo\nbar\n[2017-05-25T01:20:57.126Z] \"GET /readiness HTTP/1.1\" 200 2 1000\n",
		},
		{
			name: "exclude excludes specified types",
			input: `{"type":"service.1","message":"foo\nbar"}
{"type":"request.1","time":"2017-05-25T01:20:57.126Z","method":"GET","protocol":"HTTP/1.1","path":"/readiness","pathParams":{},"queryParams":{},"headerParams":{"Accept":"*/*","Host":"myhost.mydomain.com:8901","User-Agent":"curl/7.40.0"},"status":200,"requestSize":"0","responseSize":"2","duration":1000,"uid":null,"sid":null,"tokenId":null,"traceId":"64181265c6a8daf6","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}
{"type":"custom-type"}`,
			excludeFlags: []string{
				"request.1",
				"custom-type",
			},
			want: "      [0001-01-01T00:00:00Z]     foo\nbar\n",
		},
		{
			name: "only and exclude can both be specified",
			input: `{"type":"service.1","message":"foo\nbar"}
{"type":"request.1","time":"2017-05-25T01:20:57.126Z","method":"GET","protocol":"HTTP/1.1","path":"/readiness","pathParams":{},"queryParams":{},"headerParams":{"Accept":"*/*","Host":"myhost.mydomain.com:8901","User-Agent":"curl/7.40.0"},"status":200,"requestSize":"0","responseSize":"2","duration":1000,"uid":null,"sid":null,"tokenId":null,"traceId":"64181265c6a8daf6","unsafeParams":{"pathParams":{},"queryParams":{},"headerParams":{}}}
{"type":"custom-type"}`,
			onlyFlags: []string{
				"service.1",
				"request.1",
			},
			excludeFlags: []string{
				"request.1",
			},
			want: "      [0001-01-01T00:00:00Z]     foo\nbar\n",
		},
	} {
		args := []string{"--input", currCase.input}
		for _, curr := range currCase.onlyFlags {
			args = append(args, "--only", curr)
		}
		for _, curr := range currCase.excludeFlags {
			args = append(args, "--exclude", curr)
		}
		cmd := exec.Command(cli, args...)
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Case %d: %s\nOutput: %s", i, currCase.name, string(output))
		assert.Equal(t, currCase.want, string(output), "Case %d (%s): %s", i, currCase.name, currCase.input)
	}
}

func TestColorFlag(t *testing.T) {
	for i, currCase := range []struct {
		name      string
		input     string
		flags     []string
		want      string
		wantError string
	}{
		{
			name:  "default output to non-terminal stream does not output color",
			input: `{"type":"service.1","level":"WARN","message":"foo"}`,
			want:  "WARN  [0001-01-01T00:00:00Z]     foo\n",
		},
		{
			name:  "--color=false flag disables color output",
			input: `{"type":"service.1","level":"WARN","message":"foo"}`,
			flags: []string{
				"--color=false",
			},
			want: "WARN  [0001-01-01T00:00:00Z]     foo\n",
		},
		{
			name:  "--color flag enables color output",
			input: `{"type":"service.1","level":"WARN","message":"foo"}`,
			flags: []string{
				"--color",
			},
			want: "\x1b[33mWARN  [0001-01-01T00:00:00Z]     foo\x1b[0m\n",
		},
		{
			name:  "INFO has no color",
			input: `{"type":"service.1","level":"INFO","message":"foo"}`,
			flags: []string{
				"--color",
			},
			want: "INFO  [0001-01-01T00:00:00Z]     foo\n",
		},
	} {
		args := []string{"--input", currCase.input}
		args = append(args, currCase.flags...)
		cmd := exec.Command(cli, args...)
		output, err := cmd.CombinedOutput()

		if currCase.wantError == "" {
			require.NoError(t, err, "Case %d: %s\nOutput: %s", i, currCase.name, string(output))
			assert.Equal(t, currCase.want, string(output), "Case %d (%s): %s", i, currCase.name, currCase.input)
		} else {
			require.Error(t, err, fmt.Sprintf("Case %d: %s\nOutput: %s", i, currCase.name, string(output)))
			assert.Equal(t, currCase.wantError, string(output), "Case %d (%s): %s", i, currCase.name, currCase.input)
		}
	}
}

func TestTemplatesFlag(t *testing.T) {
	for i, currCase := range []struct {
		name          string
		input         string
		templateFlags []string
		want          string
	}{
		{
			name:  "template for service.1 log",
			input: `{"type":"service.1","message":"foo\nbar"}`,
			templateFlags: []string{
				"service.1:Output {{.Type}}",
			},
			want: "Output service.1\n",
		},
		{
			name: "template for multiple log formats",
			input: `{"type":"service.1","message":"foo\nbar"}
{"type":"request.1","method":"GET"}`,
			templateFlags: []string{
				"service.1:Output {{.Type}}",
				"request.1:{{.Method}} {{.Type}}",
			},
			want: "Output service.1\nGET request.1\n",
		},
		{
			name:  "template for custom log format",
			input: `{"type":"service.9999","message":"foo"}`,
			templateFlags: []string{
				"service.9999:{{.type}} {{.message}}",
			},
			want: "service.9999 foo\n",
		},
	} {
		args := []string{"--input", currCase.input}
		for _, curr := range currCase.templateFlags {
			args = append(args, "--template", curr)
		}
		cmd := exec.Command(cli, args...)
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Case %d: %s\nOutput: %s", i, currCase.name, string(output))
		assert.Equal(t, currCase.want, string(output), "Case %d (%s): %s", i, currCase.name, currCase.input)
	}
}

func TestNoSubstituteFlag(t *testing.T) {
	for i, currCase := range []struct {
		name               string
		input              string
		wantDefault        string
		wantNoSubstitution string
	}{
		{
			name:               "request.1 substitution",
			input:              `{"type":"request.1","path":"/foo/{id}","pathParams":{"id":"test-id"}}`,
			wantDefault:        "[0001-01-01T00:00:00Z]     \"/foo/test-id \" 0  0\n",
			wantNoSubstitution: "[0001-01-01T00:00:00Z]     \"/foo/{id} \" 0  0\n",
		},
		{
			name:               "request.2 substitution",
			input:              `{"type":"request.2","path":"/foo/{id}","params":{"id":"test-id"}}`,
			wantDefault:        "[0001-01-01T00:00:00Z]     \"/foo/test-id \" 0 0 0\n",
			wantNoSubstitution: "[0001-01-01T00:00:00Z]     \"/foo/{id} \" 0 0 0\n",
		},
		{
			name:               "service.1 substitution",
			input:              `{"type":"service.1","time":"2017-05-25T18:49:10.652Z","level":"INFO","message":"Special node for '{}' already exists","origin":"com.palantir.example.NodeCreator","params":{},"unsafeParams":{"0":"my-special-node"}}`,
			wantDefault:        "INFO  [2017-05-25T18:49:10.652Z] com.palantir.example.NodeCreator: Special node for 'my-special-node' already exists (0: my-special-node)\n",
			wantNoSubstitution: "INFO  [2017-05-25T18:49:10.652Z] com.palantir.example.NodeCreator: Special node for '{}' already exists (0: my-special-node)\n",
		},
	} {
		args := []string{"--input", currCase.input}
		cmd := exec.Command(cli, args...)
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Case %d: %s\nOutput: %s", i, currCase.name, string(output))
		assert.Equal(t, currCase.wantDefault, string(output), "Case %d (%s): %s", i, currCase.name, currCase.input)

		args = []string{"--input", currCase.input, "--no-substitution"}
		cmd = exec.Command(cli, args...)
		output, err = cmd.CombinedOutput()
		require.NoError(t, err, "Case %d: %s\nOutput: %s", i, currCase.name, string(output))
		assert.Equal(t, currCase.wantNoSubstitution, string(output), "Case %d (%s): %s", i, currCase.name, currCase.input)
	}
}

func TestFileFlag(t *testing.T) {
	for i, currCase := range []struct {
		input string
		want  string
	}{
		{
			input: `{"type":"service.1","message":"foo\nbar"}`,
			want:  "      [0001-01-01T00:00:00Z]     foo\nbar\n",
		},
		{
			input: `{"type":"service.1","message":"foo\nbar"}
`,
			want: "      [0001-01-01T00:00:00Z]     foo\nbar\n",
		},
		{
			input: `foo`,
			want:  "foo\n",
		},
		{
			input: `foo
`,
			want: "foo\n",
		},
	} {
		currCaseDir, err := ioutil.TempDir("", fmt.Sprintf("case-%d", i))
		require.NoError(t, err, "Case %d")

		inputFile := path.Join(currCaseDir, "input.txt")
		err = ioutil.WriteFile(inputFile, []byte(currCase.input), 0644)
		require.NoError(t, err, "Case %d")

		cmd := exec.Command(cli, "--file", inputFile)
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Case %d\nOutput: %s", i, string(output))
		assert.Equal(t, currCase.want, string(output), "Case %d: %s", i, currCase.input)
	}
}

func TestLargeInput(t *testing.T) {
	const largeSize = 1024 * 1024
	largeBytes := make([]byte, largeSize)
	_, _ = rand.Read(largeBytes)
	largeStr := base64.StdEncoding.EncodeToString(largeBytes)

	input := fmt.Sprintf(`{"type":"service.1","message":%q}`, largeStr)
	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err, "Case %d")

	inputFile := path.Join(tmpDir, "input.txt")
	err = ioutil.WriteFile(inputFile, []byte(input), 0644)
	require.NoError(t, err, "Case %d")

	cmd := exec.Command(cli, "--file", inputFile)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Output: %s", string(output))
}

func TestInvalidFileError(t *testing.T) {
	cmd := exec.Command(cli, "--file", "_invalid_file_")
	output, err := cmd.CombinedOutput()
	assert.Error(t, err)
	assert.Equal(t, "failed to open file: open _invalid_file_: no such file or directory\n", string(output))
}

func TestCannotSpecifyInputAndFile(t *testing.T) {
	cmd := exec.Command(cli, "--input", "foo-bar-baz", "--file", "_invalid_file_")
	output, err := cmd.CombinedOutput()
	assert.Error(t, err)
	assert.Equal(t, "input and file cannot both be specified\n", string(output))
}

func runWithStdin(t *testing.T, cmd *exec.Cmd, input string) string {
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)

	err = cmd.Start()
	require.NoError(t, err)

	_, err = io.Copy(stdin, bytes.NewBufferString(input))
	require.NoError(t, err)

	err = stdin.Close()
	require.NoError(t, err)

	content, err := ioutil.ReadAll(stdout)
	require.NoError(t, err)

	err = cmd.Wait()
	require.NoError(t, err)

	return string(content)
}
