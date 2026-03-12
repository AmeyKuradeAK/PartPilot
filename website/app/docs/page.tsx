import CodeBlock from '@/components/CodeBlock';

const uploadRequest = `POST /boms/upload HTTP/1.1
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data`;

const confirmRequest = `POST /jobs/:id/confirm-parts HTTP/1.1
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "confirmations": [
    {
      "partId": "75e0624e-45f5-40c8-bdbf-...",
      "confirmedName": "RC0603FR-0710KL"
    }
  ]
}`;

export default function Docs() {
  return (
    <div style={{ backgroundColor: 'var(--bg-primary)' }}>
      {/* Header */}
      <section style={{ padding: '6rem 2rem', borderBottom: '1px solid var(--border-light)' }}>
        <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
          <div className="tech-heading" style={{ color: 'var(--accent-orange)' }}>
            <span style={{ display: 'inline-block', width: '8px', height: '8px', backgroundColor: 'var(--accent-orange)', marginRight: '8px' }}></span>
            REST API
          </div>
          <h1 style={{ fontSize: '4rem', fontWeight: 800, letterSpacing: '-0.04em', margin: 0, color: 'var(--text-primary)' }}>
            Documentation.
          </h1>
          <p style={{ fontSize: '1.25rem', color: 'var(--text-secondary)', maxWidth: '800px', lineHeight: 1.6, marginTop: '1rem' }}>
            Interact securely with the PartPilot backend using standard HTTP verbs. All endpoints require a valid JWT bearer format.
          </p>
        </div>
      </section>

      {/* Docs Grid */}
      <section style={{ maxWidth: '1400px', margin: '0 auto', borderLeft: '1px solid var(--border-light)', borderRight: '1px solid var(--border-light)' }}>
        <div style={{ display: 'grid', gridTemplateColumns: '300px 1fr' }}>
          
          {/* Sidebar */}
          <div style={{ borderRight: '1px solid var(--border-light)', padding: '3rem 2rem', display: 'flex', flexDirection: 'column', gap: '1rem', position: 'sticky', top: '70px', height: 'calc(100vh - 70px)', overflowY: 'auto' }}>
            <div className="tech-heading">RESOURCES</div>
            <a href="#authentication" style={{ color: 'var(--text-secondary)', fontSize: '0.95rem' }}>Authentication</a>
            <a href="#upload" style={{ color: 'var(--text-secondary)', fontSize: '0.95rem' }}>Upload BOM</a>
            <a href="#confirm" style={{ color: 'var(--text-secondary)', fontSize: '0.95rem' }}>Commit Confirmations</a>
            <a href="#results" style={{ color: 'var(--text-secondary)', fontSize: '0.95rem' }}>Poll Results</a>
            <a href="#po" style={{ color: 'var(--text-secondary)', fontSize: '0.95rem' }}>Retrieve PDF PO</a>
          </div>

          {/* Main Content */}
          <div style={{ padding: '4rem 5rem', display: 'flex', flexDirection: 'column', gap: '6rem' }}>
            
            <div id="upload" style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
              <h2 style={{ fontSize: '2rem', fontWeight: 700, paddingBottom: '0.5rem', borderBottom: '1px solid var(--border-light)' }}>1. Upload BOM</h2>
              <p style={{ color: 'var(--text-secondary)', lineHeight: 1.6 }}>Ingests `multipart/form-data` and issues a job lock. Asynchronous workers decode the file buffer and populate horizontal DB nodes.</p>
              <CodeBlock code={uploadRequest} language="http" />
            </div>

            <div id="confirm" style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
              <h2 style={{ fontSize: '2rem', fontWeight: 700, paddingBottom: '0.5rem', borderBottom: '1px solid var(--border-light)' }}>2. Commit Confirmations</h2>
              <p style={{ color: 'var(--text-secondary)', lineHeight: 1.6 }}>Patch deterministic flags on pending states to un-halt AI normalized jobs. Requires `application/json` payload containing valid GUIDs.</p>
              <CodeBlock code={confirmRequest} language="http" />
            </div>
            
             <div id="results" style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
              <h2 style={{ fontSize: '2rem', fontWeight: 700, paddingBottom: '0.5rem', borderBottom: '1px solid var(--border-light)' }}>3. Poll Results</h2>
              <p style={{ color: 'var(--text-secondary)', lineHeight: 1.6 }}>Returns the executed JSON payload resolving to aggregated, filtered topology from global distributors mapped to the original IDs.</p>
              <CodeBlock code={`GET /jobs/:id/results HTTP/1.1\nAuthorization: Bearer <jwt_token>`} language="http" />
            </div>

          </div>
        </div>
      </section>
    </div>
  );
}
