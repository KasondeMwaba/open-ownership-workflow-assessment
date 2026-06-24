import { ArrowLeft, Plus, Save, Trash2 } from 'lucide-react';
import { FormEvent, useEffect, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { getSubmission, saveSubmission } from '../api/client';
import type { User } from '../types/domain';
import { hasPermission } from '../utils/permissions';

type OwnerForm = {
  name: string;
  ownershipPercent: string;
  controlType: string;
};

type SubmissionForm = {
  title: string;
  summary: string;
  company: string;
  jurisdiction: string;
  registrationNumber: string;
  beneficialOwners: OwnerForm[];
};

const emptyOwner: OwnerForm = { name: '', ownershipPercent: '', controlType: 'shares' };

const initialForm: SubmissionForm = {
  title: 'Beneficial ownership declaration',
  summary: 'Declaration prepared for approval by the review team.',
  company: '',
  jurisdiction: '',
  registrationNumber: '',
  beneficialOwners: [{ ...emptyOwner }],
};

export default function SubmissionEditor({ user }: { user: User }) {
  const { id } = useParams();
  const navigate = useNavigate();
  const [form, setForm] = useState<SubmissionForm>(initialForm);
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    if (!id) return;
    getSubmission(id).then((item) => {
      setForm({
        title: item.title,
        summary: item.summary,
        company: stringValue(item.data.company),
        jurisdiction: stringValue(item.data.jurisdiction),
        registrationNumber: stringValue(item.data.registrationNumber),
        beneficialOwners: ownersFromData(item.data.beneficialOwners),
      });
    });
  }, [id]);

  async function submit(event: FormEvent) {
    event.preventDefault();
    setBusy(true);
    setError('');
    try {
      const owners = form.beneficialOwners.map((owner) => ({
        name: owner.name.trim(),
        ownershipPercent: Number(owner.ownershipPercent),
        controlType: owner.controlType.trim(),
      }));
      if (owners.some((owner) => !owner.name || Number.isNaN(owner.ownershipPercent) || owner.ownershipPercent < 0 || owner.ownershipPercent > 100)) {
        setError('Each beneficial owner needs a name and an ownership percentage between 0 and 100.');
        return;
      }
      const saved = await saveSubmission({
        title: form.title,
        summary: form.summary,
        data: {
          company: form.company,
          jurisdiction: form.jurisdiction,
          registrationNumber: form.registrationNumber,
          beneficialOwners: owners,
        },
      }, id);
      navigate(`/submissions/${saved.id}`);
    } catch {
      setError('Could not save submission.');
    } finally {
      setBusy(false);
    }
  }

  function updateOwner(index: number, patch: Partial<OwnerForm>) {
    setForm({
      ...form,
      beneficialOwners: form.beneficialOwners.map((owner, ownerIndex) => ownerIndex === index ? { ...owner, ...patch } : owner),
    });
  }

  function addOwner() {
    setForm({ ...form, beneficialOwners: [...form.beneficialOwners, { ...emptyOwner }] });
  }

  function removeOwner(index: number) {
    const nextOwners = form.beneficialOwners.filter((_, ownerIndex) => ownerIndex !== index);
    setForm({ ...form, beneficialOwners: nextOwners.length ? nextOwners : [{ ...emptyOwner }] });
  }

  if (!hasPermission(user, 'submissions:create')) {
    return <main className="mx-auto max-w-7xl px-5 py-8 text-sm text-slate-500">You do not have permission to create or edit submissions.</main>;
  }

  return (
    <main className="mx-auto max-w-5xl px-5 py-8">
      <Link to={id ? `/submissions/${id}` : '/'} className="focus-ring inline-flex items-center gap-2 rounded-md px-2 py-1 text-sm font-medium text-slate-600 hover:bg-white">
        <ArrowLeft size={16} />
        Back
      </Link>
      <form onSubmit={submit} className="mt-5 rounded-md border border-slate-200 bg-white p-6 shadow-panel">
        <div className="flex flex-wrap items-start justify-between gap-4">
          <div>
            <p className="text-xs font-bold uppercase tracking-[0.18em] text-gold">Submission request</p>
            <h2 className="mt-1 text-2xl font-black text-ink">{id ? 'Edit submission' : 'New submission'}</h2>
            <p className="mt-1 text-sm text-slate-500">Fill in the declaration details. The system stores it as structured JSON for review.</p>
          </div>
          <button
            className="btn-primary"
            disabled={busy}
          >
            <Save size={18} />
            {busy ? 'Saving...' : 'Save'}
          </button>
        </div>

        <section className="mt-6 grid gap-5 rounded-md border border-slate-200 p-4 md:grid-cols-2">
          <TextField label="Title" value={form.title} onChange={(value) => setForm({ ...form, title: value })} required />
          <TextField label="Company name" value={form.company} onChange={(value) => setForm({ ...form, company: value })} required />
          <TextField label="Jurisdiction" value={form.jurisdiction} onChange={(value) => setForm({ ...form, jurisdiction: value })} required />
          <TextField label="Registration number" value={form.registrationNumber} onChange={(value) => setForm({ ...form, registrationNumber: value })} required />
          <label className="block md:col-span-2">
            <span className="text-sm font-semibold text-slate-700">Summary</span>
            <textarea
              className="focus-ring mt-1 min-h-24 w-full rounded-md border border-slate-300 px-3 py-2"
              value={form.summary}
              onChange={(event) => setForm({ ...form, summary: event.target.value })}
              required
            />
          </label>
        </section>

        <section className="mt-6 overflow-hidden rounded-md border border-slate-200">
          <div className="table-header-band flex items-center justify-between">
            <div>
              <h3 className="font-black text-ink">Beneficial owners</h3>
              <p className="text-xs text-slate-500">Add each person or entity with ownership or control.</p>
            </div>
            <button type="button" onClick={addOwner} className="btn-secondary text-sm">
              <Plus size={16} />
              Add owner
            </button>
          </div>
          <div className="grid gap-3 p-4">
            {form.beneficialOwners.map((owner, index) => (
              <article key={index} className="grid gap-3 rounded-md border border-slate-200 p-3 md:grid-cols-[1fr_150px_1fr_auto]">
                <TextField label="Owner name" value={owner.name} onChange={(value) => updateOwner(index, { name: value })} required />
                <TextField label="Ownership %" type="number" value={owner.ownershipPercent} onChange={(value) => updateOwner(index, { ownershipPercent: value })} required />
                <TextField label="Control type" value={owner.controlType} onChange={(value) => updateOwner(index, { controlType: value })} required />
                <button
                  type="button"
                  disabled={form.beneficialOwners.length === 1}
                  onClick={() => removeOwner(index)}
                  className="focus-ring self-end rounded-md border border-slate-300 p-2 text-slate-500 hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-40"
                  aria-label="Remove owner"
                >
                  <Trash2 size={18} />
                </button>
              </article>
            ))}
          </div>
        </section>

        {error && <p className="mt-4 rounded-md border border-rose-200 bg-rose-50 px-3 py-2 text-sm font-semibold text-rose-700">{error}</p>}
      </form>
    </main>
  );
}

function TextField({ label, value, onChange, type = 'text', required = false }: { label: string; value: string; onChange: (value: string) => void; type?: string; required?: boolean }) {
  return (
    <label className="block">
      <span className="text-sm font-semibold text-slate-700">{label}</span>
      <input
        className="focus-ring mt-1 w-full rounded-md border border-slate-300 px-3 py-2"
        type={type}
        value={value}
        min={type === 'number' ? 0 : undefined}
        max={type === 'number' ? 100 : undefined}
        step={type === 'number' ? '0.01' : undefined}
        onChange={(event) => onChange(event.target.value)}
        required={required}
      />
    </label>
  );
}

function stringValue(value: unknown) {
  return typeof value === 'string' ? value : '';
}

function ownersFromData(value: unknown): OwnerForm[] {
  if (!Array.isArray(value)) {
    return [{ ...emptyOwner }];
  }
  const owners = value.map((item) => {
    const owner = item as Record<string, unknown>;
    return {
      name: stringValue(owner.name),
      ownershipPercent: typeof owner.ownershipPercent === 'number' ? String(owner.ownershipPercent) : '',
      controlType: stringValue(owner.controlType) || 'shares',
    };
  });
  return owners.length ? owners : [{ ...emptyOwner }];
}
