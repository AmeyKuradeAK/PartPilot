'use client';

interface CodeBlockProps {
  code: string;
  language?: string;
}

export default function CodeBlock({ code, language = 'json' }: CodeBlockProps) {
  return (
    <div className="code-block" style={{
      backgroundColor: '#050508',
      border: '1px solid var(--border)',
      borderRadius: '8px',
      overflow: 'hidden',
      fontFamily: 'var(--font-geist-mono), monospace',
      fontSize: '0.85rem',
      lineHeight: '1.6'
    }}>
      <div style={{ 
        display: 'flex', 
        alignItems: 'center', 
        backgroundColor: '#0a0a0f', 
        padding: '0.5rem 1rem', 
        borderBottom: '1px solid var(--border)',
        color: 'var(--text-secondary)'
      }}>
        <div style={{ display: 'flex', gap: '6px' }}>
          <div style={{ width: '10px', height: '10px', borderRadius: '50%', backgroundColor: '#ff5f56' }} />
          <div style={{ width: '10px', height: '10px', borderRadius: '50%', backgroundColor: '#ffbd2e' }} />
          <div style={{ width: '10px', height: '10px', borderRadius: '50%', backgroundColor: '#27c93f' }} />
        </div>
        <span style={{ marginLeft: '1rem', fontSize: '0.75rem' }}>{language.toUpperCase()} Response</span>
      </div>
      <pre className="code-block-pre" style={{ 
        padding: '1.5rem', 
        overflowX: 'auto', 
        color: '#d4d4d8', 
        margin: 0
      }}>
        <code>{code}</code>
      </pre>
    </div>
  );
}
