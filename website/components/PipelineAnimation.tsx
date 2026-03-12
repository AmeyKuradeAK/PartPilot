'use client';

import { motion } from 'framer-motion';
import { Database, Cpu, Network, CheckCircle, FileText } from 'lucide-react';
import { useEffect, useState } from 'react';

const steps = [
  { id: 1, title: 'Parse BOM', icon: FileText, desc: 'Ingest CSV/Excel' },
  { id: 2, title: 'AI Normalize', icon: Cpu, desc: 'Convert to exact MPN' },
  { id: 3, title: 'Fan-Out', icon: Network, desc: '50x Concurrent API queries' },
  { id: 4, title: 'Rank & Rank', icon: Database, desc: 'Filter MOQ/Lead times' },
  { id: 5, title: 'Output PO', icon: CheckCircle, desc: 'Generate Purchase Order' },
];

export default function PipelineAnimation() {
  const [activeStep, setActiveStep] = useState(1);

  useEffect(() => {
    const timer = setInterval(() => {
      setActiveStep((prev) => (prev >= steps.length ? 1 : prev + 1));
    }, 2000);
    return () => clearInterval(timer);
  }, []);

  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      gap: '2rem',
      padding: '3rem',
      backgroundColor: 'var(--bg-elevated)',
      borderRadius: '16px',
      border: '1px solid var(--border)',
      position: 'relative',
      overflow: 'hidden'
    }}>
      {/* Background Grid Pattern */}
      <div style={{
        position: 'absolute', inset: 0, opacity: 0.1,
        backgroundImage: 'linear-gradient(var(--border) 1px, transparent 1px), linear-gradient(90deg, var(--border) 1px, transparent 1px)',
        backgroundSize: '20px 20px'
      }} />

      <h3 style={{ fontSize: '1.25rem', fontWeight: 600, color: 'var(--text-primary)', zIndex: 1 }}>
        Live Architecture Simulation
      </h3>

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', position: 'relative', zIndex: 1 }}>
        {/* Connecting Line */}
        <div style={{ position: 'absolute', top: '50%', left: 0, right: 0, height: '2px', backgroundColor: 'var(--border)', zIndex: -1 }} />

        {steps.map((step, index) => {
          const isActive = index + 1 === activeStep;
          const isPast = index + 1 < activeStep;
          const Icon = step.icon;

          return (
            <div key={step.id} style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '1rem', width: '120px' }}>
              <motion.div 
                animate={{
                  backgroundColor: isActive ? 'var(--accent)' : isPast ? 'var(--bg-surface)' : 'var(--bg-primary)',
                  borderColor: isActive ? 'var(--accent)' : 'var(--border)',
                  scale: isActive ? 1.1 : 1,
                  boxShadow: isActive ? '0 0 20px var(--accent-glow)' : 'none'
                }}
                transition={{ duration: 0.3 }}
                style={{
                  width: '64px', height: '64px', borderRadius: '50%', border: '2px solid',
                  display: 'flex', alignItems: 'center', justifyContent: 'center',
                  color: isActive ? '#fff' : 'var(--text-secondary)'
                }}
              >
                <Icon size={28} />
              </motion.div>
              
              <div style={{ textAlign: 'center' }}>
                <div style={{ 
                  fontWeight: 600, fontSize: '0.9rem', 
                  color: isActive ? 'var(--text-primary)' : 'var(--text-secondary)' 
                }}>
                  {step.title}
                </div>
                <div style={{ fontSize: '0.75rem', color: 'var(--text-secondary)', marginTop: '4px' }}>
                  {step.desc}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
