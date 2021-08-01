package db

import (
	"sync"

	"gorm.io/gorm"
)

var DB *gorm.DB

type Currency = int
type Money = int64
type UserID = uint64

const (
	CurrBTC = iota
	CurrXLM
	CurrXRP
	CurrETH
	CurrUSD
	CurrEUR
)

type env struct {
	db *gorm.DB
	//logger? and stuff too
}

type Transaction struct {
	*env
	ID   uint64  `gorm:"primaryKey;"` // commit id, one commit can *should* have multiple transactions, ie + in one curr - in another
	Diff Money   //current balance, after the change, in the currency's base units
	From Address `gorm:"primaryKey;"`
	To   Address `gorm:"primaryKey;"`
}

func (t *Transaction) Insert() error {
	return nil
}

func (t *Transaction) Equal(t2 *Transaction) bool {
	return false
}

//Opposite transaction when negative diff
func (t *Transaction) Opposite(t2 *Transaction) bool {
	return false
}

//Pair transaction when opposite address and same other stuff
func (t *Transaction) Pair(t2 *Transaction) bool {
	return false
}

type Address struct {
	*env
	lock     sync.Locker
	Addr     []byte   `gorm:"primaryKey;"`
	Currency Currency `gorm:"primaryKey;"`
}

//GetBalance gets a sum of all transactions to the address
func (a *Address) GetBalance() Money {
	return 0
}

//Lock processing on address
func (a *Address) Lock() {
}

//Unlock processing on address
func (a *Address) Unlock() {
}
