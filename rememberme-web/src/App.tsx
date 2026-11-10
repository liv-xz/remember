import { lazy, Suspense } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { Spin } from 'antd'
import ProtectedRoute from './components/layout/ProtectedRoute'
import AppLayout from './components/layout/AppLayout'
import LoginPage from './pages/auth/LoginPage'
import RegisterPage from './pages/auth/RegisterPage'
import ShareJoinPage from './pages/grave/ShareJoinPage'
import DashboardPage from './pages/dashboard/DashboardPage'

const GraveDetailPage = lazy(() => import('./pages/grave/GraveDetailPage'))
const MemberPage = lazy(() => import('./pages/members/MemberPage'))
const ProfilePage = lazy(() => import('./pages/profile/ProfilePage'))

function PageLoader() {
  return <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', padding: 80 }}><Spin size="large" /></div>
}

export default function App() {
  return (
    <Suspense fallback={<PageLoader />}>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/join/:shareCode" element={<ShareJoinPage />} />
        <Route element={<ProtectedRoute />}>
          <Route element={<AppLayout />}>
            <Route path="/" element={<DashboardPage />} />
            <Route path="/grave/:id" element={<GraveDetailPage />} />
            <Route path="/grave/:id/members" element={<MemberPage />} />
            <Route path="/profile" element={<ProfilePage />} />
          </Route>
        </Route>
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Suspense>
  )
}
