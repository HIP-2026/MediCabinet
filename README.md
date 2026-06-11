# MediCabinet


# MediCabinet — A Drug Inventory Tracker

## Problem Statement

In many households — especially those with elderly relatives, chronic conditions, or shared caregiving — medications quickly turn into a chaotic mess. Painkillers, prescriptions, vitamins, and over-the-counter remedies pile up across kitchen cabinets, bedside drawers, and travel bags. Half of them are nearly empty, some have expired without anyone noticing, and a few should never have been taken together in the first place. The consequences of losing track range from the merely annoying — a missed dose, a late-night drive to the pharmacy — to genuinely dangerous: an expired antibiotic taken in good faith, a forgotten contraindication between two innocent-looking pills, or an accidental double-dose during a confusing morning.
We want to build an application — working title MediCabinet — that brings order to this. The application should help users maintain an accurate, up-to-date inventory of every medication in their home, with quantities, dosage instructions, and expiration dates always at their fingertips. Adding a new medication should be effortless: scanning the packaging barcode, taking a photo of the box, or typing it in by hand. Searching and filtering the inventory should feel as natural as searching contacts on a phone. When something is about to expire or run low, the app should warn the user with enough advance notice to restock without panic. When two medications shouldn't be combined, the app should make that visible before the next dose, not after.



## Analysis

The Analysis Object Model captures the problem domain. Key decisions: **Medication** (the abstract drug — name, barcode, active ingredient) is separate from **InventoryItem** (the concrete package at home — quantity, expiry, location). **Contraindication** is an association class between two Medications. Every household member can have a **MedicationPlan** consisting of **PlanEntries** (dose + schedule per medication) — the personal dose lives in the plan, not on the drug.

![Analysis Object Model](docs/diagrams/medicabinet-aom.svg)

## Design Goals

Starting point: the quality attributes from the requirements process. Design goals must be **measurable** and **prioritized** (MoSCoW).

| # | Design Goal | Measurable Formulation | Priority |
|---|-------------|------------------------|----------|
| DG1 | Safety | An interaction warning is shown *before* a dose is confirmed, for 100% of dose confirmations | Must |
| DG2 | Usability | Adding a medication via barcode scan takes ≤ 30 seconds and ≤ 3 interactions | Must |
| DG3 | Shared access | A caregiver can view and update a patient's inventory from a different device/location | Won't (v1) |
| DG4 | Reliability | Expiry and low-stock warnings fire at least 7 days before the event | Should |
| DG5 | Privacy | Health data is transport-encrypted and not shared with third parties | Must |
| DG6 | Modifiability | The drug data provider (barcode lookup, interaction data) can be swapped without changing other subsystems | Should |

Architectural decisions derived from these goals:

