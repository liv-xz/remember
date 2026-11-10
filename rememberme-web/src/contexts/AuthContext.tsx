import { createContext, useContext, useState, useCallback, useEffect, type ReactNode } from 'react'
import { useNavigate } from 'react-router-dom'
import * as authApi from '../api/auth'
import * as userApi from '../api/user'
import { getAccessToken, getRefreshToken, setTokens, clearTokens } from '../utils/token'
import type { UserInfo } from '../types/auth'
import type { UserProfileResp } from '../types/user'

interface AuthState {
  user: UserInfo | null
  profile: UserProfileResp | null
  isAuthenticated: boolean
  loading: boolean
  login: (username: string, password: string, redirect?: string) => Promise<void>
  register: (username: string, email: string, password: string) => Promise<void>
  logout: () => Promise<void>
  updateProfile: (profile: UserProfileResp) => void
}

const AuthContext = createContext<AuthState | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<UserInfo | null>(null)
  const [profile, setProfile] = useState<UserProfileResp | null>(null)
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()

  useEffect(() => {
    const token = getAccessToken()
    if (!token) {
      setLoading(false)
      return
    }
    userApi.getProfile()
      .then((res) => {
        if (res.code === 0 && res.data) {
          setProfile(res.data)
          setUser({
            id: res.data.id,
            username: res.data.username,
            nickname: res.data.nickname,
            avatar: res.data.avatar,
            email: res.data.email,
          })
        } else {
          clearTokens()
        }
      })
      .catch(() => {
        clearTokens()
      })
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    const handler = () => {
      clearTokens()
      setUser(null)
      setProfile(null)
      navigate('/login')
    }
    window.addEventListener('auth:unauthorized', handler)
    return () => window.removeEventListener('auth:unauthorized', handler)
  }, [navigate])

  const login = useCallback(async (username: string, password: string, redirect?: string) => {
    const res = await authApi.login({ username, password })
    if (res.code !== 0 || !res.data) {
      throw new Error(res.message || '登录失败')
    }
    const { access_token, refresh_token, user: userInfo } = res.data
    setTokens(access_token, refresh_token)
    setUser(userInfo)
    navigate(redirect || '/')
  }, [navigate])

  const register = useCallback(async (username: string, email: string, password: string) => {
    const res = await authApi.register({ username, email, password })
    if (res.code !== 0) {
      throw new Error(res.message || '注册失败')
    }
    navigate('/login')
  }, [navigate])

  const logout = useCallback(async () => {
    const refreshToken = getRefreshToken()
    if (refreshToken) {
      try { await authApi.logout({ refresh_token: refreshToken }) } catch { /* ignore */ }
    }
    clearTokens()
    setUser(null)
    setProfile(null)
    navigate('/login')
  }, [navigate])

  const updateProfile = useCallback((p: UserProfileResp) => {
    setProfile(p)
    setUser({
      id: p.id,
      username: p.username,
      nickname: p.nickname,
      avatar: p.avatar,
      email: p.email,
    })
  }, [])

  return (
    <AuthContext.Provider value={{
      user, profile, isAuthenticated: !!user, loading,
      login, register, logout, updateProfile,
    }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthState {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
