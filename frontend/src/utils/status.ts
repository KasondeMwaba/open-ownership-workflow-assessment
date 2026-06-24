import type { Role, Status, User } from '../types/domain';
import { hasPermission } from './permissions';

export const statusLabels: Record<Status, string> = {
  draft: 'Draft',
  submitted: 'Submitted',
  changes_required: 'Changes required',
  approved: 'Approved',
  rejected: 'Rejected',
  withdrawn: 'Withdrawn',
};

export const statusTone: Record<Status, string> = {
  draft: 'bg-slate-100 text-slate-700 ring-slate-200 dark:bg-slate-800 dark:text-slate-200 dark:ring-slate-700',
  submitted: 'bg-sky-50 text-sky-700 ring-sky-200 dark:bg-sky-950 dark:text-sky-200 dark:ring-sky-900',
  changes_required: 'bg-amber-50 text-amber-700 ring-amber-200 dark:bg-amber-950 dark:text-amber-200 dark:ring-amber-900',
  approved: 'bg-emerald-50 text-emerald-700 ring-emerald-200 dark:bg-emerald-950 dark:text-emerald-200 dark:ring-emerald-900',
  rejected: 'bg-rose-50 text-rose-700 ring-rose-200 dark:bg-rose-950 dark:text-rose-200 dark:ring-rose-900',
  withdrawn: 'bg-zinc-100 text-zinc-700 ring-zinc-200 dark:bg-zinc-900 dark:text-zinc-200 dark:ring-zinc-700',
};

const transitions: Record<Status, Partial<Record<Role, Status[]>>> = {
  draft: { requester: ['submitted', 'withdrawn'], admin: ['submitted', 'withdrawn'] },
  submitted: { reviewer: ['changes_required', 'approved', 'rejected'], admin: ['changes_required', 'approved', 'rejected'] },
  changes_required: { requester: ['submitted', 'withdrawn'], admin: ['submitted', 'withdrawn'] },
  approved: {},
  rejected: {},
  withdrawn: {},
};

export function allowedTransitions(status: Status, role: Role) {
  return transitions[status][role] ?? [];
}

export function allowedTransitionsForUser(status: Status, user: User): Status[] {
  if ((status === 'draft' || status === 'changes_required') && hasPermission(user, 'submissions:create')) {
    return ['submitted', 'withdrawn'];
  }
  if (status === 'submitted' && hasPermission(user, 'submissions:review')) {
    return ['changes_required', 'approved', 'rejected'];
  }
  return [];
}
