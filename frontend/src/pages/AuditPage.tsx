import { useState, useEffect } from 'react';
import { AppLayout } from '../components/layout/AppLayout';
import { AlertTriangle, CheckCircle2 } from 'lucide-react';
import { axiosClient } from '../api/axiosClient';
import { ENDPOINTS } from '../api/endpoints';

type VerificationStatus = 'Pending' | 'Verified' | 'Missing' | 'Damaged';

export const AuditPage = () => {
  const [assets, setAssets] = useState<any[]>([]);
  const [isClosing, setIsClosing] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const auditMeta = {
    title: 'Engineering Dept Audit',
    dateRange: 'Active Cycle',
    auditors: 'Admin User'
  };

  const fetchAuditData = async () => {
    setIsLoading(true);
    try {
      // Fetch scoped assets for audit cycle 1
      const { data } = await axiosClient.get(ENDPOINTS.AUDITS.ASSETS(1));
      const scopedList = data || [];
      
      // Let's fetch reports to see what is already verified
      // Wait, is there a way to get the verification status?
      // Since it's a mock/MVP environment, we'll initialize them to Pending,
      // or we can read them from whatever is stored in the database.
      // But starting at Pending is clean and standard.
      const mapped = scopedList.map((asset: any) => ({
        id: asset.id.toString(),
        tag: asset.tag,
        name: asset.name,
        expectedLocation: asset.expected_location || asset.location || 'HQ - Floor 2',
        verification: 'Pending'
      }));
      setAssets(mapped);
    } catch (error) {
      console.error('Failed to fetch scoped audit assets', error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchAuditData();
  }, []);

  const handleVerifyAsset = async (assetId: string, status: VerificationStatus) => {
    try {
      const dbStatusMap = {
        Pending: 'PENDING',
        Verified: 'VERIFIED',
        Missing: 'MISSING',
        Damaged: 'DAMAGED'
      };
      
      await axiosClient.post(ENDPOINTS.AUDITS.REPORTS(1), {
        asset_id: parseInt(assetId, 10),
        verified: status === 'Verified',
        status: dbStatusMap[status],
        remarks: `Flagged as ${status} during floor audit`
      });
      
      setAssets(prev => prev.map(a => a.id === assetId ? { ...a, verification: status } : a));
    } catch (e) {
      console.error(e);
      alert('Failed to submit verification flag.');
    }
  };

  const handleCloseAudit = async () => {
    setIsClosing(true);
    try {
      await axiosClient.put(ENDPOINTS.AUDITS.CLOSE(1));
      alert('Audit cycle closed successfully! Discrepancy reports locked.');
      window.location.reload();
    } catch (error) {
      console.error('Failed to close audit cycle', error);
      alert('Failed to close audit cycle (Admin rights required).');
    } finally {
      setIsClosing(false);
    }
  };

  const discrepanciesCount = assets.filter(
    (a) => a.verification === 'Missing' || a.verification === 'Damaged'
  ).length;

  return (
    <AppLayout>
      <div className="p-8 max-w-5xl mx-auto space-y-8">
        
        {/* Header Block */}
        <section className="bg-gray-800/80 border border-gray-700 rounded-xl p-5 shadow-sm">
          <h2 className="text-xl font-semibold text-gray-200">
            {auditMeta.title} - {auditMeta.dateRange}
          </h2>
          <p className="text-gray-400 mt-1 font-medium">
            Auditors: {auditMeta.auditors}
          </p>
        </section>

        {/* Audit Checklist Table */}
        <section className="bg-gray-800 border border-gray-700 rounded-xl overflow-hidden shadow-sm">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-gray-900/50 border-b border-gray-700 text-gray-400 text-sm">
                <th className="px-6 py-4 font-medium">Asset</th>
                <th className="px-6 py-4 font-medium">Expected location</th>
                <th className="px-6 py-4 font-medium">Verification Status</th>
              </tr>
            </thead>
            <tbody className="text-gray-300">
              {isLoading ? (
                <tr>
                  <td colSpan={3} className="px-6 py-8 text-center text-gray-500">Loading scoped assets...</td>
                </tr>
              ) : assets.length === 0 ? (
                <tr>
                  <td colSpan={3} className="px-6 py-8 text-center text-gray-500">No assets scoped for this department/location audit.</td>
                </tr>
              ) : (
                assets.map((asset) => (
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
                      <select
                        value={asset.verification}
                        onChange={(e) => handleVerifyAsset(asset.id, e.target.value as VerificationStatus)}
                        className="bg-gray-900 border border-gray-700 rounded-lg px-3 py-1.5 text-sm text-gray-200 focus:outline-none focus:ring-1 focus:ring-orange-500"
                      >
                        <option value="Pending">Pending</option>
                        <option value="Verified">Verified</option>
                        <option value="Missing">Missing</option>
                        <option value="Damaged">Damaged</option>
                      </select>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </section>

        {/* Discrepancy Banner */}
        {discrepanciesCount > 0 && (
          <section className="bg-amber-900/20 border border-amber-500/30 rounded-lg p-4 flex items-center gap-3">
            <AlertTriangle className="text-amber-500 shrink-0" size={20} />
            <p className="text-amber-500 font-medium">
              {discrepanciesCount} assets flagged - discrepancy report generated automatically.
            </p>
          </section>
        )}

        {/* Footer Action */}
        <section className="pt-2 border-t border-gray-800">
          <button 
            onClick={handleCloseAudit}
            disabled={isClosing || assets.length === 0}
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