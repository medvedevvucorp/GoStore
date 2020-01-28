package GoStore

import (
	"github.com/youricorocks/shop_competition"
	"reflect"
	"testing"
)

func TestStore_AddProduct(t *testing.T) {
	var shop shop_competition.Shop
	shop = NewStore()

	product := shop_competition.Product{
		Name:  "Banana",
		Price: 100,
		Type:  shop_competition.ProductNormal,
	}

	if err := shop.AddProduct(product); err != nil && err.(StoreError).errStatus == E {
		t.Error(err)
	}
	if product != *shop.(*Store).ProductList[product.Name] {
		t.Error(StoreError{ error:"AddProduct(product) result doesn't match with \"product\" " })
	}
	expected := StoreError{error: "product has already declared"}
	if err := shop.AddProduct(product); !reflect.DeepEqual(expected, err) {
		t.Errorf("Looking for '%v', got '%v' ", expected.error, err.(StoreError).error)
	}
}

func TestStore_ModifyProduct(t *testing.T) {
	var shop shop_competition.Shop
	shop = NewStore()

	product := shop_competition.Product{
		Name:  "Banana",
		Price: 100,
		Type:  shop_competition.ProductNormal,
	}
	if err := shop.AddProduct(product); err != nil && err.(StoreError).errStatus == E {
		t.Error(err)
	}

	product.Price = 200

	if err := shop.ModifyProduct(product); err != nil && err.(StoreError).errStatus == E {
		t.Error(err)
	}

	if product != *shop.(*Store).ProductList[product.Name] {
		t.Error(StoreError{ error:"ModifyProduct(product) result doesn't match with \"product\" " })
	}

	product.Name = "Apple"
	expected := StoreError{error: "can't modify product " + product.Name + ", that does not exist"}
	if err := shop.ModifyProduct(product); !reflect.DeepEqual(expected, err) {
		t.Errorf("Looking for '%v', got '%v' ", expected.error, err.(StoreError).error)
	}

	product.Name = "Banana"
	product.Price = -200
	expected = StoreError{error: "price must contain positive value"}
	if err := shop.ModifyProduct(product); !reflect.DeepEqual(expected, err) {
		t.Errorf("Looking for '%v', got '%v' ", expected.error, err.(StoreError).error)
	}
}

func TestStore_RemoveProduct(t *testing.T) {
	var shop shop_competition.Shop
	shop = NewStore()

	product := shop_competition.Product{
		Name:  "Banana",
		Price: 100,
		Type:  shop_competition.ProductNormal,
	}
	if err := shop.AddProduct(product); err != nil && err.(StoreError).errStatus == E {
		t.Error(err)
	}

	if err := shop.RemoveProduct(product.Name); err != nil && err.(StoreError).errStatus == E {
		t.Error(err)
	}

	if _, ok := shop.(*Store).ProductList["Banana"]; ok != false {
		t.Error("Deleting did not give result")
	}

	expected := StoreError{error: "can't delete product " + product.Name + ", that does not exist"}
	if err := shop.RemoveProduct(product.Name); !reflect.DeepEqual(expected, err) {
		t.Errorf("Looking for '%v', got '%v' ", expected.error, err.(StoreError).error)
	}
}

func TestStore_CalculateOrder(t *testing.T) {
	var shop shop_competition.Shop
	shop = NewStore()

	shop.Register("Dimas")
	shop.AddBalance("Dimas", 100)

	shop.(*Store).Accounts.EditType("Dimas", shop_competition.AccountPremium)

	products := []shop_competition.Product{
		shop_competition.Product{
			Name :  "banana",
			Price: 40,
			Type:  0,
		},
		shop_competition.Product{
			Name :  "apple",
			Price: 30,
			Type:  shop_competition.ProductPremium,
		},
	}

	order := shop_competition.Order{
		Products: products,
		Bundles:  nil,
	}

	if val, err := shop.CalculateOrder("Dimas", order); val != 60.5 {
		t.Error(val, err)
	} else { t.Log(val, err) }

	if val, err := shop.CalculateOrder("Dimas", order); val != 60.5 {
		t.Error(val, err)
	} else { t.Log(val, err) }

	products = []shop_competition.Product{
		shop_competition.Product{
			Name :  "banana",
			Price: 40,
			Type:  0,
		},
		shop_competition.Product{
			Name :  "sapre",
			Price: 30,
			Type:  shop_competition.ProductSample,
		},
		shop_competition.Product{
			Name :  "apple",
			Price: 30,
			Type:  shop_competition.ProductPremium,
		},
	}

	order = shop_competition.Order{
		Products: products,
		Bundles:  nil,
	}

	if val, err := shop.CalculateOrder("Dimas", order); val != 60.5 {
		t.Error(val, err)
	} else { t.Log(val, err) }

	bundles := []shop_competition.Bundle {
		{ Products: products, Type: shop_competition.BundleSample, Discount: 32 },
	}

	order = shop_competition.Order{
		Products: nil,
		Bundles:  bundles,
	}

	if val, err := shop.CalculateOrder("Dimas", order); err == nil {
		t.Error(val, err)
	} else { t.Log(val, err) }
}

