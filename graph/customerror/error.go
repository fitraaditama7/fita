package customerror

import "errors"

var ErrProductNotFound = errors.New("Product not found")
var ErrProductOutOfStock = errors.New("Product out of stock")
