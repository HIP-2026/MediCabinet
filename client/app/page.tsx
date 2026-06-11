"use client";

import { useEffect, useState } from "react";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

type Health = { status: string; db: string };

export default function Home() {
  const [health, setHealth] = useState<Health | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch(`${API_URL}/health`)
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data: Health) => setHealth(data))
      .catch((err: unknown) =>
        setError(err instanceof Error ? err.message : "unknown error"),
      );
  }, []);

  return (
    <main style={{ fontFamily: "var(--font-geist-sans)", padding: "3rem" }}>
      <h1>MediCabinet</h1>
      <p>Household drug inventory tracker.</p>
      <section style={{ marginTop: "2rem" }}>
        <h2>Backend status</h2>
        {error && <p style={{ color: "crimson" }}>Unreachable: {error}</p>}
        {!error && !health && <p>Checking…</p>}
        {health && (
          <p>
            API: <strong>{health.status}</strong> · Database:{" "}
            <strong>{health.db}</strong>
          </p>
        )}
      </section>
    </main>
  );
}
