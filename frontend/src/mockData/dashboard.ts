export const mockDashboardData = {
  metrics: {
    assetsAvailable: 128,
    assetsAllocated: 76,
    maintenanceToday: 4,
    activeBookings: 9,
    pendingTransfers: 3,
    upcomingReturns: 12,
  },
  alerts: {
    overdueAllocations: [
      { allocation_id: 101, asset_tag: 'AF-0112', held_by: 'Alex', expected_return: '2023-10-20', days_overdue: 5 },
      { allocation_id: 102, asset_tag: 'AF-0014', held_by: 'Priya Singh', expected_return: '2023-10-25', days_overdue: 3 },
      { allocation_id: 103, asset_tag: 'AF-0088', held_by: 'John Doe', expected_return: '2023-10-26', days_overdue: 2 },
    ],
  },
  recentActivity: [
    { id: 1, text: 'Laptop AF-0114 - allocated to Priya Shah - IT dept', time: '2 hours ago' },
    { id: 2, text: 'Room B2 - booking confirmed - 2:00 to 3:00 PM', time: '3 hours ago' },
    { id: 3, text: 'Projector AF-0062 - maintenance resolved', time: '5 hours ago' },
  ]
};