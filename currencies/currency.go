package currencies

import (
	"k.malhotra.cc/go/payment-thing/currencies/stellar"
	"k.malhotra.cc/go/payment-thing/currencies/types"
)

func init() {
	// Assert types
	var _ = []types.CurrencyIn{&stellar.Stellar{}}

}
