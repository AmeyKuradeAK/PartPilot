'use client';

import Link from 'next/link';
import { useState } from 'react';

export default function Navbar() {
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <nav className="site-nav" style={{
      borderBottom: '1px solid var(--border-light)',
      backgroundColor: 'var(--bg-primary)',
      fontSize: '0.85rem'
    }}>
      {/* Top Banner (Utility) */}
      <div className="nav-utility" style={{ backgroundColor: 'var(--bg-surface)', borderBottom: '1px solid var(--border-dim)', padding: '0.25rem 2rem', color: 'var(--text-tertiary)', display: 'flex', justifyContent: 'flex-end', gap: '1rem' }}>
        <Link href="#">Support</Link>
        <Link href="#">Log In</Link>
      </div>

      {/* Main Nav */}
      <div className="nav-main" style={{ padding: '0.75rem 2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div className="nav-left" style={{ display: 'flex', alignItems: 'center', gap: '2.5rem' }}>
          {/* Logo */}
          <Link href="/" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontWeight: 700, fontSize: '1.1rem', letterSpacing: '-0.03em' }}>
            <div style={{ width: '12px', height: '12px', backgroundColor: 'var(--accent-orange)' }}></div>
            PARTPILOT
          </Link>
        </div>

        <div className={`nav-menu ${menuOpen ? 'open' : ''}`}>
          <div className="nav-links" style={{ display: 'flex', gap: '1.5rem', color: 'var(--text-secondary)', fontWeight: 500 }}>
            <Link href="/how-it-works" style={{ transition: 'color 0.1s' }}>Products</Link>
            <Link href="/pricing">Pricing</Link>
            <Link href="/docs">Developers</Link>
            <Link href="#">Network</Link>
          </div>

          <div className="nav-actions" style={{ display: 'flex', gap: '1rem' }}>
            <Link href="/docs" className="btn-secondary">
              Documentation
            </Link>
            <Link href="#" className="btn-primary">
              Sign Up
            </Link>
          </div>
        </div>

        <button
          className="nav-toggle"
          type="button"
          aria-label="Toggle navigation"
          aria-expanded={menuOpen}
          onClick={() => setMenuOpen((open) => !open)}
        >
          <span />
          <span />
          <span />
        </button>
      </div>
    </nav>
  );
}
