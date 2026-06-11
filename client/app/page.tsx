"use client";

import { useEffect, useState } from "react";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

type InventoryItem = {
  id: number;
  medication_id: number;
  medication_name: string;
  quantity: number;
  expiry_date?: string;
  location?: string;
  created_at: string;
};

export default function Home() {
  const [items, setItems] = useState<InventoryItem[]>([]);
  const [fetchError, setFetchError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const [name, setName] = useState("");
  const [quantity, setQuantity] = useState<string>("");
  const [expiryDate, setExpiryDate] = useState("");
  const [location, setLocation] = useState("");
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    fetch(`${API_URL}/inventory`)
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data: InventoryItem[]) => setItems(data))
      .catch((err: unknown) =>
        setFetchError(err instanceof Error ? err.message : "unknown error"),
      )
      .finally(() => setLoading(false));
  }, []);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setFormError(null);

    if (!name.trim()) {
      setFormError("Medication name is required.");
      return;
    }
    const qty = parseInt(quantity, 10);
    if (isNaN(qty) || qty < 0) {
      setFormError("Quantity must be a non-negative whole number.");
      return;
    }

    const body: Record<string, unknown> = { name: name.trim(), quantity: qty };
    if (expiryDate) body.expiry_date = expiryDate;
    if (location.trim()) body.location = location.trim();

    setSubmitting(true);
    try {
      const res = await fetch(`${API_URL}/inventory`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });
      if (!res.ok) {
        const err = await res.json().catch(() => ({ error: `HTTP ${res.status}` }));
        setFormError((err as { error: string }).error ?? `HTTP ${res.status}`);
        return;
      }
      const item: InventoryItem = await res.json();
      setItems((prev) => [item, ...prev]);
      setName("");
      setQuantity("");
      setExpiryDate("");
      setLocation("");
    } catch (err: unknown) {
      setFormError(err instanceof Error ? err.message : "unknown error");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <main style={{ fontFamily: "var(--font-geist-sans)", padding: "3rem", maxWidth: "640px" }}>
      <h1>MediCabinet</h1>
      <p style={{ color: "#666" }}>Household drug inventory tracker.</p>

      <section style={{ marginTop: "2rem" }}>
        <h2>Add medication</h2>
        <form onSubmit={handleSubmit} style={{ display: "flex", flexDirection: "column", gap: "0.75rem" }}>
          <label>
            Name *
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. Ibuprofen 400mg"
              style={inputStyle}
            />
          </label>
          <label>
            Quantity *
            <input
              type="number"
              min={0}
              step={1}
              value={quantity}
              onChange={(e) => setQuantity(e.target.value)}
              placeholder="e.g. 24"
              style={inputStyle}
            />
          </label>
          <label>
            Expiry date
            <input
              type="date"
              value={expiryDate}
              onChange={(e) => setExpiryDate(e.target.value)}
              style={inputStyle}
            />
          </label>
          <label>
            Storage location
            <input
              type="text"
              value={location}
              onChange={(e) => setLocation(e.target.value)}
              placeholder="e.g. bathroom cabinet"
              style={inputStyle}
            />
          </label>
          {formError && <p style={{ color: "crimson", margin: 0 }}>{formError}</p>}
          <button type="submit" disabled={submitting} style={buttonStyle}>
            {submitting ? "Adding…" : "Add to cabinet"}
          </button>
        </form>
      </section>

      <section style={{ marginTop: "3rem" }}>
        <h2>Your cabinet</h2>
        {loading && <p>Loading…</p>}
        {fetchError && <p style={{ color: "crimson" }}>Could not load inventory: {fetchError}</p>}
        {!loading && !fetchError && items.length === 0 && (
          <p style={{ color: "#666" }}>No medications in your cabinet yet.</p>
        )}
        <ul style={{ listStyle: "none", padding: 0, margin: 0, display: "flex", flexDirection: "column", gap: "0.75rem" }}>
          {items.map((item) => (
            <li key={item.id} style={cardStyle}>
              <strong style={{ fontSize: "1.05rem" }}>{item.medication_name}</strong>
              <span style={{ color: "#555" }}>Qty: {item.quantity}</span>
              {item.expiry_date && <span style={{ color: "#555" }}>Expires: {item.expiry_date}</span>}
              {item.location && <span style={{ color: "#555" }}>Location: {item.location}</span>}
            </li>
          ))}
        </ul>
      </section>
    </main>
  );
}

const inputStyle: React.CSSProperties = {
  display: "block",
  width: "100%",
  marginTop: "0.25rem",
  padding: "0.4rem 0.6rem",
  fontSize: "1rem",
  boxSizing: "border-box",
  border: "1px solid #ccc",
  borderRadius: "4px",
};

const buttonStyle: React.CSSProperties = {
  padding: "0.5rem 1.25rem",
  fontSize: "1rem",
  cursor: "pointer",
  alignSelf: "flex-start",
  background: "#0070f3",
  color: "#fff",
  border: "none",
  borderRadius: "4px",
};

const cardStyle: React.CSSProperties = {
  display: "flex",
  flexDirection: "column",
  gap: "0.2rem",
  padding: "0.75rem 1rem",
  border: "1px solid #e0e0e0",
  borderRadius: "6px",
};
