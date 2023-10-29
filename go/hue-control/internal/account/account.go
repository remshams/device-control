package account

type AccountStore interface {
	Save(account Account) error
	Load() (Account error)
}

type Account struct {
	apiKey string
}

func Init(apiKey string) Account {
	return Account{apiKey: apiKey}
}
