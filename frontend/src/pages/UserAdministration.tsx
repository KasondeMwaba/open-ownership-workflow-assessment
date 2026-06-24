import { CheckCircle2, KeyRound, Plus, Power, Save, SearchX, ShieldCheck, UserCog, XCircle } from 'lucide-react';
import { FormEvent, ReactNode, useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  createPermission,
  createRole,
  createUser,
  listPermissions,
  listRoles,
  listUsers,
  setUserStatus,
  updateRole,
  updateUser,
} from '../api/client';
import ConfirmDialog from '../components/ConfirmDialog';
import EmptyState from '../components/EmptyState';
import PageHeader from '../components/PageHeader';
import PaginationControls from '../components/PaginationControls';
import type {
  AccessRole,
  CreatePermissionPayload,
  CreateRolePayload,
  CreateUserPayload,
  Permission,
  Role,
  UpdateUserPayload,
  User,
} from '../types/domain';

type AdminTab = 'users' | 'roles';
const PAGE_SIZE = 8;
const PERMISSION_PAGE_SIZE = 6;

const emptyUserForm: CreateUserPayload = {
  name: '',
  email: '',
  password: 'password123',
  role: 'requester',
  isActive: true,
};

const emptyRoleForm: CreateRolePayload = {
  name: '',
  description: '',
  permissionIds: [],
};

const emptyPermissionForm: CreatePermissionPayload = {
  name: '',
  description: '',
};

function sectionTitle(section: AdminTab) {
  switch (section) {
    case 'users':
      return 'User Management';
    case 'roles':
      return 'Role Management';
  }
}

