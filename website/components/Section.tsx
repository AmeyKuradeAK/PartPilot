'use client';

import { ReactNode } from 'react';
import { motion } from 'framer-motion';

interface SectionProps {
  children: ReactNode;
  className?: string;
  id?: string;
  style?: React.CSSProperties;
}

export default function Section({ children, className = '', id, style }: SectionProps) {
  return (
    <motion.section 
      id={id}
      className={`section-padding ${className}`}
      style={style}
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, margin: "-100px" }}
      transition={{ duration: 0.5, ease: 'easeOut' }}
    >
      <div className="container">
        {children}
      </div>
    </motion.section>
  );
}
