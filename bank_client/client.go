package bankclient

import (
	"fmt"

	"github.com/google/uuid"
)

var ErrInsufficientFunds = fmt.Errorf("balance is less then withdrawal amount")
var ErrNegativeAmountOperation = fmt.Errorf("operation with negative amount")

type BankClient interface {
	// Deposit deposits given amount to clients account
	Deposit(amount int) error

	// Withdrawal withdraws given amount from clients account.
	// return error if clients balance less the withdrawal amount
	Withdrawal(amount int) error

	// Balance returns clients balance
	Balance() int
}

type bankClient struct {
	id    uuid.UUID
	owner *person
}

func (bc *bankClient) Deposit(amount int) error {
	if amount < 0 {
		return ErrNegativeAmountOperation
	}

	bc.owner.balanceMutex.Lock()
	bc.owner.balance += amount
	bc.owner.balanceMutex.Unlock()

	return nil
}

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

func (bc *bankClient) Balance() int {
	bc.owner.balanceMutex.RLock()
	balance := bc.owner.balance
	bc.owner.balanceMutex.RUnlock()

	return balance
}

func New(owner *person) *bankClient {
	return &bankClient{
		id:    uuid.New(),
		owner: owner,
	}
}
