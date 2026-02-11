import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { useUserStore } from './store'
import AppLayout from './components/Layout'
import Login from './pages/Login'
import RoomList from './pages/RoomList'
import Chat from './pages/Chat'
import './App.css'

// 受保护的路由
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { isLoggedIn } = useUserStore()
  return isLoggedIn ? <>{children}</> : <Navigate to="/login" />
}

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <AppLayout>
                <RoomList />
              </AppLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/rooms"
          element={
            <ProtectedRoute>
              <AppLayout>
                <RoomList />
              </AppLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/chat/:roomId"
          element={
            <ProtectedRoute>
              <AppLayout>
                <Chat />
              </AppLayout>
            </ProtectedRoute>
          }
        />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  )
}

export default App
