package types

import "k.malhotra.cc/go/payment-thing/db"

type Key struct {
	Priv string
	Pub  string
}

// TODO check that it's positive
type TxCallback = func(TXID string, UID uint64, Amount db.Money)
type CurrencyIn interface {
	//starts the listener or a polling loop and calls back
	Init(pubKey Key, cb TxCallback)
	GenAddr(UserID uint64) string
}
