package gojsonschema

import (
	"testing"
)

func TestJsonContextReferenceTokens(t *testing.T) {
	cases := []struct {
		input    *JsonContext
		expected string
	}{
		{
			nil,
			"",
		},
		{
			&JsonContext{
				head: STRING_CONTEXT_ROOT,
				tail: nil,
			},
			"",
		},
		{
			&JsonContext{
				head: "a",
				tail: &JsonContext{
					head: STRING_CONTEXT_ROOT,
					tail: nil,
				},
			},
			"/a",
		},
		{
			&JsonContext{
				head: "b",
				tail: &JsonContext{
					head: "a",
					tail: &JsonContext{
						head: STRING_CONTEXT_ROOT,
						tail: nil,
					},
				},
			},
			"/a/b",
		},
		{
			&JsonContext{
				head: "~/",
				tail: &JsonContext{
					head: "/~",
					tail: &JsonContext{
						head: STRING_CONTEXT_ROOT,
						tail: nil,
					},
				},
			},
			"/~1~0/~0~1",
		},
	}
	for _, c := range cases {
		actual := c.input.JSONPointer()
		if actual != c.expected {
			t.Errorf("expected %s but got %s", c.expected, actual)
		}
	}
}
