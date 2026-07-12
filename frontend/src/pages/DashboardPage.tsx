import { useEffect, useState } from 'react';
import { AppLayout } from '../components/layout/AppLayout';

import { AlertCircle, Plus, CalendarPlus, Wrench } from 'lucide-react';
import { Modal } from '../components/common/Modal';
import { FormInput } from '../components/common/FormInput';
import { Button } from '../components/common/Button';
import { axiosClient } from '../api/axiosClient';
import { ENDPOINTS } from '../api/endpoints';

export const DashboardPage = () => {
  const [metrics, setMetrics] = useState({
    assetsAvailable: 0,
    assetsAllocated: 0,
    maintenanceToday: 0,
    activeBookings: 0,
    pendingTransfers: 0,
    upcomingReturns: 0
  });
  const [overdueCount, setOverdueCount] = useState(0);
  const [activities, setActivities] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  
  // Dynamic lists from backend
  const [categories, setCategories] = useState<{ id: number; name: string }[]>([]);
  const [assets, setAssets] = useState<any[]>([]);

  // Modal State
  const [activeModal, setActiveModal] = useState<'register' | 'book' | 'request' | null>(null);

  // Form States
  // Register Asset
  const [regName, setRegName] = useState('');
  const [regCategoryId, setRegCategoryId] = useState('');
  const [regLocation, setRegLocation] = useState('');
  const [regCost, setRegCost] = useState('');
  const [regCondition, setRegCondition] = useState('New');
  const [isSharable, setIsSharable] = useState(true);
  const [isBookable, setIsBookable] = useState(false);

  // Book Asset
  const [bookAssetId, setBookAssetId] = useState('');
  const [bookStart, setBookStart] = useState('');
  const [bookEnd, setBookEnd] = useState('');

  // Maintenance Request
  const [maintAssetTag, setMaintAssetTag] = useState('');
  const [maintPriority, setMaintPriority] = useState('MEDIUM');
  const [maintDescription, setMaintDescription] = useState('');

  useEffect(() => {
    const fetchDashboardData = async () => {
      setIsLoading(true);
      try {
        const [metricsRes, alertsRes, logsRes, categoriesRes, assetsRes] = await Promise.all([
          axiosClient.get(ENDPOINTS.DASHBOARD.METRICS),
          axiosClient.get(ENDPOINTS.DASHBOARD.ALERTS),
          axiosClient.get(ENDPOINTS.LOGS.RECENT),
          axiosClient.get(ENDPOINTS.ORGANIZATION.ASSET_CATEGORIES),
          axiosClient.get(ENDPOINTS.ASSETS.DIRECTORY)
        ]);
        
        setMetrics(metricsRes.data);
        setOverdueCount(alertsRes.data.overdueAllocations?.length || 0);
        setCategories(categoriesRes.data || []);
        
        const mappedLogs = (logsRes.data || []).map((log: any, idx: number) => ({
          id: idx,
          text: `${log.user || 'System'} - ${log.action}`,
          time: log.timestamp ? new Date(log.timestamp).toLocaleTimeString() : 'Recent'
        }));
        setActivities(mappedLogs);

        const assetsList = assetsRes.data?.data || assetsRes.data || [];
        setAssets(assetsList);
      } catch (error) {
        console.error("Failed to fetch dashboard data", error);
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchDashboardData();
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
      };
      await axiosClient.post(ENDPOINTS.ASSETS.REGISTER, payload);
      alert('Asset registered successfully!');
      setActiveModal(null);
      // Reset form
      setRegName('');
      setRegCategoryId('');
      setRegLocation('');
      setRegCost('');
      // Reload
      window.location.reload();
    } catch (error) {
      console.error('Failed to register asset:', error);
      alert('Failed to register asset.');
    }
  };

  const handleBookSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload = {
        asset_id: parseInt(bookAssetId, 10),
        start_time: new Date(bookStart).toISOString(),
        end_time: new Date(bookEnd).toISOString()
      };
      await axiosClient.post(ENDPOINTS.BOOKINGS.CREATE, payload);
      alert('Booking confirmed!');
      setActiveModal(null);
      setBookAssetId('');
      setBookStart('');
      setBookEnd('');
      window.location.reload();
    } catch (error) {
      console.error('Failed to book resource:', error);
      alert('Failed to book resource (overlap conflict or server error).');
    }
  };

  const handleMaintenanceSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const assetsRes = await axiosClient.get(`${ENDPOINTS.ASSETS.DIRECTORY}?tag=${maintAssetTag}`);
      const matchedAssets = assetsRes.data?.data || [];
      if (matchedAssets.length === 0) {
        alert(`Asset with tag ${maintAssetTag} not found.`);
        return;
      }
      const asset = matchedAssets[0];
      const payload = {
        asset_id: asset.id,
        priority: maintPriority,
        description: maintDescription
      };
      await axiosClient.post(ENDPOINTS.MAINTENANCE.CREATE, payload);
      alert('Maintenance request submitted!');
      setActiveModal(null);
      setMaintAssetTag('');
      setMaintDescription('');
      window.location.reload();
    } catch (error) {
      console.error('Failed to raise maintenance request:', error);
      alert('Failed to raise maintenance request.');
    }
  };

  const KpiCard = ({ title, value }: { title: string; value: number | string }) => (
    <div className="bg-gray-800 border border-gray-700 rounded-xl p-5 shadow-sm hover:border-gray-600 transition-colors">
      <h3 className="text-gray-400 text-sm font-medium mb-2">{title}</h3>
      <p className="text-3xl font-bold text-gray-50">{value}</p>
    </div>
  );

  return (
    <AppLayout>
      <div className="p-8 max-w-7xl mx-auto space-y-8">
        
        {/* Header Section */}
        <div>
          <h2 className="text-2xl font-semibold text-gray-100 mb-6">Today's Overview</h2>
          
          {/* KPI Grid */}
          {isLoading ? (
            <div className="flex justify-center items-center h-32">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500"></div>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
              <KpiCard title="Available" value={metrics.assetsAvailable} />
              <KpiCard title="Allocated" value={metrics.assetsAllocated} />
              <KpiCard title="Maintenance Today" value={metrics.maintenanceToday} />
              <KpiCard title="Active Bookings" value={metrics.activeBookings} />
              <KpiCard title="Pending Transfers" value={metrics.pendingTransfers} />
              <KpiCard title="Upcoming returns" value={metrics.upcomingReturns} />
            </div>
          )}

          {/* Overdue Alert Banner */}
          {overdueCount > 0 && (
            <div className="flex items-center gap-3 bg-red-900/20 border border-red-500/30 rounded-lg p-4 mb-6">
              <AlertCircle className="text-red-400" size={20} />
              <p className="text-red-400 font-medium">
                {overdueCount} assets overdue for return - flagged for follow-up
              </p>
            </div>
          )}

          {/* Quick Actions */}
          <div className="flex gap-4">
            <button 
              onClick={() => setActiveModal('register')}
              className="flex items-center gap-2 bg-gray-800 border-2 border-gray-700 hover:border-orange-500 text-gray-200 px-5 py-2.5 rounded-lg font-medium transition-all shadow-sm"
            >
              <Plus size={18} className="text-orange-500" />
              register asset
            </button>
            <button 
              onClick={() => setActiveModal('book')}
              className="flex items-center gap-2 bg-gray-800 border border-gray-700 hover:bg-gray-700 text-gray-200 px-5 py-2.5 rounded-lg font-medium transition-colors shadow-sm"
            >
              <CalendarPlus size={18} className="text-gray-400" />
              Book resource
            </button>
            <button 
              onClick={() => setActiveModal('request')}
              className="flex items-center gap-2 bg-gray-800 border border-gray-700 hover:bg-gray-700 text-gray-200 px-5 py-2.5 rounded-lg font-medium transition-colors shadow-sm"
            >
              <Wrench size={18} className="text-gray-400" />
              Raise requests
            </button>
          </div>
        </div>

        {/* Recent Activity Section */}
        <div className="pt-4">
          <h2 className="text-xl font-semibold text-gray-100 mb-4">Recent Activity</h2>
          {isLoading ? (
            <p className="text-gray-500">Loading activity logs...</p>
          ) : (
            <div className="space-y-3">
              {activities.length === 0 ? (
                <p className="text-gray-500 text-sm">No recent activities.</p>
              ) : (
                activities.map((activity) => (
                  <div 
                    key={activity.id} 
                    className="flex items-center text-gray-300 bg-gray-800/50 p-3 rounded-lg border border-gray-800"
                  >
                    <div className="w-2 h-2 rounded-full bg-gray-500 mr-4"></div>
                    <p className="flex-1 text-sm">{activity.text}</p>
                    <span className="text-xs text-gray-500">{activity.time}</span>
                  </div>
                ))
              )}
            </div>
          )}
        </div>

      </div>

      {/* Modals */}
      <Modal isOpen={activeModal === 'register'} onClose={() => setActiveModal(null)} title="Register New Asset">
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

      <Modal isOpen={activeModal === 'book'} onClose={() => setActiveModal(null)} title="Book Shared Resource">
        <form onSubmit={handleBookSubmit} className="space-y-4">
          <FormInput 
            label="Select Resource" 
            as="select" 
            options={assets.filter(a => a.is_bookable).map(a => ({ label: `${a.tag} - ${a.name}`, value: a.id.toString() }))}
            value={bookAssetId}
            onChange={(e) => setBookAssetId(e.target.value)}
            required 
          />
          <div className="grid grid-cols-2 gap-4">
            <FormInput 
              label="Start Time" 
              type="datetime-local" 
              value={bookStart}
              onChange={(e) => setBookStart(e.target.value)}
              required 
            />
            <FormInput 
              label="End Time" 
              type="datetime-local" 
              value={bookEnd}
              onChange={(e) => setBookEnd(e.target.value)}
              required 
            />
          </div>
          <div className="pt-4">
            <Button type="submit">Confirm Booking</Button>
          </div>
        </form>
      </Modal>

      <Modal isOpen={activeModal === 'request'} onClose={() => setActiveModal(null)} title="Raise Maintenance Request">
        <form onSubmit={handleMaintenanceSubmit} className="space-y-4">
          <FormInput 
            label="Asset Tag" 
            placeholder="e.g. AF-0115" 
            value={maintAssetTag}
            onChange={(e) => setMaintAssetTag(e.target.value)}
            required 
          />
          <FormInput 
            label="Priority" 
            as="select" 
            options={[
              { label: 'Low', value: 'LOW' },
              { label: 'Medium', value: 'MEDIUM' },
              { label: 'High', value: 'HIGH' },
              { label: 'Critical', value: 'CRITICAL' }
            ]} 
            value={maintPriority}
            onChange={(e) => setMaintPriority(e.target.value)}
            required 
          />
          <FormInput 
            label="Issue Description" 
            as="textarea" 
            placeholder="Describe the problem..." 
            value={maintDescription}
            onChange={(e) => setMaintDescription(e.target.value)}
            required 
          />
          <div className="pt-4">
            <Button type="submit">Submit Request</Button>
          </div>
        </form>
      </Modal>
    </AppLayout>
  );
};