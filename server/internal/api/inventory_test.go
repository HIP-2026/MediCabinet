package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/HIP-2026/MediCabinet/server/internal/inventory"
)

// fakeInventoryStore is a test double for InventoryStore.
type fakeInventoryStore struct {
	addItemFn  func(context.Context, inventory.AddItemParams) (inventory.Item, error)
	listItemsFn func(context.Context) ([]inventory.Item, error)
}

func (f *fakeInventoryStore) AddItem(ctx context.Context, p inventory.AddItemParams) (inventory.Item, error) {
	return f.addItemFn(ctx, p)
}

func (f *fakeInventoryStore) ListItems(ctx context.Context) ([]inventory.Item, error) {
	return f.listItemsFn(ctx)
}

func ptr[T any](v T) *T { return &v }

func setupRouter(store InventoryStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return NewRouter(nil, store)
}

func TestCreateInventoryItem_HappyPath(t *testing.T) {
	expiry := "2026-11-01"
	store := &fakeInventoryStore{
		addItemFn: func(_ context.Context, p inventory.AddItemParams) (inventory.Item, error) {
			return inventory.Item{
				ID:             1,
				MedicationID:   1,
				MedicationName: p.Name,
				Quantity:       p.Quantity,
				ExpiryDate:     p.ExpiryDate,
				Location:       p.Location,
				CreatedAt:      "2026-06-11T12:00:00Z",
			}, nil
		},
	}

	body, _ := json.Marshal(map[string]any{
		"name":        "Ibuprofen",
		"quantity":    24,
		"expiry_date": expiry,
		"location":    "cabinet",
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/inventory", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var item inventory.Item
	if err := json.Unmarshal(rec.Body.Bytes(), &item); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if item.MedicationName != "Ibuprofen" {
		t.Errorf("expected medication_name Ibuprofen, got %q", item.MedicationName)
	}
	if item.Quantity != 24 {
		t.Errorf("expected quantity 24, got %d", item.Quantity)
	}
	if item.ExpiryDate == nil || *item.ExpiryDate != expiry {
		t.Errorf("expected expiry_date %q, got %v", expiry, item.ExpiryDate)
	}
}

func TestCreateInventoryItem_MissingName(t *testing.T) {
	store := &fakeInventoryStore{}
	body, _ := json.Marshal(map[string]any{"quantity": 5})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/inventory", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestCreateInventoryItem_MissingQuantity(t *testing.T) {
	store := &fakeInventoryStore{}
	body, _ := json.Marshal(map[string]any{"name": "Aspirin"})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/inventory", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestCreateInventoryItem_NegativeQuantity(t *testing.T) {
	store := &fakeInventoryStore{}
	body, _ := json.Marshal(map[string]any{"name": "Aspirin", "quantity": -1})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/inventory", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestCreateInventoryItem_InvalidExpiryFormat(t *testing.T) {
	store := &fakeInventoryStore{}
	body, _ := json.Marshal(map[string]any{"name": "Aspirin", "quantity": 10, "expiry_date": "11/2026"})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/inventory", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestCreateInventoryItem_OptionalFieldsAbsent(t *testing.T) {
	store := &fakeInventoryStore{
		addItemFn: func(_ context.Context, p inventory.AddItemParams) (inventory.Item, error) {
			return inventory.Item{
				ID:             2,
				MedicationID:   2,
				MedicationName: p.Name,
				Quantity:       p.Quantity,
				CreatedAt:      "2026-06-11T12:00:00Z",
			}, nil
		},
	}

	body, _ := json.Marshal(map[string]any{"name": "Paracetamol", "quantity": 10})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/inventory", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setupRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var item inventory.Item
	json.Unmarshal(rec.Body.Bytes(), &item)
	if item.ExpiryDate != nil {
		t.Errorf("expected nil expiry_date, got %v", item.ExpiryDate)
	}
	if item.Location != nil {
		t.Errorf("expected nil location, got %v", item.Location)
	}
}

func TestListInventoryItems_EmptyList(t *testing.T) {
	store := &fakeInventoryStore{
		listItemsFn: func(_ context.Context) ([]inventory.Item, error) {
			return []inventory.Item{}, nil
		},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/inventory", nil)
	setupRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var items []inventory.Item
	if err := json.Unmarshal(rec.Body.Bytes(), &items); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if items == nil {
		t.Error("expected empty array, got null")
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

func TestListInventoryItems_MultipleItems(t *testing.T) {
	store := &fakeInventoryStore{
		listItemsFn: func(_ context.Context) ([]inventory.Item, error) {
			return []inventory.Item{
				{ID: 1, MedicationID: 1, MedicationName: "Ibuprofen", Quantity: 10, CreatedAt: "2026-06-11T12:00:00Z"},
				{ID: 2, MedicationID: 2, MedicationName: "Aspirin", Quantity: 5, CreatedAt: "2026-06-10T10:00:00Z"},
			}, nil
		},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/inventory", nil)
	setupRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var items []inventory.Item
	json.Unmarshal(rec.Body.Bytes(), &items)
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].MedicationName != "Ibuprofen" || items[1].MedicationName != "Aspirin" {
		t.Errorf("unexpected medication names: %v", items)
	}
}

func TestListInventoryItems_StoreError(t *testing.T) {
	store := &fakeInventoryStore{
		listItemsFn: func(_ context.Context) ([]inventory.Item, error) {
			return nil, errors.New("db error")
		},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/inventory", nil)
	setupRouter(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rec.Code)
	}
}
