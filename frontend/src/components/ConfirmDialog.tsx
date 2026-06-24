import type { ReactNode } from 'react';

export default function ConfirmDialog({
  title,
  description,
  confirmLabel,
  tone = 'default',
  children,
  onCancel,
  onConfirm,
}: {
  title: string;
  description: string;
  confirmLabel: string;
  tone?: 'default' | 'danger';
  children?: ReactNode;
  onCancel: () => void;
  onConfirm: () => void;
}) {
  return (
    <div className="fixed inset-0 z-50 grid place-items-center bg-black/50 px-4 py-6">
      <section className="w-full max-w-md rounded-md bg-white p-5 shadow-2xl dark:bg-slate-900">
        <h3 className="text-lg font-black text-ink dark:text-slate-100">{title}</h3>
        <p className="mt-2 text-sm leading-6 text-slate-500 dark:text-slate-400">{description}</p>
        {children && <div className="mt-4">{children}</div>}
        <div className="mt-6 flex justify-end gap-2">
          <button type="button" onClick={onCancel} className="btn-secondary">
            Cancel
          </button>
          <button
            type="button"
            onClick={onConfirm}
            className={tone === 'danger' ? 'focus-ring inline-flex items-center justify-center rounded-md bg-rose-600 px-4 py-2 font-bold text-white hover:bg-rose-700' : 'btn-ink'}
          >
            {confirmLabel}
          </button>
        </div>
      </section>
    </div>
  );
}
