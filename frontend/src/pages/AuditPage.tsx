import React, { useState } from 'react';
import { AppLayout } from '../components/layout/AppLayout';
import { mockAuditCycle, mockAuditAssets, type VerificationStatus } from '../mockData/audit';
import { AlertTriangle, CheckCircle2 } from 'lucide-react';
// import { axiosClient } from '../api/axiosClient';
// import { ENDPOINTS } from '../api/endpoints';

export const AuditPage = () => {
  const [assets, setAssets] = useState(mockAuditAssets);
  const [isClosing, setIsClosing] = useState(false);

  // Calculate discrepancies (anything not 'Verified' or 'Pending')
  const discrepanciesCount = assets.filter(
    (a) => a.verification === 'Missing' || a.verification === 'Damaged'
  ).length;

  const handleCloseAudit = async () => {
    setIsClosing(true);
    /*
    // TODO: Uncomment when backend is ready
    try {
      await axiosClient.put(ENDPOINTS.AUDITS.CLOSE(mockAuditCycle.id));
      alert('Audit cycle closed successfully!');
    } catch (error) {
      console.error('Failed to close audit cycle', error);
    } finally {
      setIsClosing(false);
    }
    */
    setTimeout(() => {
      setIsClosing(false);
      alert('Mock Audit Cycle Closed!');
    }, 800);
  };

  // Helper component for the specific audit status pills
  const VerificationBadge = ({ status }: { status: VerificationStatus }) => {
    let styles = '';
    switch (status) {
      case 'Verified':
        styles = 'border-emerald-500/50 text-emerald-400 bg-emerald-500/10';
        break;
      case 'Missing':
        styles = 'border-red-500/50 text-red-400 bg-red-500/10';
        break;
      case 'Damaged':
        styles = 'border-amber-500/50 text-amber-400 bg-amber-500/10';
        break;
      default:
        styles = 'border-gray-500/50 text-gray-400 bg-gray-500/10';
    }
    
    return (
      <span className={`px-4 py-1.5 text-xs font-medium rounded-full border ${styles} inline-block min-w-[90px] text-center`}>
        {status}
      </span>
    );
  };

  return (
    <AppLayout>
      <div className="p-8 max-w-5xl mx-auto space-y-8">
        
        {/* Header Block */}
        <section className="bg-gray-800/80 border border-gray-700 rounded-xl p-5 shadow-sm">
          <h2 className="text-xl font-semibold text-gray-200">
            {mockAuditCycle.title} - {mockAuditCycle.dateRange}
          </h2>
          <p className="text-gray-400 mt-1">
            Auditors: {mockAuditCycle.auditors}
          </p>
        </section>

        {/* Audit Checklist Table */}
        <section className="bg-gray-800 border border-gray-700 rounded-xl overflow-hidden shadow-sm">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-gray-900/50 border-b border-gray-700 text-gray-400 text-sm">
                <th className="px-6 py-4 font-medium">Asset</th>
                <th className="px-6 py-4 font-medium">Expected location</th>
                <th className="px-6 py-4 font-medium">Verification</th>
              </tr>
            </thead>
            <tbody className="text-gray-300">
              {assets.map((asset) => (
                <tr 
                  key={asset.id} 
                  className="border-b border-gray-700/50 last:border-0 hover:bg-gray-700/20 transition-colors"
                >
                  <td className="px-6 py-4">
                    <span className="font-mono text-gray-400 mr-2">{asset.tag}</span>
                    <span className="font-medium text-gray-200">{asset.name}</span>
                  </td>
                  <td className="px-6 py-4 text-gray-400">{asset.expectedLocation}</td>
                  <td className="px-6 py-4">
                    <VerificationBadge status={asset.verification} />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>

        {/* Discrepancy Banner */}
        {discrepanciesCount > 0 && (
          <section className="bg-amber-900/20 border border-amber-500/30 rounded-lg p-4 flex items-center gap-3">
            <AlertTriangle className="text-amber-500 shrink-0" size={20} />
            <p className="text-amber-500 font-medium">
              {discrepanciesCount} assets flagged - discrepancy report generated automatically
            </p>
          </section>
        )}

        {/* Footer Action */}
        <section className="pt-2 border-t border-gray-800">
          <button 
            onClick={handleCloseAudit}
            disabled={isClosing}
            className="flex items-center gap-2 px-6 py-2.5 bg-gray-800 border-2 border-gray-700 hover:border-gray-500 text-gray-300 hover:text-white rounded-lg font-medium transition-all shadow-sm disabled:opacity-50"
          >
            <CheckCircle2 size={18} className={isClosing ? 'text-gray-500' : 'text-gray-400'} />
            {isClosing ? 'Closing Cycle...' : 'Close audit cycle'}
          </button>
        </section>

      </div>
    </AppLayout>
  );
};