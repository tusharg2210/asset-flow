import React, { useState, useEffect } from 'react';
import { AppLayout } from '../components/layout/AppLayout';
import { Badge } from '../components/common/Badge';
import { Plus, Search, ChevronDown } from 'lucide-react';
import { Modal } from '../components/common/Modal';
import { FormInput } from '../components/common/FormInput';
import { Button } from '../components/common/Button';
import { axiosClient } from '../api/axiosClient';
import { ENDPOINTS } from '../api/endpoints';

export const AssetsPage = () => {
  const [assets, setAssets] = useState<any[]>([]);
  const [categories, setCategories] = useState<any[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isRegisterModalOpen, setIsRegisterModalOpen] = useState(false);

  // Form States
  const [regName, setRegName] = useState('');
  const [regCategoryId, setRegCategoryId] = useState('');
  const [regLocation, setRegLocation] = useState('');
  const [regCost, setRegCost] = useState('');
  const [regCondition, setRegCondition] = useState('New');
  const [isSharable, setIsSharable] = useState(true);
  const [isBookable, setIsBookable] = useState(false);

  const fetchAssetsAndCategories = async () => {
    setIsLoading(true);
    try {
      const [assetsRes, categoriesRes] = await Promise.all([
        axiosClient.get(ENDPOINTS.ASSETS.DIRECTORY),
        axiosClient.get(ENDPOINTS.ORGANIZATION.ASSET_CATEGORIES)
      ]);
      const assetsList = assetsRes.data?.data || assetsRes.data || [];
      setAssets(assetsList);
      setCategories(categoriesRes.data || []);
    } catch (error) {
      console.error("Failed to fetch assets and categories", error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchAssetsAndCategories();
  }, []);

  const handleRegisterSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload = {
        name: regName,
        category_id: regCategoryId ? parseInt(regCategoryId, 10) : null,
        location: regLocation,
        condition: regCondition,
        acquisition_cost: regCost ? parseFloat(regCost) : 0,
        is_sharable: isSharable,
        is_bookable: isBookable,
        serial_number: "SN-" + Math.floor(Math.random() * 1000000)
      };
      await axiosClient.post(ENDPOINTS.ASSETS.REGISTER, payload);
      alert('Asset registered successfully!');
      setIsRegisterModalOpen(false);
      // Reset form
      setRegName('');
      setRegCategoryId('');
      setRegLocation('');
      setRegCost('');
      fetchAssetsAndCategories();
    } catch (error) {
      console.error('Failed to register asset:', error);
      alert('Failed to register asset.');
    }
  };

  const getCategoryName = (categoryId: number) => {
    const matched = categories.find(c => c.id === categoryId);
    return matched ? matched.name : 'Electronics';
  };

  const filteredAssets = assets.filter(asset => {
    const query = searchQuery.toLowerCase();
    return (
      asset.name?.toLowerCase().includes(query) ||
      asset.tag?.toLowerCase().includes(query) ||
      asset.serial_number?.toLowerCase().includes(query) ||
      asset.location?.toLowerCase().includes(query)
    );
  });

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
              placeholder="Search by name, tag, location, serial..." 
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
                filteredAssets.map((asset) => (
                  <tr 
                    key={asset.id} 
                    className="border-b border-gray-700/50 last:border-0 hover:bg-gray-700/20 transition-colors cursor-pointer"
                  >
                    <td className="px-6 py-4 font-mono text-sm text-gray-400">{asset.tag}</td>
                    <td className="px-6 py-4 font-medium text-gray-200">{asset.name}</td>
                    <td className="px-6 py-4">{getCategoryName(asset.category_id)}</td>
                    <td className="px-6 py-4">
                      <Badge status={asset.status} />
                    </td>
                    <td className="px-6 py-4 text-gray-400">{asset.location}</td>
                  </tr>
                ))
              )}
              {filteredAssets.length === 0 && !isLoading && (
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
          <FormInput 
            label="Asset Name" 
            placeholder="e.g. MacBook Pro 16" 
            value={regName}
            onChange={(e) => setRegName(e.target.value)}
            required 
          />
          <FormInput 
            label="Category" 
            as="select" 
            options={categories.map(c => ({ label: c.name, value: c.id.toString() }))}
            value={regCategoryId}
            onChange={(e) => setRegCategoryId(e.target.value)}
            required 
          />
          <FormInput 
            label="Location" 
            placeholder="e.g. HQ - Floor 2" 
            value={regLocation}
            onChange={(e) => setRegLocation(e.target.value)}
            required 
          />
          <FormInput 
            label="Acquisition Cost ($)" 
            type="number" 
            placeholder="2400.00" 
            value={regCost}
            onChange={(e) => setRegCost(e.target.value)}
          />
          <FormInput 
            label="Condition" 
            as="select"
            options={[
              { label: 'New', value: 'New' },
              { label: 'Good', value: 'Good' },
              { label: 'Fair', value: 'Fair' },
              { label: 'Poor', value: 'Poor' }
            ]}
            value={regCondition}
            onChange={(e) => setRegCondition(e.target.value)}
            required
          />
          <div className="flex gap-4 pt-2">
            <label className="flex items-center gap-2 text-sm text-gray-300">
              <input 
                type="checkbox" 
                checked={isSharable} 
                onChange={(e) => setIsSharable(e.target.checked)}
                className="rounded border-gray-700 bg-gray-900 text-orange-500 focus:ring-orange-500"
              />
              Is Sharable
            </label>
            <label className="flex items-center gap-2 text-sm text-gray-300">
              <input 
                type="checkbox" 
                checked={isBookable} 
                onChange={(e) => setIsBookable(e.target.checked)}
                className="rounded border-gray-700 bg-gray-900 text-orange-500 focus:ring-orange-500"
              />
              Is Bookable (Resource)
            </label>
          </div>
          <div className="pt-4">
            <Button type="submit">Submit Registration</Button>
          </div>
        </form>
      </Modal>
    </AppLayout>
  );
};