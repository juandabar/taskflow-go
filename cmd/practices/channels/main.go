package main

import (
	"fmt"
	"time"
)

type Result struct {
	name  string
	value string
}

func getUser(ch chan Result) {
	time.Sleep(300 * time.Millisecond)
	ch <- Result{"user", "Juan"}
}

func getBalance(ch chan Result) {
	time.Sleep(500 * time.Millisecond)
	ch <- Result{"balance", "$1000"}
}

func getTransactions(ch chan Result) {
	time.Sleep(700 * time.Millisecond)
	ch <- Result{"transactions", "10 ops"}
}

func main() {
	ch := make(chan Result)

	go getUser(ch)
	go getBalance(ch)
	go getTransactions(ch)

	infoUser := <-ch
	infoBalance := <-ch
	infoTransaction := <-ch

	fmt.Println("Pintando los datos en dashboard: ", infoUser, infoBalance, infoTransaction)

}
