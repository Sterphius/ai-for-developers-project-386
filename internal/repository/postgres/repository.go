package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Sterphius/ai-for-developers-project-386/internal/domain"
	"github.com/Sterphius/ai-for-developers-project-386/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (repository *Repository) CreateEventType(ctx context.Context, eventType domain.EventType) (domain.EventType, error) {
	_, err := repository.pool.Exec(ctx, `
		insert into event_types (id, owner_id, name, description, duration_minutes)
		values ($1, $2, $3, $4, $5)
	`, eventType.ID, eventType.OwnerID, eventType.Name, eventType.Description, eventType.DurationMinutes)
	if err != nil {
		return domain.EventType{}, mapDatabaseError(err)
	}

	return eventType, nil
}

func (repository *Repository) ListEventTypes(ctx context.Context) ([]domain.EventType, error) {
	rows, err := repository.pool.Query(ctx, `
		select id, owner_id, name, description, duration_minutes
		from event_types
		order by name asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]domain.EventType, 0)
	for rows.Next() {
		var eventType domain.EventType
		if err := rows.Scan(&eventType.ID, &eventType.OwnerID, &eventType.Name, &eventType.Description, &eventType.DurationMinutes); err != nil {
			return nil, err
		}
		result = append(result, eventType)
	}

	return result, rows.Err()
}

func (repository *Repository) GetEventType(ctx context.Context, eventTypeID string) (domain.EventType, error) {
	var eventType domain.EventType
	err := repository.pool.QueryRow(ctx, `
		select id, owner_id, name, description, duration_minutes
		from event_types
		where id = $1
	`, eventTypeID).Scan(&eventType.ID, &eventType.OwnerID, &eventType.Name, &eventType.Description, &eventType.DurationMinutes)
	if err != nil {
		return domain.EventType{}, mapDatabaseError(err)
	}

	return eventType, nil
}

func (repository *Repository) CreateBooking(ctx context.Context, booking domain.Booking) error {
	_, err := repository.pool.Exec(ctx, `
		insert into bookings (id, event_type_id, start_at, end_at, created_at)
		values ($1, $2, $3, $4, $5)
	`, booking.ID, booking.EventTypeID, booking.Slot.StartAt, booking.Slot.EndAt, booking.CreatedAt)
	if err != nil {
		return mapDatabaseError(err)
	}

	return nil
}

func (repository *Repository) ListBookings(ctx context.Context) ([]domain.Booking, error) {
	rows, err := repository.pool.Query(ctx, `
		select id, event_type_id, start_at, end_at, created_at
		from bookings
		order by start_at asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanBookings(rows)
}

func (repository *Repository) ListBookingsBetween(ctx context.Context, from time.Time, to time.Time) ([]domain.Booking, error) {
	rows, err := repository.pool.Query(ctx, `
		select id, event_type_id, start_at, end_at, created_at
		from bookings
		where tstzrange(start_at, end_at, '[)') && tstzrange($1, $2, '[)')
		order by start_at asc
	`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanBookings(rows)
}

func scanBookings(rows pgxRows) ([]domain.Booking, error) {
	result := make([]domain.Booking, 0)
	for rows.Next() {
		var booking domain.Booking
		if err := rows.Scan(&booking.ID, &booking.EventTypeID, &booking.Slot.StartAt, &booking.Slot.EndAt, &booking.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, booking)
	}

	return result, rows.Err()
}

type pgxRows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close()
}

func mapDatabaseError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return service.ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch pgError.Code {
		case "23505", "23P01":
			return service.ErrConflict
		case "23503":
			return service.ErrNotFound
		}
	}

	return err
}
