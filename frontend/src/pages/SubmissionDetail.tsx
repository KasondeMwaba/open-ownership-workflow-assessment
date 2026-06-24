import { ArrowLeft, Check, Edit3, RotateCcw, X } from 'lucide-react';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { getAudit, getSubmission, transitionSubmission } from '../api/client';
import ConfirmDialog from '../components/ConfirmDialog';
import PaginationControls from '../components/PaginationControls';
import StatusBadge from '../components/StatusBadge';
import type { AuditEvent, Status, Submission, User } from '../types/domain';
import { hasPermission } from '../utils/permissions';
import { allowedTransitionsForUser, statusLabels } from '../utils/status';

type BeneficialOwnerView = {
  name: string;
  ownershipPercent: string;
  controlType: string;
};

const AUDIT_PAGE_SIZE = 5;

const actionIcon: Record<Status, typeof Check> = {
  draft: Edit3,
  submitted: RotateCcw,
  changes_required: RotateCcw,
  approved: Check,
  rejected: X,
  withdrawn: X,
};

export default function SubmissionDetail({ user }: { user: User }) {
  const { id = '' } = useParams();
  const [submission, setSubmission] = useState<Submission | null>(null);
  const [audit, setAudit] = useState<AuditEvent[]>([]);
  const [comment, setComment] = useState('');
  const [pendingStatus, setPendingStatus] = useState<Status | null>(null);
  const [auditPage, setAuditPage] = useState(1);
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);

  const transitions = useMemo(
    () => (submission ? allowedTransitionsForUser(submission.status, user) : []),
    [submission, user],
  );
  const auditPages = Math.max(1, Math.ceil(audit.length / AUDIT_PAGE_SIZE));
  const pagedAudit = useMemo(() => audit.slice((auditPage - 1) * AUDIT_PAGE_SIZE, auditPage * AUDIT_PAGE_SIZE), [audit, auditPage]);

  const load = useCallback(async () => {
    const [item, events] = await Promise.all([getSubmission(id), getAudit(id)]);
    setSubmission(item);
    setAudit(events);
    setAuditPage(1);
  }, [id]);

  useEffect(() => {
    load();
  }, [load]);

  async function move(status: Status) {
    setBusy(true);
    setError('');
    try {
      await transitionSubmission(id, status, comment);
      setComment('');
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Transition failed.');
    } finally {
      setBusy(false);
    }
  }

  if (!submission) {
    return <main className="mx-auto max-w-7xl px-5 py-8 text-sm text-slate-500">Loading submission...</main>;
  }

  const owners = ownersFromData(submission.data.beneficialOwners);

  return (
    <main className="mx-auto max-w-7xl px-5 py-8">
      <Link to="/" className="focus-ring inline-flex items-center gap-2 rounded-md px-2 py-1 text-sm font-medium text-slate-600 hover:bg-white">
        <ArrowLeft size={16} />
        Back
      </Link>

      <div className="mt-5 grid gap-6 lg:grid-cols-[1fr_380px]">
        <section className="border border-slate-200 bg-white p-6 shadow-panel">
          <div className="flex flex-wrap items-start justify-between gap-4">
            <div>
              <h2 className="text-3xl font-semibold text-ink">{submission.title}</h2>
              <p className="mt-2 max-w-2xl text-slate-600">{submission.summary}</p>
            </div>
            <StatusBadge status={submission.status} />
          </div>

          <dl className="mt-6 grid gap-4 border-y border-slate-200 py-4 sm:grid-cols-3">
            <div>
              <dt className="text-xs font-semibold uppercase tracking-wide text-slate-400">Owner</dt>
              <dd className="mt-1 text-sm text-slate-700">{submission.ownerName || 'Current user'}</dd>
            </div>
            <div>
              <dt className="text-xs font-semibold uppercase tracking-wide text-slate-400">Version</dt>
              <dd className="mt-1 text-sm text-slate-700">v{submission.version}</dd>
            </div>
            <div>
              <dt className="text-xs font-semibold uppercase tracking-wide text-slate-400">Updated</dt>
              <dd className="mt-1 text-sm text-slate-700">{new Date(submission.updatedAt).toLocaleString()}</dd>
            </div>
          </dl>

          <section className="mt-6 overflow-hidden rounded-md border border-slate-200">
            <div className="table-header-band">
              <h3 className="font-black text-ink dark:text-slate-100">Declaration details</h3>
              <p className="mt-1 text-xs font-medium normal-case tracking-normal text-slate-500 dark:text-slate-300">
                Structured information captured from the submission form.
              </p>
            </div>
            <dl className="grid gap-4 p-4 sm:grid-cols-3">
              <DetailItem label="Company" value={textValue(submission.data.company)} />
              <DetailItem label="Jurisdiction" value={textValue(submission.data.jurisdiction)} />
              <DetailItem label="Registration number" value={textValue(submission.data.registrationNumber)} />
            </dl>
          </section>

          <section className="mt-5 overflow-hidden rounded-md border border-slate-200">
            <div className="table-header-band flex flex-wrap items-center justify-between gap-3">
              <div>
                <h3 className="font-black text-ink dark:text-slate-100">Beneficial owners</h3>
                <p className="mt-1 text-xs font-medium normal-case tracking-normal text-slate-500 dark:text-slate-300">
                  Ownership and control details submitted for review.
                </p>
              </div>
              <span className="rounded bg-white px-2 py-1 text-xs font-bold text-deepgreen ring-1 ring-emerald-100 dark:bg-slate-900 dark:text-emerald-100 dark:ring-emerald-900">
                {owners.length} listed
              </span>
            </div>
            <div className="overflow-x-auto">
              <table className="app-table">
                <thead>
                  <tr>
                    <th>Name</th>
                    <th>Ownership</th>
                    <th>Control type</th>
                  </tr>
                </thead>
                <tbody>
                  {owners.map((owner, index) => (
                    <tr key={`${owner.name}-${index}`}>
                      <td className="px-4 py-3 font-semibold text-ink dark:text-slate-100">{owner.name}</td>
                      <td className="px-4 py-3 text-slate-600 dark:text-slate-300">{owner.ownershipPercent}</td>
                      <td className="px-4 py-3 capitalize text-slate-600 dark:text-slate-300">{owner.controlType}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          <details className="mt-5 rounded-md border border-slate-200 bg-slate-50 p-4 text-sm dark:border-slate-800 dark:bg-slate-900">
            <summary className="cursor-pointer font-bold text-slate-700 dark:text-slate-200">View raw payload</summary>
            <pre className="mt-3 max-h-[360px] overflow-auto rounded-md bg-slate-950 p-4 text-sm leading-6 text-slate-100">
              {JSON.stringify(submission.data, null, 2)}
            </pre>
          </details>

          {(submission.status === 'draft' || submission.status === 'changes_required') &&
            hasPermission(user, 'submissions:create') && (
              <Link
                to={`/submissions/${submission.id}/edit`}
                className="btn-secondary mt-5 px-4"
              >
                <Edit3 size={17} />
                Edit submission
              </Link>
            )}
        </section>

        <aside className="space-y-6">
          <section className="border border-slate-200 bg-white p-5 shadow-panel">
            <h3 className="text-lg font-semibold text-ink">Workflow action</h3>
            {transitions.length === 0 ? (
              <p className="mt-3 text-sm text-slate-500">No actions are available for your role at this status.</p>
            ) : (
              <>
                <textarea
                  className="focus-ring mt-4 min-h-24 w-full rounded-md border border-slate-300 px-3 py-2 text-sm"
                  placeholder="Decision note"
                  value={comment}
                  onChange={(event) => setComment(event.target.value)}
                />
                <div className="mt-3 grid gap-2">
                  {transitions.map((status) => {
                    const Icon = actionIcon[status];
                    return (
                      <button
                        key={status}
                        disabled={busy}
                        onClick={() => setPendingStatus(status)}
                        className="btn-ink"
                      >
                        <Icon size={17} />
                        Move to {statusLabels[status]}
                      </button>
                    );
                  })}
                </div>
              </>
            )}
            {error && <p className="mt-3 text-sm font-medium text-rose-600">{error}</p>}
          </section>

          <section className="border border-slate-200 bg-white p-5 shadow-panel">
            <h3 className="text-lg font-semibold text-ink">Audit trail</h3>
            <div className="mt-4 space-y-4">
              {pagedAudit.map((event) => (
                <article key={event.id} className="border-l-2 border-accent pl-4">
                  <div className="flex items-center gap-2">
                    <StatusBadge status={event.toStatus} />
                    <span className="text-xs text-slate-400">{new Date(event.createdAt).toLocaleString()}</span>
                  </div>
                  <p className="mt-2 text-sm font-medium text-ink">{event.actorName}</p>
                  <p className="text-xs capitalize text-slate-500">{event.actorRole}</p>
                  {event.comment && <p className="mt-2 text-sm text-slate-600">{event.comment}</p>}
                </article>
              ))}
            </div>
            {audit.length > 0 && (
              <div className="mt-4 overflow-hidden rounded-md border border-slate-200 dark:border-slate-800">
                <PaginationControls page={auditPage} totalPages={auditPages} totalItems={audit.length} onPage={setAuditPage} />
              </div>
            )}
          </section>
        </aside>
      </div>
      {pendingStatus && (
        <ConfirmDialog
          title={`Move to ${statusLabels[pendingStatus]}?`}
          description="This workflow decision will update the submission status and create an audit trail event."
          confirmLabel={`Move to ${statusLabels[pendingStatus]}`}
          tone={pendingStatus === 'rejected' || pendingStatus === 'withdrawn' ? 'danger' : 'default'}
          onCancel={() => setPendingStatus(null)}
          onConfirm={() => {
            const status = pendingStatus;
            setPendingStatus(null);
            move(status);
          }}
        >
          {comment ? (
            <p className="rounded-md bg-slate-50 px-3 py-2 text-sm text-slate-600 dark:bg-slate-800 dark:text-slate-300">{comment}</p>
          ) : (
            <p className="rounded-md bg-amber-50 px-3 py-2 text-sm font-semibold text-amber-700 ring-1 ring-amber-200">No decision note has been entered.</p>
          )}
        </ConfirmDialog>
      )}
    </main>
  );
}

function DetailItem({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-xs font-semibold uppercase tracking-wide text-slate-400">{label}</dt>
      <dd className="mt-1 text-sm font-semibold text-slate-700 dark:text-slate-200">{value}</dd>
    </div>
  );
}

function textValue(value: unknown, fallback = '-') {
  return typeof value === 'string' && value.trim() ? value : fallback;
}

function ownersFromData(value: unknown): BeneficialOwnerView[] {
  if (!Array.isArray(value)) {
    return [{ name: 'Not provided', ownershipPercent: '-', controlType: '-' }];
  }

  const owners = value.map((item) => {
    const owner = item as Record<string, unknown>;
    const percent = typeof owner.ownershipPercent === 'number' ? `${owner.ownershipPercent}%` : textValue(owner.ownershipPercent);
    return {
      name: textValue(owner.name, 'Unnamed owner'),
      ownershipPercent: percent,
      controlType: textValue(owner.controlType),
    };
  });

  return owners.length ? owners : [{ name: 'Not provided', ownershipPercent: '-', controlType: '-' }];
}
