package bankclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

// person represents client of bank. Contains clients personal information
// and balance.
type person struct {
	ID        uuid.UUID
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`

	balanceMutex sync.RWMutex
	balance      int
}

func (p *person) String() string {
	return fmt.Sprintf("%s %s [id:%s...]", p.FirstName, p.LastName, p.ID.String()[:6])
}

// NewTestPerson reads data from JSON file and returns pointer to a person
// struct with the data specified in file. Intends for testing purposes.
func NewTestPerson(filePath string) *person {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	p := person{}
	p.ID = uuid.New()
	err = json.Unmarshal(bytes, &p)
	if err != nil {
		log.Fatal(err)
	}

	return &p
}
