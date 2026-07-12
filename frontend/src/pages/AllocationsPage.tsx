import React, { useState, useEffect } from 'react';
import { AppLayout } from '../components/layout/AppLayout';
import { AlertCircle, History } from 'lucide-react';
import { axiosClient } from '../api/axiosClient';
import { ENDPOINTS } from '../api/endpoints';

export const AllocationsPage = () => {
  const [assets, setAssets] = useState<any[]>([]);
  const [employees, setEmployees] = useState<any[]>([]);
  const [assetSearch, setAssetSearch] = useState('AF-0004'); // default tag matching our seeded asset
  const [assetDetails, setAssetDetails] = useState<any>(null);
  const [allocationHistory, setAllocationHistory] = useState<any[]>([]);
  
  const [selectedEmployee, setSelectedEmployee] = useState('');
  const [transferReason, setTransferReason] = useState('');
  const [expectedReturnDate, setExpectedReturnDate] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    const initData = async () => {
      try {
        const [assetsRes, employeesRes] = await Promise.all([
          axiosClient.get(ENDPOINTS.ASSETS.DIRECTORY),
          axiosClient.get(ENDPOINTS.ORGANIZATION.USERS)
        ]);
        const assetsList = assetsRes.data?.data || assetsRes.data || [];
        setAssets(assetsList);
        setEmployees(employeesRes.data || []);
      } catch (err) {
        console.error("Failed to load initial data", err);
      }
    };
    initData();
  }, []);

  useEffect(() => {
    const match = assets.find(a => a.tag.toUpperCase() === assetSearch.trim().toUpperCase());
    if (match) {
      const loadDetails = async () => {
        try {
          const [detailRes, historyRes] = await Promise.all([
            axiosClient.get(ENDPOINTS.ASSETS.DETAILS(match.id)),
            axiosClient.get(`/api/assets/${match.id}/history`)
          ]);
          setAssetDetails(detailRes.data);
          
          const allocations = historyRes.data?.allocations || [];
          const mapped = allocations.map((item: any, idx: number) => ({
            id: idx,
            date: item.allotted_date ? new Date(item.allotted_date).toLocaleDateString() : 'Recent',
            action: `Allocated to user #${item.to_user_id}. ${item.actual_return_date ? 'Returned.' : 'Currently active.'}`
          }));
          setAllocationHistory(mapped);
        } catch (e) {
          console.error(e);
        }
      };
      loadDetails();
    } else {
      setAssetDetails(null);
      setAllocationHistory([]);
    }
  }, [assetSearch, assets]);

  const handleTransferSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!assetDetails) return;
    setIsSubmitting(true);
    try {
      const payload = {
        asset_id: assetDetails.id,
        from_user_id: assetDetails.current_holder?.user_id || 0,
        to_user_id: parseInt(selectedEmployee, 10),
        reason: transferReason
      };
      await axiosClient.post(ENDPOINTS.TRANSFERS.REQUEST, payload);
      alert('Transfer request submitted successfully!');
      setTransferReason('');
      setSelectedEmployee('');
      // Refresh asset details
      setAssetSearch(assetSearch + ' ');
      setAssetSearch(assetSearch.trim());
    } catch (error) {
      console.error('Failed to submit transfer request', error);
      alert('Failed to submit transfer request.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDirectAllocateSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!assetDetails) return;
    setIsSubmitting(true);
    try {
      const payload = {
        asset_id: assetDetails.id,
        to_user_id: parseInt(selectedEmployee, 10),
        expected_return_date: new Date(expectedReturnDate).toISOString(),
        reason: transferReason
      };
      await axiosClient.post(ENDPOINTS.ALLOCATIONS.CREATE, payload);
      alert('Asset allocated successfully!');
      setTransferReason('');
      setSelectedEmployee('');
      setExpectedReturnDate('');
      // Refresh details
      setAssetSearch(assetSearch + ' ');
      setAssetSearch(assetSearch.trim());
    } catch (error) {
      console.error('Failed to allocate asset', error);
      alert('Failed to allocate asset.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const isAllocated = assetDetails?.status === 'ALLOCATED';
  const holder = assetDetails?.current_holder;

  return (
    <AppLayout>
      <div className="p-8 max-w-4xl mx-auto space-y-8">
        
        {/* Main Form Section */}
        <section>
          <div className="space-y-6">
            
            {/* Asset Selection */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium text-gray-400">Search Asset Tag</label>
              <input 
                type="text" 
                value={assetSearch}
                onChange={(e) => setAssetSearch(e.target.value)}
                placeholder="e.g. AF-0004"
                className="w-full px-4 py-2.5 bg-gray-900 border border-gray-700 rounded-lg text-gray-200 focus:outline-none focus:ring-1 focus:ring-orange-500 focus:border-orange-500 transition-all"
              />
              {assetDetails ? (
                <p className="text-xs text-emerald-400 mt-1">
                  Matched asset: {assetDetails.name} ({assetDetails.status})
                </p>
              ) : (
                <p className="text-xs text-gray-500 mt-1">
                  Type a registered asset tag (e.g. AF-0004) to load details.
                </p>
              )}
            </div>

            {/* Double-Allocation Block Alert */}
            {isAllocated && holder && (
              <div className="bg-red-500/10 border border-red-500/30 rounded-lg p-4 flex gap-3">
                <AlertCircle className="text-red-400 shrink-0 mt-0.5" size={20} />
                <div>
                  <p className="text-red-400 font-medium">
                    Already Allocated to User ID #{holder.user_id}
                  </p>
                  <p className="text-red-400/80 text-sm mt-1">
                    Direct re-allocation is blocked - submit a transfer request below.
                  </p>
                </div>
              </div>
            )}

            {/* Transfer Form (Rendered conditionally if allocated) */}
            {isAllocated && assetDetails && (
              <form onSubmit={handleTransferSubmit} className="space-y-6 mt-4">
                <h3 className="text-lg font-medium text-gray-200 border-b border-gray-800 pb-2">Transfer Request</h3>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {/* From */}
                  <div className="flex flex-col gap-1.5">
                    <label className="text-sm font-medium text-gray-400">From User ID</label>
                    <input 
                      type="text" 
                      readOnly
                      value={holder?.user_id || ''}
                      className="w-full px-4 py-2.5 bg-slate-800/40 border border-slate-700/60 rounded-lg text-slate-400 focus:outline-none"
                    />
                  </div>

                  {/* To */}
                  <div className="flex flex-col gap-1.5">
                    <label className="text-sm font-medium text-gray-400">To Employee</label>
                    <select 
                      value={selectedEmployee}
                      onChange={(e) => setSelectedEmployee(e.target.value)}
                      required
                      className="w-full px-4 py-2.5 bg-gray-900 border border-gray-700 rounded-lg text-gray-200 focus:outline-none focus:ring-1 focus:ring-orange-500 focus:border-orange-500 transition-all appearance-none"
                    >
                      <option value="" disabled>Select Employee...</option>
                      {employees.map(emp => (
                        <option key={emp.id} value={emp.id}>{emp.name} ({emp.email})</option>
                      ))}
                    </select>
                  </div>
                </div>

                {/* Reason Textarea */}
                <div className="flex flex-col gap-1.5">
                  <label className="text-sm font-medium text-gray-400">Reason</label>
                  <textarea 
                    rows={4}
                    value={transferReason}
                    onChange={(e) => setTransferReason(e.target.value)}
                    required
                    placeholder="Provide transfer justification..."
                    className="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-lg text-gray-200 focus:outline-none focus:ring-1 focus:ring-orange-500 focus:border-orange-500 transition-all resize-none"
                  />
                </div>

                {/* Submit Button */}
                <button 
                  type="submit" 
                  disabled={isSubmitting}
                  className="px-6 py-2.5 bg-emerald-600/10 border border-emerald-600 text-emerald-500 hover:bg-emerald-600 hover:text-white rounded-lg font-medium transition-all shadow-sm disabled:opacity-50"
                >
                  {isSubmitting ? 'Submitting...' : 'Submit Request'}
                </button>
              </form>
            )}

            {/* Direct Allocation Form (Rendered if available) */}
            {assetDetails && !isAllocated && (
              <form onSubmit={handleDirectAllocateSubmit} className="space-y-6 mt-4">
                <h3 className="text-lg font-medium text-gray-200 border-b border-gray-800 pb-2">Direct Allocation</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {/* To Employee */}
                  <div className="flex flex-col gap-1.5">
                    <label className="text-sm font-medium text-gray-400">Allocate To</label>
                    <select 
                      value={selectedEmployee}
                      onChange={(e) => setSelectedEmployee(e.target.value)}
                      required
                      className="w-full px-4 py-2.5 bg-gray-900 border border-gray-700 rounded-lg text-gray-200 focus:outline-none focus:ring-1 focus:ring-orange-500 focus:border-orange-500 transition-all appearance-none"
                    >
                      <option value="" disabled>Select Employee...</option>
                      {employees.map(emp => (
                        <option key={emp.id} value={emp.id}>{emp.name} ({emp.email})</option>
                      ))}
                    </select>
                  </div>
                  {/* Expected Return Date */}
                  <div className="flex flex-col gap-1.5">
                    <label className="text-sm font-medium text-gray-400">Expected Return Date</label>
                    <input 
                      type="date" 
                      value={expectedReturnDate}
                      onChange={(e) => setExpectedReturnDate(e.target.value)}
                      required
                      className="w-full px-4 py-2.5 bg-gray-900 border border-gray-700 rounded-lg text-gray-200 focus:outline-none focus:ring-1 focus:ring-orange-500 focus:border-orange-500 transition-all"
                    />
                  </div>
                </div>
                {/* Reason */}
                <div className="flex flex-col gap-1.5">
                  <label className="text-sm font-medium text-gray-400">Reason</label>
                  <textarea 
                    rows={4}
                    value={transferReason}
                    onChange={(e) => setTransferReason(e.target.value)}
                    required
                    placeholder="State the purpose of this allocation..."
                    className="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-lg text-gray-200 focus:outline-none focus:ring-1 focus:ring-orange-500 focus:border-orange-500 transition-all resize-none"
                  />
                </div>
                <button 
                  type="submit" 
                  disabled={isSubmitting}
                  className="px-6 py-2.5 bg-orange-600/10 border border-orange-600 text-orange-500 hover:bg-orange-600 hover:text-white rounded-lg font-medium transition-all shadow-sm disabled:opacity-50"
                >
                  {isSubmitting ? 'Allocating...' : 'Allocate Asset'}
                </button>
              </form>
            )}
          </div>
        </section>

        {/* Allocation History Section */}
        <section className="pt-8">
          <div className="flex items-center gap-2 mb-4">
            <History size={18} className="text-gray-400" />
            <h3 className="text-lg font-medium text-gray-200">Allocation history</h3>
          </div>
          
          <div className="border-t border-gray-800 pt-4 space-y-3">
            {allocationHistory.length === 0 ? (
              <p className="text-gray-500 text-sm">No allocation history for this asset.</p>
            ) : (
              allocationHistory.map((history) => (
                <div key={history.id} className="flex text-sm">
                  <span className="w-24 text-gray-500 shrink-0 font-mono">{history.date}</span>
                  <span className="text-gray-300">- {history.action}</span>
                </div>
              ))
            )}
          </div>
        </section>

      </div>
    </AppLayout>
  );
};