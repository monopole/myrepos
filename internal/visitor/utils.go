package visitor

import "fmt"

func indent(s int) {
	for i := 0; i < s; i++ {
		fmt.Print("  ")
	}
}
