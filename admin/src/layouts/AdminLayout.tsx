import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { ProLayout } from '@ant-design/pro-components'
import { DashboardOutlined, BookOutlined, AppstoreOutlined, TagsOutlined, UserOutlined, ShoppingCartOutlined, MessageOutlined, CommentOutlined, SettingOutlined, BarChartOutlined } from '@ant-design/icons'
import { Button, Dropdown } from 'antd'
import { LogoutOutlined } from '@ant-design/icons'

const menuRoutes = {
  routes: [
    { path: '/dashboard', name: '仪表盘', icon: <DashboardOutlined /> },
    { path: '/courses', name: '课程管理', icon: <BookOutlined /> },
    { path: '/categories', name: '分类管理', icon: <AppstoreOutlined /> },
    { path: '/tags', name: '标签管理', icon: <TagsOutlined /> },
    { path: '/users', name: '用户管理', icon: <UserOutlined /> },
    { path: '/orders', name: '订单管理', icon: <ShoppingCartOutlined /> },
    { path: '/guestbook', name: '留言管理', icon: <MessageOutlined /> },
    { path: '/comments', name: '评论管理', icon: <CommentOutlined /> },
    { path: '/config', name: '首页配置', icon: <SettingOutlined /> },
    { path: '/data', name: '数据看板', icon: <BarChartOutlined /> },
  ],
}

export default function AdminLayout() {
  const navigate = useNavigate()
  const location = useLocation()

  const handleLogout = () => {
    localStorage.removeItem('admin_token')
    navigate('/login', { replace: true })
  }

  return (
    <ProLayout
      title="知猿管理后台"
      logo={null}
      layout="mix"
      fixSiderbar
      fixedHeader
      route={menuRoutes}
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
