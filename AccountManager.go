package GoStore

import (
	"github.com/youricorocks/shop_competition"
	sorting "sort"
)

func (account Accounts) Register(username string) error {
	if _, ok := account[username]; ok {
		return StoreError{error: "user " + username + " already registered"}
	}

	account[username] = &shop_competition.Account{
		Name:        username,
		Balance:     0,
		AccountType: shop_competition.AccountNormal,
	}

	return nil
}

func (account Accounts) EditType(username string, accountType shop_competition.AccountType) error {
	account[username].AccountType = accountType

	return nil
}

func (account Accounts) AddBalance(username string, sum float32) error {
	if _, ok := account[username]; !ok {
		return StoreError{error: "can't AddBalance to user " + username + ", that does not exist"}
	}
	if sum <= 0 {
		return StoreError{error: "no positive value can't be added"}
	}

	account[username].Balance += sum

	return nil
}

func (account Accounts) Balance(username string) (float32, error) {
	if _, ok := account[username]; !ok {
		return 0, StoreError{error: "can't read balance from user " + username + ", that does not exist"}
	}

	return account[username].Balance, nil
}

func (account Accounts) GetAccounts(sort shop_competition.AccountSortType) []shop_competition.Account {
	accounts := make([]shop_competition.Account, len(account))

	var i int
	for _, val := range account {
		accounts[i] = *val
		i++
	}

	switch sort {
	case shop_competition.SortByName:
		sorting.Slice(accounts, func(i, j int) bool { return accounts[i].Name < accounts[j].Name })
	case shop_competition.SortByNameReverse:
		sorting.Slice(accounts, func(i, j int) bool { return accounts[i].Name > accounts[j].Name })
	case shop_competition.SortByBalance:
		sorting.Slice(accounts, func(i, j int) bool { return accounts[i].Balance < accounts[j].Balance })
	}

	return accounts
}
