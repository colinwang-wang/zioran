'use client';

import { useEffect, useState } from 'react';
import { changePassword, getProfile, updateProfile } from '@/lib/services';
import { useAuth } from '@/contexts/AuthContext';

export default function SettingsPage() {
  const { user, setAuth, token } = useAuth();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [oldPwd, setOldPwd] = useState('');
  const [newPwd, setNewPwd] = useState('');
  const [loading, setLoading] = useState(false);
  const [profileLoading, setProfileLoading] = useState(false);

  useEffect(() => {
    getProfile().then((profile) => {
      setUsername(profile.username || '');
      setEmail(profile.email || '');
    }).catch(() => {});
  }, []);

  const handleProfileSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setProfileLoading(true);
    try {
      const profile = await updateProfile({ username, email });
      setUsername(profile.username || '');
      setEmail(profile.email || '');
      // 即时更新头像处名称
      if (user && token) {
        setAuth(token, { ...user, username: profile.username || user.username });
      }
      alert('资料已保存');
    } catch { alert('保存失败'); }
    setProfileLoading(false);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!oldPwd || !newPwd) return;
    setLoading(true);
    try {
      await changePassword({ old_password: oldPwd, new_password: newPwd });
      alert('修改成功');
      setOldPwd(''); setNewPwd('');
    } catch { alert('修改失败'); }
    setLoading(false);
  };

  return (
    <div className="space-y-8">
      <section>
        <h2 className="text-lg font-bold mb-4">账号资料</h2>
        <form onSubmit={handleProfileSubmit} className="max-w-sm space-y-4">
          <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} placeholder="用户名" className="w-full px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
          <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="邮箱（用于找回密码）" className="w-full px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
          <button type="submit" disabled={profileLoading} className="px-6 py-3 bg-primary text-white text-sm font-bold rounded-card disabled:opacity-50">保存资料</button>
        </form>
      </section>

      <section>
        <h2 className="text-lg font-bold mb-4">修改密码</h2>
        <p className="text-xs text-mute mb-4">邮箱和密码修改只对邮箱注册的用户有效，微信登录用户无需设置。</p>
      <form onSubmit={handleSubmit} className="max-w-sm space-y-4">
        <input type="password" value={oldPwd} onChange={(e) => setOldPwd(e.target.value)} placeholder="当前密码" className="w-full px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
        <input type="password" value={newPwd} onChange={(e) => setNewPwd(e.target.value)} placeholder="新密码（至少6位）" className="w-full px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
        <button type="submit" disabled={loading} className="px-6 py-3 bg-primary text-white text-sm font-bold rounded-card disabled:opacity-50">修改密码</button>
      </form>
      </section>
    </div>
  );
}
