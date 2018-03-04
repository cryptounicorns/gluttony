package main

import (
	"fmt"

	"github.com/corpix/reflect"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	var (
		numbers           = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
		mapNumberToString = map[int]string{
			1: "1",
			2: "2",
			3: "3",
			4: "4",
			5: "5",
			6: "6",
			7: "7",
			8: "8",
			9: "9",
			0: "0",
		}

		targetType reflect.Type
	)

	targetType = reflect.TypeOf([]string{})
	fmt.Printf(
		"Convert %s to %s\n",
		reflect.TypeOf(numbers),
		targetType,
	)
	spew.Dump(
		reflect.ConvertToType(
			numbers,
			targetType,
		),
	)

	fmt.Println("")

	targetType = reflect.TypeOf(map[string]int{})
	fmt.Printf(
		"Convert %s to %s\n",
		reflect.TypeOf(mapNumberToString),
		targetType,
	)
	spew.Dump(
		reflect.ConvertToType(
			mapNumberToString,
			targetType,
		),
	)

}
