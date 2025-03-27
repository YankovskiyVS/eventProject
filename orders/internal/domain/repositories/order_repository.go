package repositories

import (
	"github.com/YankovskiyVS/eventProject/orders/internal/domain/entities"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(order *entities.Order) (entities.Order, error)
	Update(order *entities.Order) (*entities.Order, error)
	FindAll(userID int) ([]*entities.Order, error)
	FindByID(userID int, id uuid.UUID) (*entities.Order, error)
	GetCurrent(userID int) (*entities.Order, error)
	PayCurrent(order *entities.Order) (*entities.Order, error)
}
