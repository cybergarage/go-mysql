// Copyright (C) 2025 The go-mysql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stmt

// ParameterOption is the option of the parameter.
type ParameterOption func(*parameter)

type parameter struct {
	name string
	typ  FieldType
	v    any
}

// WithParameterName sets the name of the parameter.
func WithParameterName(name string) ParameterOption {
	return func(p *parameter) {
		p.name = name
	}
}

// WithParameterType sets the type of the parameter.
func WithParameterType(typ FieldType) ParameterOption {
	return func(p *parameter) {
		p.typ = typ
	}
}

// WithParameterValue sets the value of the parameter.
func WithParameterValue(v any) ParameterOption {
	return func(p *parameter) {
		p.v = v
	}
}

// NewParameter creates a new parameter with the options.
func NewParameter(opts ...ParameterOption) Parameter {
	p := &parameter{
		name: "",
		typ:  0,
		v:    nil,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Name returns the name of the parameter.
func (p *parameter) Name() string {
	return p.name
}

// Type returns the type of the parameter.
func (p *parameter) Type() FieldType {
	return p.typ
}

// Value returns the value of the parameter.
func (p *parameter) Value() any {
	return p.v
}
