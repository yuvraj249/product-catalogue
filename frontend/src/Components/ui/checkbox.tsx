import React from 'react';

interface CheckboxProps extends React.InputHTMLAttributes<HTMLInputElement> {
  indeterminate?: boolean;
}

export const Checkbox: React.FC<CheckboxProps> = ({
  indeterminate,
  className = '',
  ...props
}) => {
  const ref = React.useRef<HTMLInputElement>(null);

  React.useEffect(() => {
    if (ref.current) {
      ref.current.indeterminate = Boolean(indeterminate);
    }
  }, [indeterminate]);

  return (
    <input
      type="checkbox"
      ref={ref}
      className={`h-4 w-4 rounded border-slate-700 bg-slate-900 text-indigo-600 focus:ring-indigo-500 transition cursor-pointer ${className}`}
      {...props}
    />
  );
};
