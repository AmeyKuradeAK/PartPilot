-- PartPilot initial schema

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-----------------------------------------------------------
-- Users
-----------------------------------------------------------
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-----------------------------------------------------------
-- BOMs (uploaded files)
-----------------------------------------------------------
CREATE TABLE boms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    filename VARCHAR(512) NOT NULL,
    raw_headers JSONB,          -- stored for column mapping fallback
    column_mapping JSONB,       -- user-provided mapping if auto-detect fails
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_boms_user_id ON boms(user_id);

-----------------------------------------------------------
-- BOM Parts (individual rows extracted from a BOM)
-----------------------------------------------------------
CREATE TABLE bom_parts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bom_id UUID NOT NULL REFERENCES boms(id) ON DELETE CASCADE,
    row_index INTEGER NOT NULL,
    raw_name VARCHAR(1024) NOT NULL,
    normalized_name VARCHAR(1024),
    quantity INTEGER NOT NULL DEFAULT 1,
    is_ai_normalized BOOLEAN NOT NULL DEFAULT FALSE,
    ai_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bom_parts_bom_id ON bom_parts(bom_id);

-----------------------------------------------------------
-- Jobs (processing status — the bridge between API and Engine)
-----------------------------------------------------------
-- Status flow: pending → awaiting_confirmation → processing → done | failed
CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bom_id UUID NOT NULL REFERENCES boms(id) ON DELETE CASCADE,
    status VARCHAR(32) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'awaiting_confirmation', 'processing', 'done', 'failed')),
    claimed_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_bom_id ON jobs(bom_id);

-----------------------------------------------------------
-- Supplier Results (what came back from each supplier per part)
-----------------------------------------------------------
CREATE TABLE supplier_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bom_part_id UUID NOT NULL REFERENCES bom_parts(id) ON DELETE CASCADE,
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    supplier VARCHAR(64) NOT NULL,
    part_number VARCHAR(512),
    unit_price NUMERIC(12,4),
    stock_qty INTEGER,
    lead_time_days INTEGER,
    moq INTEGER,
    product_url VARCHAR(2048),
    rank INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_supplier_results_bom_part ON supplier_results(bom_part_id);
CREATE INDEX idx_supplier_results_job ON supplier_results(job_id);

-----------------------------------------------------------
-- Purchase Orders
-----------------------------------------------------------
CREATE TABLE purchase_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_name VARCHAR(512),
    logo_path VARCHAR(1024),       -- nullable branding slot, empty for V1
    pdf_path VARCHAR(1024),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_purchase_orders_job ON purchase_orders(job_id);
CREATE INDEX idx_purchase_orders_user ON purchase_orders(user_id);

-----------------------------------------------------------
-- Row Level Security (RLS)
-----------------------------------------------------------
-- Enable RLS on all tables to secure them in Supabase Cloud
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE boms ENABLE ROW LEVEL SECURITY;
ALTER TABLE bom_parts ENABLE ROW LEVEL SECURITY;
ALTER TABLE jobs ENABLE ROW LEVEL SECURITY;
ALTER TABLE supplier_results ENABLE ROW LEVEL SECURITY;
ALTER TABLE purchase_orders ENABLE ROW LEVEL SECURITY;

-- Note: The Go engine and Node API connect using the 'postgres' role
-- (or a service_role key), which automatically bypasses RLS.
-- This prevents any direct unauthorized access via the Supabase Data API (anon key).

