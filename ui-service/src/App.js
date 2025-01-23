import { Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import LiveFeed from './pages/LiveFeed'
import './App.css';
import './styles.css'

function App() {
  return (
    <div className="App">
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/live-feed" element={<LiveFeed />} />
      </Routes>
    </div>
  );
}

export default App;
