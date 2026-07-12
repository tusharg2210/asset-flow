import React, { useState, useEffect } from 'react';
import { AppLayout } from '../components/layout/AppLayout';
import { Badge } from '../components/common/Badge';
import { mockDepartments } from '../mockData/organization';
import { Plus } from 'lucide-react';
import { Modal } from '../components/common/Modal';
import { FormInput } from '../components/common/FormInput';
import { Button } from '../components/common/Button';
// import { axiosClient } from '../../api/axiosClient';
// import { ENDPOINTS } from '../../api/endpoints';

type TabType = 'Departments' | 'Categories' | 'Employee';

export const OrganizationPage = () => {
  const [activeTab, setActiveTab] = useState<TabType>('Departments');
  const [departments, setDepartments] = useState(mockDepartments);
  const [isLoading, setIsLoading] = useState(false);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);

  /*
  // TODO: Uncomment when backend is ready
  useEffect(() => {
    const fetchTabContent = async () => {
      setIsLoading(true);
      try {
        if (activeTab === 'Departments') {
          const { data } = await axiosClient.get(ENDPOINTS.ORGANIZATION.DEPARTMENTS);
          setDepartments(data);
        }
        // Handle Categories and Employee fetching here based on activeTab
      } catch (error) {
        console.error(`Failed to fetch ${activeTab}`, error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchTabContent();
  }, [activeTab]);
  */

  const tabs: TabType[] = ['Departments', 'Categories', 'Employee'];

  const handleAddSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    alert(`Added new ${activeTab.slice(0, -1).toLowerCase()} successfully!`);
    setIsAddModalOpen(false);
  };

  return (
    <AppLayout>
      <div className="p-8 max-w-7xl mx-auto space-y-6">
        
        {/* Top Controls: Tabs & Add Button */}
        <div className="flex justify-between items-center">
          <div className="flex gap-3">
            {tabs.map((tab) => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={`px-5 py-2 rounded-lg font-medium transition-all duration-200 border ${
                  activeTab === tab
                    ? 'border-orange-500 text-orange-400 bg-orange-500/10 shadow-sm'
                    : 'border-gray-700 text-gray-400 hover:border-gray-500 hover:text-gray-200 bg-gray-800'
                }`}
              >
                {tab}
              </button>
            ))}
          </div>
          
          <button 
            onClick={() => setIsAddModalOpen(true)}
            className="flex items-center gap-2 px-5 py-2 rounded-lg font-medium transition-all duration-200 border border-gray-700 bg-gray-800 text-gray-200 hover:border-orange-500 hover:text-orange-400 shadow-sm"
          >
            <Plus size={18} />
            Add
          </button>
        </div>

        {/* Data Table */}
        <div className="bg-gray-800 border border-gray-700 rounded-xl overflow-hidden shadow-sm mt-6">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-gray-900/50 border-b border-gray-700 text-gray-400 text-sm">
                <th className="px-6 py-4 font-medium">Department</th>
                <th className="px-6 py-4 font-medium">Head</th>
                <th className="px-6 py-4 font-medium">Parent Dept</th>
                <th className="px-6 py-4 font-medium">Status</th>
              </tr>
            </thead>
            <tbody className="text-gray-300">
              {isLoading ? (
                <tr>
                  <td colSpan={4} className="px-6 py-8 text-center text-gray-500">
                    Loading data...
                  </td>
                </tr>
              ) : (
                departments.map((dept) => (
                  <tr 
                    key={dept.id} 
                    className="border-b border-gray-700/50 last:border-0 hover:bg-gray-700/20 transition-colors"
                  >
                    <td className="px-6 py-4">{dept.name}</td>
                    <td className="px-6 py-4 capitalize">{dept.head}</td>
                    <td className="px-6 py-4 text-gray-500">{dept.parent}</td>
                    <td className="px-6 py-4">
                      <Badge status={dept.status} />
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Contextual Helper Text */}
        <div className="pt-4 border-t border-gray-800">
          <p className="text-sm text-gray-500">
            Editing a department here also drives the picklist in Asset Registration and Allocations.
          </p>
        </div>

      </div>

      <Modal 
        isOpen={isAddModalOpen} 
        onClose={() => setIsAddModalOpen(false)} 
        title={`Add New ${activeTab === 'Categories' ? 'Category' : activeTab.slice(0, -1)}`}
      >
        <form onSubmit={handleAddSubmit} className="space-y-4">
          {activeTab === 'Departments' && (
            <>
              <FormInput label="Department Name" placeholder="e.g. Marketing" required />
              <FormInput label="Head of Department (Employee ID)" placeholder="e.g. 104" required />
              <FormInput label="Parent Department (Optional)" placeholder="e.g. Operations" />
            </>
          )}

          {activeTab === 'Categories' && (
            <>
              <FormInput label="Category Name" placeholder="e.g. Furniture" required />
              <FormInput label="Description" as="textarea" placeholder="Optional details..." />
            </>
          )}

          {activeTab === 'Employee' && (
            <>
              <FormInput label="Full Name" placeholder="e.g. Jane Smith" required />
              <FormInput label="Email Address" type="email" placeholder="jane@company.com" required />
              <FormInput label="Department ID" placeholder="e.g. 4" required />
            </>
          )}

          <div className="pt-4">
            <Button type="submit">Create Record</Button>
          </div>
        </form>
      </Modal>
    </AppLayout>
  );
};