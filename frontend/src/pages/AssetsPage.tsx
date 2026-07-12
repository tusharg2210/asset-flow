import React, { useState, useEffect } from 'react';
import { AppLayout } from '../components/layout/AppLayout';
import { Badge } from '../components/common/Badge';
import { mockAssets } from '../mockData/assets';
import { Plus, Search, ChevronDown } from 'lucide-react';
import { Modal } from '../components/common/Modal';
import { FormInput } from '../components/common/FormInput';
import { Button } from '../components/common/Button';
// import { axiosClient } from '../api/axiosClient';
// import { ENDPOINTS } from '../api/endpoints';

export const AssetsPage = () => {
  const [assets, setAssets] = useState(mockAssets);
  const [searchQuery, setSearchQuery] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isRegisterModalOpen, setIsRegisterModalOpen] = useState(false);

  /*
  // TODO: Uncomment when backend is ready
  useEffect(() => {
    const fetchAssets = async () => {
      setIsLoading(true);
      try {
        const { data } = await axiosClient.get(ENDPOINTS.ASSETS.DIRECTORY);
        setAssets(data);
      } catch (error) {
        console.error("Failed to fetch assets", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchAssets();
  }, []);
  */

  const handleRegisterSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    alert('Asset registered successfully! (Mock)');
    setIsRegisterModalOpen(false);
  };

  // Quick reusable component for the filter buttons
  const FilterButton = ({ label }: { label: string }) => (
    <button className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-gray-300 bg-gray-800 border border-gray-700 rounded-lg hover:border-gray-500 transition-colors">
      {label}
      <ChevronDown size={14} className="text-gray-500" />
    </button>
  );

  return (
    <AppLayout>
      <div className="p-8 max-w-7xl mx-auto space-y-6">
        
        {/* Top Header: Search and Action */}
        <div className="flex justify-between items-center gap-4">
          <div className="relative flex-1 max-w-2xl">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500" size={18} />
            <input 
              type="text" 
              placeholder="Search by tag, serial, or QR code.." 
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-2.5 bg-gray-800 border border-gray-700 rounded-lg text-gray-200 placeholder-gray-500 focus:outline-none focus:ring-1 focus:ring-orange-500 focus:border-orange-500 transition-all shadow-sm"
            />
          </div>
          
          <button 
            onClick={() => setIsRegisterModalOpen(true)}
            className="flex items-center gap-2 px-5 py-2.5 rounded-lg font-semibold transition-all duration-200 border border-orange-600 bg-orange-600/10 text-orange-500 hover:bg-orange-600 hover:text-white shadow-sm whitespace-nowrap"
          >
            <Plus size={18} />
            Register Asset
          </button>
        </div>

        {/* Filters Row */}
        <div className="flex gap-3">
          <FilterButton label="Category" />
          <FilterButton label="Status" />
          <FilterButton label="Department" />
        </div>

        {/* Data Table */}
        <div className="bg-gray-800 border border-gray-700 rounded-xl overflow-hidden shadow-sm mt-4">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-gray-900/50 border-b border-gray-700 text-gray-400 text-sm">
                <th className="px-6 py-4 font-medium">Tag</th>
                <th className="px-6 py-4 font-medium">Name</th>
                <th className="px-6 py-4 font-medium">Category</th>
                <th className="px-6 py-4 font-medium">Status</th>
                <th className="px-6 py-4 font-medium">Location</th>
              </tr>
            </thead>
            <tbody className="text-gray-300">
              {isLoading ? (
                <tr>
                  <td colSpan={5} className="px-6 py-8 text-center text-gray-500">
                    Loading directory...
                  </td>
                </tr>
              ) : (
                assets.map((asset) => (
                  <tr 
                    key={asset.id} 
                    className="border-b border-gray-700/50 last:border-0 hover:bg-gray-700/20 transition-colors cursor-pointer"
                  >
                    <td className="px-6 py-4 font-mono text-sm text-gray-400">{asset.tag}</td>
                    <td className="px-6 py-4 font-medium text-gray-200">{asset.name}</td>
                    <td className="px-6 py-4">{asset.category}</td>
                    <td className="px-6 py-4">
                      <Badge status={asset.status} />
                    </td>
                    <td className="px-6 py-4 text-gray-400">{asset.location}</td>
                  </tr>
                ))
              )}
              {assets.length === 0 && !isLoading && (
                <tr>
                  <td colSpan={5} className="px-6 py-8 text-center text-gray-500">
                    No assets found.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

      </div>

      <Modal isOpen={isRegisterModalOpen} onClose={() => setIsRegisterModalOpen(false)} title="Register New Asset">
        <form onSubmit={handleRegisterSubmit} className="space-y-4">
          <FormInput label="Asset Name" placeholder="e.g. MacBook Pro 16" required />
          <FormInput 
            label="Category" 
            as="select" 
            options={[
              { label: 'Electronics', value: 'electronics' },
              { label: 'Furniture', value: 'furniture' },
              { label: 'Vehicles', value: 'vehicles' }
            ]} 
            required 
          />
          <FormInput label="Location" placeholder="e.g. HQ - Floor 2" required />
          <FormInput label="Acquisition Cost ($)" type="number" placeholder="2400.00" />
          <div className="pt-4">
            <Button type="submit">Submit Registration</Button>
          </div>
        </form>
      </Modal>
    </AppLayout>
  );
};