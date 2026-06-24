import { Eye, FileSearch, ShieldCheck, X } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { listAuditTrail } from '../api/client';
import EmptyState from '../components/EmptyState';
import PageHeader from '../components/PageHeader';
import PaginationControls from '../components/PaginationControls';
import StatusBadge from '../components/StatusBadge';
import type { ActivityAuditEvent, AuditEvent, AuditTrailResult, SessionAuditEvent, SystemAuditEvent, User } from '../types/domain';

const PAGE_SIZE = 8;

type DetailEvent =
  | { kind: 'submission'; event: AuditEvent }
  | { kind: 'system'; event: SystemAuditEvent }
  | { kind: 'session'; event: SessionAuditEvent }
  | { kind: 'activity'; event: ActivityAuditEvent };

export type AuditMode = 'activity' | 'session' | 'system' | 'submission';

const auditTitles: Record<AuditMode, { title: string; description: string }> = {
  activity: { title: 'Activity Audit', description: 'Authenticated user movements across workflow and administration endpoints.' },
  session: { title: 'Session Audit', description: 'Login, failed login, disabled-account attempts, and logout events.' },
  system: { title: 'System Audit', description: 'User, role, and permission administration events.' },
  submission: { title: 'Submission Audit', description: 'Submission creation, edit, and workflow transition events.' },
};

export default function AuditTrail({ user, mode = 'submission' }: { user: User; mode?: AuditMode }) {
  const [audit, setAudit] = useState<AuditTrailResult>({ submissionEvents: [], systemEvents: [], sessionEvents: [], activityEvents: [] });
  const [busy, setBusy] = useState(true);
  const [error, setError] = useState('');
  const [detail, setDetail] = useState<DetailEvent | null>(null);
  const [search, setSearch] = useState('');
  const [eventFilter, setEventFilter] = useState('');
  const [sessionResult, setSessionResult] = useState<'all' | 'success' | 'failed'>('all');

  useEffect(() => {
    listAuditTrail()
      .then(setAudit)
      .catch(() => setError('Could not load audit trail.'))
      .finally(() => setBusy(false));
  }, []);

  const filteredSessionEvents = useMemo(
    () =>
      audit.sessionEvents.filter((event) => {
        const text = `${event.eventType} ${event.email} ${event.actorName} ${event.browser} ${event.ipAddress} ${event.reason}`.toLowerCase();
        const matchesSearch = text.includes(search.toLowerCase());
        const matchesEvent = !eventFilter || event.eventType === eventFilter;
        const matchesResult = sessionResult === 'all' || (sessionResult === 'success' ? event.success : !event.success);
        return matchesSearch && matchesEvent && matchesResult;
      }),
    [audit.sessionEvents, eventFilter, search, sessionResult],
  );
  const filteredActivityEvents = useMemo(
    () =>
      audit.activityEvents.filter((event) => {
        const text = `${event.method} ${event.path} ${event.query} ${event.actorName} ${event.actorRole} ${event.browser} ${event.ipAddress} ${event.statusCode}`.toLowerCase();
        const matchesSearch = text.includes(search.toLowerCase());
        const matchesResult = sessionResult === 'all' || (sessionResult === 'success' ? event.success : !event.success);
        return matchesSearch && matchesResult;
      }),
    [audit.activityEvents, search, sessionResult],
  );
  const filteredSystemEvents = useMemo(
    () =>
      audit.systemEvents.filter((event) => {
        const text = `${event.eventType} ${event.resourceType} ${event.resourceName} ${event.actorName} ${event.summary}`.toLowerCase();
        const matchesSearch = text.includes(search.toLowerCase());
        const matchesEvent = !eventFilter || event.eventType === eventFilter;
        return matchesSearch && matchesEvent;
      }),
    [audit.systemEvents, eventFilter, search],
  );
  const filteredSubmissionEvents = useMemo(
    () =>
      audit.submissionEvents.filter((event) => {
        const text = `${event.submissionTitle} ${event.actorName} ${event.actorRole} ${event.comment} ${event.toStatus}`.toLowerCase();
        return text.includes(search.toLowerCase());
      }),
    [audit.submissionEvents, search],
  );
  const eventOptions = useMemo(() => {
    const names = new Set<string>();
    audit.sessionEvents.forEach((event) => names.add(event.eventType));
    audit.systemEvents.forEach((event) => names.add(event.eventType));
    return Array.from(names).sort();
  }, [audit.sessionEvents, audit.systemEvents]);

  return (
    <main className="px-4 py-6 md:px-6 md:py-8">
      <PageHeader
        eyebrow="Governance"
        title={auditTitles[mode].title}
        description={user.role === 'admin' ? auditTitles[mode].description : 'Submission activity for records you are allowed to view.'}
        action={
          <div className="inline-flex items-center gap-2 rounded-md border border-emerald-200 bg-emerald-50 px-3 py-2 text-sm font-semibold text-deepgreen">
            <ShieldCheck size={17} />
            Permission-aware log
          </div>
        }
      />

      {busy ? (
        <p className="mt-6 rounded-md border border-slate-200 bg-white p-6 text-sm text-slate-500 shadow-panel dark:border-slate-800 dark:bg-slate-900">Loading audit trail...</p>
      ) : error ? (
        <p className="mt-6 rounded-md border border-rose-200 bg-rose-50 p-6 text-sm font-semibold text-rose-600">{error}</p>
      ) : (
        <div className="mt-6 space-y-6">
          <AuditFilters mode={mode} search={search} eventFilter={eventFilter} sessionResult={sessionResult} eventOptions={eventOptions} onSearch={setSearch} onEventFilter={setEventFilter} onSessionResult={setSessionResult} />
          {mode === 'activity' && <ActivityAuditSection events={filteredActivityEvents} onView={(event) => setDetail({ kind: 'activity', event })} />}
          {mode === 'session' && <SessionAuditSection events={filteredSessionEvents} onView={(event) => setDetail({ kind: 'session', event })} />}
          {mode === 'system' && <SystemAuditSection events={filteredSystemEvents} onView={(event) => setDetail({ kind: 'system', event })} />}
          {mode === 'submission' && <SubmissionAuditSection events={filteredSubmissionEvents} onView={(event) => setDetail({ kind: 'submission', event })} />}
        </div>
      )}

      {detail && <AuditDetailModal detail={detail} onClose={() => setDetail(null)} />}
    </main>
  );
}

