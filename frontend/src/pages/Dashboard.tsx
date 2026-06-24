import { CheckCircle2, Clock3, FilePlus2, Inbox, ListFilter, SearchX, RotateCcw } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { getDashboard, listSubmissionsPage } from '../api/client';
import EmptyState from '../components/EmptyState';
import PageHeader from '../components/PageHeader';
import PaginationControls from '../components/PaginationControls';
import StatusBadge from '../components/StatusBadge';
import StatTile from '../components/StatTile';
import type { DashboardStats, Status, Submission, User } from '../types/domain';
import { hasPermission } from '../utils/permissions';
import { statusLabels } from '../utils/status';

const filters: Array<'' | Status> = ['', 'submitted', 'changes_required', 'draft', 'approved', 'rejected', 'withdrawn'];
const PAGE_SIZE = 8;

export default function Dashboard({ user }: { user: User }) {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [items, setItems] = useState<Submission[]>([]);
  const [searchParams, setSearchParams] = useSearchParams();
  const urlStatus = searchParams.get('status') as Status | null;
  const filter: '' | Status = urlStatus && filters.includes(urlStatus) ? urlStatus : '';
  const [busy, setBusy] = useState(true);
  const [page, setPage] = useState(1);
  const [totalItems, setTotalItems] = useState(0);
  const totalPages = Math.max(1, Math.ceil(totalItems / PAGE_SIZE));

  useEffect(() => {
    setBusy(true);
    Promise.all([getDashboard(), listSubmissionsPage(filter || undefined, page, PAGE_SIZE)])
      .then(([nextStats, result]) => {
        setStats(nextStats);
        setItems(result.items);
        setTotalItems(result.total);
      })
      .finally(() => setBusy(false));
  }, [filter, page]);

  useEffect(() => {
    setPage(1);
  }, [filter]);

  return (
    <main className="px-4 py-6 md:px-6 md:py-8">
      <PageHeader
        eyebrow="Command Center"
        title="Submission work queue"
        description={hasPermission(user, 'submissions:create') && !hasPermission(user, 'submissions:review') ? 'Your declarations and review outcomes.' : 'Submissions waiting for reviewer action.'}
        action={hasPermission(user, 'submissions:create') && (
          <Link
            to="/submissions/new"
            className="btn-primary"
          >
            <FilePlus2 size={18} />
            New submission
          </Link>
        )}
      />

      <section className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatTile label="Total" value={stats?.total ?? '-'} icon={Inbox} />
        <StatTile label="Awaiting review" value={stats?.awaitingReview ?? '-'} icon={Clock3} />
        <StatTile label="Needs requester" value={stats?.needsRequester ?? '-'} icon={RotateCcw} />
        <StatTile label="Completed" value={stats?.completed ?? '-'} icon={CheckCircle2} />
      </section>

      <div className="mt-6 flex flex-wrap items-center gap-2">
        <ListFilter size={18} className="text-slate-500" />
        {filters.map((item) => (
          <button
            key={item || 'all'}
            onClick={() => (item ? setSearchParams({ status: item }) : setSearchParams({}))}
            className={`focus-ring rounded-md border px-3 py-2 text-sm font-medium ${
              filter === item ? 'border-accent bg-teal-50 text-accent dark:bg-emerald-950 dark:text-emerald-200' : 'border-slate-200 bg-white text-slate-600 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-300'
            }`}
          >
            {item ? statusLabels[item] : 'All'}
          </button>
        ))}
        {stats && <span className="ml-auto text-xs text-slate-500">Redis cache: {stats.redisCacheState}</span>}
      </div>

      <section className="mt-4 overflow-hidden rounded-md border border-slate-200 bg-white shadow-panel dark:border-slate-800 dark:bg-slate-900">
        <div className="grid grid-cols-[1fr_150px_130px_160px] border-b border-emerald-100 bg-emerald-50 px-4 py-3 text-xs font-black uppercase tracking-wide text-deepgreen dark:border-emerald-900 dark:bg-emerald-950/70 dark:text-emerald-100 max-lg:hidden">
          <span>Submission</span>
          <span>Status</span>
          <span>Version</span>
          <span>Updated</span>
        </div>
        {busy ? (
          <p className="p-6 text-sm text-slate-500">Loading submissions...</p>
        ) : items.length === 0 ? (
          <EmptyState
            icon={SearchX}
            title="No submissions found"
            description="There are no records for the selected status. Try another filter or create a new submission if your role allows it."
          />
        ) : (
          items.map((item) => (
            <Link
              key={item.id}
              to={`/submissions/${item.id}`}
              className="grid gap-3 border-b border-slate-100 px-4 py-4 hover:bg-slate-50 dark:border-slate-800 dark:hover:bg-slate-800/70 lg:grid-cols-[1fr_150px_130px_160px]"
            >
              <div>
                <p className="font-semibold text-ink dark:text-slate-100">{item.title}</p>
                <p className="mt-1 line-clamp-1 text-sm text-slate-500">{item.summary}</p>
                <p className="mt-1 text-xs text-slate-400">Owner: {item.ownerName || 'Current user'}</p>
              </div>
              <div><StatusBadge status={item.status} /></div>
              <p className="text-sm text-slate-600">v{item.version}</p>
              <p className="text-sm text-slate-500">{new Date(item.updatedAt).toLocaleString()}</p>
            </Link>
          ))
        )}
        {!busy && <PaginationControls page={page} totalPages={totalPages} totalItems={totalItems} onPage={setPage} />}
      </section>
    </main>
  );
}
