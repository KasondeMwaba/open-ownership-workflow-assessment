import type { LucideIcon } from 'lucide-react';
import type { ReactNode } from 'react';

export default function EmptyState({
  icon: Icon,
  title,
  description,
  action,
}: {
  icon: LucideIcon;
  title: string;
  description: string;
  action?: ReactNode;
}) {
  return (
    <div className="grid place-items-center px-6 py-12 text-center">
      <div className="flex h-12 w-12 items-center justify-center rounded-md bg-emerald-50 text-deepgreen ring-1 ring-emerald-100 dark:bg-emerald-950 dark:text-emerald-100 dark:ring-emerald-900">
        <Icon size={22} />
      </div>
      <h3 className="mt-4 text-base font-black text-ink dark:text-slate-100">{title}</h3>
      <p className="mt-1 max-w-md text-sm leading-6 text-slate-500 dark:text-slate-400">{description}</p>
      {action && <div className="mt-5">{action}</div>}
    </div>
  );
}