function ActivityAuditSection({ events, onView }: { events: ActivityAuditEvent[]; onView: (event: ActivityAuditEvent) => void }) {
  const [page, setPage] = useState(1);
  const rows = usePaginatedRows(events, page);
  const pages = Math.max(1, Math.ceil(events.length / PAGE_SIZE));
  useEffect(() => setPage(1), [events.length]);

  return (
    <section className="overflow-hidden rounded-md border border-slate-200 bg-white shadow-panel dark:border-slate-800 dark:bg-slate-900">
      <SectionHeader title="Activity Audit" subtitle="Authenticated user movements across workflow and administration endpoints." count={events.length} />
      <div className="overflow-x-auto">
        <table className="app-table">
          <thead>
            <tr>
              <th>Time</th>
              <th>User</th>
              <th>Action</th>
              <th>Status</th>
              <th>Browser</th>
              <th>Duration</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((event) => (
              <tr key={event.id}>
                <td className="px-4 py-3 text-slate-500">{new Date(event.createdAt).toLocaleString()}</td>
                <td className="px-4 py-3">
                  <p className="font-semibold text-slate-700 dark:text-slate-200">{event.actorName || 'Unknown user'}</p>
                  <p className="text-xs capitalize text-slate-400">{event.actorRole || 'unknown'}</p>
                </td>
                <td className="px-4 py-3">
                  <p className="font-black text-ink dark:text-slate-100">{event.method}</p>
                  <p className="text-xs text-slate-500">{event.path}{event.query ? `?${event.query}` : ''}</p>
                </td>
                <td className="px-4 py-3">
                  <span className={`inline-flex rounded px-2 py-1 text-xs font-bold ring-1 ${event.success ? 'bg-emerald-50 text-emerald-700 ring-emerald-200' : 'bg-rose-50 text-rose-700 ring-rose-200'}`}>
                    {event.statusCode}
                  </span>
                </td>
                <td className="px-4 py-3 text-slate-600 dark:text-slate-300">{event.browser || 'Unknown'}</td>
                <td className="px-4 py-3 text-slate-600 dark:text-slate-300">{event.durationMs}ms</td>
                <td className="px-4 py-3">
                  <button onClick={() => onView(event)} className="focus-ring inline-flex items-center gap-2 rounded-md border border-slate-300 px-3 py-2 text-xs font-bold text-slate-600 hover:bg-slate-50 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800">
                    <Eye size={14} />
                    View
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <PaginationControls page={page} totalPages={pages} totalItems={events.length} onPage={setPage} />
    </section>
  );
}

function AuditFilters({
  mode,
  search,
  eventFilter,
  sessionResult,
  eventOptions,
  onSearch,
  onEventFilter,
  onSessionResult,
}: {
  mode: AuditMode;
  search: string;
  eventFilter: string;
  sessionResult: 'all' | 'success' | 'failed';
  eventOptions: string[];
  onSearch: (value: string) => void;
  onEventFilter: (value: string) => void;
  onSessionResult: (value: 'all' | 'success' | 'failed') => void;
}) {
  return (
    <section className="grid gap-3 rounded-md border border-slate-200 bg-white p-4 shadow-panel dark:border-slate-800 dark:bg-slate-900 md:grid-cols-[1fr_220px_180px]">
      <input className="focus-ring rounded-md border border-slate-300 px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-950 dark:text-slate-100" placeholder="Search audit events, actors, emails, IPs..." value={search} onChange={(event) => onSearch(event.target.value)} />
      {(mode === 'session' || mode === 'system') ? (
        <select className="focus-ring rounded-md border border-slate-300 px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-950 dark:text-slate-100" value={eventFilter} onChange={(event) => onEventFilter(event.target.value)}>
          <option value="">All event types</option>
          {eventOptions.map((eventType) => (
            <option key={eventType} value={eventType}>
              {eventType}
            </option>
          ))}
        </select>
      ) : (
        <span className="hidden md:block" />
      )}
      <select className={`focus-ring rounded-md border border-slate-300 px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-950 dark:text-slate-100 ${mode === 'system' || mode === 'submission' ? 'opacity-50' : ''}`} value={sessionResult} onChange={(event) => onSessionResult(event.target.value as 'all' | 'success' | 'failed')} disabled={mode === 'system' || mode === 'submission'}>
        <option value="all">All results</option>
        <option value="success">Success only</option>
        <option value="failed">Failed only</option>
      </select>
    </section>
  );
}

function SessionAuditSection({ events, onView }: { events: SessionAuditEvent[]; onView: (event: SessionAuditEvent) => void }) {
  const [page, setPage] = useState(1);
  const rows = usePaginatedRows(events, page);
  const pages = Math.max(1, Math.ceil(events.length / PAGE_SIZE));
  useEffect(() => setPage(1), [events.length]);

  return (
    <section className="overflow-hidden rounded-md border border-slate-200 bg-white shadow-panel dark:border-slate-800 dark:bg-slate-900">
      <SectionHeader title="Session Audit" subtitle="Login, failed login, and logout attempts with browser and network details." count={events.length} />
      <div className="overflow-x-auto">
        <table className="app-table">
          <thead>
            <tr>
              <th>Time</th>
              <th>Event</th>
              <th>Email</th>
              <th>Browser</th>
              <th>IP address</th>
              <th>Result</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((event) => (
              <tr key={event.id}>
                <td className="px-4 py-3 text-slate-500">{new Date(event.createdAt).toLocaleString()}</td>
                <td className="px-4 py-3 font-black text-ink dark:text-slate-100">{event.eventType}</td>
                <td className="px-4 py-3">
                  <p className="font-semibold text-slate-700 dark:text-slate-200">{event.email || 'Unknown email'}</p>
                  <p className="text-xs text-slate-400">{event.actorName || 'No matched account'}</p>
                </td>
                <td className="px-4 py-3 text-slate-600 dark:text-slate-300">{event.browser || 'Unknown'}</td>
                <td className="px-4 py-3 text-slate-600 dark:text-slate-300">{event.ipAddress || '-'}</td>
                <td className="px-4 py-3">
                  <span className={`inline-flex rounded px-2 py-1 text-xs font-bold ring-1 ${event.success ? 'bg-emerald-50 text-emerald-700 ring-emerald-200' : 'bg-rose-50 text-rose-700 ring-rose-200'}`}>
                    {event.success ? 'Success' : 'Failed'}
                  </span>
                </td>
                <td className="px-4 py-3">
                  <button onClick={() => onView(event)} className="focus-ring inline-flex items-center gap-2 rounded-md border border-slate-300 px-3 py-2 text-xs font-bold text-slate-600 hover:bg-slate-50 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800">
                    <Eye size={14} />
                    View
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <PaginationControls page={page} totalPages={pages} totalItems={events.length} onPage={setPage} />
    </section>
  );
}

function SubmissionAuditSection({ events, onView }: { events: AuditEvent[]; onView: (event: AuditEvent) => void }) {
  const [page, setPage] = useState(1);
  const rows = usePaginatedRows(events, page);
  const pages = Math.max(1, Math.ceil(events.length / PAGE_SIZE));
  useEffect(() => setPage(1), [events.length]);

  return (
    <section className="overflow-hidden rounded-md border border-slate-200 bg-white shadow-panel dark:border-slate-800 dark:bg-slate-900">
      <SectionHeader title="Submission Audit" subtitle="Events for submissions visible to this account." count={events.length} />
      <div className="overflow-x-auto">
        <table className="app-table">
          <thead>
            <tr>
              <th className="px-4 py-3">Time</th>
              <th className="px-4 py-3">Submission</th>
              <th className="px-4 py-3">Status</th>
              <th className="px-4 py-3">Actor</th>
              <th className="px-4 py-3">Comment</th>
              <th className="px-4 py-3">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
            {rows.map((event) => (
              <tr key={event.id}>
                <td className="px-4 py-3 text-slate-500">{new Date(event.createdAt).toLocaleString()}</td>
                <td className="px-4 py-3">
                  <Link to={`/submissions/${event.submissionId}`} className="font-semibold text-ink hover:text-accent dark:text-slate-100">
                    {event.submissionTitle || event.submissionId}
                  </Link>
                </td>
                <td className="px-4 py-3"><StatusBadge status={event.toStatus} /></td>
                <td className="px-4 py-3">
                  <p className="font-semibold text-slate-700 dark:text-slate-200">{event.actorName || 'Unknown actor'}</p>
                  <p className="text-xs capitalize text-slate-400">{event.actorRole || 'unknown'}</p>
                </td>
                <td className="max-w-sm truncate px-4 py-3 text-slate-600 dark:text-slate-300">{event.comment || 'No comment'}</td>
                <td className="px-4 py-3">
                  <button onClick={() => onView(event)} className="focus-ring inline-flex items-center gap-2 rounded-md border border-slate-300 px-3 py-2 text-xs font-bold text-slate-600 hover:bg-slate-50 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800">
                    <Eye size={14} />
                    View
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {events.length === 0 ? (
        <EmptyState icon={FileSearch} title="No submission audit events" description="There are no visible submission events for this account yet." />
      ) : (
        <PaginationControls page={page} totalPages={pages} totalItems={events.length} onPage={setPage} />
      )}
    </section>
  );
}

function SystemAuditSection({ events, onView }: { events: SystemAuditEvent[]; onView: (event: SystemAuditEvent) => void }) {
  const [page, setPage] = useState(1);
  const rows = usePaginatedRows(events, page);
  const pages = Math.max(1, Math.ceil(events.length / PAGE_SIZE));
  useEffect(() => setPage(1), [events.length]);

  return (
    <section className="overflow-hidden rounded-md border border-slate-200 bg-white shadow-panel dark:border-slate-800 dark:bg-slate-900">
      <SectionHeader title="System Audit" subtitle="Admin, user, role, and permission changes." count={events.length} />
      <div className="overflow-x-auto">
        <table className="app-table">
          <thead>
            <tr>
              <th className="px-4 py-3">Time</th>
              <th className="px-4 py-3">Event</th>
              <th className="px-4 py-3">Resource</th>
              <th className="px-4 py-3">Actor</th>
              <th className="px-4 py-3">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
            {rows.map((event) => (
              <tr key={event.id}>
                <td className="px-4 py-3 text-slate-500">{new Date(event.createdAt).toLocaleString()}</td>
                <td className="px-4 py-3 font-black text-ink dark:text-slate-100">{event.eventType}</td>
                <td className="px-4 py-3">
                  <p className="font-semibold text-slate-700 dark:text-slate-200">{event.resourceName}</p>
                  <p className="text-xs uppercase tracking-wide text-slate-400">{event.resourceType}</p>
                </td>
                <td className="px-4 py-3">
                  <p className="font-semibold text-slate-700 dark:text-slate-200">{event.actorName || 'Unknown actor'}</p>
                  <p className="text-xs capitalize text-slate-400">{event.actorRole || 'unknown'}</p>
                </td>
                <td className="px-4 py-3">
                  <button onClick={() => onView(event)} className="focus-ring inline-flex items-center gap-2 rounded-md border border-slate-300 px-3 py-2 text-xs font-bold text-slate-600 hover:bg-slate-50 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800">
                    <Eye size={14} />
                    View
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <PaginationControls page={page} totalPages={pages} totalItems={events.length} onPage={setPage} />
    </section>
  );
}

function AuditDetailModal({ detail, onClose }: { detail: DetailEvent; onClose: () => void }) {
  const event = detail.event;
  const title = detail.kind === 'system' ? detail.event.eventType : detail.kind === 'session' ? `${detail.event.eventType} attempt` : detail.kind === 'activity' ? `${detail.event.method} ${detail.event.path}` : detail.event.submissionTitle || 'Submission audit event';
  const metadata = detail.event.metadata;

  return (
    <div className="fixed inset-0 z-50 grid place-items-center bg-black/50 px-4 py-6">
      <section className="max-h-[88vh] w-full max-w-3xl overflow-hidden rounded-md bg-white shadow-2xl dark:bg-slate-900">
        <div className="flex items-start justify-between gap-4 border-b border-slate-200 p-5 dark:border-slate-800">
          <div>
            <p className="text-xs font-bold uppercase tracking-[0.18em] text-gold">{detail.kind === 'system' ? 'System audit' : detail.kind === 'session' ? 'Session audit' : detail.kind === 'activity' ? 'Activity audit' : 'Submission audit'}</p>
            <h3 className="mt-1 text-xl font-black text-ink dark:text-slate-100">{title}</h3>
            <p className="mt-1 text-sm text-slate-500">{new Date(event.createdAt).toLocaleString()}</p>
          </div>
          <button onClick={onClose} className="focus-ring rounded-md border border-slate-300 p-2 text-slate-500 hover:bg-slate-50 dark:border-slate-700 dark:hover:bg-slate-800">
            <X size={18} />
          </button>
        </div>
        <div className="grid max-h-[70vh] gap-4 overflow-y-auto p-5 md:grid-cols-2">
          <DetailBlock label="Actor" value={`${event.actorName || 'Unknown actor'} (${event.actorRole || 'unknown'})`} />
          <DetailBlock label="Actor ID" value={event.actorId || 'No matched account'} />
          {detail.kind === 'system' ? (
            <>
              <DetailBlock label="Resource" value={`${detail.event.resourceType}: ${detail.event.resourceName}`} />
              <DetailBlock label="Resource ID" value={detail.event.resourceId} />
              <DetailBlock label="Summary" value={detail.event.summary || 'No summary'} wide />
            </>
          ) : detail.kind === 'session' ? (
            <>
              <DetailBlock label="Email" value={detail.event.email || 'Unknown email'} />
              <DetailBlock label="Result" value={detail.event.success ? 'Success' : 'Failed'} />
              <DetailBlock label="Reason" value={detail.event.reason || 'No reason recorded'} />
              <DetailBlock label="Browser" value={detail.event.browser || 'Unknown browser'} />
              <DetailBlock label="IP Address" value={detail.event.ipAddress || '-'} />
              <DetailBlock label="User Agent" value={detail.event.userAgent || '-'} wide />
            </>
          ) : detail.kind === 'activity' ? (
            <>
              <DetailBlock label="Request" value={`${detail.event.method} ${detail.event.path}${detail.event.query ? `?${detail.event.query}` : ''}`} wide />
              <DetailBlock label="Status" value={`${detail.event.statusCode} (${detail.event.success ? 'success' : 'failed'})`} />
              <DetailBlock label="Duration" value={`${detail.event.durationMs}ms`} />
              <DetailBlock label="Browser" value={detail.event.browser || 'Unknown browser'} />
              <DetailBlock label="IP Address" value={detail.event.ipAddress || '-'} />
              <DetailBlock label="User Agent" value={detail.event.userAgent || '-'} wide />
            </>
          ) : (
            <>
              <DetailBlock label="Submission" value={detail.event.submissionTitle || detail.event.submissionId} />
              <DetailBlock label="Submission ID" value={detail.event.submissionId} />
              <DetailBlock label="Status Change" value={`${detail.event.fromStatus || 'none'} -> ${detail.event.toStatus}`} />
              <DetailBlock label="Comment" value={detail.event.comment || 'No comment'} wide />
            </>
          )}
          <div className="md:col-span-2">
            <p className="text-xs font-bold uppercase tracking-wide text-slate-400">Metadata</p>
            <pre className="mt-2 max-h-80 overflow-auto rounded-md bg-slate-950 p-4 text-xs leading-5 text-slate-100">
              {JSON.stringify(metadata, null, 2)}
            </pre>
          </div>
        </div>
      </section>
    </div>
  );
}

function SectionHeader({ title, subtitle, count }: { title: string; subtitle: string; count: number }) {
  return (
    <div className="table-header-band flex flex-wrap items-center justify-between gap-3">
      <div>
        <h3 className="font-black text-ink dark:text-slate-100">{title}</h3>
        <p className="mt-1 text-xs text-slate-500">{subtitle}</p>
      </div>
      <span className="rounded bg-white px-2 py-1 text-xs font-black text-slate-500 ring-1 ring-slate-200 dark:bg-slate-900 dark:ring-slate-800">{count} records</span>
    </div>
  );
}

function DetailBlock({ label, value, wide = false }: { label: string; value: string; wide?: boolean }) {
  return (
    <div className={wide ? 'md:col-span-2' : ''}>
      <p className="text-xs font-bold uppercase tracking-wide text-slate-400">{label}</p>
      <p className="mt-1 break-words text-sm font-semibold text-slate-700 dark:text-slate-200">{value}</p>
    </div>
  );
}

function usePaginatedRows<T>(rows: T[], page: number) {
  return useMemo(() => rows.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE), [rows, page]);
}
