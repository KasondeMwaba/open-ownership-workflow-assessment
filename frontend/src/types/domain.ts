export type Role = string;
export type Status = 'draft' | 'submitted' | 'changes_required' | 'approved' | 'rejected' | 'withdrawn';

export type Permission = {
  id: string;
  createdAt: string;
  name: string;
  description: string;
};

export type AccessRole = {
  id: string;
  createdAt: string;
  name: string;
  description: string;
  permissions: Permission[];
};

export type User = {
  id: string;
  createdAt: string;
  name: string;
  email: string;
  role: Role;
  isActive: boolean;
  permissions: string[];
};

export type CreateUserPayload = {
  name: string;
  email: string;
  password: string;
  role: Role;
  isActive: boolean;
};

export type UpdateUserPayload = {
  name: string;
  email: string;
  role: Role;
  isActive: boolean;
};

export type CreatePermissionPayload = {
  name: string;
  description: string;
};

export type CreateRolePayload = {
  name: string;
  description: string;
  permissionIds: string[];
};

export type UpdateRolePayload = CreateRolePayload;

export type Submission = {
  id: string;
  title: string;
  summary: string;
  data: Record<string, unknown>;
  status: Status;
  ownerId: string;
  ownerName: string;
  reviewerId?: string;
  version: number;
  createdAt: string;
  updatedAt: string;
};

export type AuditEvent = {
  id: string;
  submissionId: string;
  submissionTitle: string;
  actorId: string;
  actorName: string;
  actorRole: Role;
  fromStatus?: Status;
  toStatus: Status;
  comment: string;
  metadata: Record<string, unknown>;
  createdAt: string;
};

export type SystemAuditEvent = {
  id: string;
  createdAt: string;
  actorId: string;
  actorName: string;
  actorRole: Role;
  eventType: string;
  resourceType: string;
  resourceId: string;
  resourceName: string;
  summary: string;
  metadata: Record<string, unknown>;
};

export type SessionAuditEvent = {
  id: string;
  createdAt: string;
  actorId?: string;
  actorName: string;
  actorRole: Role;
  email: string;
  eventType: string;
  success: boolean;
  ipAddress: string;
  userAgent: string;
  browser: string;
  reason: string;
  metadata: Record<string, unknown>;
};

export type ActivityAuditEvent = {
  id: string;
  createdAt: string;
  actorId: string;
  actorName: string;
  actorRole: Role;
  method: string;
  path: string;
  query: string;
  statusCode: number;
  success: boolean;
  durationMs: number;
  ipAddress: string;
  userAgent: string;
  browser: string;
  metadata: Record<string, unknown>;
};

export type AuditTrailResult = {
  submissionEvents: AuditEvent[];
  systemEvents: SystemAuditEvent[];
  sessionEvents: SessionAuditEvent[];
  activityEvents: ActivityAuditEvent[];
};

export type DashboardStats = {
  total: number;
  byStatus: Record<Status, number>;
  awaitingReview: number;
  needsRequester: number;
  completed: number;
  generatedAt: string;
  redisCacheState: string;
};
