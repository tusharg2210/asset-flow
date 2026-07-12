import React from 'react';

interface BadgeProps {
  status: string;
}

export const Badge: React.FC<BadgeProps> = ({ status }) => {
  const normalizedStatus = status.toLowerCase();
  
  let colorStyles = '';

  switch (normalizedStatus) {
    case 'available':
    case 'active':
      colorStyles = 'border-emerald-500/50 text-emerald-400 bg-emerald-500/10';
      break;
    case 'allocated':
      colorStyles = 'border-blue-500/50 text-blue-400 bg-blue-500/10';
      break;
    case 'maintenance':
    case 'under_maintenance':
      colorStyles = 'border-amber-500/50 text-amber-400 bg-amber-500/10';
      break;
    default:
      // Inactive, Lost, Retired, etc.
      colorStyles = 'border-gray-500/50 text-gray-400 bg-gray-500/10';
  }
  
  return (
    <span className={`px-3 py-1 text-xs font-medium rounded-full border ${colorStyles}`}>
      {status}
    </span>
  );
};