package domain

// IsPositiveQuantity reports whether a quantity can be used in a stock operation.
func IsPositiveQuantity(quantity int64) bool {
	return quantity > 0
}

// IsNonNegativePriceCents reports whether a price represented in cents is valid.
func IsNonNegativePriceCents(priceCents int64) bool {
	return priceCents >= 0
}

// IsLowStock reports whether a positive quantity has reached its minimum threshold.
func IsLowStock(quantity int64, minimumThreshold int64) bool {
	return quantity > 0 && minimumThreshold > 0 && quantity <= minimumThreshold
}

// IsOutOfStock reports whether a stock quantity is depleted.
func IsOutOfStock(quantity int64) bool {
	return quantity <= 0
}

// EvaluateStockState classifies a stock quantity against the minimum threshold.
func EvaluateStockState(quantity int64, minimumThreshold int64) StockState {
	if IsOutOfStock(quantity) {
		return StockStateOut
	}

	if IsLowStock(quantity, minimumThreshold) {
		return StockStateLow
	}

	return StockStateInStock
}