func TestStore_PlaceOrder(t *testing.T) {
	var shop shop_competition.Shop
	shop = NewStore()

	shop.Register("Dimas")
	shop.AddBalance("Dimas", 100)

	shop.(*Store).Accounts.EditType("Dimas", shop_competition.AccountPremium)

	products := []shop_competition.Product{
		shop_competition.Product{
			Name :  "banana",
			Price: 40,
			Type:  0,
		},
		shop_competition.Product{
			Name :  "sapre",
			Price: 30,
			Type:  shop_competition.ProductSample,
		},
		shop_competition.Product{
			Name :  "apple",
			Price: 30,
			Type:  shop_competition.ProductPremium,
		},
	}

	bundles := []shop_competition.Bundle {
		{ Products: products, Type: shop_competition.BundleSample, Discount: 0.1 },
	}

	order := shop_competition.Order{
		Products: products,
		Bundles:  bundles,
	}

	if err := shop.PlaceOrder("Dimas", order); err != nil && err.(StoreError).errStatus == E {
		t.Error(err)
	} else if shop.(*Store).Accounts["Dimas"].Balance == 39.5 {
		t.Log(err)
	}

	//if shop.(*Store).Accounts["Dimas"].Balance
}

func TestStore_Register(t *testing.T) {
	var shop shop_competition.Shop
	shop = NewStore()

	if err := shop.Register("Pavel_007"); err != nil && err.(StoreError).errStatus == E  {
		t.Error(err)
	}

	if _, ok := shop.(*Store).Accounts["Pavel_007"]; ok == false {
		t.Error("Registration unsuccessful")
	}

	expected := StoreError{error: "user Pavel_007 already registered"}
	if err := shop.Register("Pavel_007"); !reflect.DeepEqual(expected, err) {
		t.Errorf("Looking for '%v', got '%v' ", expected.error, err.(StoreError).error)
	}
}

func TestStore_AddBalance(t *testing.T) {
	var shop shop_competition.Shop
	shop = NewStore()

	if err := shop.Register("Pavel_007"); err != nil && err.(StoreError).errStatus == E  {
		t.Error(err)
	}

	if err := shop.AddBalance("Pavel_007", 30); err != nil && err.(StoreError).errStatus == E {
		t.Error(err)
	}

	if shop.(*Store).Accounts["Pavel_007"].Balance != 30 {
		t.Error("AddBalance result is wrong")
	}

	expected := StoreError{error: "no positive value can't be added"}
	if err := shop.AddBalance("Pavel_007", -19); !reflect.DeepEqual(expected, err) {
		t.Errorf("Looking for '%v', got '%v' ", expected.error, err.(StoreError).error)
	}
}

func TestStore_Balance(t *testing.T) {
	var shop shop_competition.Shop
	shop = NewStore()

	if err := shop.Register("Pavel_007"); err != nil && err.(StoreError).errStatus == E  {
		t.Error(err)
	}

	if err := shop.AddBalance("Pavel_007", 30); err != nil && err.(StoreError).errStatus == E {
		t.Error(err)
	}

	if res, err := shop.Balance("Pavel_007"); err != nil &&  err.(StoreError).errStatus == E {
		t.Error(err)
	} else if res != 30 {
		t.Error("Balance result is wrong")
	}
}

func TestStore_GetAccounts(t *testing.T) {
	var shop shop_competition.Shop
	shop = NewStore()

	if err := shop.Register("Pavel_007"); err != nil && err.(StoreError).errStatus == E  {
		t.Error(err)
	}
	if err := shop.Register("Mashka"); err != nil && err.(StoreError).errStatus == E  {
		t.Error(err)
	}
	if err := shop.AddBalance("Pavel_007", 30); err != nil && err.(StoreError).errStatus == E {
		t.Error(err)
	}
	if err := shop.Register("Dimas"); err != nil && err.(StoreError).errStatus == E  {
		t.Error(err)
	}
	if err := shop.AddBalance("Dimas", 150); err != nil && err.(StoreError).errStatus == E {
		t.Error(err)
	}

	accounts := shop.GetAccounts(shop_competition.SortByName)
	for i := 1; i < len(accounts); i++ {
		if accounts[i - 1].Name > accounts[i].Name {
			t.Error("wrong sorting by Name: ", accounts)
		}
	}

	accounts = shop.GetAccounts(shop_competition.SortByNameReverse)
	for i := 1; i < len(accounts); i++ {
		if accounts[i - 1].Name < accounts[i].Name {
			t.Error("wrong sorting by NameReverse: ", accounts)
		}
	}

	accounts = shop.GetAccounts(shop_competition.SortByBalance)
	for i := 1; i < len(accounts); i++ {
		if accounts[i - 1].Balance > accounts[i].Balance {
			t.Error("wrong sorting by Balance: ", accounts)
		}
	}
}