export default function UserAdministration({ user, section = 'roles' }: { user: User; section?: AdminTab }) {
  const navigate = useNavigate();
  const [tab, setTab] = useState<AdminTab>(section);
  const [users, setUsers] = useState<User[]>([]);
  const [roles, setRoles] = useState<AccessRole[]>([]);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [userForm, setUserForm] = useState<CreateUserPayload>(emptyUserForm);
  const [roleForm, setRoleForm] = useState<CreateRolePayload>(emptyRoleForm);
  const [permissionForm, setPermissionForm] = useState<CreatePermissionPayload>(emptyPermissionForm);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [editingRole, setEditingRole] = useState<AccessRole | null>(null);
  const [statusTarget, setStatusTarget] = useState<User | null>(null);
  const [roleModalOpen, setRoleModalOpen] = useState(false);
  const [roleStep, setRoleStep] = useState<1 | 2>(1);
  const [roleSearch, setRoleSearch] = useState('');
  const [userPage, setUserPage] = useState(1);
  const [rolePage, setRolePage] = useState(1);
  const [permissionPage, setPermissionPage] = useState(1);
  const [busy, setBusy] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const roleOptions = useMemo(() => {
    const names = roles.map((role) => role.name);
    return names.length ? names : ['requester', 'reviewer', 'admin'];
  }, [roles]);

  useEffect(() => {
    setTab(section);
  }, [section]);

  function selectSection(nextSection: AdminTab) {
    setTab(nextSection);
    navigate(`/admin/${nextSection}`);
  }

  async function refreshAll() {
    const [nextUsers, nextRoles, nextPermissions] = await Promise.all([listUsers(), listRoles(), listPermissions()]);
    setUsers(nextUsers);
    setRoles(nextRoles);
    setPermissions(nextPermissions);
  }

  useEffect(() => {
    if (user.role === 'admin') {
      refreshAll().catch(() => setError('Could not load administration data.'));
    }
  }, [user.role]);

  const filteredRoles = useMemo(() => {
    const query = roleSearch.toLowerCase();
    return roles.filter((role) => {
      const text = `${role.name} ${role.description}`.toLowerCase();
      return text.includes(query);
    });
  }, [roles, roleSearch]);

  const userPages = Math.max(1, Math.ceil(users.length / PAGE_SIZE));
  const rolePages = Math.max(1, Math.ceil(filteredRoles.length / PAGE_SIZE));
  const permissionPages = Math.max(1, Math.ceil(permissions.length / PERMISSION_PAGE_SIZE));
  const pagedUsers = useMemo(() => users.slice((userPage - 1) * PAGE_SIZE, userPage * PAGE_SIZE), [users, userPage]);
  const pagedRoles = useMemo(() => filteredRoles.slice((rolePage - 1) * PAGE_SIZE, rolePage * PAGE_SIZE), [filteredRoles, rolePage]);
  const pagedPermissions = useMemo(
    () => permissions.slice((permissionPage - 1) * PERMISSION_PAGE_SIZE, permissionPage * PERMISSION_PAGE_SIZE),
    [permissions, permissionPage],
  );

  useEffect(() => {
    setUserPage(1);
  }, [users.length]);

  useEffect(() => {
    setRolePage(1);
  }, [roleSearch, roles.length]);

  useEffect(() => {
    setPermissionPage(1);
  }, [permissions.length]);

  if (user.role !== 'admin') {
    return (
      <main className="px-4 py-6 md:px-6 md:py-8">
        <section className="rounded-md border border-slate-200 bg-white p-6 shadow-panel dark:border-slate-800 dark:bg-slate-900">
          <h2 className="text-xl font-black text-ink dark:text-slate-100">Administration</h2>
          <p className="mt-2 text-sm text-slate-500">Only admins can manage users, roles, and permissions.</p>
        </section>
      </main>
    );
  }

  function startEditUser(account: User) {
    setEditingUser(account);
    setUserForm({
      name: account.name,
      email: account.email,
      password: '',
      role: account.role,
      isActive: account.isActive,
    });
    setMessage('');
    setError('');
    selectSection('users');
  }

  function resetUserForm() {
    setEditingUser(null);
    setUserForm({ ...emptyUserForm, role: roleOptions.includes('requester') ? 'requester' : roleOptions[0] });
  }

  function startEditRole(role: AccessRole) {
    setEditingRole(role);
    setRoleForm({
      name: role.name,
      description: role.description,
      permissionIds: role.permissions.map((permission) => permission.id),
    });
    setMessage('');
    setError('');
    setRoleStep(1);
    setRoleModalOpen(true);
    selectSection('roles');
  }

  function resetRoleForm() {
    setEditingRole(null);
    setRoleForm(emptyRoleForm);
    setRoleStep(1);
  }

  function openCreateRole() {
    resetRoleForm();
    setMessage('');
    setError('');
    setRoleModalOpen(true);
  }

  function closeRoleModal() {
    setRoleModalOpen(false);
    resetRoleForm();
  }

  async function submitUser(event: FormEvent) {
    event.preventDefault();
    setBusy(true);
    setError('');
    setMessage('');
    try {
      if (editingUser) {
        const payload: UpdateUserPayload = {
          name: userForm.name,
          email: userForm.email,
          role: userForm.role,
          isActive: userForm.isActive,
        };
        await updateUser(editingUser.id, payload);
        setMessage('Account updated.');
      } else {
        await createUser(userForm);
        setMessage('Account created.');
      }
      resetUserForm();
      await refreshAll();
    } catch {
      setError('Could not save account. Check the email, password, and role.');
    } finally {
      setBusy(false);
    }
  }

  async function submitRole(event: FormEvent) {
    event.preventDefault();
    if (roleStep === 1) {
      if (!roleForm.name.trim()) {
        setError('Role name is required before choosing permissions.');
        return;
      }
      setError('');
      setRoleStep(2);
      return;
    }
    setBusy(true);
    setError('');
    setMessage('');
    try {
      if (editingRole) {
        await updateRole(editingRole.id, roleForm);
        setMessage('Role updated.');
      } else {
        await createRole(roleForm);
        setMessage('Role created.');
      }
      resetRoleForm();
      setRoleModalOpen(false);
      await refreshAll();
    } catch {
      setError('Could not save role. Role names must be unique and permissions must exist.');
    } finally {
      setBusy(false);
    }
  }

  async function submitPermission(event: FormEvent) {
    event.preventDefault();
    setBusy(true);
    setError('');
    setMessage('');
    try {
      await createPermission(permissionForm);
      setPermissionForm(emptyPermissionForm);
      setMessage('Permission created.');
      await refreshAll();
    } catch {
      setError('Could not create permission. Permission names must be unique.');
    } finally {
      setBusy(false);
    }
  }

  async function toggleStatus(account: User) {
    setBusy(true);
    setError('');
    setMessage('');
    try {
      await setUserStatus(account.id, !account.isActive);
      setMessage(account.isActive ? 'Account disabled.' : 'Account enabled.');
      await refreshAll();
    } catch {
      setError('Could not update account status.');
    } finally {
      setBusy(false);
    }
  }

  function togglePermission(permissionId: string) {
    const exists = roleForm.permissionIds.includes(permissionId);
    setRoleForm({
      ...roleForm,
      permissionIds: exists ? roleForm.permissionIds.filter((id) => id !== permissionId) : [...roleForm.permissionIds, permissionId],
    });
  }

  return (
    <main className="px-4 py-6 md:px-6 md:py-8">
      <PageHeader
        eyebrow="Administration"
        title={sectionTitle(tab)}
        description="Create accounts, define roles, and control the permissions each role receives."
        action={
          <div className="inline-flex items-center gap-2 rounded-md border border-emerald-200 bg-emerald-50 px-3 py-2 text-sm font-semibold text-deepgreen">
            <ShieldCheck size={17} />
            Role-based access control
          </div>
        }
      />

      {message && <p className="mt-4 rounded-md border border-emerald-200 bg-emerald-50 px-3 py-2 text-sm font-semibold text-emerald-700">{message}</p>}
      {error && <p className="mt-4 rounded-md border border-rose-200 bg-rose-50 px-3 py-2 text-sm font-semibold text-rose-700">{error}</p>}

      {tab === 'users' && (
        <div className="mt-6 grid gap-6 xl:grid-cols-[420px_1fr]">
          <form onSubmit={submitUser} className="rounded-md border border-slate-200 bg-white p-5 shadow-panel dark:border-slate-800 dark:bg-slate-900">
            <PanelHeading icon={editingUser ? <Save size={20} /> : <Plus size={20} />} title={editingUser ? 'Edit user account' : 'Create user account'} subtitle={editingUser ? editingUser.email : 'Create requesters, reviewers, admins, or custom role users.'} />

            <div className="mt-5 space-y-4">
              <TextField label="Name" value={userForm.name} onChange={(value) => setUserForm({ ...userForm, name: value })} required />
              <TextField label="Email" type="email" value={userForm.email} onChange={(value) => setUserForm({ ...userForm, email: value })} required />
              {!editingUser && <TextField label="Temporary password" value={userForm.password} onChange={(value) => setUserForm({ ...userForm, password: value })} required minLength={8} />}
              <div className="grid grid-cols-2 gap-3">
                <label className="block">
                  <span className="text-sm font-semibold text-slate-700 dark:text-slate-300">Role</span>
                  <select className="focus-ring mt-1 w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm capitalize dark:border-slate-700 dark:bg-slate-950 dark:text-slate-100" value={userForm.role} onChange={(event) => setUserForm({ ...userForm, role: event.target.value as Role })}>
                    {roleOptions.map((role) => (
                      <option key={role} value={role}>
                        {role}
                      </option>
                    ))}
                  </select>
                </label>
                <label className="block">
                  <span className="text-sm font-semibold text-slate-700 dark:text-slate-300">Status</span>
                  <select className="focus-ring mt-1 w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-950 dark:text-slate-100" value={userForm.isActive ? 'active' : 'disabled'} onChange={(event) => setUserForm({ ...userForm, isActive: event.target.value === 'active' })}>
                    <option value="active">Active</option>
                    <option value="disabled">Disabled</option>
                  </select>
                </label>
              </div>
            </div>

            <FormActions busy={busy} editing={Boolean(editingUser)} createLabel="Create account" updateLabel="Save changes" onCancel={resetUserForm} />
          </form>

          <section className="overflow-hidden rounded-md border border-slate-200 bg-white shadow-panel dark:border-slate-800 dark:bg-slate-900">
            <div className="flex items-center justify-between border-b border-slate-100 px-4 py-3 text-sm dark:border-slate-800">
              <p className="font-black text-ink dark:text-slate-100">Accounts</p>
              <span className="rounded bg-slate-100 px-2 py-1 text-xs font-bold text-slate-500 dark:bg-slate-800">{users.length} users</span>
            </div>
            <div className="grid grid-cols-[1fr_130px_120px_170px] border-b border-emerald-100 bg-emerald-50 px-4 py-3 text-xs font-black uppercase tracking-wide text-deepgreen dark:border-emerald-900 dark:bg-emerald-950/70 dark:text-emerald-100 max-lg:hidden">
              <span>User</span>
              <span>Role</span>
              <span>Status</span>
              <span>Actions</span>
            </div>
            {pagedUsers.map((account) => (
              <article key={account.id} className="grid gap-3 border-b border-slate-100 px-4 py-4 hover:bg-emerald-50/40 dark:border-slate-800 dark:hover:bg-slate-800/70 lg:grid-cols-[1fr_130px_120px_170px]">
                <div className="flex items-center gap-3">
                  <div className="flex h-10 w-10 items-center justify-center rounded-md bg-emerald-50 text-deepgreen">
                    <UserCog size={18} />
                  </div>
                  <div>
                    <p className="font-bold text-ink dark:text-slate-100">{account.name}</p>
                    <p className="text-sm text-slate-500">{account.email}</p>
                  </div>
                </div>
                <RoleBadge role={account.role} />
                <StatusPill active={account.isActive} />
                <div className="flex items-center gap-2 self-center">
                  <button onClick={() => startEditUser(account)} className="btn-secondary text-sm">
                    Edit
                  </button>
                  <button disabled={busy || account.id === user.id} onClick={() => setStatusTarget(account)} className="btn-secondary text-sm">
                    <Power size={15} />
                    {account.isActive ? 'Disable' : 'Enable'}
                  </button>
                </div>
              </article>
            ))}
            <PaginationControls page={userPage} totalPages={userPages} totalItems={users.length} onPage={setUserPage} />
          </section>
        </div>
      )}

      {tab === 'roles' && (
        <div className="mt-6 grid gap-6 xl:grid-cols-[1fr_380px]">
          <section className="overflow-hidden rounded-md border border-slate-200 bg-white shadow-panel dark:border-slate-800 dark:bg-slate-900">
            <div className="flex flex-wrap items-center gap-3 border-b border-slate-200 px-4 py-3 dark:border-slate-800">
              <input
                className="focus-ring min-w-[240px] flex-1 rounded-md border border-slate-300 bg-slate-50 px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-950 dark:text-slate-100"
                placeholder="search for roles"
                value={roleSearch}
                onChange={(event) => setRoleSearch(event.target.value)}
              />
              <button className="btn-secondary text-sm">
                Filters
              </button>
              <button onClick={openCreateRole} className="btn-ink ml-auto text-sm">
                <Plus size={15} />
                Add Role
              </button>
            </div>
            <div className="grid grid-cols-[1fr_1fr_110px] border-b border-emerald-100 bg-emerald-50 px-6 py-3 text-xs font-black uppercase tracking-wide text-deepgreen dark:border-emerald-900 dark:bg-emerald-950/70 dark:text-emerald-100 max-md:hidden">
              <span>Name</span>
              <span>Description</span>
              <span>Actions</span>
            </div>
            {pagedRoles.map((role) => (
              <article key={role.id} className="grid gap-3 border-t border-slate-100 px-6 py-4 hover:bg-emerald-50/40 dark:border-slate-800 dark:hover:bg-slate-800/70 md:grid-cols-[1fr_1fr_110px]">
                <p className="self-center font-bold uppercase tracking-wide text-ink dark:text-slate-100">{role.name}</p>
                <p className="self-center text-sm text-slate-500">{role.description || 'No description provided.'}</p>
                <button onClick={() => startEditRole(role)} className="focus-ring justify-self-start rounded-md px-3 py-2 text-sm font-black text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800">
                  ...
                </button>
                <div className="md:col-span-3 flex flex-wrap gap-2">
                  {role.permissions.map((permission) => (
                    <span key={permission.id} className="rounded bg-emerald-50 px-2 py-1 text-xs font-bold text-deepgreen ring-1 ring-emerald-100 dark:bg-emerald-950 dark:text-emerald-100 dark:ring-emerald-900">
                      {permission.name}
                    </span>
                  ))}
                </div>
              </article>
            ))}
            {filteredRoles.length === 0 && <EmptyState icon={SearchX} title="No roles found" description="No role matches the current search. Clear the search or create a new role." />}
            <PaginationControls page={rolePage} totalPages={rolePages} totalItems={filteredRoles.length} onPage={setRolePage} />
          </section>

          <aside className="space-y-4">
            <form onSubmit={submitPermission} className="rounded-md border border-slate-200 bg-white p-5 shadow-panel dark:border-slate-800 dark:bg-slate-900">
              <PanelHeading icon={<KeyRound size={20} />} title="Create permission" subtitle="New permissions become available inside the role wizard." />
              <div className="mt-5 space-y-4">
                <TextField label="Permission name" value={permissionForm.name} onChange={(value) => setPermissionForm({ ...permissionForm, name: value })} required />
                <TextField label="Description" value={permissionForm.description} onChange={(value) => setPermissionForm({ ...permissionForm, description: value })} />
              </div>
              <div className="mt-5">
                <button disabled={busy} className="btn-ink w-full py-2.5">
                  <Plus size={17} />
                  Create permission
                </button>
              </div>
            </form>
            <section className="overflow-hidden rounded-md border border-slate-200 bg-white shadow-panel dark:border-slate-800 dark:bg-slate-900">
              <div className="table-header-band flex items-center justify-between">
                <h3 className="font-black text-ink dark:text-slate-100">Permissions</h3>
                <span className="rounded bg-white px-2 py-1 text-xs font-bold text-deepgreen ring-1 ring-emerald-100 dark:bg-slate-900 dark:text-emerald-100 dark:ring-emerald-900">{permissions.length}</span>
              </div>
              <div className="grid gap-2 p-4">
                {pagedPermissions.map((permission) => (
                  <div key={permission.id} className="rounded-md border border-slate-200 px-3 py-2 hover:border-emerald-200 hover:bg-emerald-50/50 dark:border-slate-800 dark:hover:bg-slate-800/70">
                    <p className="text-sm font-black text-ink dark:text-slate-100">{permission.name}</p>
                    <p className="mt-1 text-xs text-slate-500">{permission.description || 'No description provided.'}</p>
                  </div>
                ))}
              </div>
              <PaginationControls page={permissionPage} totalPages={permissionPages} totalItems={permissions.length} onPage={setPermissionPage} />
            </section>
          </aside>
          {roleModalOpen && (
            <div className="fixed inset-0 z-50 grid place-items-center bg-black/50 px-4 py-6">
              <form onSubmit={submitRole} className="flex max-h-[88vh] w-full max-w-3xl flex-col overflow-hidden rounded-md bg-white shadow-2xl dark:bg-slate-900">
                <div className="border-b border-slate-200 p-6 dark:border-slate-800">
                  <h3 className="text-2xl font-black text-ink dark:text-slate-100">{editingRole ? 'Edit Role' : 'Create New Role'}</h3>
                  <div className="mt-5 flex items-center gap-2">
                    <WizardStep active={roleStep === 1} complete={roleStep > 1} label="Role Details" number="1" />
                    <div className="h-px flex-1 bg-slate-200 dark:bg-slate-800" />
                    <WizardStep active={roleStep === 2} complete={false} label="Permissions" number="2" />
                  </div>
                </div>
                <div className="flex-1 overflow-y-auto p-6">
                  {roleStep === 1 ? (
                    <div className="space-y-5">
                      <TextField label="Role Name" value={roleForm.name} onChange={(value) => setRoleForm({ ...roleForm, name: value })} required />
                      <TextField label="Description" value={roleForm.description} onChange={(value) => setRoleForm({ ...roleForm, description: value })} />
                    </div>
                  ) : (
                    <div>
                      <div className="flex items-center justify-between gap-3">
                        <h4 className="font-black text-ink dark:text-slate-100">Select Permissions</h4>
                        <span className="rounded bg-emerald-50 px-2 py-1 text-xs font-bold text-deepgreen">{roleForm.permissionIds.length} selected</span>
                      </div>
                      <div className="mt-4 grid gap-3 md:grid-cols-2">
                        {permissions.map((permission) => (
                          <label key={permission.id} className="flex min-h-[82px] items-start gap-3 rounded-md border border-slate-200 px-3 py-3 text-sm hover:bg-slate-50 dark:border-slate-800 dark:hover:bg-slate-800/70">
                            <input className="mt-1 h-4 w-4 rounded border-slate-300 text-deepgreen focus:ring-deepgreen" type="checkbox" checked={roleForm.permissionIds.includes(permission.id)} onChange={() => togglePermission(permission.id)} />
                            <span>
                              <span className="block font-bold uppercase text-ink dark:text-slate-100">{permission.name}</span>
                              <span className="text-slate-500">{permission.description || 'No description provided.'}</span>
                            </span>
                          </label>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
                <div className="flex justify-between border-t border-slate-200 p-4 dark:border-slate-800">
                  {roleStep === 1 ? (
                    <button type="button" onClick={closeRoleModal} className="btn-secondary px-4 py-2.5">
                      Cancel
                    </button>
                  ) : (
                    <button type="button" onClick={() => setRoleStep(1)} className="btn-secondary px-4 py-2.5">
                      Back
                    </button>
                  )}
                  {roleStep === 1 ? (
                    <button type="submit" className="btn-ink px-5 py-2.5">
                      Next
                    </button>
                  ) : (
                    <button type="submit" disabled={busy} className="btn-ink px-5 py-2.5">
                      {editingRole ? 'Update Role' : 'Create Role'}
                    </button>
                  )}
                </div>
              </form>
            </div>
          )}
        </div>
      )}
      {statusTarget && (
        <ConfirmDialog
          title={statusTarget.isActive ? 'Disable account?' : 'Enable account?'}
          description={`${statusTarget.name} will ${statusTarget.isActive ? 'no longer be able to sign in' : 'be able to sign in again'}. This action will be recorded in the audit trail.`}
          confirmLabel={statusTarget.isActive ? 'Disable account' : 'Enable account'}
          tone={statusTarget.isActive ? 'danger' : 'default'}
          onCancel={() => setStatusTarget(null)}
          onConfirm={() => {
            const account = statusTarget;
            setStatusTarget(null);
            toggleStatus(account);
          }}
        />
      )}
    </main>
  );
}

function WizardStep({ active, complete, label, number }: { active: boolean; complete: boolean; label: string; number: string }) {
  return (
    <div className={`flex items-center gap-2 text-sm font-bold ${active || complete ? 'text-deepgreen' : 'text-slate-400'}`}>
      <span className={`flex h-7 w-7 items-center justify-center rounded-full text-xs ${active || complete ? 'bg-emerald-100 text-deepgreen' : 'bg-slate-100 text-slate-400 dark:bg-slate-800'}`}>
        {complete ? <CheckCircle2 size={14} /> : number}
      </span>
      <span>{label}</span>
    </div>
  );
}

function PanelHeading({ icon, title, subtitle }: { icon: ReactNode; title: string; subtitle: string }) {
  return (
    <div className="flex items-center gap-3">
      <div className="flex h-11 w-11 items-center justify-center rounded-md bg-emerald-50 text-deepgreen">{icon}</div>
      <div>
        <h3 className="text-lg font-black text-ink dark:text-slate-100">{title}</h3>
        <p className="text-xs text-slate-500">{subtitle}</p>
      </div>
    </div>
  );
}

function TextField({ label, value, onChange, type = 'text', required = false, minLength }: { label: string; value: string; onChange: (value: string) => void; type?: string; required?: boolean; minLength?: number }) {
  return (
    <label className="block">
      <span className="text-sm font-semibold text-slate-700 dark:text-slate-300">{label}</span>
      <input className="focus-ring mt-1 w-full rounded-md border border-slate-300 px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-950 dark:text-slate-100" type={type} value={value} onChange={(event) => onChange(event.target.value)} required={required} minLength={minLength} />
    </label>
  );
}

function FormActions({ busy, editing, createLabel, updateLabel, onCancel }: { busy: boolean; editing: boolean; createLabel: string; updateLabel: string; onCancel: () => void }) {
  return (
    <div className="mt-5 flex gap-2">
      <button disabled={busy} className="btn-ink flex-1 py-2.5">
        {editing ? <Save size={17} /> : <Plus size={17} />}
        {editing ? updateLabel : createLabel}
      </button>
      {editing && (
        <button type="button" onClick={onCancel} className="btn-secondary px-4 py-2.5">
          Cancel
        </button>
      )}
    </div>
  );
}

function StatusPill({ active }: { active: boolean }) {
  return (
    <div className="self-center">
      <span className={`inline-flex items-center gap-1 rounded px-2 py-1 text-xs font-bold ring-1 ${active ? 'bg-emerald-50 text-emerald-700 ring-emerald-200' : 'bg-rose-50 text-rose-700 ring-rose-200'}`}>
        {active ? <CheckCircle2 size={13} /> : <XCircle size={13} />}
        {active ? 'Active' : 'Disabled'}
      </span>
    </div>
  );
}

function RoleBadge({ role }: { role: Role }) {
  return (
    <div className="self-center">
      <span className="inline-flex max-w-full rounded bg-slate-100 px-2 py-1 text-xs font-black uppercase tracking-wide text-slate-600 ring-1 ring-slate-200 dark:bg-slate-800 dark:text-slate-200 dark:ring-slate-700">
        {role}
      </span>
    </div>
  );
}
