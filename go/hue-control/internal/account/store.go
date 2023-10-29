package account

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type AccountDto struct {
	ApiKey string
}

func (accountDto AccountDto) toJson() ([]byte, error) {
	return json.Marshal(accountDto)
}

func dtoFromAccount(account Account) AccountDto {
	return AccountDto{
		ApiKey: account.apiKey,
	}
}

func jsonFromAccount(account Account) ([]byte, error) {
	return dtoFromAccount(account).toJson()
}

type AccountJsonStore struct {
	FilePath string
}

func (store AccountJsonStore) Save(account Account) error {
	accountJson, err := jsonFromAccount(account)
	if err != nil {
		log.Error().Msg("Could create json for account")
		return err
	}
	return store.createOrUpdateFile(accountJson)
}

func (store AccountJsonStore) createOrUpdateFile(accountJson []byte) error {
	dir := filepath.Dir(store.FilePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Error().Msg("Could not create account file directory")
		return err
	}
	file, err := os.Create(store.FilePath)
	if err != nil {
		log.Error().Msg("Could not create account file")
		return err
	}
	defer file.Close()
	_, err = file.Write(accountJson)
	if err != nil {
		log.Error().Msg("Could not write account file")
	}
	return err
}
