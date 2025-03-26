package entities

type ValidatedOrder struct {
	Order
	isValidated bool
}

func NewValidateddOrder(order *Order) (*ValidatedOrder, error) {
	if err := order.validate(); err != nil {
		return nil, err
	}
	return &ValidatedOrder{
		Order:       *order,
		isValidated: true,
	}, nil
}
