'use client';

import { useState } from 'react';
import { changePassword } from '@/lib/services';

export default function SettingsPage() {
  const [oldPwd, setOldPwd] = useState('');
  const [newPwd, setNewPwd] = useState('');
  const [loading, setLoading] = useState(false);

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
    <div>
      <h2 className="text-lg font-bold mb-4">账号设置</h2>
      <form onSubmit={handleSubmit} className="max-w-sm space-y-4">
        <input type="password" value={oldPwd} onChange={(e) => setOldPwd(e.target.value)} placeholder="当前密码" className="w-full px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
        <input type="password" value={newPwd} onChange={(e) => setNewPwd(e.target.value)} placeholder="新密码（至少6位）" className="w-full px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
        <button type="submit" disabled={loading} className="px-6 py-3 bg-primary text-white text-sm font-bold rounded-card disabled:opacity-50">修改密码</button>
      </form>
    </div>
  );
}
