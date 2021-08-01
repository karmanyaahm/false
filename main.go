package main

import (
	"fmt"

	"k.malhotra.cc/go/payment-thing/config"
	"k.malhotra.cc/go/payment-thing/currencies/stellar"
	"k.malhotra.cc/go/payment-thing/currencies/types"
	"k.malhotra.cc/go/payment-thing/db"
)

func main() {
	fmt.Println("vim-go")
}

//resolve on addresses with 'to' currency of balances with 'from' currency
//afterCommit and beforeCommit are inclusive and refer to the commit number of transactions
//any required currency exchanges inside this function will from 'from' to 'to'
func resolve(from db.Currency, to db.Currency, afterCommit int64, beforeCommit int64) {

}

type payment struct {
	Addr db.Address
	Amt  db.Money
}

//payoff pays off ppl with applicable balances
//make sure it isn't too lossy transaction and there's enough balance and the last payoff isn't too recent
func payoff(c db.Currency) {}

func payoffWho(addrs []payment) {}
func init() {
	db.Init()
	s := stellar.New(db.DB, types.Key{Pub: config.Get().Stellar.PubKey, Priv: config.Get().Stellar.PrivKey}, func(_ string, _ uint64, _ db.Money) {})
	s.In()
}
