package inventory

import (
	"context"
	"time"

	"github.com/HIP-2026/MediCabinet/server/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type Store struct {
	q *db.Queries
}

func New(q *db.Queries) *Store {
	return &Store{q: q}
}

type AddItemParams struct {
	Name       string
	Quantity   int32
	ExpiryDate *string // "YYYY-MM-DD" or nil
	Location   *string
}

type Item struct {
	ID             int64   `json:"id"`
	MedicationID   int64   `json:"medication_id"`
	MedicationName string  `json:"medication_name"`
	Quantity       int32   `json:"quantity"`
	ExpiryDate     *string `json:"expiry_date,omitempty"`
	Location       *string `json:"location,omitempty"`
	CreatedAt      string  `json:"created_at"`
}

func (s *Store) AddItem(ctx context.Context, p AddItemParams) (Item, error) {
	med, err := s.q.CreateMedication(ctx, db.CreateMedicationParams{
		Name: p.Name,
	})
	if err != nil {
		return Item{}, err
	}

	var expiryDate pgtype.Date
	if p.ExpiryDate != nil {
		t, err := time.Parse("2006-01-02", *p.ExpiryDate)
		if err != nil {
			return Item{}, err
		}
		expiryDate = pgtype.Date{Time: t, Valid: true}
	}

	row, err := s.q.CreateInventoryItem(ctx, db.CreateInventoryItemParams{
		MedicationID: med.ID,
		Quantity:     p.Quantity,
		ExpiryDate:   expiryDate,
		Location:     p.Location,
	})
	if err != nil {
		return Item{}, err
	}

	return rowToItem(med.Name, row.ID, row.MedicationID, row.Quantity, row.ExpiryDate, row.Location, row.CreatedAt), nil
}

func (s *Store) ListItems(ctx context.Context) ([]Item, error) {
	rows, err := s.q.ListInventoryItems(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]Item, 0, len(rows))
	for _, r := range rows {
		items = append(items, rowToItem(r.MedicationName, r.ID, r.MedicationID, r.Quantity, r.ExpiryDate, r.Location, r.CreatedAt))
	}
	return items, nil
}

func rowToItem(medName string, id, medID int64, qty int32, expiry pgtype.Date, location *string, createdAt pgtype.Timestamptz) Item {
	item := Item{
		ID:             id,
		MedicationID:   medID,
		MedicationName: medName,
		Quantity:       qty,
		Location:       location,
		CreatedAt:      createdAt.Time.UTC().Format(time.RFC3339),
	}
	if expiry.Valid {
		s := expiry.Time.Format("2006-01-02")
		item.ExpiryDate = &s
	}
	return item
}
