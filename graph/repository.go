package graph

import (
	"test-fite/graph/customerror"
	"test-fite/graph/model"
)

var products = make(map[string]model.Product)

func InitProduct() {
	products["120P90"] = model.Product{
		Sku:   "120P90",
		Name:  "Google Home",
		Price: 49.99,
		Qty:   10,
	}

	products["43N23P"] = model.Product{
		Sku:   "43N23P",
		Name:  "MacBook Pro",
		Price: 5399.99,
		Qty:   5,
	}

	products["A304SD"] = model.Product{
		Sku:   "A304SD",
		Name:  "Alexa Speaker",
		Price: 109.50,
		Qty:   10,
	}

	products["234234"] = model.Product{
		Sku:   "234234",
		Name:  "Raspberry Pi B",
		Price: 30.0,
		Qty:   2,
	}
}

func getProductBySKU(sku string) (*model.Product, error) {
	if _, ok := products[sku]; ok {
		result := products[sku]
		return &result, nil
	} else {
		return nil, customerror.ErrProductNotFound
	}
}

func processCheckout(checkoutData []*model.NewCheckout) (*model.ProductTotal, error) {
	var result model.ProductTotal
	var total = 0.0
	productData := make(map[string]int)
	for _, value := range checkoutData {
		if _, ok := productData[value.Sku]; !ok {
			productData[value.Sku] = value.Qty
		} else {
			productData[value.Sku] += value.Qty
		}

		if data, ok := products[value.Sku]; ok {
			if data.Qty < productData[value.Sku] && productData[value.Sku] != 0 {
				return nil, customerror.ErrProductOutOfStock
			}
		} else {
			return nil, customerror.ErrProductNotFound
		}
	}

	// Google Home
	if _, ok := productData["120P90"]; ok {
		if data, ok := products["120P90"]; ok {
			data.Qty = data.Qty - productData["120P90"]
			products["120P90"] = data

			result.Products = append(result.Products, &model.ProductQuantity{
				Product: &model.Product{
					Sku:   data.Sku,
					Name:  data.Name,
					Price: data.Price,
					Qty:   productData["120P90"],
				},
				Qty: productData["120P90"],
			})
		}

		counter := 1
		for i := 0; i < productData["120P90"]; i++ {
			if counter == 3 {
				counter = 0
			} else {
				total += products["120P90"].Price
			}

			counter++
		}
	}

	// Alexa Speaker
	if _, ok := productData["A304SD"]; ok {
		if data, ok := products["A304SD"]; ok {
			data.Qty = data.Qty - productData["A304SD"]
			products["A304SD"] = data
			discount := data.Price * 0.1

			for i := 0; i < productData["A304SD"]; i++ {
				prices := data.Price
				if productData["A304SD"] > 3 {
					prices = prices - discount
				}

				total += prices
			}

			result.Products = append(result.Products, &model.ProductQuantity{
				Product: &model.Product{
					Sku:   data.Sku,
					Name:  data.Name,
					Price: data.Price,
					Qty:   productData["A304SD"],
				},
				Qty: productData["A304SD"],
			})
		}
	}

	// Macbook
	if _, ok := productData["43N23P"]; ok {
		if data, ok := products["43N23P"]; ok {

			data.Qty = data.Qty - productData["43N23P"]
			products["43N23P"] = data
			result.Products = append(result.Products, &model.ProductQuantity{
				Product: &model.Product{
					Sku:   data.Sku,
					Name:  data.Name,
					Price: data.Price,
					Qty:   productData["43N23P"],
				},
				Qty: productData["43N23P"],
			})

			for i := 0; i < productData["43N23P"]; i++ {
				total += data.Price
			}
		}

		if data, ok := products["234234"]; ok {
			if data.Qty != 0 {
				counter := 0
				for i := 0; i < productData["43N23P"]; i++ {
					data.Qty--
					productData["43N23P"]--
					counter++
				}
				products["234234"] = data
				result.Products = append(result.Products, &model.ProductQuantity{
					Product: &model.Product{
						Sku:   data.Sku,
						Name:  data.Name,
						Price: data.Price,
						Qty:   counter,
					},
					Qty: counter,
				})
			}
		}
	}

	if _, ok := productData["234234"]; ok {
		if data, ok := products["234234"]; ok {
			data.Qty = data.Qty - productData["234234"]
			products["234234"] = data
			total += data.Price * float64(productData["234234"])

			isFound := false
			for index, product := range result.Products {
				if product.Product.Sku == "234234" {
					isFound = true
					result.Products[index].Qty += productData["234234"]
					result.Products[index].Product.Qty += productData["234234"]
					break
				}
			}
			if !isFound {
				result.Products = append(result.Products, &model.ProductQuantity{
					Product: &model.Product{
						Sku:   data.Sku,
						Name:  data.Name,
						Price: data.Price,
						Qty:   productData["234234"],
					},
					Qty: productData["234234"],
				})
			}
		}
	}

	return &model.ProductTotal{
		Products: result.Products,
		Totals:   total,
	}, nil
}
