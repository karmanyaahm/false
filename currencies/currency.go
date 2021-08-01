package currencies

import (
	"github.com/karmanyaahm/payment-thing/currencies/stellar"
	"github.com/karmanyaahm/payment-thing/currencies/types"
)

func init() {
	// Assert types
	var _ = []types.CurrencyIn{&stellar.Stellar{}}

}
