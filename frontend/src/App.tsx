import { useState, useEffect } from 'react'; // Add useEffect import
import logo from './assets/images/logo-universal.png';
import './App.css';
import { Greet } from "../wailsjs/go/main/App";
import * as RiotClient from '../wailsjs/go/riot/RiotClient';

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
    RiotClient.AddAccount("americas", "the thirsty rock", "NA1").then((res: string) => {
      setAccount(res);
    }).catch((err: any) => {
      setAccount(`Error: ${err}`);
    });
  }, []); // Empty dependency array = run once when component mounts

  return (
    <div id="App">
      <img src={logo} id="logo" alt="logo" />
      <div id="result" className="result">{resultText}</div>
      <div id="result" className="result">{account}</div>
      <div id="input" className="input-box">
        <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text" />
        <button className="btn" onClick={greet}>Greet</button>
      </div>
    </div>
  )
}

export default App