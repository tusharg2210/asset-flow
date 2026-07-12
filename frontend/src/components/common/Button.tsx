import React, { type ButtonHTMLAttributes } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'outline';
  isLoading?: boolean;
}

export const Button: React.FC<ButtonProps> = ({ 
  children, 
  variant = 'primary', 
  isLoading, 
  className = '', 
  ...props 
}) => {
  const baseStyles = "w-full px-4 py-2.5 rounded-lg font-semibold transition-all duration-200 flex justify-center items-center";
  const variants = {
    primary: "bg-orange-600 hover:bg-orange-700 text-white shadow-md hover:shadow-lg",
    outline: "bg-transparent border-2 border-slate-200 hover:border-slate-300 text-slate-700",
  };

  return (
    <button 
      className={`${baseStyles} ${variants[variant]} ${isLoading ? 'opacity-70 cursor-not-allowed' : ''} ${className}`}
      disabled={isLoading || props.disabled}
      {...props}
    >
      {isLoading ? 'Loading...' : children}
    </button>
  );
};