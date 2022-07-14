package graph

import (
	"test-fite/graph/customerror"
	"test-fite/graph/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetProductSKU(t *testing.T) {
	t.Run("be able to get products", func(t *testing.T) {
		InitProduct()

		data, err := getProductBySKU("120P90")
		require.Nil(t, err)
		require.NotNil(t, data)
		require.Equal(t, data.Name, "Google Home")
		require.Equal(t, data.Price, 49.99)
		require.Equal(t, data.Qty, 10)
	})

	t.Run("NOT be able to get products. data not found", func(t *testing.T) {
		InitProduct()

		data, err := getProductBySKU("120P900")
		require.Nil(t, data)
		require.Error(t, err)
		require.Equal(t, err, customerror.ErrProductNotFound)
	})
}

func TestProcessCheckout(t *testing.T) {
	t.Run("be able to checkout products", func(t *testing.T) {
		InitProduct()
		request := []*model.NewCheckout{
			{
				Sku: "120P90",
				Qty: 2,
			},
		}

		dataBefore, _ := getProductBySKU(request[0].Sku)

		data, err := processCheckout(request)
		require.Nil(t, err)
		require.NoError(t, err)
		require.Equal(t, data.Totals, 49.99*2)
		require.Equal(t, data.Products[0].Qty, request[0].Qty)
		require.Equal(t, data.Products[0].Product.Sku, request[0].Sku)

		dataAfter, _ := getProductBySKU(request[0].Sku)
		require.Equal(t, dataAfter.Qty, dataBefore.Qty-request[0].Qty)
	})

	t.Run("be able to checkout products when Google Home Discount", func(t *testing.T) {
		InitProduct()
		request := []*model.NewCheckout{
			{
				Sku: "120P90",
				Qty: 3,
			},
		}

		dataBefore, _ := getProductBySKU(request[0].Sku)

		data, err := processCheckout(request)
		require.Nil(t, err)
		require.NoError(t, err)
		require.Equal(t, data.Totals, 49.99*2)
		require.Equal(t, data.Products[0].Qty, request[0].Qty)
		require.Equal(t, data.Products[0].Product.Sku, request[0].Sku)

		dataAfter, _ := getProductBySKU(request[0].Sku)
		require.Equal(t, dataAfter.Qty, dataBefore.Qty-request[0].Qty)
	})

	t.Run("be able to checkout products when Google Home Discount with multiple same products", func(t *testing.T) {
		InitProduct()
		request := []*model.NewCheckout{
			{
				Sku: "120P90",
				Qty: 3,
			},
			{
				Sku: "120P90",
				Qty: 2,
			},
		}

		dataBefore, _ := getProductBySKU(request[0].Sku)

		data, err := processCheckout(request)
		require.Nil(t, err)
		require.NoError(t, err)
		require.Equal(t, data.Totals, 49.99*4)
		require.Equal(t, data.Products[0].Qty, request[0].Qty+request[1].Qty)
		require.Equal(t, data.Products[0].Product.Sku, request[0].Sku)

		dataAfter, _ := getProductBySKU(request[0].Sku)
		require.Equal(t, dataAfter.Qty, dataBefore.Qty-request[0].Qty-request[1].Qty)
	})

	t.Run("be able to checkout products when Order Macbook", func(t *testing.T) {
		InitProduct()
		request := []*model.NewCheckout{
			{
				Sku: "43N23P",
				Qty: 1,
			},
		}

		dataMacbookBefore, _ := getProductBySKU(request[0].Sku)
		dataRaspberryBefore, _ := getProductBySKU("234234")

		data, err := processCheckout(request)
		require.Nil(t, err)
		require.NoError(t, err)
		require.Equal(t, data.Totals, 5399.99)
		require.Equal(t, len(data.Products), 2)
		require.Equal(t, data.Products[0].Product.Sku, request[0].Sku)
		require.Equal(t, data.Products[0].Qty, 1)
		require.Equal(t, data.Products[1].Product.Sku, "234234")
		require.Equal(t, data.Products[1].Qty, 1)

		dataMacbookAfter, _ := getProductBySKU(request[0].Sku)
		dataRaspberryAfter, _ := getProductBySKU("234234")
		require.Equal(t, dataMacbookAfter.Qty, dataMacbookBefore.Qty-request[0].Qty)
		require.Equal(t, dataRaspberryAfter.Qty, dataRaspberryBefore.Qty-1)
	})

	t.Run("be able to checkout products when order Alexa Speaker", func(t *testing.T) {
		InitProduct()
		request := []*model.NewCheckout{
			{
				Sku: "A304SD",
				Qty: 4,
			},
		}

		discount := 109.50 * 0.1

		data, err := processCheckout(request)
		require.Nil(t, err)
		require.NoError(t, err)
		require.Equal(t, data.Totals, (109.50-discount)*4)
		require.Equal(t, data.Products[0].Product.Sku, request[0].Sku)
		require.Equal(t, data.Products[0].Qty, request[0].Qty)
	})

	t.Run("be able to checkout products when order Raspberry Pi B", func(t *testing.T) {
		InitProduct()
		request := []*model.NewCheckout{
			{
				Sku: "234234",
				Qty: 2,
			},
		}

		data, err := processCheckout(request)
		require.Nil(t, err)
		require.NoError(t, err)
		require.Equal(t, data.Totals, 30.0*2)
		require.Equal(t, data.Products[0].Product.Sku, request[0].Sku)
		require.Equal(t, data.Products[0].Qty, request[0].Qty)
	})

	t.Run("be able to checkout products when Order Macbook & Raspberry P", func(t *testing.T) {
		InitProduct()
		request := []*model.NewCheckout{
			{
				Sku: "43N23P",
				Qty: 1,
			},
			{
				Sku: "234234",
				Qty: 1,
			},
		}

		dataMacbookBefore, _ := getProductBySKU(request[0].Sku)
		dataRaspberryBefore, _ := getProductBySKU(request[1].Sku)

		data, err := processCheckout(request)
		require.Nil(t, err)
		require.NoError(t, err)
		require.Equal(t, data.Totals, 5399.99+30.0)
		require.Equal(t, len(data.Products), 2)
		require.Equal(t, data.Products[0].Product.Sku, request[0].Sku)
		require.Equal(t, data.Products[0].Qty, 1)
		require.Equal(t, data.Products[1].Product.Sku, "234234")
		require.Equal(t, data.Products[1].Qty, 2)

		dataMacbookAfter, _ := getProductBySKU(request[0].Sku)
		dataRaspberryAfter, _ := getProductBySKU("234234")
		require.Equal(t, dataMacbookAfter.Qty, dataMacbookBefore.Qty-request[0].Qty)
		require.Equal(t, dataRaspberryAfter.Qty, dataRaspberryBefore.Qty-2)
	})

	t.Run("NOT be able to checkout products when quantity over limit", func(t *testing.T) {
		InitProduct()
		request := []*model.NewCheckout{
			{
				Sku: "120P90",
				Qty: 20,
			},
		}

		data, err := processCheckout(request)
		require.Nil(t, data)
		require.Error(t, err)
		require.Equal(t, err, customerror.ErrProductOutOfStock)
	})
}
