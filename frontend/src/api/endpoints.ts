export const ENDPOINTS = {
  AUTH: {
    LOGIN: '/api/auth/login',
    SIGNUP: '/api/auth/signup',
    ME: '/api/auth/me',
  },
  DASHBOARD: {
    METRICS: '/api/dashboard/metrics',
    ALERTS: '/api/dashboard/alerts',
  },
  LOGS: {
    RECENT: '/api/logs',
  },
  ORGANIZATION: {
    DEPARTMENTS: '/api/departments',
    USERS: '/api/users', // Employee directory
    ASSET_CATEGORIES: '/api/categories', // Assuming this endpoint will exist
  },
  ASSETS: {
    DIRECTORY: '/api/assets',
    REGISTER: '/api/assets',
    DETAILS: (id: number | string) => `/api/assets/${id}`,
  },
  ALLOCATIONS: {
    CREATE: '/api/allocations',
    RETURN: (id: number | string) => `/api/allocations/${id}/return`,
  },
  TRANSFERS: {
    REQUEST: '/api/transfers',
    STATUS: (id: number | string) => `/api/transfers/${id}/status`,
  },
  BOOKINGS: {
    ASSET_BOOKINGS: (assetId: number | string) => `/api/assets/${assetId}/bookings`,
    CREATE: '/api/bookings',
    STATUS: (id: number | string) => `/api/bookings/${id}/status`,
  },
  MAINTENANCE: {
    LIST: '/api/maintenance',
    CREATE: '/api/maintenance',
    WORKFLOW: (id: number | string) => `/api/maintenance/${id}/workflow`,
  },
  AUDITS: {
    CREATE: '/api/audits',
    ASSETS: (id: number | string) => `/api/audits/${id}/assets`,
    REPORTS: (id: number | string) => `/api/audits/${id}/reports`,
    CLOSE: (id: number | string) => `/api/audits/${id}/close`,
  },
  REPORTS: {
    ANALYTICS: '/api/reports',
  },
  NOTIFICATIONS: {
    ALL: '/api/notifications',
    MARK_READ: (id: number | string) => `/api/notifications/${id}/read`,
  },
};