import React, { type InputHTMLAttributes } from 'react';

interface FormInputProps extends InputHTMLAttributes<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement> {
  label: string;
  as?: 'input' | 'select' | 'textarea';
  options?: { label: string; value: string }[];
}

export const FormInput: React.FC<FormInputProps> = ({ 
  label, 
  as = 'input', 
  options, 
  className = '', 
  ...props 
}) => {
  const baseStyles = "w-full px-4 py-2.5 bg-gray-900/50 border border-gray-700 rounded-lg text-gray-200 placeholder-gray-600 focus:outline-none focus:border-orange-500 focus:ring-1 focus:ring-orange-500 transition-all";

  return (
    <div className="flex flex-col gap-1.5 w-full">
      <label className="text-sm font-medium text-gray-400">{label}</label>
      
      {as === 'select' ? (
        <select className={`${baseStyles} appearance-none ${className}`} {...(props as any)}>
          <option value="" disabled>Select an option...</option>
          {options?.map(opt => (
            <option key={opt.value} value={opt.value}>{opt.label}</option>
          ))}
        </select>
      ) : as === 'textarea' ? (
        <textarea className={`${baseStyles} resize-none ${className}`} rows={3} {...(props as any)} />
      ) : (
        <input className={`${baseStyles} ${className}`} {...(props as any)} />
      )}
    </div>
  );
};