// src/App.tsx
import './App.css';
import { Route, Routes } from 'react-router';
import { Layout } from './components/layout';
import HomePage from './pages/home';
import AccountsPage from './pages/accounts';
import MatchesPage from './pages/matches';
import { ThemeProvider } from "./components/themeProvider"
import { StreamerProvider } from './contexts/streamerContext';
import StreamersPage from './pages/streamers';

function App() {
  return (
    <ThemeProvider defaultTheme="dark">
      <StreamerProvider>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<HomePage />} />
            <Route path="accounts" element={<AccountsPage />} />
            <Route path="matches" element={<MatchesPage />} />
            <Route path="streamers" element={<StreamersPage />} />
          </Route>
        </Routes>
      </StreamerProvider>
    </ThemeProvider>
  );
}

export default App;