package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Type = reflect.Type

var (
	Types = []Type{
		TypeBool,
		TypeInt,
		TypeInt8,
		TypeInt16,
		TypeInt32,
		TypeInt64,
		TypeUint,
		TypeUint8,
		TypeUint16,
		TypeUint32,
		TypeUint64,
		TypeFloat32,
		TypeFloat64,
		TypeComplex64,
		TypeComplex128,
		TypeUintptr,
		TypeString,
	}

	StringToType = map[string]Type{}

	TypeInvalid    = Type(nil)
	TypeBool       = reflect.TypeOf(false)
	TypeInt        = reflect.TypeOf(int(0))
	TypeInt8       = reflect.TypeOf(int8(0))
	TypeInt16      = reflect.TypeOf(int16(0))
	TypeInt32      = reflect.TypeOf(int32(0))
	TypeInt64      = reflect.TypeOf(int64(0))
	TypeUint       = reflect.TypeOf(uint(0))
	TypeUint8      = reflect.TypeOf(uint8(0))
	TypeUint16     = reflect.TypeOf(uint16(0))
	TypeUint32     = reflect.TypeOf(uint32(0))
	TypeUint64     = reflect.TypeOf(uint64(0))
	TypeFloat32    = reflect.TypeOf(float32(0))
	TypeFloat64    = reflect.TypeOf(float64(0))
	TypeComplex64  = reflect.TypeOf(complex64(0))
	TypeComplex128 = reflect.TypeOf(complex128(0))
	TypeUintptr    = reflect.TypeOf(uintptr(0))
	TypeString     = reflect.TypeOf(string(""))
)

func init() {
	var (
		t   Type
		buf = []Type{}
	)

	for _, v := range Types {
		StringToType[v.String()] = v

		t = reflect.SliceOf(v)
		buf = append(buf, t)
		StringToType[t.String()] = t

		t = reflect.PtrTo(v)
		buf = append(buf, t)
		StringToType[t.String()] = t

		for _, vv := range Types {
			t = reflect.MapOf(v, vv)
			buf = append(buf, t)
			StringToType[t.String()] = t
		}

	}

	Types = append(Types, buf...)
}

func representPrimitiveConstructor(name string) string {
	switch name {
	case "bool":
		return "false"
	case "int":
		return "int(0)"
	case "int8":
		return "int8(0)"
	case "int16":
		return "int16(0)"
	case "int32":
		return "int32(0)"
	case "int64":
		return "int64(0)"
	case "uint":
		return "uint(0)"
	case "uint8":
		return "uint8(0)"
	case "uint16":
		return "uint16(0)"
	case "uint32":
		return "uint32(0)"
	case "uint64":
		return "uint64(0)"
	case "float32":
		return "float32(0)"
	case "float64":
		return "float64(0)"
	case "complex64":
		return "complex64(0)"
	case "complex128":
		return "complex128(0)"
	case "uintptr":
		return "uintptr(0)"
	case "string":
		return "string(\"\")"
	default:
		panic(
			fmt.Sprintf(
				"PrimitUnsupported primitive '%s'",
				name,
			),
		)
	}
}

func representConstructor(name string) string {
	var (
		buf = ""
	)

	for k, c := range name {
		buf += string(c)

		switch buf {
		case "map":
			first := k + 2
			last := first
			for kk, vv := range name[first:] {
				switch vv {
				case ']':
					last = first + kk
				default:
					continue
				}
			}

			return fmt.Sprintf(
				"reflect.MapOf(%s, %s)",
				representConstructor(string(name[first])+name[first+1:last]),
				representConstructor(name[last+1:]),
			)
		case "[]":
			return fmt.Sprintf(
				"reflect.SliceOf(%s)",
				representConstructor(name[2:]),
			)
		case "*":
			return fmt.Sprintf(
				"reflect.PtrTo(%s)",
				representConstructor(name[1:]),
			)
		default:
			continue
		}
	}

	return fmt.Sprintf(
		"reflect.TypeOf(%s)",
		representPrimitiveConstructor(name),
	)
}

func representName(name string) string {
	var (
		buf = ""
	)

	for k, c := range name {
		buf += string(c)

		switch buf {
		case "map":
			first := k + 2
			last := first
			for kk, vv := range name[first:] {
				switch vv {
				case ']':
					last = first + kk
				default:
					continue
				}
			}

			return fmt.Sprintf(
				"MapOf%s",
				representName(string(name[first])+name[first+1:last])+
					representName(name[last+1:]),
			)
		case "[]":
			return fmt.Sprintf(
				"SliceOf%s",
				representName(name[2:]),
			)
		case "*":
			return fmt.Sprintf(
				"PtrTo%s",
				representName(name[1:]),
			)
		default:
			continue
		}
	}

	return strings.ToUpper(string(name[0])) + name[1:]
}

func main() {
	var (
		types = [][]string{}
	)

	for _, t := range Types {
		types = append(
			types,
			[]string{
				representName(t.String()),
				representConstructor(t.String()),
			},
		)
	}

	fmt.Println("package reflect")
	fmt.Println("")
	fmt.Println("import (\n\t\"reflect\"\n)")
	fmt.Println("")
	fmt.Println("var (")
	for _, v := range types {
		fmt.Printf("\tType%s = %s\n", v[0], v[1])
	}
	fmt.Println("")
	fmt.Println("\tTypes = []Type{")
	for _, v := range types {
		fmt.Printf("\t\tType%s,\n", v[0])
	}
	fmt.Println("\t}")
	fmt.Println("")
	fmt.Println("\tStringToType = map[string]Type{")
	for _, v := range types {
		fmt.Printf("\t\tType%s.String(): Type%s,\n", v[0], v[0])
	}
	fmt.Println("\t}")

	fmt.Println(")")
}
