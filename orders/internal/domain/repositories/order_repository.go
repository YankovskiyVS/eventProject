package repositories

import (
	"context"

	"github.com/YankovskiyVS/eventProject/orders/internal/domain/entities"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) error
	Update(ctx context.Context, order *entities.Order) (*entities.Order, error)
	FindAll(ctx context.Context, userID int) ([]*entities.Order, error)
	FindByID(ctx context.Context, userID int, uuID string) (*entities.Order, error)
	GetCurrent(ctx context.Context, userID int) (*entities.Order, error)
	PayCurrent(ctx context.Context, order *entities.Order) (*entities.Order, error)
}
