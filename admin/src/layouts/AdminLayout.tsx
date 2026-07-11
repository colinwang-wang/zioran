import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { ProLayout } from '@ant-design/pro-components'
import { DashboardOutlined, BookOutlined, AppstoreOutlined, TagsOutlined, UserOutlined, ShoppingCartOutlined, MessageOutlined, CommentOutlined, SettingOutlined, BarChartOutlined, CustomerServiceOutlined, TeamOutlined, ToolOutlined, SafetyOutlined } from '@ant-design/icons'
import { Button, Dropdown } from 'antd'
import { LogoutOutlined } from '@ant-design/icons'

const allRoutes = [
  { path: '/dashboard', name: '仪表盘', icon: <DashboardOutlined />, perm: 'dashboard' },
  { path: '/courses', name: '课程管理', icon: <BookOutlined />, perm: 'courses' },
  { path: '/categories', name: '分类管理', icon: <AppstoreOutlined />, perm: 'categories' },
  { path: '/tags', name: '标签管理', icon: <TagsOutlined />, perm: 'tags' },
  { path: '/users', name: '用户管理', icon: <UserOutlined />, perm: 'users' },
  { path: '/orders', name: '订单管理', icon: <ShoppingCartOutlined />, perm: 'orders' },
  { path: '/guestbook', name: '留言管理', icon: <MessageOutlined />, perm: 'guestbook' },
  { path: '/comments', name: '评论管理', icon: <CommentOutlined />, perm: 'comments' },
  { path: '/config', name: '首页配置', icon: <SettingOutlined />, perm: 'home_config' },
  { path: '/data', name: '数据看板', icon: <BarChartOutlined />, perm: 'data' },
  { path: '/tickets', name: '工单管理', icon: <CustomerServiceOutlined />, perm: 'tickets' },
  { path: '/settings', name: '系统设置', icon: <ToolOutlined />, perm: 'settings' },
  { path: '/admins', name: '管理员', icon: <TeamOutlined />, perm: 'admins' },
  { path: '/permissions', name: '权限配置', icon: <SafetyOutlined />, perm: 'admins' },
]

function getMenuRoutes() {
  const role = localStorage.getItem('admin_role') || 'admin'
  if (role === 'super_admin') return { routes: allRoutes }
  let perms: string[] = []
  try { perms = JSON.parse(localStorage.getItem('admin_permissions') || '[]') } catch { /* */ }
  const routes = allRoutes.filter(r => perms.includes(r.perm))
  return { routes }
}

export default function AdminLayout() {
  const navigate = useNavigate()
  const location = useLocation()

  const handleLogout = () => {
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_role')
    localStorage.removeItem('admin_permissions')
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
