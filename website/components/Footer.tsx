import Link from 'next/link';

export default function Footer() {
  return (
    <footer style={{
      borderTop: '1px solid var(--border-light)',
      backgroundColor: 'var(--bg-surface)',
      fontFamily: 'var(--font-geist-mono), monospace',
      fontSize: '0.85rem'
    }}>
      <div style={{
        maxWidth: '1400px',
        margin: '0 auto',
        padding: '4rem 2rem',
        display: 'grid',
        gridTemplateColumns: '1fr auto',
        gap: '4rem'
      }}>
        
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontWeight: 700, color: 'var(--text-primary)', marginBottom: '1rem', fontSize: '1rem' }}>
            <div style={{ width: '12px', height: '12px', backgroundColor: 'var(--accent-orange)' }}></div>
            PARTPILOT
          </div>
          <p style={{ color: 'var(--text-secondary)', maxWidth: '400px', lineHeight: 1.6 }}>
            The global programmatic API for electronic component sourcing and supply chain automation.
          </p>
        </div>

        <div style={{ display: 'flex', gap: '4rem' }}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
            <div style={{ color: 'var(--text-primary)', fontWeight: 600 }}>Infrastructure</div>
            <Link href="#" style={{ color: 'var(--text-secondary)' }}>Network Map</Link>
            <Link href="#" style={{ color: 'var(--text-secondary)' }}>System Status</Link>
            <Link href="#" style={{ color: 'var(--text-secondary)' }}>Pricing</Link>
          </div>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
            <div style={{ color: 'var(--text-primary)', fontWeight: 600 }}>Developers</div>
            <Link href="#" style={{ color: 'var(--text-secondary)' }}>Documentation</Link>
            <Link href="#" style={{ color: 'var(--text-secondary)' }}>API Reference</Link>
            <Link href="#" style={{ color: 'var(--text-secondary)' }}>GitHub</Link>
          </div>
        </div>

      </div>

      <div style={{
        borderTop: '1px solid var(--border-light)',
      }}>
        <div style={{
          maxWidth: '1400px',
          margin: '0 auto',
          padding: '1.5rem 2rem',
          display: 'flex',
          justifyContent: 'space-between',
          color: 'var(--text-tertiary)',
        }}>
          <div>&copy; {new Date().getFullYear()} PartPilot Inc.</div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            <span style={{ display: 'inline-block', width: '8px', height: '8px', backgroundColor: '#32cd32' }}></span>
            All systems operational (99.99%)
          </div>
        </div>
      </div>
    </footer>
  );
}
