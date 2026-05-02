package domain

// ID identifies a StockWise domain entity.
type ID uint64

// IsZero reports whether the identifier has not been assigned yet.
func (id ID) IsZero() bool {
	return id == 0
}
