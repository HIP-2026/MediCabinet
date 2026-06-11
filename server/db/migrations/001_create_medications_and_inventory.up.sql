CREATE TABLE medications (
    id                BIGSERIAL PRIMARY KEY,
    name              TEXT NOT NULL,
    barcode           TEXT,
    active_ingredient TEXT
);

CREATE TABLE inventory_items (
    id            BIGSERIAL PRIMARY KEY,
    medication_id BIGINT NOT NULL REFERENCES medications(id) ON DELETE CASCADE,
    quantity      INT NOT NULL CHECK (quantity >= 0),
    expiry_date   DATE,
    location      TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX inventory_items_medication_id_idx ON inventory_items(medication_id);
