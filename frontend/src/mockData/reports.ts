export const mockReportsData = {
  utilization: [
    { label: 'Eng', value: 40 },
    { label: 'HR', value: 65 },
    { label: 'IT', value: 85 },
    { label: 'Ops', value: 55 },
    { label: 'Sales', value: 35 },
    { label: 'Mktg', value: 75 },
  ],
  mostUsed: [
    { asset: 'Room B2', stat: '34 booking this month' },
    { asset: 'Van AF-343', stat: '21 trips this month' },
    { asset: 'Projector AF-335', stat: '18 uses' },
  ],
  idle: [
    { asset: 'Camera AF-0301', stat: 'unused 60+ days' },
    { asset: 'Chair AF-0410', stat: 'unused 45 days' },
  ],
  actionNeeded: [
    { asset: 'Forklift AF-0087', stat: 'service due in 5 days' },
    { asset: 'Laptop AF-0020', stat: '4 years old : nearing retirement' },
  ],
};