- **Client-server architecture with a web app** (thin client) — no installation, works on any device with a browser, one codebase. A central backend keeps the door open for shared access (DG3) later; offline caching and native apps are later refinements.
- **Access / household management is deferred (DG3 = Won't for v1).** The webapp serves one household without login. Household and HouseholdMember stay in the analysis model — the domain doesn't change, we just don't build the subsystem yet.
- **DG1 → safety-critical checks run server-side** in one place (Interaction Checker), not scattered across UI code.
- **DG6 → external drug data sits behind our own Medication Catalog component** (adapter), so the provider is replaceable.



## Component Architecture

The system is decomposed into four backend subsystems, each offering exactly one service, and a webapp consuming them. The Interaction Checker uses the Catalog Service for contraindication data; Reminder & Alerts reads the inventory to detect upcoming expiry and low stock. This keeps coupling low (the webapp only talks to services) and cohesion high (everything about stock levels lives in Inventory).

| Component | Responsibility |
|-----------|----------------|
| Inventory | Manages the household's inventory items: quantities, expiry dates, locations |
| Medication Catalog | Drug master data: barcode lookup, ingredients, contraindication data (adapter to an external provider) |
| Interaction Checker | Checks a planned dose against the inventory/plan for dangerous combinations |
| Reminder & Alerts | Medication plans, dose reminders, expiry and low-stock alerts |
| MediCabinet Webapp | UI: scan & add, browse/search inventory, dose & reminder interaction |

![Top-level component diagram](docs/diagrams/medicabinet-top-level.svg)

## Deployment

Adapting to the design goals (client-server, replaceable external providers, safety checks in one place server-side): all backend components are deployed together on a single **MediCabinet Backend Server** — a modular monolith, adequate for the team size and avoiding microservice overhead. The webapp runs in the user's browser. Two external systems are attached behind our own components, so each provider can be swapped without touching the rest: the **Drug Database Provider** (consumed only by the Medication Catalog) and a **Notification Provider** (Web Push / e-mail, consumed only by Reminder & Alerts).

![Deployment view](docs/diagrams/medicabinet-deployment.svg)


## Implementation Decisions

- The client should be in Typescript with React 
- The server should be in go with a REST API 
- The database should be Postgres, with a simple schema reflecting the analysis model (Medication, InventoryItem, Household, etc.)
- Use SQLc for type-safe database access in Go, generating code from SQL queries
- Docker for containerization, with a docker-compose setup for local development (backend + Postgres)

## Roadmap (Feature Increments)

Work proceeds one user story at a time, each a **vertical slice** (migration + sqlc query + Go service + REST endpoint + Next.js UI + tests). The order is dependency-aware: every slice de-risks the next, and the two `Must` goals (DG1 safety, DG2 barcode usability) land as early as their dependencies allow.

| # | Increment | What | Subsystem | Goal |
|---|-----------|------|-----------|------|
| 1 | Manual add + browse inventory | Add an item by hand (name, quantity, expiry, location) and list it. First real schema (`Medication` + `InventoryItem`); first thing to exercise sqlc/migrations end to end. | Inventory | Baseline |
| 2 | Search & filter inventory | Browsing that feels "as natural as searching contacts." Builds on #1. | Inventory | DG2 |
| 3 | Barcode scan to add | Introduces the Medication Catalog adapter and external Drug Database Provider (stub the provider first, then swap in a real one - also validates DG6). Photo/OCR add is a later refinement of this flow. | Catalog | DG2 (Must), DG6 |
| 4 | Expiry & low-stock alerts | Warn ≥ 7 days ahead by reading the inventory. Start in-app; external Notification Provider is a later swap behind the adapter. | Reminder & Alerts | DG4 |
| 5 | Medication plans + dose tracking | `MedicationPlan` / `PlanEntry` with `markAsTaken`. Must precede #6: DG1 is defined relative to a dose confirmation. | Reminder & Alerts | Enables DG1 |
| 6 | Interaction check before dose | The headline safety feature; runs server-side in the Interaction Checker, never only in the UI. Depends on contraindication data (#3) and the dose trigger (#5). | Interaction Checker | DG1 (Must) |

**Cross-cutting (not increments):**

- **DG5 (privacy / transport encryption)** is deployment config (HTTPS/TLS), handled at the infra layer before any real health data flows.
- **DG3 (shared access)** is `Won't` for v1 - the no-login single-household assumption holds, so Household/Member management is not built yet.

**Sequencing tension:** DG1 is a `Must` but the most dependency-heavy slice, so it lands at #6. To demonstrate safety value sooner, a stripped check ("warn when two contraindicated meds both exist in inventory") could be pulled forward after #3, then enriched with dose-time confirmation once #5 exists - at the cost of some rework.

## Development Setup

### Layout

```
client/   Next.js (TypeScript, App Router) thin client
server/   Go + gin REST backend (modular monolith)
```

The backend `internal/` packages mirror the subsystem decomposition (`inventory`, `catalog`, `interaction`, `reminder`). Database access is generated by [sqlc](https://sqlc.dev) into `server/internal/db` from SQL in `server/db/migrations` (schema) and `server/db/queries`. Both directories are empty until the first user story adds tables.

### Prerequisites

- Go 1.26+
- Node 20+
- Docker (with Compose)
- Optional: [sqlc](https://sqlc.dev) and [golang-migrate](https://github.com/golang-migrate/migrate) CLIs for DB codegen and migrations

### Run the backend + database

```sh
docker compose up --build      # or: make up
```

Postgres comes up on `localhost:5432`, the API on `localhost:8080`. Verify:

```sh
curl localhost:8080/health     # {"status":"ok","db":"up"}
```

### Run the client

```sh
cd client
cp .env.local.example .env.local
npm install
npm run dev                     # http://localhost:3000
```

The home page pings the backend `/health` and shows its status.

### Database code (sqlc & migrations)

Migrations use golang-migrate naming (`db/migrations/<version>_<name>.up.sql` / `.down.sql`). After adding a migration and matching queries in `db/queries`, regenerate the typed access layer:

```sh
cd server && sqlc generate      # or: make sqlc
make migrate-up                 # applies migrations (needs golang-migrate)
```

### Tests

```sh
cd server && go test ./...      # or: make test
```

