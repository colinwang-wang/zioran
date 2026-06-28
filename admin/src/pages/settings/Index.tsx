import { useEffect, useState } from 'react'
import { Button, Card, Form, Input, InputNumber, message, Space, Spin, Switch } from 'antd'
import { getSettings, getVipPackagesAdmin, updateSettings, updateVipPackage } from '@/api'
import type { Settings, VipPackageConfig } from '@/types'

export default function SettingsPage() {
  const [form] = Form.useForm()
  const [vipForm] = Form.useForm()
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [vipSaving, setVipSaving] = useState(false)
  const [vipPackages, setVipPackages] = useState<VipPackageConfig[]>([])

  useEffect(() => {
    Promise.all([getSettings(), getVipPackagesAdmin()]).then(([settingsRes, vipRes]) => {
      form.setFieldsValue(settingsRes.data)
      setVipPackages(vipRes.data)
      vipForm.setFieldsValue({ packages: vipRes.data })
      setLoading(false)
    })
  }, [])

  const handleSave = async () => {
    const values = await form.validateFields()
    setSaving(true)
    try {
      await updateSettings(values as Partial<Settings>)
      message.success('保存成功')
    } finally { setSaving(false) }
  }

  const handleSaveVip = async () => {
    const values = await vipForm.validateFields()
    setVipSaving(true)
    try {
      await Promise.all((values.packages || []).map((pkg: VipPackageConfig, index: number) =>
        updateVipPackage(vipPackages[index].id, { ...vipPackages[index], ...pkg })
      ))
      message.success('VIP套餐已保存')
    } finally { setVipSaving(false) }
  }

  if (loading) return <Spin />

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Card title="系统设置">
        <Form form={form} labelCol={{ span: 4 }} wrapperCol={{ span: 12 }}>
          <Form.Item label="站点名称" name="siteName" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item label="站点描述" name="siteDescription"><Input.TextArea rows={2} /></Form.Item>
          <Form.Item label="联系电话" name="contactPhone"><Input /></Form.Item>
          <Form.Item label="联系邮箱" name="contactEmail"><Input /></Form.Item>
          <Form.Item label="最低提现(元)" name="withdrawMinAmount"><InputNumber min={0} precision={2} /></Form.Item>
          <Form.Item label="佣金比例(%)" name="commissionRate"><InputNumber min={0} max={100} /></Form.Item>
          <Form.Item label="充值比例" name="coinRechargeRatio" rules={[{ required: true, message: '请输入充值比例' }]} extra="填写 1 元可兑换多少金币">
            <InputNumber min={1} precision={0} addonBefore="1 元 =" addonAfter="金币" />
          </Form.Item>
          <Form.Item label="充值档位(元)" name="coinRechargeAmounts" extra="多个金额用英文逗号分隔，例如：10,50,100,200">
            <Input placeholder="10,50,100,200,500,1000" />
          </Form.Item>
          <Form.Item wrapperCol={{ offset: 4 }}><Button type="primary" loading={saving} onClick={handleSave}>保存设置</Button></Form.Item>
        </Form>
      </Card>

      <Card title="VIP套餐">
        <Form form={vipForm} labelCol={{ span: 4 }} wrapperCol={{ span: 12 }}>
          {vipPackages.map((pkg, index) => (
            <div key={pkg.id} style={{ borderBottom: index === vipPackages.length - 1 ? 'none' : '1px solid #f0f0f0', marginBottom: 24, paddingBottom: 8 }}>
              <Form.Item label="套餐名称" name={['packages', index, 'name']} rules={[{ required: true, message: '请输入套餐名称' }]}><Input /></Form.Item>
              <Form.Item label="现价(金币)" name={['packages', index, 'price']} rules={[{ required: true, message: '请输入现价' }]}><InputNumber min={0} precision={0} /></Form.Item>
              <Form.Item label="原价(金币)" name={['packages', index, 'originalPrice']}><InputNumber min={0} precision={0} /></Form.Item>
              <Form.Item label="权益说明" name={['packages', index, 'benefits']} extra="前台会按现有格式读取，可填写 JSON 字符串或文案。">
                <Input.TextArea rows={3} />
              </Form.Item>
              <Form.Item label="排序" name={['packages', index, 'sortOrder']}><InputNumber min={0} precision={0} /></Form.Item>
              <Form.Item label="启用" name={['packages', index, 'isActive']} valuePropName="checked"><Switch /></Form.Item>
            </div>
          ))}
          <Form.Item wrapperCol={{ offset: 4 }}><Button type="primary" loading={vipSaving} onClick={handleSaveVip}>保存VIP套餐</Button></Form.Item>
        </Form>
      </Card>
    </Space>
  )
}
