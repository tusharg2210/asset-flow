export type FilterTab = 'All' | 'Alerts' | 'Approvals' | 'Bookings';

export interface AppLog {
  id: string;
  text: string;
  time: string;
  category: FilterTab;
  indicatorClass: string;
}

export const mockLogs: AppLog[] = [
  { 
    id: '1', 
    text: 'Laptop AF-0014 assigned to Priya shah', 
    time: '2m ago', 
    category: 'All', // General
    indicatorClass: 'bg-blue-500/80' 
  },
  { 
    id: '2', 
    text: 'Maintenance request AF-0055 approved', 
    time: '18m ago', 
    category: 'Approvals', 
    indicatorClass: 'border-2 border-emerald-500/80 bg-transparent' 
  },
  { 
    id: '3', 
    text: 'Booking confirmed : Room B2 : 2:00 to 3:00 PM', 
    time: '1h ago', 
    category: 'Bookings', 
    indicatorClass: 'bg-blue-500/80' 
  },
  { 
    id: '4', 
    text: 'Transfer approved : AF-0033 to facilities dept', 
    time: '3h ago', 
    category: 'Approvals', 
    indicatorClass: 'bg-rose-500/80' 
  },
  { 
    id: '5', 
    text: 'Overdue return : AF-0021 was due 3 days ago', 
    time: '1d ago', 
    category: 'Alerts', 
    indicatorClass: 'bg-amber-600/80' 
  },
  { 
    id: '6', 
    text: 'audit discrepancy flagged : AF-0088 damaged', 
    time: '2d ago', 
    category: 'Alerts', 
    indicatorClass: 'bg-rose-500/80' 
  },
];