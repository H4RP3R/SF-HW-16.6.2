package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	bc "github.com/H4RP3R/bankclient"
)

func main() {
	user := bc.NewTestPerson("bank_client/test_person.json")
	var bcli bc.BankClient = bc.New(user)
	done := make(chan struct{})
	wg := &sync.WaitGroup{}

	fmt.Println("Account:", user)
	fmt.Println("Supported commands: balance, deposit, withdrawal, exit")

	for i := range 10 {
		wg.Add(1)
		go func(i int, done chan struct{}) {
			for {
				select {
				case <-done:
					log.Printf("[g deposit %d] finished", i)
					wg.Done()
					return
				default:
					time.Sleep(time.Millisecond * time.Duration(randInRange(500, 1000)))
					bcli.Deposit(randInRange(1, 10))
				}
			}
		}(i, done)
	}

	for i := range 5 {
		wg.Add(1)
		go func(i int, done chan struct{}) {
			for {
				select {
				case <-done:
					log.Printf("[g withdrawal %d] finished", i)
					wg.Done()
					return
				default:
					time.Sleep(time.Millisecond * time.Duration(randInRange(500, 1000)))
					err := bcli.Withdrawal(randInRange(1, 5))
					if err != nil {
						log.Println(err)
					}
				}
			}
		}(i, done)
	}

	for {
		command, err := readString()
		if err != nil {
			log.Println(err)
		}

		switch command {
		case "balance":
			fmt.Println(bcli.Balance())

		case "deposit":
			fmt.Print("Enter a value to deposit: ")
			amount, err := readInt()
			if err != nil {
				log.Println(fmt.Errorf("amount must be a number"))
				continue
			}
			err = bcli.Deposit(amount)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println("Successful operation. Balance:", bcli.Balance())
			}

		case "withdrawal":
			fmt.Print("Enter a value to withdraw: ")
			amount, err := readInt()
			if err != nil {
				log.Println(fmt.Errorf("amount must be a number"))
				continue
			}
			err = bcli.Withdrawal(amount)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println("Successful operation. Balance:", bcli.Balance())
			}

		case "exit":
			close(done)
			wg.Wait()
			fmt.Println("Bye!")
			os.Exit(0)

		default:
			fmt.Println("Unsupported command. You can use commands: balance, deposit, withdrawal, exit")
		}
	}
}

// readString reads from os.Stdin and return normalized command and error.
func readString() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	command, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.ToLower(strings.TrimSpace(command)), nil
}

// readInt reads from os.Stdin and returns int value and error.
func readInt() (int, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	input = strings.ToLower(strings.TrimSpace(input))
	amount, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

// randInRange returns pseudo-random number n. min<=n<=max
func randInRange(min, max int) int {
	max += 1
	return rand.Intn(max-min) + min
}
