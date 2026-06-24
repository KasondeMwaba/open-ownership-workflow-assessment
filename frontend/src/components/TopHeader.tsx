import { Bell, LogOut, Search, ShieldCheck } from 'lucide-react';
import { useMemo } from 'react';
import type { User } from '../types/domain';
import ThemeToggle from './ThemeToggle';

type Props = {
  user: User;
  onLogout: () => void;
};

export default function TopHeader({ user, onLogout }: Props) {
  const initials = useMemo(
    () =>
      user.name
        .split(' ')
        .map((part) => part[0])
        .join('')
        .slice(0, 2)
        .toUpperCase(),
    [user.name],
  );

  return (
    <header className="sticky top-0 z-30 flex h-16 shrink-0 items-center justify-between border-b border-slate-200 bg-white/90 px-4 shadow-sm backdrop-blur dark:border-slate-800 dark:bg-slate-950/88 md:px-6">
      <div className="flex items-center gap-3">
        <span className="hidden rounded border border-emerald-200 bg-emerald-50 px-2.5 py-1 text-[10px] font-bold uppercase tracking-[0.18em] text-deepgreen sm:inline-flex">
          Workflow Beta
        </span>
        <div className="relative hidden w-72 md:block">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-slate-400" />
          <input
            type="text"
            placeholder="Search submissions, owners, status..."
            className="focus-ring h-9 w-full rounded-md border border-slate-300 bg-white pl-9 pr-3 text-sm dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
          />
        </div>
      </div>

      <div className="flex items-center gap-3">
        <div className="hidden items-center gap-2 rounded-md border border-slate-200 bg-slate-50 px-3 py-2 text-xs font-semibold text-emerald-700 md:flex">
          <ShieldCheck size={15} />
          Operational
        </div>
        <button className="focus-ring relative inline-flex h-9 w-9 items-center justify-center rounded-md text-slate-500 hover:bg-slate-100" title="Notifications">
          <Bell size={18} />
          <span className="absolute right-2 top-2 h-2 w-2 rounded-full bg-rose-500 ring-2 ring-white" />
        </button>
        <ThemeToggle />
        <div className="hidden items-center gap-2 border-l border-slate-200 pl-3 dark:border-slate-800 sm:flex">
          <div className="text-right">
            <p className="text-sm font-semibold leading-none text-ink dark:text-slate-100">{user.name}</p>
            <p className="mt-1 text-xs capitalize text-slate-500">{user.role}</p>
          </div>
          <div className="flex h-9 w-9 items-center justify-center rounded-md bg-emerald-100 text-sm font-black text-deepgreen">
            {initials}
          </div>
        </div>
        <button
          className="focus-ring inline-flex h-9 w-9 items-center justify-center rounded-md text-slate-500 hover:bg-rose-50 hover:text-rose-600 dark:text-slate-300 dark:hover:bg-rose-950/40"
          onClick={onLogout}
          title="Sign out"
        >
          <LogOut size={17} />
        </button>
      </div>
    </header>
  );
}
