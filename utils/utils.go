package utils

import (
	"strconv"

	"scalper/config"
)

func Abs(a string) string {
	if a[0:1] == "-" {
		return a[1:]
	}
	return a
}

func FormatPrice(price float64) string {
	return strconv.FormatFloat(price, 'f', config.Param.Symbol.Precision.Price, 64)
}

func FormatQuantity(price float64) string {
	return strconv.FormatFloat(price, 'f', config.Param.Symbol.Precision.Quantity, 64)
}

func QuantityEqual(a, b float64) bool {
	return FormatQuantity(a) == FormatQuantity(b)
}

func PriceEqual(a, b float64) bool {
	return FormatPrice(a) == FormatPrice(b)
}

func PriceZero() string {
	return FormatPrice(0.0)
}

func QuantityZero() string {
	return FormatQuantity(0.0)
}
