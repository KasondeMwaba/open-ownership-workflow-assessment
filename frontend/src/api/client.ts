import axios from 'axios';
import type {
  AccessRole,
  AuditEvent,
  AuditTrailResult,
  CreatePermissionPayload,
  CreateRolePayload,
  CreateUserPayload,
  DashboardStats,
  Permission,
  Status,
  Submission,
  SystemAuditEvent,
  UpdateRolePayload,
  UpdateUserPayload,
  User,
} from '../types/domain';

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';
const TOKEN_KEY = 'oo-workflow-token';
export type PaginatedResult<T> = {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
};

export const api = axios.create({
  baseURL: API_URL,
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY);
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export function saveToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
}

export async function login(email: string, password: string) {
  const { data } = await api.post<{ token: string; user: User }>('/api/auth/login', { email, password });
  saveToken(data.token);
  return data.user;
}

export async function logout() {
  await api.post('/api/auth/logout');
}

export async function me() {
  const { data } = await api.get<User>('/api/me');
  return data;
}

export async function getDashboard() {
  const { data } = await api.get<DashboardStats>('/api/dashboard');
  return data;
}

export async function listSubmissions(status?: string) {
  const { data } = await api.get<Submission[]>('/api/submissions', { params: { status: status || undefined } });
  return data;
}

export async function listSubmissionsPage(status: string | undefined, page: number, pageSize: number) {
  const { data } = await api.get<PaginatedResult<Submission>>('/api/submissions', {
    params: { status: status || undefined, page, pageSize },
  });
  return data;
}

export async function getSubmission(id: string) {
  const { data } = await api.get<Submission>(`/api/submissions/${id}`);
  return data;
}

export async function saveSubmission(payload: Pick<Submission, 'title' | 'summary' | 'data'>, id?: string) {
  const { data } = id
    ? await api.put<Submission>(`/api/submissions/${id}`, payload)
    : await api.post<Submission>('/api/submissions', payload);
  return data;
}

export async function transitionSubmission(id: string, status: Status, comment: string) {
  const { data } = await api.post<Submission>(`/api/submissions/${id}/transition`, { status, comment });
  return data;
}

export async function getAudit(id: string) {
  const { data } = await api.get<AuditEvent[]>(`/api/submissions/${id}/audit`);
  return data;
}

export async function listAuditTrail() {
  const { data } = await api.get<AuditTrailResult>('/api/audit');
  return data;
}

export async function listSystemAudit() {
  const { data } = await api.get<SystemAuditEvent[]>('/api/admin/audit');
  return data;
}

export async function listUsers() {
  const { data } = await api.get<User[]>('/api/admin/users');
  return data;
}

export async function createUser(payload: CreateUserPayload) {
  const { data } = await api.post<User>('/api/admin/users', payload);
  return data;
}

export async function updateUser(id: string, payload: UpdateUserPayload) {
  const { data } = await api.put<User>(`/api/admin/users/${id}`, payload);
  return data;
}

export async function setUserStatus(id: string, isActive: boolean) {
  const { data } = await api.post<User>(`/api/admin/users/${id}/status`, { isActive });
  return data;
}

export async function listPermissions() {
  const { data } = await api.get<Permission[]>('/api/admin/permissions');
  return data;
}

export async function createPermission(payload: CreatePermissionPayload) {
  const { data } = await api.post<Permission>('/api/admin/permissions', payload);
  return data;
}

export async function listRoles() {
  const { data } = await api.get<AccessRole[]>('/api/admin/roles');
  return data;
}

export async function createRole(payload: CreateRolePayload) {
  const { data } = await api.post<AccessRole>('/api/admin/roles', payload);
  return data;
}

export async function updateRole(id: string, payload: UpdateRolePayload) {
  const { data } = await api.put<AccessRole>(`/api/admin/roles/${id}`, payload);
  return data;
}
