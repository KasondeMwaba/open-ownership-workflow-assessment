import clsx from 'clsx';
import type { Status } from '../types/domain';
import { statusLabels, statusTone } from '../utils/status';

export default function StatusBadge({ status }: { status: Status }) {
  return (
    <span className={clsx('inline-flex items-center gap-1.5 rounded px-2 py-1 text-xs font-bold ring-1', statusTone[status])}>
      <span className="h-1.5 w-1.5 rounded-full bg-current" />
      {statusLabels[status]}
    </span>
  );
}
