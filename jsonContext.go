// Copyright 2013 MongoDB, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// author           tolsen
// author-github    https://github.com/tolsen
//
// repository-name  gojsonschema
// repository-desc  An implementation of JSON Schema, based on IETF's draft v4 - Go language.
//
// description      Implements a persistent (immutable w/ shared structure) singly-linked list of strings for the purpose of storing a json context
//
// created          04-09-2013

package gojsonschema

import (
	"bytes"
	"fmt"
	"strings"
)

// JsonContext implements a persistent linked-list of strings
type JsonContext struct {
	head string
	tail *JsonContext
}

// NewJsonContext creates a new JsonContext
func NewJsonContext(head string, tail *JsonContext) *JsonContext {
	return &JsonContext{head, tail}
}

// String displays the context in reverse.
// This plays well with the data structure's persistent nature with
// Cons and a json document's tree structure.
func (c *JsonContext) String(del ...string) string {
	byteArr := make([]byte, 0, c.stringLen())
	buf := bytes.NewBuffer(byteArr)
	c.writeStringToBuffer(buf, del)

	return buf.String()
}

func (c *JsonContext) stringLen() int {
	length := 0
	if c.tail != nil {
		length = c.tail.stringLen() + 1 // add 1 for "."
	}

	length += len(c.head)
	return length
}

func (c *JsonContext) writeStringToBuffer(buf *bytes.Buffer, del []string) {
	if c.tail != nil {
		c.tail.writeStringToBuffer(buf, del)

		if len(del) > 0 {
			buf.WriteString(del[0])
		} else {
			buf.WriteString(".")
		}
	}

	buf.WriteString(c.head)
}

func (c *JsonContext) ReferenceTokens() []string {
	var output []string

	curr := c
	for {
		if curr == nil {
			break
		}
		if curr.tail == nil && curr.head == STRING_CONTEXT_ROOT {
			// Special case: skip (root)
		} else {
			output = append(output, curr.head)
		}
		curr = curr.tail
	}

	// Reverse
	for i, j := 0, len(output)-1; i < j; i, j = i+1, j-1 {
		output[i], output[j] = output[j], output[i]
	}

	return output
}

func (c *JsonContext) JSONPointer() string {
	tokens := c.ReferenceTokens()
	builder := &strings.Builder{}
	for _, token := range tokens {
		fmt.Fprintf(builder, "/%s", c.EncodeReferenceToken(token))
	}
	return builder.String()
}

func (c *JsonContext) EncodeReferenceToken(token string) string {
	token = strings.Replace(token, "~", "~0", -1)
	token = strings.Replace(token, "/", "~1", -1)
	return token
}
