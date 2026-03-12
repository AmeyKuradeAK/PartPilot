# PartPilot

Automated hardware component sourcing. Upload a BOM → system finds best suppliers → select → generate PO.

## Architecture

- **`api/`** — Node/TypeScript API layer (auth, uploads, job status, results)
- **`core-engine/`** — Go processing engine (supplier queries, ranking, PO generation)
- **`supabase/`** — Managed Supabase database and config

Services communicate through Postgres only. No direct service-to-service calls.

## Local Development

```bash
# Start Supabase (Postgres, Auth, Studio)
npx supabase start

# Start API and Engine
docker-compose up
```

## Environment

Each service has its own `.env` file. See `.env.example` in each service directory.
