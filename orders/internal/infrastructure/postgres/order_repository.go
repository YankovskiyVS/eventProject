package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/YankovskiyVS/eventProject/orders/internal/domain/entities"
	"github.com/YankovskiyVS/eventProject/orders/internal/domain/repositories"
	"github.com/google/uuid"
)

type OrderRepository struct {
	orderDB  *sql.DB
	ticketDB *sql.DB
}

func NewOrderRepository(orderDB *sql.DB, ticketDB *sql.DB) repositories.OrderRepository {
	return &OrderRepository{
		orderDB:  orderDB,
		ticketDB: ticketDB,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *entities.Order) error {
	tx, err := r.orderDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Save order
	orderPG := OrderPostgres{
		UUID:        order.ID().String(),
		OrderStatus: string(order.OrderStatus()),
		UserID:      uint(order.UserID()),
		CreatedAt:   order.CreatedAt(),
		UpdatedAt:   order.UpdatedAt(),
	}

	err = tx.QueryRowContext(ctx, `
        INSERT INTO orders (uuid, order_status, user_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, orderPG.UUID, orderPG.OrderStatus, orderPG.UserID,
		orderPG.CreatedAt, orderPG.UpdatedAt).Scan(&orderPG.ID)
	if err != nil {
		return fmt.Errorf("error saving order: %w", err)
	}

	// Save tickets
	for _, ticket := range order.Tickets() {
		ticketPG := TicketPostgres{
			UUID:    ticket.ID().String(),
			EventID: uint(ticket.EventID()),
			Price:   ticket.Price(),
			OrderID: orderPG.ID,
		}

		_, err = tx.ExecContext(ctx, `
            INSERT INTO tickets (uuid, event_id, price, order_id)
            VALUES ($1, $2, $3, $4)
        `, ticketPG.UUID, ticketPG.EventID, ticketPG.Price, ticketPG.OrderID)
		if err != nil {
			return fmt.Errorf("error saving ticket: %w", err)
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) Update(ctx context.Context, order *entities.Order) (*entities.Order, error) {
	return nil, nil
}

func (r *OrderRepository) GetCurrent(ctx context.Context, userID int) (*entities.Order, error) {
	return nil, nil
}

func (r *OrderRepository) PayCurrent(ctx context.Context, order *entities.Order) (*entities.Order, error) {
	return nil, nil
}

func (r *OrderRepository) FindByID(ctx context.Context, userID int, orderUUID string) (*entities.Order, error) {
	// Get order
	var orderPG OrderPostgres
	err := r.orderDB.QueryRowContext(ctx, `
        SELECT id, uuid, order_status, user_id, created_at, updated_at
        FROM orders
        WHERE uuid = $1 AND user_id = $2
    `, orderUUID, userID).Scan(
		&orderPG.ID,
		&orderPG.UUID,
		&orderPG.OrderStatus,
		&orderPG.UserID,
		&orderPG.CreatedAt,
		&orderPG.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching order: %w", err)
	}

	// Get tickets
	rows, err := r.orderDB.QueryContext(ctx, `
        SELECT uuid, event_id, price
        FROM tickets
        WHERE order_id = $1
    `, orderPG.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching tickets: %w", err)
	}
	defer rows.Close()

	var tickets []*entities.Ticket
	for rows.Next() {
		var ticketPG TicketPostgres
		err := rows.Scan(
			&ticketPG.UUID,
			&ticketPG.EventID,
			&ticketPG.Price,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning ticket: %w", err)
		}

		ticket, err := r.convertToDomainTicket(ticketPG)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	return r.convertToDomainOrder(orderPG, tickets)
}

func (r *OrderRepository) FindAll(ctx context.Context, userID int) ([]*entities.Order, error) {
	return nil, nil
}

// Private conversion methods
func (r *OrderRepository) convertToDomainOrder(pg OrderPostgres, tickets []*entities.Ticket) (*entities.Order, error) {
	id, err := uuid.Parse(pg.UUID)
	if err != nil {
		return nil, fmt.Errorf("invalid order UUID: %w", err)
	}

	return entities.SetOrderWithDetails(
		id,
		int(pg.UserID),
		tickets,
		pg.OrderStatus,
		pg.CreatedAt,
		pg.UpdatedAt,
	)
}

func (r *OrderRepository) convertToDomainTicket(pg TicketPostgres) (*entities.Ticket, error) {
	id, err := uuid.Parse(pg.UUID)
	if err != nil {
		return nil, fmt.Errorf("invalid ticket UUID: %w", err)
	}

	return entities.SetTicketWithDetails(
		id,
		float32(pg.Price),
		int(pg.EventID),
	), nil
}
