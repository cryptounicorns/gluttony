package reflect

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseType(t *testing.T) {
	var (
		samples = []struct {
			name   string
			input  []string
			output []Type
		}{
			{
				name:  "primitive",
				input: []string{"int", "int32", "float32", "uintptr"},
				output: []Type{
					TypeInt,
					TypeInt32,
					TypeFloat32,
					TypeUintptr,
				},
			},
			{
				name: "composite",
				input: []string{
					"[]int",
					"[]int32",
					"[]float32",
					"[]uintptr",
					"map[string]int",
					"map[bool]string",
				},
				output: []Type{
					reflect.SliceOf(TypeInt),
					reflect.SliceOf(TypeInt32),
					reflect.SliceOf(TypeFloat32),
					reflect.SliceOf(TypeUintptr),
					reflect.MapOf(TypeString, TypeInt),
					reflect.MapOf(TypeBool, TypeString),
				},
			},
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				var (
					tt  Type
					err error
				)
				for k, v := range sample.input {
					tt, err = ParseType(v)
					assert.Equal(
						t,
						sample.output[k],
						tt,
						fmt.Sprintf(
							"Expected type %s",
							sample.output[k].String(),
						),
					)
					assert.Equal(t, nil, err)
				}
			},
		)
	}
}
