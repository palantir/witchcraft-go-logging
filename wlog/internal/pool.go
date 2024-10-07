// Copyright (c) 2024 Palantir Technologies. All rights reserved.
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

package wloginternal

import (
	"sync"

	"github.com/palantir/witchcraft-go-logging/wapi/logging"
)

// SyncPool exposes a sync.Pool with a reset function for the type T.
// The reset function is called on every object before it is returned to the pool.
type SyncPool[T logging.LogTypes] struct {
	pool  sync.Pool
	reset func(*T)
}

// NewPool returns a new SyncPool whose Get() returns a new pointer to a zero-value T
// and a closure to reset the object and return it to the pool.
func NewPool[T logging.LogTypes](reset func(*T)) *SyncPool[T] {
	return &SyncPool[T]{pool: sync.Pool{New: func() any { return new(T) }}, reset: reset}
}

func (p *SyncPool[T]) Get() *T {
	return p.pool.Get().(*T)
}

func (p *SyncPool[T]) Put(obj *T) {
	p.reset(obj)
	p.pool.Put(obj)
}
