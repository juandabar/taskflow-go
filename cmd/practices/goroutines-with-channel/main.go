package main

import (
	"fmt"
	"time"
)

type Result struct {
	Name string
}

func tarea1(ch chan Result) {
	time.Sleep(6 * time.Second)
	ch <- Result{Name: "tarea1"}
}

func tarea2(ch chan Result) {
	time.Sleep(6 * time.Second)
	ch <- Result{Name: "tarea2"}
}

func tarea3(ch chan Result) {
	time.Sleep(6 * time.Second)
	ch <- Result{Name: "tarea3"}
}

func main() {
	ch := make(chan Result)

	go tarea1(ch)
	go tarea2(ch)
	go tarea3(ch)

	r := <-ch
	fmt.Println(r.Name)

	r = <-ch
	fmt.Println(r.Name)

	r = <-ch
	fmt.Println(r.Name)

	fmt.Println("finalizo todo...")

}
