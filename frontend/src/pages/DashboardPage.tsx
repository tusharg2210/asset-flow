import { useEffect, useState } from 'react';
import { AppLayout } from '../components/layout/AppLayout';
import { mockDashboardData } from '../mockData/dashboard.ts';
import { AlertCircle, Plus, CalendarPlus, Wrench } from 'lucide-react';
// import { axiosClient } from '../api/axiosClient';
// import { ENDPOINTS } from '../api/endpoints';

export const DashboardPage = () => {
  const [metrics, setMetrics] = useState(mockDashboardData.metrics);
  const [overdueCount, setOverdueCount] = useState(mockDashboardData.alerts.overdueAllocations.length);
  const [activities, setActivities] = useState(mockDashboardData.recentActivity);
  const [isLoading, setIsLoading] = useState(false);

  /* 
  // TODO: Uncomment when backend is ready
  useEffect(() => {
    const fetchDashboardData = async () => {
      setIsLoading(true);
      try {
        const [metricsRes, alertsRes, logsRes] = await Promise.all([
          axiosClient.get(ENDPOINTS.DASHBOARD.METRICS),
          axiosClient.get(ENDPOINTS.DASHBOARD.ALERTS),
          axiosClient.get(ENDPOINTS.LOGS.RECENT) // Assuming this returns recent logs
        ]);
        
        setMetrics(metricsRes.data);
        setOverdueCount(alertsRes.data.overdueAllocations?.length || 0);
        // Map backend logs to frontend activity format here
      } catch (error) {
        console.error("Failed to fetch dashboard data", error);
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchDashboardData();
  }, []);
  */

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
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
            <KpiCard title="Available" value={metrics.assetsAvailable} />
            <KpiCard title="Allocated" value={metrics.assetsAllocated} />
            <KpiCard title="Maintenance Today" value={metrics.maintenanceToday} />
            <KpiCard title="Active Bookings" value={metrics.activeBookings} />
            <KpiCard title="Pending Transfers" value={metrics.pendingTransfers} />
            <KpiCard title="Upcoming returns" value={metrics.upcomingReturns} />
          </div>

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
            <button className="flex items-center gap-2 bg-gray-800 border-2 border-gray-700 hover:border-orange-500 text-gray-200 px-5 py-2.5 rounded-lg font-medium transition-all shadow-sm">
              <Plus size={18} className="text-orange-500" />
              register asset
            </button>
            <button className="flex items-center gap-2 bg-gray-800 border border-gray-700 hover:bg-gray-700 text-gray-200 px-5 py-2.5 rounded-lg font-medium transition-colors shadow-sm">
              <CalendarPlus size={18} className="text-gray-400" />
              Book resource
            </button>
            <button className="flex items-center gap-2 bg-gray-800 border border-gray-700 hover:bg-gray-700 text-gray-200 px-5 py-2.5 rounded-lg font-medium transition-colors shadow-sm">
              <Wrench size={18} className="text-gray-400" />
              Raise requests
            </button>
          </div>
        </div>

        {/* Recent Activity Section */}
        <div className="pt-4">
          <h2 className="text-xl font-semibold text-gray-100 mb-4">Recent Activity</h2>
          <div className="space-y-3">
            {activities.map((activity) => (
              <div 
                key={activity.id} 
                className="flex items-center text-gray-300 bg-gray-800/50 p-3 rounded-lg border border-gray-800"
              >
                <div className="w-2 h-2 rounded-full bg-gray-500 mr-4"></div>
                <p className="flex-1 text-sm">{activity.text}</p>
                <span className="text-xs text-gray-500">{activity.time}</span>
              </div>
            ))}
          </div>
        </div>

      </div>
    </AppLayout>
  );
};