import React, { useState, useEffect } from 'react';
import { AppLayout } from '../components/layout/AppLayout';
export type MaintenanceStatus = 'Pending' | 'Approved' | 'Technician assigned' | 'In progress' | 'Resolved';

export interface MaintenanceTicket {
  id: string;
  tag: string;
  issue: string;
  subtext?: string;
  status: MaintenanceStatus;
}

export const kanbanColumns: MaintenanceStatus[] = [
  'Pending', 
  'Approved', 
  'Technician assigned', 
  'In progress', 
  'Resolved'
];
import { axiosClient } from '../api/axiosClient';
import { ENDPOINTS } from '../api/endpoints';
import { Modal } from '../components/common/Modal';
import { FormInput } from '../components/common/FormInput';
import { Button } from '../components/common/Button';

export const MaintenancePage = () => {
  const [tickets, setTickets] = useState<MaintenanceTicket[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [selectedTicket, setSelectedTicket] = useState<any | null>(null);
  
  // Workflow update state
  const [workflowStatus, setWorkflowStatus] = useState('APPROVED');
  const [workflowDesc, setWorkflowDesc] = useState('');

  const toDbStatus = (status: MaintenanceStatus): string => {
    switch (status) {
      case 'Pending': return 'PENDING';
      case 'Approved': return 'APPROVED';
      case 'Technician assigned': return 'TECHNICIAN_ASSIGNED';
      case 'In progress': return 'IN_PROGRESS';
      case 'Resolved': return 'RESOLVED';
      default: return 'PENDING';
    }
  };

  const fromDbStatus = (status: string): MaintenanceStatus => {
    switch (status) {
      case 'PENDING': return 'Pending';
      case 'APPROVED': return 'Approved';
      case 'TECHNICIAN_ASSIGNED': return 'Technician assigned';
      case 'IN_PROGRESS': return 'In progress';
      case 'RESOLVED': return 'Resolved';
      default: return 'Pending';
    }
  };

  const fetchMaintenanceTickets = async () => {
    setIsLoading(true);
    try {
      const [ticketsRes, assetsRes] = await Promise.all([
        axiosClient.get(ENDPOINTS.MAINTENANCE.LIST),
        axiosClient.get(ENDPOINTS.ASSETS.DIRECTORY)
      ]);
      const rawTickets = ticketsRes.data?.data || ticketsRes.data || [];
      const assetsList = assetsRes.data?.data || assetsRes.data || [];
      
      const mapped = rawTickets.map((t: any) => {
        const asset = assetsList.find((a: any) => a.id === t.asset_id);
        return {
          id: t.id.toString(),
          tag: asset ? asset.tag : `Asset #${t.asset_id}`,
          issue: t.description,
          subtext: `Priority: ${t.priority}`,
          status: fromDbStatus(t.status)
        };
      });
      setTickets(mapped);
    } catch (error) {
      console.error("Failed to fetch maintenance data", error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchMaintenanceTickets();
  }, []);

  const getTicketsByStatus = (status: MaintenanceStatus) => {
    return tickets.filter(ticket => ticket.status === status);
  };

  const handleCardClick = (ticket: MaintenanceTicket) => {
    setSelectedTicket(ticket);
    setWorkflowStatus(toDbStatus(ticket.status));
  };

  const handleWorkflowSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedTicket) return;
    try {
      await axiosClient.post(ENDPOINTS.MAINTENANCE.WORKFLOW(selectedTicket.id), {
        status: workflowStatus,
        description: workflowDesc
      });
      alert('Maintenance status advanced!');
      setSelectedTicket(null);
      setWorkflowDesc('');
      fetchMaintenanceTickets();
    } catch (err) {
      console.error(err);
      alert('Failed to advance maintenance workflow (Insufficient roles or server error).');
    }
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
                      onClick={() => handleCardClick(ticket)}
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
            Click on any card to advance the ticket's workflow status. Approving moves the asset to <span className="text-amber-500/80 border-b border-amber-500/30">under maintenance</span>, resolving returns it to <span className="text-emerald-500/80 border-b border-emerald-500/30">available</span>.
          </p>
        </div>

      </div>

      <Modal isOpen={!!selectedTicket} onClose={() => setSelectedTicket(null)} title="Update Maintenance Ticket">
        {selectedTicket && (
          <form onSubmit={handleWorkflowSubmit} className="space-y-4">
            <div>
              <p className="text-sm text-gray-400">Asset: <span className="font-mono text-gray-200">{selectedTicket.tag}</span></p>
              <p className="text-sm text-gray-400 mt-1">Issue: <span className="text-gray-200">{selectedTicket.issue}</span></p>
            </div>
            <FormInput 
              label="Workflow Action / Status" 
              as="select"
              options={[
                { label: 'Pending', value: 'PENDING' },
                { label: 'Approved', value: 'APPROVED' },
                { label: 'Assign Technician', value: 'TECHNICIAN_ASSIGNED' },
                { label: 'In Progress', value: 'IN_PROGRESS' },
                { label: 'Resolved', value: 'RESOLVED' }
              ]}
              value={workflowStatus}
              onChange={(e) => setWorkflowStatus(e.target.value)}
              required
            />
            <FormInput 
              label="Action Notes" 
              as="textarea"
              placeholder="e.g. Assigned technician, diagnostic results, parts ordered..."
              value={workflowDesc}
              onChange={(e) => setWorkflowDesc(e.target.value)}
            />
            <div className="pt-4">
              <Button type="submit">Update Status</Button>
            </div>
          </form>
        )}
      </Modal>
    </AppLayout>
  );
};