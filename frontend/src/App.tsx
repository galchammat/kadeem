// src/App.tsx
import './App.css';
import { Route, Routes } from 'react-router';
import { Layout } from './components/layout';
import Home from './pages/home';
import About from './pages/about';
import Accounts from './pages/accounts';
import Matches from './pages/matches';

function App() {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<Home />} />
        <Route path="accounts" element={<Accounts />} />
        <Route path="matches" element={<Matches />} />
        <Route path="about" element={<About />} />
      </Route>
    </Routes>
  );
}

export default App;