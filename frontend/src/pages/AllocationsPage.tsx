import React, { useState } from 'react';
import { AppLayout } from '../components/layout/AppLayout';
import { AlertCircle, History } from 'lucide-react';
import { mockAssetDetails, mockAllocationHistory, mockEmployees } from '../mockData/allocations';
// import { axiosClient } from '../api/axiosClient';
// import { ENDPOINTS } from '../api/endpoints';

export const AllocationsPage = () => {
  // In a real app, 'assetSearch' would trigger an API call on debounce to fetch 'assetDetails'
  const [assetSearch, setAssetSearch] = useState(`${mockAssetDetails.tag} - ${mockAssetDetails.name}`);
  const [selectedEmployee, setSelectedEmployee] = useState('');
  const [transferReason, setTransferReason] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Using mock data to simulate an already allocated asset
  const isAllocated = mockAssetDetails.status === 'Allocated';
  const holder = mockAssetDetails.currentHolder;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    
    /*
    // TODO: Uncomment when backend is ready
    try {
      const payload = {
        asset_id: mockAssetDetails.id, // Need actual ID from search result
        from_user_id: holder.id,
        to_user_id: selectedEmployee,
        reason: transferReason
      };
      await axiosClient.post(ENDPOINTS.TRANSFERS.REQUEST, payload);
      alert('Transfer request submitted successfully!');
      // Reset form or redirect
    } catch (error) {
      console.error('Failed to submit transfer request', error);
    } finally {
      setIsSubmitting(false);
    }
    */
    
    setTimeout(() => {
      setIsSubmitting(false);
      alert("Mock Transfer Request Submitted!");
    }, 800);
  };

  return (
    <AppLayout>
      <div className="p-8 max-w-4xl mx-auto space-y-8">
        
        {/* Main Form Section */}
        <section>
          <div className="space-y-6">
            
            {/* Asset Selection */}
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium text-gray-400">Asset</label>
              <input 
                type="text" 
                value={assetSearch}
                onChange={(e) => setAssetSearch(e.target.value)}
                className="w-full px-4 py-2.5 bg-gray-900 border border-gray-700 rounded-lg text-gray-200 focus:outline-none focus:ring-1 focus:ring-orange-500 focus:border-orange-500 transition-all"
              />
            </div>

            {/* Double-Allocation Block Alert */}
            {isAllocated && (
              <div className="bg-red-500/10 border border-red-500/30 rounded-lg p-4 flex gap-3">
                <AlertCircle className="text-red-400 shrink-0 mt-0.5" size={20} />
                <div>
                  <p className="text-red-400 font-medium">
                    Already Allocated to {holder.name} ({holder.department})
                  </p>
                  <p className="text-red-400/80 text-sm mt-1">
                    Direct re-allocation is blocked - submit a transfer request below
                  </p>
                </div>
              </div>
            )}

            {/* Transfer Form (Rendered conditionally if allocated) */}
            {isAllocated && (
              <form onSubmit={handleSubmit} className="space-y-6 mt-4">
                <h3 className="text-lg font-medium text-gray-200 border-b border-gray-800 pb-2">Transfer Request</h3>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {/* From (Active Input matching screenshot aesthetic) */}
                  <div className="flex flex-col gap-1.5">
                    <label className="text-sm font-medium text-gray-400">From</label>
                    <input 
                      type="text" 
                      defaultValue={holder.name}
                      className="w-full px-4 py-2.5 bg-slate-800/40 border border-slate-700/60 rounded-lg text-slate-400 focus:outline-none focus:ring-1 focus:ring-orange-500 focus:border-orange-500 transition-all"
                    />
                  </div>

                  {/* To */}
                  <div className="flex flex-col gap-1.5">
                    <label className="text-sm font-medium text-gray-400">To</label>
                    <select 
                      value={selectedEmployee}
                      onChange={(e) => setSelectedEmployee(e.target.value)}
                      required
                      className="w-full px-4 py-2.5 bg-gray-900 border border-gray-700 rounded-lg text-gray-200 focus:outline-none focus:ring-1 focus:ring-orange-500 focus:border-orange-500 transition-all appearance-none"
                    >
                      <option value="" disabled>Select Employee....</option>
                      {mockEmployees.map(emp => (
                        <option key={emp.id} value={emp.id}>{emp.name}</option>
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
          </div>
        </section>

        {/* Allocation History Section */}
        <section className="pt-8">
          <div className="flex items-center gap-2 mb-4">
            <History size={18} className="text-gray-400" />
            <h3 className="text-lg font-medium text-gray-200">Allocation history</h3>
          </div>
          
          <div className="border-t border-gray-800 pt-4 space-y-3">
            {mockAllocationHistory.map((history) => (
              <div key={history.id} className="flex text-sm">
                <span className="w-20 text-gray-500 shrink-0 font-mono">{history.date}</span>
                <span className="text-gray-300">- {history.action}</span>
              </div>
            ))}
          </div>
        </section>

      </div>
    </AppLayout>
  );
};