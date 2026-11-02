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
)

// syncPool is like sync.Pool generically-typed and with a reset function for the type T.
// The reset function is called on every object before it is returned to the pool.
type syncPool[T any] struct {
	pool  sync.Pool
	reset func(*T)
}

// NewPool returns a new syncPool whose Get() returns a pointer to a zero-value or reset T,
// using the provided reset function to clear objects before they are returned to the pool.
func NewPool[T any](reset func(*T)) *syncPool[T] {
	return &syncPool[T]{
		pool:  sync.Pool{New: func() any { return new(T) }},
		reset: reset,
	}
}

// Get returns a pointer to a zero-value (or reset) T from the pool.
func (p *syncPool[T]) Get() *T {
	return p.pool.Get().(*T)
}

// Put resets the object and returns it to the pool.
func (p *syncPool[T]) Put(obj *T) {
	p.reset(obj)
	p.pool.Put(obj)
}

// Log objects start with nil maps so that unset map fields are omitted, so their maps are pooled separately.
var anyMaps, stringMaps sync.Pool

func mapPool[V any]() *sync.Pool {
	switch any((*V)(nil)).(type) {
	case *any:
		return &anyMaps
	case *string:
		return &stringMaps
	}
	return nil
}

func getMap[V any]() map[string]V {
	if pool := mapPool[V](); pool != nil {
		if m, ok := pool.Get().(map[string]V); ok {
			return m
		}
	}
	return make(map[string]V)
}

// ReleaseMap clears an owned map, returns it to the pool, and sets the field to nil.
func ReleaseMap[V any](m *map[string]V) {
	if *m == nil {
		return
	}
	clear(*m)
	if pool := mapPool[V](); pool != nil {
		pool.Put(*m)
	}
	*m = nil
}
