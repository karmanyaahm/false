package config

type Config struct {
	ProductionNet bool
	Stellar       struct {
		PubKey  string
		PrivKey string
	}
}

var config Config

func init() {
	config.Stellar.PubKey = "GBAHNGNUSW52NZ5JYECOVAETFXPSKJREOEXMMF6BIRQF7GLDCKWDAKSU"
}

func Get() Config {
	return config
}
