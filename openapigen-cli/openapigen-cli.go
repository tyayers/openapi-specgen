package main

import (
	"fmt"
	"os"

	"github.com/tyayers/openapigen"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please call with a single URL as parameter, containing the GET list result.")
	} else {
		result := openapigen.GenerateSpec(os.Args[1])
		fmt.Println(result)
	}
}
