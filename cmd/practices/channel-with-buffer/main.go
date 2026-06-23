package main

import (
	"fmt"
	"time"
)

func producer(ch chan string) {
	for i := 1; i <= 5; i++ {
		fmt.Println("Produciendo:", i)
		ch <- fmt.Sprintf("event-%d", i)
		fmt.Println("Enviado:", i)
		time.Sleep(1 * time.Second)
	}
}

func consumer(ch chan string) {
	for i := 1; i <= 5; i++ {
		event := <-ch
		fmt.Println("Procesando:", event)
		time.Sleep(5 * time.Second)
	}
}

func main() {
	ch := make(chan string, 3) // 🔴 SIN BUFFER

	go producer(ch)
	go consumer(ch)

	time.Sleep(30 * time.Second)
}
