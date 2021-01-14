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

package logs_test

import (
	"testing"
)

func Test_diagnostics1LogTyper_DefaultFormatter_ThreadDump(t *testing.T) {
	RunLogTests(t, []LogTest{
		{
			name: "Thread dump diagnostic log entry",
			input: []string{
				`{ "type": "diagnostic.1", "time": "2018-09-26T18:23:09.392Z", "diagnostic": { "type": "threadDump", "threadDump": { "threads": [{ "id": 33476, "name": "{threadName33476}", "stackTrace": [{ "procedure": "sun.misc.Unsafe.park", "file": "Unsafe.java", "line": -2, "params": { "isNative": true } }, { "procedure": "java.util.concurrent.locks.LockSupport.parkNanos", "file": "LockSupport.java", "line": 215, "params": { "isNative": false } }, { "procedure": "java.util.concurrent.SynchronousQueue$TransferStack.awaitFulfill", "file": "SynchronousQueue.java", "line": 460, "params": { "isNative": false } }], "params": { "threadState": "TIMED_WAITING", "blockedTime": "-1", "blockedCount": "0", "lockOwnerId": "-1", "waitedTime": "-1", "waitedCount": "1", "isSuspended": false, "isNative": false } }] } }, "unsafeParams": { "threadName33476": "SQL execute statement-20375" } }`,
			},
			output: []string{
				`[2018-09-26T18:23:09.392Z]
(threadName33476: SQL execute statement-20375)
"SQL execute statement-20375" tid=33476 (blockedCount: 0, blockedTime: -1, isNative: false, isSuspended: false, lockOwnerId: -1, threadState: TIMED_WAITING, waitedCount: 1, waitedTime: -1)
	sun.misc.Unsafe.park(Unsafe.java:-2)
	- isNative: true
	java.util.concurrent.locks.LockSupport.parkNanos(LockSupport.java:215)
	- isNative: false
	java.util.concurrent.SynchronousQueue$TransferStack.awaitFulfill(SynchronousQueue.java:460)
	- isNative: false
`,
			},
		},
	})
}

func Test_diagnostics1LogTyper_DefaultFormatter_ThreadDump_Without_UnsafeParams(t *testing.T) {
	RunLogTests(t, []LogTest{
		{
			name: "Thread dump diagnostic log entry",
			input: []string{
				`{ "type": "diagnostic.1", "time": "2018-09-26T18:23:09.392Z", "diagnostic": { "type": "threadDump", "threadDump": { "threads": [{ "id": 33476, "name": "SQL execute statement-20375", "stackTrace": [{ "procedure": "sun.misc.Unsafe.park", "file": "Unsafe.java", "line": -2, "params": { "isNative": true } }, { "procedure": "java.util.concurrent.locks.LockSupport.parkNanos", "file": "LockSupport.java", "line": 215, "params": { "isNative": false } }, { "procedure": "java.util.concurrent.SynchronousQueue$TransferStack.awaitFulfill", "file": "SynchronousQueue.java", "line": 460, "params": { "isNative": false } }], "params": { "threadState": "TIMED_WAITING", "blockedTime": "-1", "blockedCount": "0", "lockOwnerId": "-1", "waitedTime": "-1", "waitedCount": "1", "isSuspended": false, "isNative": false } }] } }}`,
			},
			output: []string{
				`[2018-09-26T18:23:09.392Z]
"SQL execute statement-20375" tid=33476 (blockedCount: 0, blockedTime: -1, isNative: false, isSuspended: false, lockOwnerId: -1, threadState: TIMED_WAITING, waitedCount: 1, waitedTime: -1)
	sun.misc.Unsafe.park(Unsafe.java:-2)
	- isNative: true
	java.util.concurrent.locks.LockSupport.parkNanos(LockSupport.java:215)
	- isNative: false
	java.util.concurrent.SynchronousQueue$TransferStack.awaitFulfill(SynchronousQueue.java:460)
	- isNative: false
`,
			},
		},
	})
}

