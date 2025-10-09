// src/App.tsx
import './App.css';
import { Greet } from "../wailsjs/go/main/App";
import * as RiotClient from '../wailsjs/go/riot/RiotClient';
import { useState, useEffect } from 'react';
import { Route, Routes } from 'react-router';
import { Layout } from './components/layout';
import Home from './pages/home';
import About from './pages/about';

function App() {
  const [resultText, setResultText] = useState("Please enter your name below 👇");
  const [name, setName] = useState('');
  const updateName = (e: any) => setName(e.target.value);
  const updateResultText = (result: string) => setResultText(result);

  function greet() {
    Greet(name).then(updateResultText);
  }

  const [account, setAccount] = useState<string>('Loading...');

  useEffect(() => {
    RiotClient.AddAccount("americas", "the thirsty rock", "NA1").then(() => {console.log("Added account")}).catch((err: any) => {
      setAccount(`Error: ${err}`);
    });
  }, []); // Empty dependency array = run once when component mounts

  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<Home />} />
        <Route path="about" element={<About />} />
      </Route>
    </Routes>
  );
}

export default App;