import React from 'react';
import { Sidebar } from './Sidebar';

interface AppLayoutProps {
  children: React.ReactNode;
}

export const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  return (
    <div className="flex h-screen bg-gray-950 font-sans">
      <Sidebar />
      <main className="flex-1 overflow-y-auto bg-[#111827]"> {/* Strict gray-900 equivalent */}
        {children}
      </main>
    </div>
  );
};