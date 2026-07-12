import React, { useState } from 'react';
import { AppLayout } from '../components/layout/AppLayout';
import { mockLogs,type FilterTab } from '../mockData/notifications';
// import { axiosClient } from '../api/axiosClient';
// import { ENDPOINTS } from '../api/endpoints';

export const NotificationsPage = () => {
  const [activeTab, setActiveTab] = useState<FilterTab>('All');
  const [logs, setLogs] = useState(mockLogs);

  /*
  // TODO: Uncomment when backend is ready
  useEffect(() => {
    const fetchNotifications = async () => {
      try {
        const { data } = await axiosClient.get(ENDPOINTS.NOTIFICATIONS.ALL);
        // Map backend data to local state
        // setLogs(data);
      } catch (error) {
        console.error('Failed to fetch notifications', error);
      }
    };
    fetchNotifications();
  }, []);
  */

  const tabs: FilterTab[] = ['All', 'Alerts', 'Approvals', 'Bookings'];

  // Filter logic based on the selected tab
  const filteredLogs = logs.filter(log => {
    if (activeTab === 'All') return true;
    return log.category === activeTab;
  });

  return (
    <AppLayout>
      <div className="p-8 max-w-4xl mx-auto">
        
        {/* Filter Tabs */}
        <div className="flex gap-3 mb-8 border-b border-gray-800 pb-6">
          {tabs.map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`px-5 py-2 rounded-lg font-medium text-sm transition-all duration-200 border ${
                activeTab === tab
                  ? 'border-emerald-500/50 text-emerald-400 bg-emerald-500/10 shadow-sm'
                  : 'border-gray-700 text-gray-400 hover:border-gray-500 hover:text-gray-200 bg-gray-900/50'
              }`}
            >
              {tab}
            </button>
          ))}
        </div>

        {/* Logs List */}
        <div className="flex flex-col">
          {filteredLogs.length > 0 ? (
            filteredLogs.map((log) => (
              <div 
                key={log.id} 
                className="flex items-center justify-between py-4 border-b border-gray-800/80 group hover:bg-gray-800/20 px-2 rounded-lg transition-colors cursor-default"
              >
                <div className="flex items-center gap-4">
                  {/* Status Square Indicator */}
                  <div className={`w-2.5 h-2.5 rounded-sm flex-shrink-0 ${log.indicatorClass}`} />
                  
                  {/* Log Text */}
                  <span className="text-gray-300 text-sm font-medium group-hover:text-gray-200 transition-colors">
                    {log.text}
                  </span>
                </div>
                
                {/* Timestamp */}
                <span className="text-gray-500 text-sm whitespace-nowrap ml-4 font-mono">
                  {log.time}
                </span>
              </div>
            ))
          ) : (
            <div className="py-12 text-center text-gray-500">
              No {activeTab.toLowerCase()} to display.
            </div>
          )}
        </div>

      </div>
    </AppLayout>
  );
};