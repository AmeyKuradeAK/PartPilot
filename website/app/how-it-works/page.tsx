import CodeBlock from '@/components/CodeBlock';
import { ArrowDown } from 'lucide-react';

const csvInput = `rawName,quantity
STM32F103C8T6,5
10k 0603 resistor,50`;

const aiOutput = `{
  "original": "10k 0603 resistor",
  "normalizedMPN": "RC0603FR-0710KL",
  "confidence": 0.98,
  "requiresHumanConfirmation": true
}`;

export default function HowItWorks() {
  return (
    <div style={{ backgroundColor: 'var(--bg-primary)' }}>
      {/* Header */}
      <section style={{ padding: '6rem 2rem', borderBottom: '1px solid var(--border-light)' }}>
        <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
          <div className="tech-heading" style={{ color: 'var(--accent-orange)' }}>
            <span style={{ display: 'inline-block', width: '8px', height: '8px', backgroundColor: 'var(--accent-orange)', marginRight: '8px' }}></span>
            SYSTEM ARCHITECTURE
          </div>
          <h1 style={{ fontSize: '4rem', fontWeight: 800, letterSpacing: '-0.04em', margin: 0, color: 'var(--text-primary)' }}>
            Inside the Engine.
          </h1>
          <p style={{ fontSize: '1.25rem', color: 'var(--text-secondary)', maxWidth: '800px', lineHeight: 1.6, marginTop: '1rem' }}>
            PartPilot is a distributed Go pipeline orchestrating file parsing, strict LLM normalization, and high-concurrency network requests over a PostgreSQL transaction pool.
          </p>
        </div>
      </section>

      {/* Grid Pipeline Stages */}
      <section style={{ maxWidth: '1400px', margin: '0 auto', borderLeft: '1px solid var(--border-light)', borderRight: '1px solid var(--border-light)' }}>
        
        {/* Stage 1 */}
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', borderBottom: '1px solid var(--border-light)' }}>
          <div style={{ padding: '4rem', borderRight: '1px solid var(--border-light)' }}>
            <div className="tech-heading">STAGE 1 // INGESTION</div>
            <h2 style={{ fontSize: '2rem', fontWeight: 700, margin: '1rem 0' }}>Upload & Parse</h2>
            <p style={{ color: 'var(--text-secondary)', fontSize: '1.1rem', lineHeight: 1.6 }}>
              Multipart form uploads are buffered into blob storage. Background Go workers instantly claim the job utilizing `FOR UPDATE SKIP LOCKED` PostgreSQL row locks to guarantee exactly-once processing across horizontally scaled replicas.
            </p>
          </div>
          <div style={{ padding: '4rem', backgroundColor: 'var(--bg-surface)' }}>
            <CodeBlock code={csvInput} language="csv" />
          </div>
        </div>

        {/* Intersect */}
        <div style={{ display: 'flex', justifyContent: 'center', padding: '2rem 0', borderBottom: '1px solid var(--border-light)' }}>
          <ArrowDown size={32} color="var(--border-light)" />
        </div>

        {/* Stage 2 */}
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', borderBottom: '1px solid var(--border-light)' }}>
          <div style={{ padding: '4rem', backgroundColor: 'var(--bg-surface)', borderRight: '1px solid var(--border-light)' }}>
            <CodeBlock code={aiOutput} language="json" />
          </div>
          <div style={{ padding: '4rem' }}>
            <div className="tech-heading">STAGE 2 // NORMALIZATION</div>
            <h2 style={{ fontSize: '2rem', fontWeight: 700, margin: '1rem 0' }}>AI State Halts</h2>
            <p style={{ color: 'var(--text-secondary)', fontSize: '1.1rem', lineHeight: 1.6 }}>
              Regex classifiers route ambiguous component descriptions to OpenAI's structured JSON endpoints. If probabilistic matching occurs, the job state is deterministicly halted (`awaiting_confirmation`) until manual cryptographic authorization is provided by an engineer.
            </p>
          </div>
        </div>

        {/* Intersect */}
        <div style={{ display: 'flex', justifyContent: 'center', padding: '2rem 0', borderBottom: '1px solid var(--border-light)' }}>
          <ArrowDown size={32} color="var(--border-light)" />
        </div>

        {/* Stage 3 */}
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', borderBottom: '1px solid var(--border-light)' }}>
          <div style={{ padding: '4rem', borderRight: '1px solid var(--border-light)' }}>
            <div className="tech-heading">STAGE 3 // AGGREGATION</div>
            <h2 style={{ fontSize: '2rem', fontWeight: 700, margin: '1rem 0' }}>Concurrent Fan-Out</h2>
            <p style={{ color: 'var(--text-secondary)', fontSize: '1.1rem', lineHeight: 1.6 }}>
              The orchestrator spawns unblocked goroutines spanning HTTP requests to DigiKey, Mouser, and Custom ERPs simultaneously. Results are ranked against internal logic blocks (MOQ thresholds, Lead Times) and compiled into a unified Purchase Order blob.
            </p>
          </div>
          <div style={{ padding: '4rem', backgroundColor: 'var(--bg-surface)' }}>
             <ul style={{ border: '1px solid var(--border-light)', fontFamily: 'var(--font-geist-mono), monospace' }}>
              <li style={{ display: 'flex', justifyContent: 'space-between', borderBottom: '1px solid var(--border-light)', padding: '1rem' }}>
                <span style={{ color: 'var(--text-secondary)' }}>GET /api/v1/digikey/search</span>
                <span style={{ color: '#32cd32' }}>124ms</span>
              </li>
              <li style={{ display: 'flex', justifyContent: 'space-between', borderBottom: '1px solid var(--border-light)', padding: '1rem' }}>
                <span style={{ color: 'var(--text-secondary)' }}>GET /api/v1/mouser/search</span>
                <span style={{ color: '#32cd32' }}>180ms</span>
              </li>
              <li style={{ display: 'flex', justifyContent: 'space-between', padding: '1rem', backgroundColor: 'var(--bg-elevated)' }}>
                <span style={{ color: 'var(--text-primary)', fontWeight: 600 }}>Total Execution Time</span>
                <span style={{ color: 'var(--accent-orange)', fontWeight: 600 }}>181ms</span>
              </li>
            </ul>
          </div>
        </div>

      </section>
    </div>
  );
}
