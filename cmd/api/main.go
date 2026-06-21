package main

import (
	"fmt"
)

func main() {
	names := []string{"Juan", "Maria", "Pedro"}
	names = append(names, "Ana")

	for _, name := range names {
		fmt.Println(name)
	}

	ages := map[string]int{
		"Juan": 28,
		"Ana":  25,
	}

	ages["Maria"] = 26

	fmt.Println(ages)
}
