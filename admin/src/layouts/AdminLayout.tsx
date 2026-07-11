import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { ProLayout } from '@ant-design/pro-components'
import { DashboardOutlined, BookOutlined, AppstoreOutlined, TagsOutlined, UserOutlined, ShoppingCartOutlined, MessageOutlined, CommentOutlined, SettingOutlined, BarChartOutlined, CustomerServiceOutlined, TeamOutlined, ToolOutlined } from '@ant-design/icons'
import { Button, Dropdown } from 'antd'
import { LogoutOutlined } from '@ant-design/icons'

const allRoutes = [
  { path: '/dashboard', name: '仪表盘', icon: <DashboardOutlined />, roles: ['*'] },
  { path: '/courses', name: '课程管理', icon: <BookOutlined />, roles: ['super_admin', 'admin', 'operator'] },
  { path: '/categories', name: '分类管理', icon: <AppstoreOutlined />, roles: ['super_admin', 'admin', 'operator'] },
  { path: '/tags', name: '标签管理', icon: <TagsOutlined />, roles: ['super_admin', 'admin', 'operator'] },
  { path: '/users', name: '用户管理', icon: <UserOutlined />, roles: ['super_admin', 'admin'] },
  { path: '/orders', name: '订单管理', icon: <ShoppingCartOutlined />, roles: ['super_admin', 'admin', 'operator'] },
  { path: '/guestbook', name: '留言管理', icon: <MessageOutlined />, roles: ['super_admin', 'admin', 'operator', 'support'] },
  { path: '/comments', name: '评论管理', icon: <CommentOutlined />, roles: ['super_admin', 'admin', 'operator', 'support'] },
  { path: '/config', name: '首页配置', icon: <SettingOutlined />, roles: ['super_admin', 'admin', 'operator'] },
  { path: '/data', name: '数据看板', icon: <BarChartOutlined />, roles: ['super_admin', 'admin', 'operator'] },
  { path: '/tickets', name: '工单管理', icon: <CustomerServiceOutlined />, roles: ['super_admin', 'admin', 'operator', 'support'] },
  { path: '/settings', name: '系统设置', icon: <ToolOutlined />, roles: ['super_admin'] },
  { path: '/admins', name: '管理员', icon: <TeamOutlined />, roles: ['super_admin'] },
]

function getMenuRoutes() {
  const role = localStorage.getItem('admin_role') || 'admin'
  const routes = allRoutes.filter(r => r.roles.includes('*') || r.roles.includes(role))
  return { routes }
}

export default function AdminLayout() {
  const navigate = useNavigate()
  const location = useLocation()

  const handleLogout = () => {
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_role')
    navigate('/login', { replace: true })
  }

  return (
    <ProLayout
      title="知猿管理后台"
      logo={null}
      layout="mix"
      fixSiderbar
      fixedHeader
      route={getMenuRoutes()}
      location={{ pathname: location.pathname }}
      menuItemRender={(item, dom) => (
        <div onClick={() => item.path && navigate(item.path)}>{dom}</div>
      )}
      actionsRender={() => [
        <Dropdown key="user" menu={{ items: [{ key: 'logout', icon: <LogoutOutlined />, label: '退出登录', onClick: handleLogout }] }}>
          <Button type="text">管理员</Button>
        </Dropdown>,
      ]}
    >
      <Outlet />
    </ProLayout>
  )
}
