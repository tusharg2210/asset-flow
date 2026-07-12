import React, { useState, useEffect } from 'react';
import { AppLayout } from '../components/layout/AppLayout';
import { mockMaintenanceTickets, kanbanColumns,type MaintenanceTicket,type MaintenanceStatus } from '../mockData/maintenance';
// import { axiosClient } from '../api/axiosClient';
// import { ENDPOINTS } from '../api/endpoints';

export const MaintenancePage = () => {
  const [tickets, setTickets] = useState<MaintenanceTicket[]>(mockMaintenanceTickets);
  const [isLoading, setIsLoading] = useState(false);

  /*
  // TODO: Uncomment when backend is ready
  useEffect(() => {
    const fetchMaintenanceTickets = async () => {
      setIsLoading(true);
      try {
        const { data } = await axiosClient.get(ENDPOINTS.MAINTENANCE.LIST);
        setTickets(data);
      } catch (error) {
        console.error("Failed to fetch maintenance data", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchMaintenanceTickets();
  }, []);
  */

  const getTicketsByStatus = (status: MaintenanceStatus) => {
    return tickets.filter(ticket => ticket.status === status);
  };

  return (
    <AppLayout>
      <div className="p-8 h-full flex flex-col max-w-7xl mx-auto">
        
        {/* Kanban Board Container */}
        <div className="flex-1 grid grid-cols-5 gap-0 border border-gray-800 rounded-xl bg-gray-900/50 overflow-hidden shadow-sm">
          
          {kanbanColumns.map((columnStatus, index) => (
            <div 
              key={columnStatus} 
              className={`flex flex-col border-gray-800 ${
                index !== kanbanColumns.length - 1 ? 'border-r' : ''
              }`}
            >
              {/* Column Header */}
              <div className="px-4 py-3 border-b border-gray-800 bg-gray-900/80">
                <h3 className="text-sm font-medium text-gray-300 capitalize text-center">
                  {columnStatus}
                </h3>
              </div>

              {/* Column Body / Dropzone */}
              <div className="flex-1 p-3 space-y-3 overflow-y-auto min-h-[400px]">
                {isLoading ? (
                  <div className="animate-pulse flex flex-col gap-3">
                    <div className="h-24 bg-gray-800 rounded-lg w-full"></div>
                  </div>
                ) : (
                  getTicketsByStatus(columnStatus).map(ticket => (
                    <div 
                      key={ticket.id}
                      className={`p-3 rounded-lg border shadow-sm transition-colors cursor-pointer hover:shadow-md ${
                        ticket.status === 'Resolved'
                          ? 'bg-emerald-900/20 border-emerald-500/50 hover:border-emerald-500/80'
                          : 'bg-gray-800 border-gray-700 hover:border-gray-500'
                      }`}
                    >
                      <div className={`font-mono text-xs mb-1.5 ${
                        ticket.status === 'Resolved' ? 'text-emerald-400' : 'text-gray-400'
                      }`}>
                        {ticket.tag}
                      </div>
                      <div className={`text-sm font-medium ${
                        ticket.status === 'Resolved' ? 'text-emerald-300' : 'text-gray-200'
                      }`}>
                        {ticket.issue}
                      </div>
                      {ticket.subtext && (
                        <div className={`text-xs mt-2 ${
                          ticket.status === 'Resolved' ? 'text-emerald-500/80' : 'text-gray-500'
                        }`}>
                          {ticket.subtext}
                        </div>
                      )}
                    </div>
                  ))
                )}
              </div>
            </div>
          ))}

        </div>

        {/* Footer Helper Text */}
        <div className="mt-6 pt-4 border-t border-gray-800">
          <p className="text-sm text-gray-400 text-center font-medium">
            Approving a card moves the asset to <span className="text-amber-500/80 border-b border-amber-500/30">under maintenance</span>, resolving returns it to <span className="text-emerald-500/80 border-b border-emerald-500/30">available</span>.
          </p>
        </div>

      </div>
    </AppLayout>
  );
};