import { lazy, Suspense } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { Spin } from 'antd'
import AdminLayout from './layouts/AdminLayout'

const Login = lazy(() => import('./pages/Login'))
const Dashboard = lazy(() => import('./pages/Dashboard'))
const CourseList = lazy(() => import('./pages/course/List'))
const CourseForm = lazy(() => import('./pages/course/Form'))
const CategoryList = lazy(() => import('./pages/category/List'))
const TagList = lazy(() => import('./pages/tag/List'))
const UserList = lazy(() => import('./pages/user/List'))
const UserDetail = lazy(() => import('./pages/user/Detail'))
const OrderList = lazy(() => import('./pages/order/List'))
const OrderDetail = lazy(() => import('./pages/order/Detail'))
const GuestbookList = lazy(() => import('./pages/guestbook/List'))
const CommentList = lazy(() => import('./pages/comment/List'))
const HomeConfig = lazy(() => import('./pages/config/HomeConfig'))
const DataBoard = lazy(() => import('./pages/data/Board'))
const TicketList = lazy(() => import('./pages/ticket/List'))
const TicketDetail = lazy(() => import('./pages/ticket/Detail'))
const SettingsPage = lazy(() => import('./pages/settings/Index'))
const AdminList = lazy(() => import('./pages/admin/List'))
const PermissionConfig = lazy(() => import('./pages/permission/Config'))

const Loading = <div style={{ display: 'flex', justifyContent: 'center', paddingTop: 200 }}><Spin size="large" /></div>

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const token = localStorage.getItem('admin_token')
  if (!token) return <Navigate to="/login" replace />
  return <>{children}</>
}

export default function App() {
  return (
    <Suspense fallback={Loading}>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/" element={<PrivateRoute><AdminLayout /></PrivateRoute>}>
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="courses" element={<CourseList />} />
          <Route path="courses/create" element={<CourseForm />} />
          <Route path="courses/:id/edit" element={<CourseForm />} />
          <Route path="categories" element={<CategoryList />} />
          <Route path="tags" element={<TagList />} />
          <Route path="users" element={<UserList />} />
          <Route path="users/:id" element={<UserDetail />} />
          <Route path="orders" element={<OrderList />} />
          <Route path="orders/:id" element={<OrderDetail />} />
          <Route path="guestbook" element={<GuestbookList />} />
          <Route path="comments" element={<CommentList />} />
          <Route path="config" element={<HomeConfig />} />
          <Route path="data" element={<DataBoard />} />
          <Route path="tickets" element={<TicketList />} />
          <Route path="tickets/:id" element={<TicketDetail />} />
          <Route path="settings" element={<SettingsPage />} />
          <Route path="admins" element={<AdminList />} />
          <Route path="permissions" element={<PermissionConfig />} />
        </Route>
      </Routes>
    </Suspense>
  )
}
