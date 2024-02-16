package modules

import (
	"fmt"

	calculator "go-modules/src/calculator"
)

func Main() {
	fmt.Println("Hello World!");
	fmt.Println(calculator.Add(1,2));
}
