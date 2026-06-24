import { BarChart3, ChevronDown, ClipboardCheck, FilePlus2, History, Home, Layers3, LogIn, Route, ShieldCheck, UsersRound } from 'lucide-react';
import { Link, useLocation } from 'react-router-dom';
import type { User } from '../types/domain';
import { hasPermission } from '../utils/permissions';

const navItems = [
  { label: 'Command Center', href: '/', icon: Home },
  { label: 'Review Queue', href: '/?status=submitted', icon: ClipboardCheck },
  { label: 'Requester Tasks', href: '/?status=changes_required', icon: UsersRound },
  { label: 'Metrics', href: '/', icon: BarChart3 },
];

const auditItems = [
  { label: 'Submission Audit', href: '/audit/submission', icon: ClipboardCheck, adminOnly: false },
  { label: 'Activity Audit', href: '/audit/activity', icon: Route, adminOnly: true },
  { label: 'Session Audit', href: '/audit/session', icon: LogIn, adminOnly: true },
  { label: 'System Audit', href: '/audit/system', icon: ShieldCheck, adminOnly: true },
];

const adminItems = [
  { label: 'User Management', href: '/admin/users', icon: UsersRound },
  { label: 'Role Management', href: '/admin/roles', icon: Layers3 },
];

export default function PortalSidebar({ user }: { user: User }) {
  const location = useLocation();

  return (
    <aside className="hidden min-h-screen w-72 shrink-0 border-r border-white/10 bg-deepgreen text-white dark:bg-emerald-950 lg:flex lg:flex-col">
      <div className="flex h-16 items-center border-b border-white/10 bg-white px-5">
        <div>
          <p className="text-sm font-black tracking-tight text-deepgreen">Open Ownership</p>
          <p className="text-[10px] font-bold uppercase tracking-[0.22em] text-slate-500">Assessment Portal</p>
        </div>
      </div>
      <div className="flex-1 overflow-y-auto px-3 py-5">
        <p className="px-3 text-[10px] font-bold uppercase tracking-[0.22em] text-white/50">Workflow Modules</p>
        <nav className="mt-3 space-y-1">
          {navItems.map((item) => {
            const [path, query] = item.href.split('?');
            const active = location.pathname === path && (!query || location.search === `?${query}`);
            return (
              <Link
                key={item.label}
                to={item.href}
                className={`flex items-center gap-3 rounded-md px-3 py-2.5 text-sm font-semibold transition ${
                  active ? 'bg-white/15 text-white' : 'text-white/78 hover:bg-white/10 hover:text-white'
                }`}
              >
                <item.icon size={18} />
                <span>{item.label}</span>
              </Link>
            );
          })}
          <div>
            <div className={`flex items-center gap-3 rounded-md px-3 py-2.5 text-sm font-semibold ${location.pathname.startsWith('/audit') ? 'bg-white/15 text-white' : 'text-white/78'}`}>
              <History size={18} />
              <span className="flex-1">Audit Trail</span>
              <ChevronDown size={15} />
            </div>
            <div className="mt-1 space-y-1 border-l border-white/15 pl-3">
              {auditItems
                .filter((item) => user.role === 'admin' || !item.adminOnly)
                .map((item) => {
                  const active = location.pathname === item.href;
                  return (
                    <Link
                      key={item.label}
                      to={item.href}
                      className={`flex items-center gap-3 rounded-md px-3 py-2 text-sm font-semibold transition ${
                        active ? 'bg-white/15 text-white' : 'text-white/70 hover:bg-white/10 hover:text-white'
                      }`}
                    >
                      <item.icon size={16} />
                      <span>{item.label}</span>
                    </Link>
                  );
                })}
            </div>
          </div>
          {(user.role === 'admin' ? adminItems : []).map((item) => {
            const active = location.pathname === item.href;
            return (
              <Link
                key={item.label}
                to={item.href}
                className={`flex items-center gap-3 rounded-md px-3 py-2.5 text-sm font-semibold transition ${
                  active ? 'bg-white/15 text-white' : 'text-white/78 hover:bg-white/10 hover:text-white'
                }`}
              >
                <item.icon size={18} />
                <span>{item.label}</span>
              </Link>
            );
          })}
        </nav>
      </div>
      <div className="border-t border-white/10 p-4">
        <div className="rounded-md bg-white/10 p-3">
          <p className="text-xs font-bold uppercase tracking-wider text-white/55">Signed in as</p>
          <p className="mt-1 truncate text-sm font-semibold">{user.name}</p>
          <p className="mt-0.5 text-xs capitalize text-white/65">{user.role}</p>
        </div>
      </div>
    </aside>
  );
}

export function MobilePortalNav({ user }: { user: User }) {
  const canReview = hasPermission(user, 'submissions:review');
  const canCreate = hasPermission(user, 'submissions:create');
  const actionHref = canReview ? '/?status=submitted' : canCreate ? '/submissions/new' : '/';
  const actionLabel = canReview ? 'Queue' : canCreate ? 'New' : 'Home';

  return (
    <div className="border-b border-slate-200 bg-deepgreen px-4 py-3 text-white lg:hidden">
      <div className="flex items-center justify-between">
        <div>
          <p className="font-black">Open Ownership</p>
          <p className="text-xs text-white/65">Workflow portal</p>
        </div>
        <Link
          to={actionHref}
          className="inline-flex items-center gap-2 rounded-md bg-white px-3 py-2 text-sm font-semibold text-deepgreen"
        >
          {canReview ? <ClipboardCheck size={16} /> : <FilePlus2 size={16} />}
          {actionLabel}
        </Link>
      </div>
    </div>
  );
}
