package GoStore

import "github.com/youricorocks/shop_competition"

func (bundles Bundles) AddBundle(name string, main shop_competition.Product, discount float32, additional ...shop_competition.Product) error {
	if discount < 1 || discount > 99 {
		return StoreError{error: "discount is not correct"}
	}
	if _, ok := bundles[name]; ok {
		return StoreError{error: "cant't AddBundle with name " + name + ", that already exist"}
	}

	if len(additional) == 0 {
		return StoreError{error: "bundle has only one product"}
	} else if additional[0].Type != shop_competition.ProductSample && main.Type != shop_competition.ProductSample {
		for i := 1; i < len(additional); i++ {
			if additional[i].Type == shop_competition.ProductSample {
				return StoreError{error: "normal bundle has no samplers"}
			}
		}
	} else {
		if len(additional) > 1{
			return StoreError{error: "sample bundle has only two products"}
		}
		bundles[name].Type = shop_competition.BundleSample
	}

	bundles[name].Products = append(additional, main)
	bundles[name].Discount = bundles.getDiscount(discount)

	return nil
}

func (bundles Bundles) ChangeDiscount(name string, discount float32) error {
	if discount < 1 || discount > 99 {
		return StoreError{error: "discount is not correct"}
	}
	bundle, ok := bundles[name]
	if !ok {
		return StoreError{error: "can't get bundle " + name + ", that does not exist"}
	}

	bundle.Discount = bundles.getDiscount(discount)

	return nil
}

func (bundles Bundles) RemoveBundle(name string) error {
	if _, ok := bundles[name]; !ok {
		return StoreError{error: "bundle " + name + " doesnt exist"}
	}
	delete(bundles, name)

	return nil
}

func (bundles Bundles) getDiscount(in float32) float32 {
	return 1 - (in / 100)
}
