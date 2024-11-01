// Copyright (C) 2024 The go-mysql Authors. All rights reserved.
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

package protocol

// AttributeMap represents a map of attributes.
type AttributeMap struct {
	keys  []string
	attrs map[string]string
}

// NewAttributeMap returns a new AttributeMap.
func NewAttributeMap() *AttributeMap {
	return &AttributeMap{
		keys:  make([]string, 0),
		attrs: make(map[string]string),
	}
}

// AddAttribute adds an attribute.
func (attrMap *AttributeMap) AddAttribute(key, value string) {
	attrMap.keys = append(attrMap.keys, key)
	attrMap.attrs[key] = value
}

// Attributes returns the attributes.
func (attrMap *AttributeMap) Attributes() map[string]string {
	return attrMap.attrs
}

func (attrMap *AttributeMap) LookupAttribute(key string) (string, bool) {
	value, hasKey := attrMap.attrs[key]
	return value, hasKey
}
