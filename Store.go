package GoStore

import (
	"bytes"
	"encoding/json"
	"github.com/youricorocks/shop_competition"
)

const (
	E errStatus = iota // Error (critical)
	I                  // Info (Warnings)
)

type errStatus uint8

type ErrStatus errStatus

type StoreError struct {
	error string
	errStatus
}

type ProductList map[string]*shop_competition.Product
type Accounts map[string]*shop_competition.Account
type Orders map[string]float32
type Bundles map[string]*shop_competition.Bundle

type Currency int64

type Store struct {
	ProductList
	Accounts
	Orders
	Bundles
}

func (err StoreError) Error() string {
	return err.error
}

func (store *Store) CalculateOrder(username string, order shop_competition.Order) (sum float32, err error) {
	var storeErr StoreError

	user, ok := store.Accounts[username]
	if !ok {
		return 0, StoreError{error: "can't CalculateOrder for user " + username}
	}

	if key, er := cache(order); er == nil {
		if v, ok := store.Orders[key]; ok {
			return v, nil
		}
	} else {
		storeErr.error += er.(StoreError).error
		storeErr.errStatus = I
	}

	var currSum Currency

	// counting products in the order
	productCount := len(order.Products)
	for i := 0; i < productCount; i++ {
		coef := float32(1)

		// if products contain sampler product
		if order.Products[i].Type == shop_competition.ProductSample {
			storeErr.error += "product " + order.Products[i].Name + " is missed cause it's invalid type\n"
			storeErr.errStatus = I

			order.Products = removeProduct(order.Products, i)
			productCount--

			if key, er := cache(order); er == nil {
				if v, ok := store.Orders[key]; ok {
					return v, storeErr
				}
			} else {
				storeErr.error += er.(StoreError).error
				storeErr.errStatus = I
			}
		}

		// get coefficient
		if user.AccountType == shop_competition.AccountNormal && order.Products[i].Type == shop_competition.ProductPremium {
			coef = 1.5
		} else if user.AccountType == shop_competition.AccountPremium && order.Products[i].Type == shop_competition.ProductNormal {
			coef = 0.80
		} else if user.AccountType == shop_competition.AccountPremium && order.Products[i].Type == shop_competition.ProductPremium {
			coef = 0.95
		}

		currSum += roundToCur(order.Products[i].Price * coef)
	}

	// now use it to count bundles
	productCount = len(order.Bundles)
	for i := 0; i < productCount; i++  {
		var bunSum Currency
		if order.Bundles[i].Type == shop_competition.BundleNormal {
			for _, p := range order.Bundles[i].Products {
				if p.Type == shop_competition.ProductSample {
					storeErr.error += "bundle with " + p.Name + " is missed cause it's invalid type\n"
					storeErr.errStatus = I

					order.Bundles = removeBundle(order.Bundles, i)
					productCount--

					if key, er := cache(order); er == nil {
						if v, ok := store.Orders[key]; ok {
							return v, storeErr
						}
					} else {
						storeErr.error += er.(StoreError).error
						storeErr.errStatus = I
					}
					bunSum = 0
					break
				}

				bunSum += roundToCur(p.Price * order.Bundles[i].Discount)
			}
		} else {
			var normal, sample uint8
			if len(order.Products) != 2 {
				storeErr.error += "one of sampleBundles has different than two products\n"
				storeErr.errStatus = I

				order.Bundles = removeBundle(order.Bundles, i)
				productCount--

				if key, er := cache(order); er == nil {
					if v, ok := store.Orders[key]; ok {
						return v, storeErr
					}
				} else {
					storeErr.error += er.(StoreError).error
					storeErr.errStatus = I
				}
				continue
			}

			for _, p := range order.Products {
				if p.Type == shop_competition.ProductNormal || p.Type == shop_competition.ProductPremium {
					normal++
				} else {
					sample++
				}

				if normal > 1 || sample > 1 {
					storeErr.error += "bundle with " + p.Name + " is missed cause it's invalid type\n"
					storeErr.errStatus = I

					order.Bundles = removeBundle(order.Bundles, i)
					productCount--

					if key, er := cache(order); er == nil {
						if v, ok := store.Orders[key]; ok {
							return v, storeErr
						}
					} else {
						storeErr.error += er.(StoreError).error
						storeErr.errStatus = I
					}
					bunSum = 0
					break
				}
				bunSum += roundToCur(p.Price * order.Bundles[i].Discount)
			}
		}
		currSum += bunSum
	}


	if key, er := cache(order); er == nil {
		if v, ok := store.Orders[key]; ok {
			return v, storeErr
		} else {
			sum = float32(currSum) / 100
			store.Orders[key] = sum
		}
	} else {
		storeErr.error += er.(StoreError).error
		storeErr.errStatus = I
	}


	return sum, storeErr
}

func (store *Store) PlaceOrder(username string, order shop_competition.Order) error {
	var stErr StoreError

	user, ok := store.Accounts[username]
	if !ok {
		return StoreError{error: "can't PlaceOrder to user " + username + ", that does not exist"}
	}
	sum, err := store.CalculateOrder(username, order)
	if err != nil && err.(StoreError).errStatus == E {
		return err
	} else if err != nil {
		stErr.error += err.(StoreError).error
		stErr.errStatus = I
	}

	if roundToCur(user.Balance) < roundToCur(sum) {
		stErr.error += "user has insufficient balance\n"
		return stErr
	}

	user.Balance = float32(roundToCur(user.Balance) - roundToCur(sum)) / 100
	return stErr
}

func (store *Store) Import(data []byte) error {
	return json.Unmarshal(data, store)
}

func (store *Store) Export() ([]byte, error) {
	reqBodyBytes := new(bytes.Buffer)
	if err := json.NewEncoder(reqBodyBytes).Encode(StoreError{}); err != nil {
		return nil, StoreError{error: err.Error()}
	}

	return reqBodyBytes.Bytes(), nil
}

func NewStore() *Store {
	return &Store{ProductList: make(map[string]*shop_competition.Product),
		Accounts: make(map[string]*shop_competition.Account),
		Orders: make(map[string]float32),
		Bundles: make(map[string]*shop_competition.Bundle)}
}

//удаления с сохранением сортировки
func removeProduct(products [] shop_competition.Product, i int) []shop_competition.Product {
	return append(products[:i], products[i+1:]...)
}

func removeBundle(bundles [] shop_competition.Bundle, i int) []shop_competition.Bundle {
	return append(bundles[:i], bundles[i+1:]...)
}

func roundToCur(fl float32) Currency {
	return Currency((fl * 100) + 0.5)
}

func cache(order shop_competition.Order) (string, error) {
	cache, err := json.Marshal(order)
	if err == nil {
		return string(cache), nil
	} else {
		return "", StoreError{error: "json parse error"}
	}
}
