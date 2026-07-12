import { useState, useEffect } from 'react';
import { AppLayout } from '../components/layout/AppLayout';
import { Badge } from '../components/common/Badge';
import { Plus } from 'lucide-react';
import { Modal } from '../components/common/Modal';
import { FormInput } from '../components/common/FormInput';
import { Button } from '../components/common/Button';
import { axiosClient } from '../api/axiosClient';
import { ENDPOINTS } from '../api/endpoints';

type TabType = 'Departments' | 'Categories' | 'Employee';

export const OrganizationPage = () => {
  const [activeTab, setActiveTab] = useState<TabType>('Departments');
  const [departments, setDepartments] = useState<any[]>([]);
  const [categories, setCategories] = useState<any[]>([]);
  const [employees, setEmployees] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);

  // Form States
  // Department Form
  const [deptName, setDeptName] = useState('');
  const [deptHeadId, setDeptHeadId] = useState('');
  const [deptParentId, setDeptParentId] = useState('');

  // Category Form
  const [catName, setCatName] = useState('');

  // Employee Form
  const [empName, setEmpName] = useState('');
  const [empEmail, setEmpEmail] = useState('');
  const [empPassword, setEmpPassword] = useState('SecurePassword123!');

  const fetchTabContent = async () => {
    setIsLoading(true);
    try {
      if (activeTab === 'Departments') {
        const { data } = await axiosClient.get(ENDPOINTS.ORGANIZATION.DEPARTMENTS);
        setDepartments(data || []);
      } else if (activeTab === 'Categories') {
        const { data } = await axiosClient.get(ENDPOINTS.ORGANIZATION.ASSET_CATEGORIES);
        setCategories(data || []);
      } else if (activeTab === 'Employee') {
        const { data } = await axiosClient.get(ENDPOINTS.ORGANIZATION.USERS);
        setEmployees(data || []);
      }
    } catch (error) {
      console.error(`Failed to fetch ${activeTab}`, error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchTabContent();
  }, [activeTab]);

  const handleAddSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (activeTab === 'Departments') {
        await axiosClient.post(ENDPOINTS.ORGANIZATION.DEPARTMENTS, {
          name: deptName,
          head_id: deptHeadId ? parseInt(deptHeadId, 10) : null,
          parent_department_id: deptParentId ? parseInt(deptParentId, 10) : null
        });
        setDeptName('');
        setDeptHeadId('');
        setDeptParentId('');
      } else if (activeTab === 'Categories') {
        await axiosClient.post(ENDPOINTS.ORGANIZATION.ASSET_CATEGORIES, {
          name: catName
        });
        setCatName('');
      } else if (activeTab === 'Employee') {
        await axiosClient.post(ENDPOINTS.AUTH.SIGNUP, {
          name: empName,
          email: empEmail,
          password: empPassword
        });
        setEmpName('');
        setEmpEmail('');
      }
      alert(`Added new ${activeTab.slice(0, -1).toLowerCase()} successfully!`);
      setIsAddModalOpen(false);
      fetchTabContent();
    } catch (error) {
      console.error('Failed to create record:', error);
      alert('Failed to create record.');
    }
  };

  const tabs: TabType[] = ['Departments', 'Categories', 'Employee'];



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
              {activeTab === 'Departments' && (
                <tr className="bg-gray-900/50 border-b border-gray-700 text-gray-400 text-sm">
                  <th className="px-6 py-4 font-medium">Department Name</th>
                  <th className="px-6 py-4 font-medium">Head of Department ID</th>
                  <th className="px-6 py-4 font-medium">Parent Dept ID</th>
                  <th className="px-6 py-4 font-medium">Status</th>
                </tr>
              )}
              {activeTab === 'Categories' && (
                <tr className="bg-gray-900/50 border-b border-gray-700 text-gray-400 text-sm">
                  <th className="px-6 py-4 font-medium">Category Name</th>
                  <th className="px-6 py-4 font-medium">Category ID</th>
                  <th className="px-6 py-4 font-medium">Created At</th>
                </tr>
              )}
              {activeTab === 'Employee' && (
                <tr className="bg-gray-900/50 border-b border-gray-700 text-gray-400 text-sm">
                  <th className="px-6 py-4 font-medium">Name</th>
                  <th className="px-6 py-4 font-medium">Email</th>
                  <th className="px-6 py-4 font-medium">Role</th>
                  <th className="px-6 py-4 font-medium">Status</th>
                </tr>
              )}
            </thead>
            <tbody className="text-gray-300">
              {isLoading ? (
                <tr>
                  <td colSpan={4} className="px-6 py-8 text-center text-gray-500">
                    Loading data...
                  </td>
                </tr>
              ) : activeTab === 'Departments' ? (
                departments.length === 0 ? (
                  <tr>
                    <td colSpan={4} className="px-6 py-8 text-center text-gray-500">No departments setup.</td>
                  </tr>
                ) : (
                  departments.map((dept) => (
                    <tr 
                      key={dept.id} 
                      className="border-b border-gray-700/50 last:border-0 hover:bg-gray-700/20 transition-colors"
                    >
                      <td className="px-6 py-4 font-medium text-gray-200">{dept.name}</td>
                      <td className="px-6 py-4 font-mono text-sm text-gray-400">{dept.head_id || '-'}</td>
                      <td className="px-6 py-4 font-mono text-sm text-gray-400">{dept.parent_department_id || '-'}</td>
                      <td className="px-6 py-4">
                        <Badge status={dept.status} />
                      </td>
                    </tr>
                  ))
                )
              ) : activeTab === 'Categories' ? (
                categories.length === 0 ? (
                  <tr>
                    <td colSpan={3} className="px-6 py-8 text-center text-gray-500">No categories setup.</td>
                  </tr>
                ) : (
                  categories.map((cat) => (
                    <tr 
                      key={cat.id} 
                      className="border-b border-gray-700/50 last:border-0 hover:bg-gray-700/20 transition-colors"
                    >
                      <td className="px-6 py-4 font-medium text-gray-200">{cat.name}</td>
                      <td className="px-6 py-4 font-mono text-sm text-gray-400">{cat.id}</td>
                      <td className="px-6 py-4 text-gray-400">{cat.created_at ? new Date(cat.created_at).toLocaleDateString() : 'N/A'}</td>
                    </tr>
                  ))
                )
              ) : (
                employees.length === 0 ? (
                  <tr>
                    <td colSpan={4} className="px-6 py-8 text-center text-gray-500">No employees directory setup.</td>
                  </tr>
                ) : (
                  employees.map((emp) => (
                    <tr 
                      key={emp.id} 
                      className="border-b border-gray-700/50 last:border-0 hover:bg-gray-700/20 transition-colors"
                    >
                      <td className="px-6 py-4 font-medium text-gray-200">{emp.name}</td>
                      <td className="px-6 py-4">{emp.email}</td>
                      <td className="px-6 py-4 text-orange-400 font-medium text-xs tracking-wider">{emp.role}</td>
                      <td className="px-6 py-4">
                        <Badge status={emp.status} />
                      </td>
                    </tr>
                  ))
                )
              )}
            </tbody>
          </table>
        </div>

        {/* Contextual Helper Text */}
        <div className="pt-4 border-t border-gray-800">
          <p className="text-sm text-gray-500">
            Adding departments, categories, and employees here feeds setup parameters dynamically to the rest of the application.
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
              <FormInput 
                label="Department Name" 
                placeholder="e.g. Marketing" 
                value={deptName}
                onChange={(e) => setDeptName(e.target.value)}
                required 
              />
              <FormInput 
                label="Head of Department (Employee User ID)" 
                placeholder="e.g. 4" 
                value={deptHeadId}
                onChange={(e) => setDeptHeadId(e.target.value)}
              />
              <FormInput 
                label="Parent Department ID (Optional)" 
                placeholder="e.g. 1" 
                value={deptParentId}
                onChange={(e) => setDeptParentId(e.target.value)}
              />
            </>
          )}

          {activeTab === 'Categories' && (
            <>
              <FormInput 
                label="Category Name" 
                placeholder="e.g. Furniture" 
                value={catName}
                onChange={(e) => setCatName(e.target.value)}
                required 
              />
            </>
          )}

          {activeTab === 'Employee' && (
            <>
              <FormInput 
                label="Full Name" 
                placeholder="e.g. Jane Smith" 
                value={empName}
                onChange={(e) => setEmpName(e.target.value)}
                required 
              />
              <FormInput 
                label="Email Address" 
                type="email" 
                placeholder="jane@company.com" 
                value={empEmail}
                onChange={(e) => setEmpEmail(e.target.value)}
                required 
              />
              <FormInput 
                label="Initial Password" 
                type="password" 
                placeholder="SecurePassword123!" 
                value={empPassword}
                onChange={(e) => setEmpPassword(e.target.value)}
                required 
              />
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