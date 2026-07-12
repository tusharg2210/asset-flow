export type MaintenanceStatus = 'Pending' | 'Approved' | 'Technician assigned' | 'In progress' | 'Resolved';

export interface MaintenanceTicket {
  id: string;
  tag: string;
  issue: string;
  subtext?: string;
  status: MaintenanceStatus;
}

export const mockMaintenanceTickets: MaintenanceTicket[] = [
  { id: '1', tag: 'AF-0062', issue: 'Projector bulb not turning on', status: 'Pending' },
  { id: '2', tag: 'AF-003', issue: 'ac unit noisy compresor', status: 'Approved' },
  { id: '3', tag: 'AF-0078', issue: 'forlift', subtext: 'tech: R varma', status: 'Technician assigned' },
  { id: '4', tag: 'AF-897', issue: 'Printer Jam', subtext: 'parts ordered', status: 'In progress' },
  { id: '5', tag: 'AF-873', issue: 'Chair repair', subtext: 'resolved 7 Jul', status: 'Resolved' },
];

export const kanbanColumns: MaintenanceStatus[] = [
  'Pending', 
  'Approved', 
  'Technician assigned', 
  'In progress', 
  'Resolved'
];