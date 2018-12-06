package order

import (
	"encoding/json"
	"fmt"
	"github.com/yuyan2077/payjs/context"
	"github.com/yuyan2077/payjs/util"
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
