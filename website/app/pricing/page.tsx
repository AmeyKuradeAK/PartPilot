import { Check } from 'lucide-react';
import Link from 'next/link';

const tiers = [
  {
    name: 'Developer',
    price: 'Free',
    description: 'Local deployment / Pipeline validation.',
    features: [
      'Mock Supplier Topology',
      '50 items / Payload limit',
      'Docker Compose scripts',
      'Community Issue Tracker'
    ]
  },
  {
    name: 'Pro',
    price: '$99',
    period: '/mo',
    description: 'Live API provisioning / High throughput.',
    highlighted: true,
    features: [
      'Production DigiKey & Mouser APIs',
      'Unlimited payload sizes',
      'LLM Normalization Engine',
      'Automated PO aggregation',
      'SLA Guaranteed Support'
    ]
  },
  {
    name: 'Enterprise',
    price: 'Custom',
    description: 'VPC deployment / ERP sync parameters.',
    features: [
      'NetSuite / SAP ERP Webhooks',
      'SAML / SSO Authorization',
      'Dedicated Network Slack',
      'Custom Supplier API Adapters',
      'On-Premise binaries'
    ]
  }
];

export default function Pricing() {
  return (
    <div style={{ backgroundColor: 'var(--bg-primary)', minHeight: 'calc(100vh - 70px)' }}>
      {/* Header */}
      <section className="section-pad" style={{ padding: '6rem 2rem', borderBottom: '1px solid var(--border-light)' }}>
        <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
          <div className="tech-heading" style={{ color: 'var(--accent-orange)' }}>
            <span style={{ display: 'inline-block', width: '8px', height: '8px', backgroundColor: 'var(--accent-orange)', marginRight: '8px' }}></span>
            NETWORK PRICING
          </div>
          <h1 className="page-title" style={{ fontSize: '4rem', fontWeight: 800, letterSpacing: '-0.04em', margin: 0, color: 'var(--text-primary)' }}>
            Transparent topology scaling.
          </h1>
          <p className="page-lead" style={{ fontSize: '1.25rem', color: 'var(--text-secondary)', maxWidth: '800px', lineHeight: 1.6, marginTop: '1rem' }}>
            Initiate locally for free. Upgrade tier parameters to unlock live API passthrough and LLM integrations.
          </p>
        </div>
      </section>

      {/* Pricing Grid */}
      <section style={{ maxWidth: '1400px', margin: '0 auto', borderLeft: '1px solid var(--border-light)', borderRight: '1px solid var(--border-light)' }}>
        <div className="pricing-grid" style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))' }}>
          {tiers.map((tier, index) => (
            <div key={tier.name} className="pricing-card" style={{
              backgroundColor: tier.highlighted ? 'var(--bg-elevated)' : 'var(--bg-surface)',
              borderRight: index !== tiers.length - 1 ? '1px solid var(--border-light)' : 'none',
              borderBottom: '1px solid var(--border-light)',
              padding: '4rem 3rem',
              display: 'flex',
              flexDirection: 'column',
              position: 'relative'
            }}>
              {tier.highlighted && (
                <div style={{ position: 'absolute', top: 0, left: 0, right: 0, height: '4px', backgroundColor: 'var(--accent-orange)' }} />
              )}
              
              <div className="mono" style={{ fontSize: '0.85rem', color: tier.highlighted ? 'var(--accent-orange)' : 'var(--text-secondary)', marginBottom: '1rem' }}>
                // {tier.name.toUpperCase()} TIER
              </div>
              
              <div style={{ margin: '1rem 0', display: 'flex', alignItems: 'baseline', gap: '4px' }}>
                <span style={{ fontSize: '3rem', fontWeight: 700, letterSpacing: '-0.02em', color: 'var(--text-primary)' }}>{tier.price}</span>
                {tier.period && <span style={{ color: 'var(--text-secondary)' }}>{tier.period}</span>}
              </div>
              
              <p style={{ color: 'var(--text-secondary)', height: '48px', fontSize: '0.95rem', marginBottom: '2rem' }}>{tier.description}</p>
              
              <Link 
                href="#" 
                className={tier.highlighted ? "btn-primary" : "btn-secondary"} 
                style={{ width: '100%', justifyContent: 'center', marginBottom: '3rem' }}
              >
                {tier.name === 'Enterprise' ? 'Contact Sales' : 'Deploy Template'}
              </Link>
              
              <ul className="pricing-features" style={{ display: 'flex', flexDirection: 'column', gap: '1rem', flex: 1 }}>
                {tier.features.map((feature) => (
                  <li key={feature} style={{ display: 'flex', gap: '0.75rem', alignItems: 'flex-start' }}>
                    <Check size={18} style={{ color: tier.highlighted ? 'var(--accent-orange)' : 'var(--text-tertiary)', flexShrink: 0, marginTop: '2px' }} />
                    <span style={{ color: 'var(--text-primary)', fontSize: '0.95rem' }}>{feature}</span>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
