export const mockAuditCycle = {
  id: 5,
  title: 'Q3 audit: Engineering dept',
  dateRange: '1-15 jul',
  auditors: 'A. Rao, S. Iqbal',
};

export type VerificationStatus = 'Verified' | 'Missing' | 'Damaged' | 'Pending';

export interface AuditAsset {
  id: string;
  tag: string;
  name: string;
  expectedLocation: string;
  verification: VerificationStatus;
}

export const mockAuditAssets: AuditAsset[] = [
  { id: '1', tag: 'AF-003', name: 'Dell laptop', expectedLocation: 'Desk E12', verification: 'Verified' },
  { id: '2', tag: 'AF-9921', name: 'Office chair', expectedLocation: 'Desk E14', verification: 'Missing' },
  { id: '3', tag: 'AF-9838', name: 'Monitor', expectedLocation: 'Desk E15', verification: 'Damaged' },
];