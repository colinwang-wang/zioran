import { useEffect, useState } from 'react'
import { Card, Form, Input, InputNumber, Button, message, Spin } from 'antd'
import { getSettings, updateSettings } from '@/api'
import type { Settings } from '@/types'

export default function SettingsPage() {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    getSettings().then(res => { form.setFieldsValue(res.data); setLoading(false) })
  }, [])

  const handleSave = async () => {
    const values = await form.validateFields()
    setSaving(true)
    try {
      await updateSettings(values as Partial<Settings>)
      message.success('保存成功')
    } finally { setSaving(false) }
  }

  if (loading) return <Spin />

  return (
    <Card title="系统设置">
      <Form form={form} labelCol={{ span: 4 }} wrapperCol={{ span: 12 }}>
        <Form.Item label="站点名称" name="siteName" rules={[{ required: true }]}><Input /></Form.Item>
        <Form.Item label="站点描述" name="siteDescription"><Input.TextArea rows={2} /></Form.Item>
        <Form.Item label="联系电话" name="contactPhone"><Input /></Form.Item>
        <Form.Item label="联系邮箱" name="contactEmail"><Input /></Form.Item>
        <Form.Item label="VIP月费(元)" name="vipMonthlyPrice"><InputNumber min={0} precision={2} /></Form.Item>
        <Form.Item label="VIP年费(元)" name="vipYearlyPrice"><InputNumber min={0} precision={2} /></Form.Item>
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
  )
}
