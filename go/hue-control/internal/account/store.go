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

func (accountDto AccountDto) toAccount() Account {
	return InitAccount(accountDto.ApiKey)
}

func dtoFromAccount(account Account) AccountDto {
	return AccountDto{
		ApiKey: account.apiKey,
	}
}

func jsonFromAccount(account Account) ([]byte, error) {
	return dtoFromAccount(account).toJson()
}

func dtoFromJson(accountJson []byte) (*AccountDto, error) {
	var accountDto AccountDto
	err := json.Unmarshal(accountJson, &accountDto)
	if err != nil {
		log.Error().Msg("Could not parsed account")
		return nil, err
	}
	return &accountDto, nil
}

func accountFromJson(accountJson []byte) (*Account, error) {
	accountDto, err := dtoFromJson(accountJson)
	if err != nil {
		return nil, err
	}
	account := accountDto.toAccount()
	return &account, nil
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

func (store AccountJsonStore) Load() (*Account, error) {
	data, err := os.ReadFile(store.FilePath)
	if err != nil {
		log.Error().Msg("Could not read account file")
		return nil, err
	}
	return accountFromJson(data)
}
