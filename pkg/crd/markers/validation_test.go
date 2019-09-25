/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package markers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

var _ = Describe("CRD Marker", func() {
	preserveUnknownFields := true

	type defaultCase struct {
		schema      *apiext.JSONSchemaProps
		value       string
		expected    string
		errExpected bool
	}

	DescribeTable("Default",
		func(c defaultCase) {
			def := Default(c.value)
			err := def.ApplyToSchema(c.schema)

			if c.errExpected {
				Expect(err).To(HaveOccurred())
				return
			}
			Expect(err).NotTo(HaveOccurred())

			expected := c.expected
			if len(expected) == 0 {
				expected = c.value
			}
			Expect(c.schema.Default).To(Equal(&apiext.JSON{Raw: []byte(expected)}))
		},
		Entry("should support boolean", defaultCase{
			schema: &apiext.JSONSchemaProps{
				Type: "boolean",
			},
			value: "true",
		}),
		Entry("should fail for invalid boolean", defaultCase{
			schema: &apiext.JSONSchemaProps{
				Type: "boolean",
			},
			value:       "foo",
			errExpected: true,
		}),
		Entry("should quote strings", defaultCase{
			schema: &apiext.JSONSchemaProps{
				Type: "string",
			},
			value:    "foo",
			expected: `"foo"`,
		}),
		Entry("should support array of primitives", defaultCase{
			schema: &apiext.JSONSchemaProps{
				Type: "array",
				Items: &apiext.JSONSchemaPropsOrArray{
					Schema: &apiext.JSONSchemaProps{
						Type: "string",
					},
				},
			},
			value:    `["a", "b"]`,
			expected: `["a","b"]`,
		}),
		Entry("should support array of objects", defaultCase{
			schema: &apiext.JSONSchemaProps{
				Type: "array",
				Items: &apiext.JSONSchemaPropsOrArray{
					Schema: &apiext.JSONSchemaProps{
						Type: "object",
						Properties: map[string]apiext.JSONSchemaProps{
							"type": {
								Type: "string",
							},
							"value": {
								Type:   "integer",
								Format: "int32",
							},
						},
					},
				},
			},
			value:    `[{"type": "magic", "value": 42}]`,
			expected: `[{"type":"magic","value":42}]`,
		}),
		Entry("should support simple object", defaultCase{
			schema: &apiext.JSONSchemaProps{
				Type: "object",
				Properties: map[string]apiext.JSONSchemaProps{
					"type": {
						Type: "string",
					},
					"value": {
						Type:   "integer",
						Format: "int32",
					},
				},
			},
			value: `{"type":"magic","value":42}`,
		}),
		Entry("should support complex object", defaultCase{
			schema: &apiext.JSONSchemaProps{
				Type: "object",
				Properties: map[string]apiext.JSONSchemaProps{
					"nested": {
						Type: "object",
						Properties: map[string]apiext.JSONSchemaProps{
							"value": {
								Type: "string",
							},
						},
					},
					"value": {
						Type:   "integer",
						Format: "int32",
					},
				},
			},
			value:    `{"type": {"nested": {"value": "magic"}}, "value": 42}`,
			expected: `{"type":{"nested":{"value":"magic"}},"value":42}`,
		}),
		Entry("should support arbitrary objects", defaultCase{
			schema: &apiext.JSONSchemaProps{
				XPreserveUnknownFields: &preserveUnknownFields,
			},
			value:    `{"type": "magic", "value": 42}`,
			expected: `{"type":"magic","value":42}`,
		}),
	)
})
