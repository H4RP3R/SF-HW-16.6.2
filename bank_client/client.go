package bankclient

import (
	"fmt"

	"github.com/google/uuid"
)

var ErrInsufficientFunds = fmt.Errorf("balance less than withdrawal amount")
var ErrNegativeAmountOperation = fmt.Errorf("operation with negative amount")

type BankClient interface {
	// Deposit deposits given amount to clients account
	Deposit(amount int) error

	// Withdrawal withdraws given amount from clients account.
	// return error if clients balance less than withdrawal amount
	Withdrawal(amount int) error

	// Balance returns clients balance
	Balance() int
}

// bankClient contains data and methods that are used to provide basic
// functionality for banking CLI app.
type bankClient struct {
	id    uuid.UUID
	owner *person
}

// Deposit deposits specified amount to clients account. If amount is
// negative returns error.
func (bc *bankClient) Deposit(amount int) error {
	if amount < 0 {
		return ErrNegativeAmountOperation
	}

	bc.owner.balanceMutex.Lock()
	bc.owner.balance += amount
	bc.owner.balanceMutex.Unlock()

	return nil
}

// Withdrawal withdraws specified amount from clients account,
// returns error if clients balance is less than withdrawal amount or
// if the amount is negative.
func (bc *bankClient) Withdrawal(amount int) error {
	if amount < 0 {
		return ErrNegativeAmountOperation
	}

	bc.owner.balanceMutex.Lock()
	defer bc.owner.balanceMutex.Unlock()

	if bc.owner.balance >= amount {
		bc.owner.balance -= amount
		return nil
	}

	return ErrInsufficientFunds
}

// Balance returns clients balance.
func (bc *bankClient) Balance() int {
	bc.owner.balanceMutex.RLock()
	balance := bc.owner.balance
	bc.owner.balanceMutex.RUnlock()

	return balance
}

// New returns new client with given owner.
func New(owner *person) *bankClient {
	return &bankClient{
		id:    uuid.New(),
		owner: owner,
	}
}
