import type { User } from '../types/domain';

export function hasPermission(user: User, permission: string) {
  return user.role === 'admin' || user.permissions?.includes(permission);
}
