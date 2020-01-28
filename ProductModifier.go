package GoStore

import "github.com/youricorocks/shop_competition"

func (list ProductList) AddProduct(product shop_competition.Product) error {
	if product.Price <= 0 && product.Type != shop_competition.ProductSample {
		return StoreError{error: "price must contain positive value"}
	} else if product.Type == shop_competition.ProductSample {
		return StoreError{error: "sample product can not have price"}
	}
	if _, ok := list[product.Name]; ok {
		return StoreError{error: "product has already declared"}
	}
	if product.Type >= shop_competition.ProductSample {
		return StoreError{error: "unknown product type"}
	}

	list[product.Name] = &product

	return nil
}

func (list ProductList) ModifyProduct(product shop_competition.Product) error {
	if product.Price <= 0 && product.Type != shop_competition.ProductSample {
		return StoreError{error: "price must contain positive value"}
	} else if product.Type == shop_competition.ProductSample {
		return StoreError{error: "sample product can not have price"}
	}
	if _, ok := list[product.Name]; !ok {
		return StoreError{error: "can't modify product " + product.Name + ", that does not exist"}
	}
	if product.Type >= shop_competition.ProductSample {
		return StoreError{error: "unknown product type"}
	}

	list[product.Name] = &product

	return nil
}

func (list ProductList) RemoveProduct(name string) error {
	if _, ok := list[name]; !ok {
		return StoreError{error: "can't delete product " + name + ", that does not exist"}
	}

	delete(list, name)

	return nil
}
