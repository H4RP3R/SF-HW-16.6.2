package bankclient

import (
	"errors"
	"math/rand"
	"sync"
	"testing"
)

func randomIntArr(size, min, max int) []int {
	arr := make([]int, size)

	for i := 0; i < size; i++ {
		arr[i] = rand.Intn(max-min) + min
	}

	return arr
}

func Test_bankClient_Deposit(t *testing.T) {
	tests := []struct {
		name            string
		amountToDeposit int
		currentBalance  int
		wantErr         error
	}{
		{"zero balance", 0, 0, nil},
		{"deposit 100", 100, 100, nil},
		{"deposit 42", 42, 142, nil},
		{"deposit -5", -5, 142, ErrNegativeAmountOperation},
	}

	user := NewTestPerson("test_person.json")
	bClient := New(user)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bClient.Deposit(tt.amountToDeposit)
			got, want := bClient.owner.balance, tt.currentBalance
			if got != want {
				t.Errorf("got balance: %d, want balance: %d", got, want)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want err: %v, got err: %v", tt.wantErr, &err)
			}
		})
	}
}

func Test_bankClient_Withdrawal(t *testing.T) {
	var startBalance = 1100
	tests := []struct {
		name             string
		amountToWithdraw int
		currentBalance   int
		wantErr          error
	}{
		{"withdraw 0", 0, 1100, nil},
		{"withdraw 100", 100, 1000, nil},
		{"withdraw -123", -123, 1000, ErrNegativeAmountOperation},
		{"withdraw -1001", 1001, 1000, ErrInsufficientFunds},
		{"withdraw 999", 999, 1, nil},
		{"withdraw 1", 1, 0, nil},
	}

	user := NewTestPerson("test_person.json")
	bClient := New(user)
	bClient.owner.balance = startBalance
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bClient.Withdrawal(tt.amountToWithdraw)
			got, want := bClient.owner.balance, tt.currentBalance
			if got != want {
				t.Errorf("got balance: %d, want balance: %d", got, want)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want err: %v, got err: %v", tt.wantErr, &err)
			}
		})
	}
}

func Test_bankClient_Deposit_MultiThreading(t *testing.T) {
	amountArr := randomIntArr(10000, -500, 10_000_000)
	wantBalance := 0
	for _, n := range amountArr {
		if n >= 0 {
			wantBalance += n
		}
	}

	user := NewTestPerson("test_person.json")
	bClient := New(user)
	gNum := 100
	step := len(amountArr) / gNum

	var wg sync.WaitGroup
	for i := range gNum {
		wg.Add(1)
		go func(i int) {
			for _, amount := range amountArr[i*step : (i+1)*step] {
				bClient.Deposit(amount)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	if wantBalance != bClient.owner.balance {
		t.Errorf("balance want: %d, balance got: %d", wantBalance, bClient.owner.balance)
	}
}

func Test_bankClient_Withdrawal_MultiThreading(t *testing.T) {
	startBalance := 1_500_000
	amountArr := randomIntArr(1000, -500, 1000)
	wantBalance := startBalance
	for _, n := range amountArr {
		if n >= 0 {
			wantBalance -= n
		}
	}

	user := NewTestPerson("test_person.json")
	bClient := New(user)
	bClient.owner.balance = startBalance
	gNum := 100
	step := len(amountArr) / gNum

	var wg sync.WaitGroup
	for i := range gNum {
		wg.Add(1)
		go func(i int) {
			for _, amount := range amountArr[i*step : (i+1)*step] {
				bClient.Withdrawal(amount)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	if wantBalance != bClient.owner.balance {
		t.Errorf("balance want: %d, balance got: %d", wantBalance, bClient.owner.balance)
	}
}
