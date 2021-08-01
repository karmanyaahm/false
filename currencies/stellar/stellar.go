package stellar

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/karmanyaahm/payment-thing/config"
	"github.com/karmanyaahm/payment-thing/currencies/types"
	"github.com/karmanyaahm/payment-thing/db"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	hProtocol "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/protocols/horizon/effects"
	"github.com/stellar/go/txnbuild"
	"gorm.io/gorm"
)

type Stellar struct {
	client   *horizonclient.Client
	db       *gorm.DB
	callback types.TxCallback
	key      keypair.KP

	passphrase string
}

type StellarUpdatedAt struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	Cursor    string
	Version   uint //will be zero by default so all good, can use in future to update pages
}

//only ever ever ever run one of these even in multi processes TODO figure out multi runs locks ever ever
var lock = sync.Mutex{}

func New(db *gorm.DB, key types.Key, cb types.TxCallback) *Stellar {
	s := &Stellar{}

	if config.Get().ProductionNet {
		s.client = horizonclient.DefaultPublicNetClient
		s.passphrase = network.PublicNetworkPassphrase
	} else {
		s.client = horizonclient.DefaultTestNetClient
		s.passphrase = network.TestNetworkPassphrase
	}

	if key.Priv != "" {
		s.key = keypair.MustParse(key.Priv)
	} else if key.Pub != "" {
		s.key = keypair.MustParseAddress(key.Pub)
	} else {
		panic("NO STELLAR KEY")
	}

	s.callback = cb
	s.db = db
	if db == nil {
		panic("AAA DB NIL STELLAR")
	}
	err := s.db.AutoMigrate(&StellarUpdatedAt{})
	if err != nil {
		panic(err)
	}
	return s
}

func (s *Stellar) Out(key string) {
	fullKey, ok := s.key.(*keypair.Full)
	if !ok {
		//TODO what if there's no privkey avaiablbale
	}

	targetKey, _ := keypair.ParseAddress(key)

	ar := horizonclient.AccountRequest{AccountID: s.key.Address()}
	sourceAccount, err := s.client.AccountDetail(ar)
	if err != nil {
		return
	}

	op := txnbuild.Payment{
		Destination: targetKey.Address(),
		Amount:      "1",
		Asset:       txnbuild.NativeAsset{},
	}

	// Construct the transaction that holds the operations to execute on the network
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Sign the transaction
	tx, err = tx.Sign(s.passphrase, fullKey)
	if err != nil {
		log.Fatalln(err)
	}

	txResult, err := s.client.SubmitTransaction(tx)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(txResult)

}

func (s *Stellar) In() {
	for {
		lastCursor := StellarUpdatedAt{}
		s.db.Order("created_at desc").Attrs(&StellarUpdatedAt{Cursor: "0"}).FirstOrCreate(&lastCursor)
		// all transactions
		transactionRequest := horizonclient.TransactionRequest{Cursor: lastCursor.Cursor, ForAccount: s.key.Address()}

		//LOG
		fmt.Println("stellar listening TXs")
		lock.Lock()
		err := s.client.StreamTransactions(context.Background(), transactionRequest, s.txHandler)
		lock.Unlock()
		println("stellar listening stopped??? sleeping and restarting")
		if err != nil {
			fmt.Println(err)
			//TODO some sort of error counter
		}
		time.Sleep(15 * time.Second)
	}
}

func (s *Stellar) txHandler(tx hProtocol.Transaction) {
	effReq := horizonclient.EffectRequest{ForTransaction: tx.ID}
	ops, err := s.client.Effects(effReq)
	if err != nil {
	}
	s.db.Create(&StellarUpdatedAt{Cursor: tx.PagingToken()})
	records := ops.Embedded.Records
	for _, op := range records {
		switch opp := op.(type) {
		case effects.AccountCredited:
			//fmt.Println(opp.Amount)
			if opp.Asset.Type == "native" {
				s.callback(op.GetID(), s.getUID(tx), int64(s.parseAmt(opp.Amount)))
			}
		default:
			fmt.Printf("%T\n", op)
		}
	}
}

func (s *Stellar) GenAddr(UserID uint64) string {
	return fmt.Sprintf("%s - Memo Type: Text - Memo Value: %d", s.key.Address(), UserID)
}
func (s *Stellar) getUID(tr hProtocol.Transaction) db.UserID {
	if tr.MemoType != "text" || tr.MemoType != "MEMO_TEXT" {
		//log error
		return 0
	}
	num, err := strconv.Atoi(tr.Memo)
	if err != nil || num < 0 {
		//log
		return 0
	}
	return uint64(num)
}

//'XLM/7'
func (*Stellar) parseAmt(s string) uint64 {
	// remove dots and parse as int. stellar always returns 7 decimal places in the API //TODO it seems like but not surez
	s = strings.ReplaceAll(s, ".", "")
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		//LOG
		return 0
	}
	return uint64(n)
}

func (*Stellar) getAmt(i db.Money) string {
	str := strconv.Itoa(int(i))
	return str[:len(str)-7] + "." + str[len(str)-7:]
	//TODO need test prob
}