func Test_diagnostics1LogTyper_DefaultFormatter_ThreadDump_WithTwoThreads(t *testing.T) {
	RunLogTests(t, []LogTest{
		{
			name: "2 threads dump diagnostic log entry",
			input: []string{
				`{ "type": "diagnostic.1", "time": "2018-09-26T18:23:09.392Z", "diagnostic": { "type": "threadDump", "threadDump": { "threads": [{ "id": 33476, "name": "{threadName33476}", "stackTrace": [{ "procedure": "sun.misc.Unsafe.park", "params": { "isNative": true } }, { "procedure": "java.util.concurrent.locks.LockSupport.parkNanos", "file": "LockSupport.java", "line": 215, "params": { "isNative": false } }, { "procedure": "java.util.concurrent.SynchronousQueue$TransferStack.awaitFulfill", "file": "SynchronousQueue.java", "line": 460, "params": { "isNative": false } }], "params": { "threadState": "TIMED_WAITING", "blockedTime": "-1", "blockedCount": "0", "lockOwnerId": "-1", "waitedTime": "-1", "waitedCount": "1", "isSuspended": false, "isNative": false } }, { "id": 33477, "name": "{threadName33477}", "stackTrace": [{ "procedure": "sun.misc.Unsafe.park", "file": "Unsafe.java", "line": -2, "params": { "isNative": true } }, { "procedure": "java.util.concurrent.locks.LockSupport.parkNanos", "file": "LockSupport.java", "line": 215, "params": { "isNative": false } }, { "procedure": "java.util.concurrent.SynchronousQueue$TransferStack.awaitFulfill", "file": "SynchronousQueue.java", "line": 460, "params": { "isNative": false } }], "params": { "threadState": "TIMED_WAITING", "blockedTime": "-1", "blockedCount": "0", "lockOwnerId": "-1", "waitedTime": "-1", "waitedCount": "1", "isSuspended": false, "isNative": false } }] } }, "unsafeParams": { "threadName33476": "SQL execute statement-20375", "threadName33477": "SQL execute statement-20377" } }`,
			},
			output: []string{
				`[2018-09-26T18:23:09.392Z]
(threadName33476: SQL execute statement-20375, threadName33477: SQL execute statement-20377)
"SQL execute statement-20375" tid=33476 (blockedCount: 0, blockedTime: -1, isNative: false, isSuspended: false, lockOwnerId: -1, threadState: TIMED_WAITING, waitedCount: 1, waitedTime: -1)
	sun.misc.Unsafe.park
	- isNative: true
	java.util.concurrent.locks.LockSupport.parkNanos(LockSupport.java:215)
	- isNative: false
	java.util.concurrent.SynchronousQueue$TransferStack.awaitFulfill(SynchronousQueue.java:460)
	- isNative: false
"SQL execute statement-20377" tid=33477 (blockedCount: 0, blockedTime: -1, isNative: false, isSuspended: false, lockOwnerId: -1, threadState: TIMED_WAITING, waitedCount: 1, waitedTime: -1)
	sun.misc.Unsafe.park(Unsafe.java:-2)
	- isNative: true
	java.util.concurrent.locks.LockSupport.parkNanos(LockSupport.java:215)
	- isNative: false
	java.util.concurrent.SynchronousQueue$TransferStack.awaitFulfill(SynchronousQueue.java:460)
	- isNative: false
`,
			},
		},
	})
}

func Test_diagnostics1LogTyper_DefaultFormatter_Generic(t *testing.T) {
	RunLogTests(t, []LogTest{
		{
			name: "Generic diagnostic log entry",
			input: []string{
				`{ "type": "diagnostic.1", "time": "2018-09-26T18:23:09.392Z", "diagnostic": { "type": "generic", "generic": {"Value": "Generic value"}}}`,
			},
			output: []string{
				`[2018-09-26T18:23:09.392Z] Generic value`,
			},
		},
	})
}

func Test_diagnostics1LogTyper_DefaultFormatter_Unknown(t *testing.T) {
	RunLogTests(t, []LogTest{
		{
			name: "Generic diagnostic log entry",
			input: []string{
				`{ "type": "diagnostic.1", "time": "2018-09-26T18:23:09.392Z", "diagnostic": { "type": "unknown", "unknown": {"Value": "Unknown value"}}}`,
			},
			output: []string{
				`[2018-09-26T18:23:09.392Z] [unknown] log type is not implemented for diagnostic.1, log line will be skipped`,
			},
		},
	})
}

func Test_diagnostics1LogTyper_DefaultFormatter_ThreadDump_Without_LineNumbers(t *testing.T) {
	RunLogTests(t, []LogTest{
		{
			name: "Thread dump diagnostic log entry",
			input: []string{
				`{ "type": "diagnostic.1", "time": "2018-09-26T18:23:09.392Z", "diagnostic": { "type": "threadDump", "threadDump": { "threads": [{ "id": 33476, "name": "SQL execute statement-20375", "stackTrace": [{ "procedure": "sun.misc.Unsafe.park", "file": "Unsafe.java", "params": { "isNative": true } }, { "procedure": "java.util.concurrent.locks.LockSupport.parkNanos", "file": "LockSupport.java", "line": 215, "params": { "isNative": false } }, { "procedure": "java.util.concurrent.SynchronousQueue$TransferStack.awaitFulfill", "file": "SynchronousQueue.java", "line": 460, "params": { "isNative": false } }], "params": { "threadState": "TIMED_WAITING", "blockedTime": "-1", "blockedCount": "0", "lockOwnerId": "-1", "waitedTime": "-1", "waitedCount": "1", "isSuspended": false, "isNative": false } }] } }}`,
			},
			output: []string{
				`[2018-09-26T18:23:09.392Z]
"SQL execute statement-20375" tid=33476 (blockedCount: 0, blockedTime: -1, isNative: false, isSuspended: false, lockOwnerId: -1, threadState: TIMED_WAITING, waitedCount: 1, waitedTime: -1)
	sun.misc.Unsafe.park(Unsafe.java)
	- isNative: true
	java.util.concurrent.locks.LockSupport.parkNanos(LockSupport.java:215)
	- isNative: false
	java.util.concurrent.SynchronousQueue$TransferStack.awaitFulfill(SynchronousQueue.java:460)
	- isNative: false
`,
			},
		},
	})
}
