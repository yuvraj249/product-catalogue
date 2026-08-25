import React from 'react';

interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: 'default' | 'success' | 'warning' | 'destructive' | 'outline';
  children: React.ReactNode;
}

export const Badge: React.FC<BadgeProps> = ({
  variant = 'default',
  children,
  className = '',
  ...props
}) => {
  const baseStyles = 'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold tracking-wide transition-colors';

  const variantStyles = {
    default: 'bg-indigo-950 text-indigo-300 border border-indigo-700/50',
    success: 'bg-emerald-950 text-emerald-300 border border-emerald-700/50',
    warning: 'bg-amber-950 text-amber-300 border border-amber-700/50',
    destructive: 'bg-rose-950 text-rose-300 border border-rose-700/50',
    outline: 'bg-transparent text-slate-300 border border-slate-700',
  };

  return (
    <span className={`${baseStyles} ${variantStyles[variant]} ${className}`} {...props}>
      {children}
    </span>
  );
};
