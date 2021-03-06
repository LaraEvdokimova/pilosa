// Copyright 2017 Pilosa Corp.
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

package tracing

import (
	"context"
	"net/http"
)

// GlobalTracer is a single, global instance of Tracer.
var GlobalTracer Tracer = NopTracer()

// StartSpanFromContext returnus a new child span and context from a given
// context using the global tracer.
func StartSpanFromContext(ctx context.Context, operationName string) (Span, context.Context) {
	return GlobalTracer.StartSpanFromContext(ctx, operationName)
}

// Tracer implements a generic distributed tracing interface.
type Tracer interface {
	// Returns a new child span and context from a given context.
	StartSpanFromContext(ctx context.Context, operationName string) (Span, context.Context)

	// Adds the required HTTP headers to pass context between nodes.
	InjectHTTPHeaders(r *http.Request)

	// Reads the HTTP headers to derive incoming context.
	ExtractHTTPHeaders(r *http.Request) (Span, context.Context)
}

// Span represents a single span in a distributed trace.
type Span interface {
	// Sets the end timestamp and finalizes Span state.
	Finish()

	// Adds key/value pairs to the span.
	LogKV(alternatingKeyValues ...interface{})
}

// NopTracer returns a tracer that doesn't do anything.
func NopTracer() Tracer {
	return &nopTracer{}
}

type nopTracer struct{}

func (t *nopTracer) StartSpanFromContext(ctx context.Context, operationName string) (Span, context.Context) {
	return &nopSpan{}, ctx
}

func (t *nopTracer) InjectHTTPHeaders(r *http.Request) {}

func (t *nopTracer) ExtractHTTPHeaders(r *http.Request) (Span, context.Context) {
	return &nopSpan{}, r.Context()
}

type nopSpan struct{}

func (s *nopSpan) Finish()                                   {}
func (s *nopSpan) LogKV(alternatingKeyValues ...interface{}) {}
