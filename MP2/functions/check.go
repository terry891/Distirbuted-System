package functions

import "fmt"

//Panic if an error is thrown
func Check(e error) {
	if e != nil {
		fmt.Printf("\n\nWe got an error \n%s\n", e)
		panic(e)
	}
}
