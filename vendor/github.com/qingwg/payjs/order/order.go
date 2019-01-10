package order

import (
	"github.com/qingwg/payjs/context"
)

// Order struct
type Order struct {
	*context.Context
}

//NewOrder init
func NewOrder(context *context.Context) *Order {
	order := new(Order)
	order.Context = context
	return order
}
