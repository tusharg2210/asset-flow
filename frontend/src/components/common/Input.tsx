import React, { type InputHTMLAttributes } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
}

export const Input: React.FC<InputProps> = ({ label, ...props }) => {
  return (
    <div className="flex flex-col gap-1.5">
      <label className="text-sm font-medium text-slate-700">{label}</label>
      <input
        className="px-4 py-2.5 rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-all text-slate-900 bg-white"
        {...props}
      />
    </div>
  );
};