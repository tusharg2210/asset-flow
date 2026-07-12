export type Role = 'ADMIN' | 'ASSET_MANAGER' | 'DEPARTMENT_HEAD' | 'EMPLOYEE';

export interface User {
  id: number;
  name: string;
  email: string;
  role: Role;
  department_id?: number;
  allotted_asset_id?: number[];
  status?: 'ACTIVE' | 'INACTIVE';
}