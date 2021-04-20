package functions

import (
	"fmt"
)

//Panic if error is thrown
func Check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}
