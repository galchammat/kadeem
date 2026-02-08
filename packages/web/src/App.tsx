// src/App.tsx
import './App.css';
import { Route, Routes, Navigate } from 'react-router';
import { Layout } from './components/layout';
import HomePage from './pages/home';
import LoginPage from './pages/login';
import AccountsPage from './pages/accounts';
import MatchesPage from './pages/matches';
import { ThemeProvider } from "./components/themeProvider"
import { AuthProvider, useAuth } from './contexts/authContext';
import { StreamerProvider } from './contexts/streamerContext';
import StreamersPage from './pages/streamers';

function RequireAuth({ children }: { children: React.ReactNode }) {
  const { session, loading } = useAuth()
  
  if (loading) return null

  if (!session) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}

function App() {
  return (
    <ThemeProvider defaultTheme="dark">
      <AuthProvider>
        <StreamerProvider>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/" element={
              <RequireAuth>
                <Layout />
              </RequireAuth>
            }>
              <Route index element={<HomePage />} />
              <Route path="accounts" element={<AccountsPage />} />
              <Route path="matches" element={<MatchesPage />} />
              <Route path="streamers" element={<StreamersPage />} />
            </Route>
          </Routes>
        </StreamerProvider>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;