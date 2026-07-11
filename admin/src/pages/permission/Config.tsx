import { useState, useEffect } from 'react'
import { Card, Table, Checkbox, Button, message, Spin, Tag } from 'antd'
import { getAllPermissions, getRolePermissions, updateRolePermissions } from '@/api'

const editableRoles = [
  { key: 'admin', label: '管理员' },
  { key: 'operator', label: '运营' },
  { key: 'support', label: '客服' },
]

interface PermDef {
  key: string
  label: string
}

export default function PermissionConfig() {
  const [allPerms, setAllPerms] = useState<PermDef[]>([])
  const [rolePerms, setRolePerms] = useState<Record<string, string[]>>({})
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState<string | null>(null)

  const fetchData = async () => {
    setLoading(true)
    try {
      const permRes = await getAllPermissions()
      setAllPerms(permRes.data.permissions || [])
      const permsMap: Record<string, string[]> = {}
      for (const role of editableRoles) {
        const res = await getRolePermissions(role.key)
        permsMap[role.key] = res.data.permissions || []
      }
      setRolePerms(permsMap)
    } catch {
      message.error('加载权限数据失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [])

  const handleToggle = (role: string, permKey: string, checked: boolean) => {
    setRolePerms(prev => {
      const current = prev[role] || []
      const next = checked ? [...current, permKey] : current.filter(k => k !== permKey)
      return { ...prev, [role]: next }
    })
  }

  const handleSave = async (role: string) => {
    setSaving(role)
    try {
      await updateRolePermissions(role, rolePerms[role] || [])
      message.success('保存成功')
    } catch {
      message.error('保存失败')
    } finally {
      setSaving(null)
    }
  }

  if (loading) return <Spin style={{ display: 'flex', justifyContent: 'center', paddingTop: 100 }} />

  const columns = [
    { title: '功能模块', dataIndex: 'label', width: 150, render: (label: string, record: PermDef) => <span>{label} <Tag style={{ fontSize: 10 }}>{record.key}</Tag></span> },
    ...editableRoles.map(role => ({
      title: () => (
        <div style={{ textAlign: 'center' }}>
          <div>{role.label}</div>
          <Button size="small" type="primary" loading={saving === role.key} onClick={() => handleSave(role.key)} style={{ marginTop: 4 }}>保存</Button>
        </div>
      ),
      width: 120,
      align: 'center' as const,
      render: (_: unknown, record: PermDef) => (
        <Checkbox
          checked={(rolePerms[role.key] || []).includes(record.key)}
          onChange={e => handleToggle(role.key, record.key, e.target.checked)}
        />
      ),
    })),
  ]

  return (
    <Card title="角色权限配置" extra={<Tag color="orange">超级管理员拥有全部权限，无需配置</Tag>}>
      <p style={{ color: '#999', marginBottom: 16 }}>勾选每个角色可以访问的功能模块。修改后点击对应角色列的"保存"按钮生效（约30秒后后端缓存刷新）。</p>
      <Table
        rowKey="key"
        dataSource={allPerms}
        columns={columns}
        pagination={false}
        size="middle"
        bordered
      />
    </Card>
  )
}
