import { useEffect, useState } from 'react';
import { Navigate, Route, Routes, useNavigate } from 'react-router-dom';
import { clearToken, logout as auditLogout, me } from './api/client';
import PortalSidebar, { MobilePortalNav } from './components/PortalSidebar';
import TopHeader from './components/TopHeader';
import AuditTrail from './pages/AuditTrail';
import Dashboard from './pages/Dashboard';
import Login from './pages/Login';
import UserAdministration from './pages/UserAdministration';
import SubmissionDetail from './pages/SubmissionDetail';
import SubmissionEditor from './pages/SubmissionEditor';
import type { User } from './types/domain';

export default function App() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    me()
      .then(setUser)
      .catch(() => clearToken())
      .finally(() => setLoading(false));
  }, []);

  async function logout() {
    try {
      await auditLogout();
    } catch {
      // The local session should still end even if audit recording fails.
    }
    clearToken();
    setUser(null);
    navigate('/login');
  }

  if (loading) {
    return <div className="grid min-h-screen place-items-center text-sm text-slate-500">Loading workflow...</div>;
  }

  if (!user) {
    return (
      <Routes>
        <Route path="/login" element={<Login onLogin={setUser} />} />
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    );
  }

  return (
    <div className="flex min-h-screen bg-paper dark:bg-slate-950">
      <PortalSidebar user={user} />
      <div className="flex min-w-0 flex-1 flex-col">
        <MobilePortalNav user={user} />
        <TopHeader user={user} onLogout={logout} />
        <div className="flex-1 overflow-auto">
          <Routes>
            <Route path="/" element={<Dashboard user={user} />} />
            <Route path="/submissions/new" element={<SubmissionEditor user={user} />} />
            <Route path="/submissions/:id/edit" element={<SubmissionEditor user={user} />} />
            <Route path="/submissions/:id" element={<SubmissionDetail user={user} />} />
            <Route path="/audit" element={<Navigate to="/audit/submission" replace />} />
            <Route path="/audit/activity" element={<AuditTrail user={user} mode="activity" />} />
            <Route path="/audit/session" element={<AuditTrail user={user} mode="session" />} />
            <Route path="/audit/system" element={<AuditTrail user={user} mode="system" />} />
            <Route path="/audit/submission" element={<AuditTrail user={user} mode="submission" />} />
            <Route path="/admin" element={<Navigate to="/admin/roles" replace />} />
            <Route path="/admin/permissions" element={<Navigate to="/admin/roles" replace />} />
            <Route path="/admin/roles" element={<UserAdministration user={user} section="roles" />} />
            <Route path="/admin/users" element={<UserAdministration user={user} section="users" />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </div>
      </div>
    </div>
  );
}
