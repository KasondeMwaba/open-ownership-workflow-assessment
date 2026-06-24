import type { LucideIcon } from 'lucide-react';

type Props = {
  label: string;
  value: number | string;
  icon: LucideIcon;
};

export default function StatTile({ label, value, icon: Icon }: Props) {
  return (
    <div className="group relative overflow-hidden rounded-md border border-slate-200 bg-white p-4 shadow-panel transition hover:-translate-y-0.5 hover:border-emerald-200 hover:shadow-lg dark:border-slate-800 dark:bg-slate-900">
      <div className="absolute inset-x-0 top-0 h-1 bg-emerald-500/70" />
      <div className="flex items-center justify-between">
        <p className="text-sm font-medium text-slate-500">{label}</p>
        <span className="relative z-10 rounded-md bg-emerald-50 p-2 text-deepgreen">
          <Icon size={18} />
        </span>
      </div>
      <p className="relative z-10 mt-3 text-3xl font-black tracking-tight text-ink dark:text-slate-100">{value}</p>
    </div>
  );
}
