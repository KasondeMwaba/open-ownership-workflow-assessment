import { ArrowRight, CheckCircle2, ClipboardCheck, KeyRound, LockKeyhole, ShieldCheck, type LucideIcon } from 'lucide-react';
import { FormEvent, useState } from 'react';
import { login } from '../api/client';
import type { User } from '../types/domain';

const portalMetrics: Array<[string, string, LucideIcon]> = [
  ['Submitted', '1', ClipboardCheck],
  ['Audit events', 'Live', CheckCircle2],
  ['Access model', 'RBAC', LockKeyhole],
];

export default function Login({ onLogin }: { onLogin: (user: User) => void }) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);

  async function submit(event: FormEvent) {
    event.preventDefault();
    setBusy(true);
    setError('');
    try {
      onLogin(await login(email, password));
    } catch {
      setError('Could not sign in with those credentials.');
    } finally {
      setBusy(false);
    }
  }

  return (
    <main className="min-h-screen bg-[#f4f7f4] px-4 py-5 text-ink sm:px-6 lg:px-8">
      <div className="mx-auto grid min-h-[calc(100vh-2.5rem)] max-w-7xl overflow-hidden rounded-lg border border-slate-200 bg-white shadow-2xl lg:grid-cols-[1.05fr_0.95fr]">
        <section className="relative flex flex-col justify-between overflow-hidden bg-deepgreen p-7 text-white md:p-10">
          <div className="absolute inset-x-0 bottom-0 h-44 bg-[linear-gradient(180deg,transparent,rgba(0,0,0,0.18))]" />
          <div className="relative z-10">
            <div className="flex items-center gap-3">
              <div className="flex h-12 w-12 items-center justify-center rounded-md bg-white text-deepgreen shadow-lg">
                <ShieldCheck size={25} />
              </div>
              <div>
                <p className="text-lg font-black tracking-tight">Open Ownership</p>
                <p className="text-[10px] font-bold uppercase tracking-[0.24em] text-white/60">Assessment Portal</p>
              </div>
            </div>

            <div className="mt-16 max-w-2xl">
              <p className="text-xs font-bold uppercase tracking-[0.24em] text-emerald-100/80">Technical assessment</p>
              <h1 className="mt-4 text-4xl font-black leading-[1.05] tracking-tight md:text-6xl">
                Submission approval workflow
              </h1>
              <p className="mt-5 max-w-xl text-base leading-7 text-emerald-50/82 md:text-lg">
                Role-aware review, strict backend transitions, audit history, and Redis-backed operational metrics in one review console.
              </p>
            </div>
          </div>

          <div className="relative z-10 mt-12 grid gap-3 sm:grid-cols-3">
            {portalMetrics.map(([label, value, Icon]) => (
              <div key={label} className="rounded-md border border-white/10 bg-white/10 p-4 backdrop-blur">
                <Icon className="text-emerald-100" size={18} />
                <p className="mt-4 text-2xl font-black">{value}</p>
                <p className="mt-1 text-[10px] font-bold uppercase tracking-[0.18em] text-white/58">{label}</p>
              </div>
            ))}
          </div>
        </section>

        <section className="flex items-center justify-center bg-paper p-5 md:p-10">
          <form onSubmit={submit} className="w-full max-w-md rounded-lg border border-slate-200 bg-white p-6 shadow-panel md:p-7">
            <div className="flex items-start justify-between gap-4">
              <div>
                <div className="flex h-12 w-12 items-center justify-center rounded-md bg-emerald-50 text-deepgreen">
                  <KeyRound size={22} />
                </div>
                <h2 className="mt-5 text-2xl font-black tracking-tight text-ink">Sign in</h2>
                <p className="mt-2 text-sm text-slate-500">Enter your email and password to access the workflow portal.</p>
              </div>
              <span className="rounded border border-emerald-200 bg-emerald-50 px-2 py-1 text-[10px] font-bold uppercase tracking-[0.18em] text-deepgreen">
                Secure
              </span>
            </div>

            <div className="mt-6 space-y-4">
              <label className="block">
                <span className="text-sm font-semibold text-slate-700">Email</span>
                <input
                  className="focus-ring mt-1 w-full rounded-md border border-slate-300 bg-white px-3 py-2.5 text-sm"
                  type="email"
                  value={email}
                  onChange={(event) => setEmail(event.target.value)}
                  required
                />
              </label>
              <label className="block">
                <span className="text-sm font-semibold text-slate-700">Password</span>
                <input
                  className="focus-ring mt-1 w-full rounded-md border border-slate-300 px-3 py-2.5 text-sm"
                  type="password"
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                  required
                />
              </label>
            </div>

            {error && <p className="mt-4 rounded-md border border-rose-200 bg-rose-50 px-3 py-2 text-sm font-medium text-rose-700">{error}</p>}

            <button
              className="btn-ink mt-6 w-full py-3"
              disabled={busy}
            >
              {busy ? 'Signing in...' : 'Continue'}
              {!busy && <ArrowRight size={18} />}
            </button>
          </form>
        </section>
      </div>
    </main>
  );
}
