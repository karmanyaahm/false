package db

import (
	"gorm.io/gorm"
)

var DB gorm.DB

type Currency = int

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
	ID      uint64  `gorm:"primaryKey;"` // commit id, one commit can *should* have multiple transactions, ie + in one curr - in another
	Balance int64   //current balance, after the change, in the currency's base units
	From    Address `gorm:"primaryKey;"`
	To      Address `gorm:"primaryKey;"`
}

func (t *Transaction) Insert() error {
	return nil
}

func (t *Transaction) Equal(t2 *Transaction) bool {
	return false
}

func (t *Transaction) Opposite(t2 *Transaction) bool {
	return false
}

type Account struct {
	*env
	ID        uint64     `gorm:"primaryKey;"`
	Addresses []*Address `gorm:"many2many:account_address;"` //addresses that this account owns
}

type Address struct {
	*env
	Addr     []byte     `gorm:"primaryKey;"`
	Currency Currency   `gorm:"primaryKey;"`
	Accounts []*Account `gorm:"many2many:account_address;"` //account that owns this addr
}
