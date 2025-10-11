// src/App.tsx
import './App.css';
import { Route, Routes } from 'react-router';
import { Layout } from './components/layout';
import Home from './pages/home';
import Accounts from './pages/accounts';
import Matches from './pages/matches';
import { ThemeProvider } from "./components/themeProvider"

function App() {
  return (
    <ThemeProvider defaultTheme="dark">
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Home />} />
          <Route path="accounts" element={<Accounts />} />
          <Route path="matches" element={<Matches />} />
        </Route>
      </Routes>
    </ThemeProvider>
  );
}

export default App;