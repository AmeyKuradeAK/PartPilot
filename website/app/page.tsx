import Link from 'next/link';
import CodeBlock from '@/components/CodeBlock';
import { ArrowLeftRight, Binary, ShieldAlert, Cpu } from 'lucide-react';

const sampleResponse = `{
  "id": "75e0624e-45f5-40c8-bdbf-a53048e7eba3",
  "rawName": "10k 0603 resistor",
  "normalizedName": "RC0603FR-0710KL",
  "quantity": 50,
  "isAiNormalized": true,
  "results": [
    {
      "supplier": "DigiKey",
      "unitPrice": 0.15,
      "stockQty": 10000,
      "leadTimeDays": 0
    }
  ]
}`;

export default function Home() {
  return (
    <>
      {/* Hero Section - Stark, Technical, Left Aligned */}
      <section style={{
        padding: '6rem 2rem',
        borderBottom: '1px solid var(--border-light)',
        backgroundColor: 'var(--bg-primary)'
      }}>
        <div style={{ maxWidth: '1200px', margin: '0 auto', display: 'flex', flexDirection: 'column', gap: '2rem' }}>
          
          <div className="tech-heading" style={{ color: 'var(--accent-orange)' }}>
            <span style={{ display: 'inline-block', width: '8px', height: '8px', backgroundColor: 'var(--accent-orange)', marginRight: '8px' }}></span>
            PARTPILOT // CORE ENGINE v1.0
          </div>
          
          <h1 style={{ fontSize: '5.5rem', fontWeight: 800, letterSpacing: '-0.04em', lineHeight: 1, maxWidth: '900px', margin: 0, color: 'var(--text-primary)' }}>
            Global Hardware Procurement Network.
          </h1>
          
          <p style={{ fontSize: '1.25rem', color: 'var(--text-secondary)', maxWidth: '650px', lineHeight: 1.6, marginTop: '1rem' }}>
            PartPilot provides programmatic, transaction-safe access to global electronic component inventory. Upload fragmented BOMs, normalize data in milliseconds, and span massive concurrent queries across global distributor APIs.
          </p>
          
          <div style={{ display: 'flex', gap: '1rem', marginTop: '1rem' }}>
            <Link href="#" className="btn-primary" style={{ padding: '0.75rem 2rem', fontSize: '1.1rem' }}>
              Deploy Infrastructure
            </Link>
            <Link href="/docs" className="btn-secondary" style={{ padding: '0.75rem 2rem', fontSize: '1.1rem' }}>
              Read the Docs
            </Link>
          </div>
        </div>
      </section>

      {/* Network Architecture Grid (replaces soft features) */}
      <section style={{ backgroundColor: 'var(--bg-primary)' }}>
        <div style={{ maxWidth: '1400px', margin: '0 auto', borderLeft: '1px solid var(--border-light)', borderRight: '1px solid var(--border-light)' }}>
          
          <div style={{ padding: '4rem 2rem', borderBottom: '1px solid var(--border-light)' }}>
            <div className="tech-heading">System Architecture</div>
            <h2 style={{ fontSize: '2.5rem', fontWeight: 700, letterSpacing: '-0.02em', margin: 0 }}>Built for extreme concurrency.</h2>
          </div>

          <div className="tech-grid">
            <div className="tech-cell">
              <Cpu size={24} color="var(--accent-orange)" style={{ marginBottom: '1.5rem' }} />
              <h3 style={{ fontSize: '1.2rem', fontWeight: 700, marginBottom: '0.75rem' }}>AI Normalization Engine</h3>
              <p style={{ color: 'var(--text-secondary)', fontSize: '0.95rem', lineHeight: 1.6 }}>Raw strings are piped through OpenAI models utilizing strict JSON schema enforcement to resolve generic descriptions into valid Manufacturer Part Numbers (MPNs).</p>
            </div>

            <div className="tech-cell">
              <ArrowLeftRight size={24} color="var(--accent-orange)" style={{ marginBottom: '1.5rem' }} />
              <h3 style={{ fontSize: '1.2rem', fontWeight: 700, marginBottom: '0.75rem' }}>Distributed Fan-Out</h3>
              <p style={{ color: 'var(--text-secondary)', fontSize: '0.95rem', lineHeight: 1.6 }}>Go routines orchestrate thousands of non-blocking HTTP transactions across DigiKey, Mouser, and custom ERP integration endpoints simultaneously.</p>
            </div>

            <div className="tech-cell">
              <ShieldAlert size={24} color="var(--accent-orange)" style={{ marginBottom: '1.5rem' }} />
              <h3 style={{ fontSize: '1.2rem', fontWeight: 700, marginBottom: '0.75rem' }}>Deterministic State Halts</h3>
              <p style={{ color: 'var(--text-secondary)', fontSize: '0.95rem', lineHeight: 1.6 }}>Jobs demanding algorithmic inference are autonomously halted (`status: awaiting_confirmation`) until deterministic engineer approval unlocks the transaction context.</p>
            </div>

            <div className="tech-cell">
              <Binary size={24} color="var(--accent-orange)" style={{ marginBottom: '1.5rem' }} />
              <h3 style={{ fontSize: '1.2rem', fontWeight: 700, marginBottom: '0.75rem' }}>ACID Compliant Processing</h3>
              <p style={{ color: 'var(--text-secondary)', fontSize: '0.95rem', lineHeight: 1.6 }}>PostgreSQL row-level locks guarantees exactly-once processing across horizontally scaled distributed workers without race conditions.</p>
            </div>
          </div>
        </div>
      </section>

      {/* Terminal Preview Section */}
      <section style={{ backgroundColor: 'var(--bg-elevated)', borderTop: '1px solid var(--border-light)', borderBottom: '1px solid var(--border-light)' }}>
        <div style={{ maxWidth: '1400px', margin: '0 auto', display: 'grid', gridTemplateColumns: 'minmax(0, 1fr) 1fr' }}>
          
          <div style={{ padding: '6rem 4rem', borderRight: '1px solid var(--border-light)' }}>
            <div className="tech-heading">API INTERFACE</div>
            <h2 style={{ fontSize: '2.5rem', fontWeight: 700, letterSpacing: '-0.02em', marginBottom: '1.5rem' }}>Raw payloads. No abstractions.</h2>
            <p style={{ color: 'var(--text-secondary)', fontSize: '1.1rem', lineHeight: 1.6, marginBottom: '2rem' }}>
              Integrate PartPilot directly into your CI/CD pipelines or internal procurement dashboards. The REST API exposes structured JSON containing exact algorithmic match flags, stock integers, and deterministic pricing arrays.
            </p>
            <ul style={{ borderTop: '1px solid var(--border-light)', paddingTop: '2rem', display: 'grid', gap: '1rem', color: 'var(--text-secondary)' }}>
              <li className="mono" style={{ fontSize: '0.85rem' }}>$ curl -X POST https://api.partpilot.net/v1/jobs</li>
              <li className="mono" style={{ fontSize: '0.85rem' }}>$ curl -X GET https://api.partpilot.net/v1/jobs/id/results</li>
            </ul>
          </div>
          
          <div style={{ padding: '4rem' }}>
            <CodeBlock code={sampleResponse} />
          </div>

        </div>
      </section>

      {/* Brutalist Footer CTA */}
      <section style={{ padding: '8rem 2rem', backgroundColor: 'var(--bg-primary)' }}>
        <div style={{ maxWidth: '800px', margin: '0 auto', textAlign: 'center' }}>
          <h2 style={{ fontSize: '3rem', fontWeight: 700, letterSpacing: '-0.02em', margin: 0, color: 'var(--text-primary)' }}>Initiate Procurement</h2>
          <p style={{ color: 'var(--text-secondary)', fontSize: '1.25rem', marginTop: '1rem', marginBottom: '2.5rem' }}>
            Scale logic, not headcount. Standardize your hardware supply chain in milliseconds.
          </p>
          <Link href="#" className="btn-primary" style={{ padding: '1rem 3rem', fontSize: '1.25rem' }}>
            Get Started
          </Link>
        </div>
      </section>
    </>
  );
}
