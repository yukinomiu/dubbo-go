/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package impl

import "maps"

type RequestPayload struct {
	Params      interface{}
	Attachments map[string]interface{}
}

func NewRequestPayload(args interface{}, atta map[string]interface{}) *RequestPayload {
	var newAtta map[string]interface{}
	if atta == nil {
		newAtta = make(map[string]interface{})
	} else {
		// dubbox fix: use cloned map
		newAtta = maps.Clone(atta)
	}
	return &RequestPayload{
		Params:      args,
		Attachments: newAtta,
	}
}

func EnsureRequestPayload(body interface{}) *RequestPayload {
	if req, ok := body.(*RequestPayload); ok {
		return req
	}
	return NewRequestPayload(body, nil)
}
