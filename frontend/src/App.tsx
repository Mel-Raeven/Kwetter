import './App.css'
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import LoginButton from './components/LoginButton';
import AuthCallback from './components/AuthCallback';
import LogoutButton from './components/LogoutButton';

function App() {

  return (
    <Router>
      <Routes>
        <Route path="/" element={<LoginButton />} />
        <Route path="/logout" element={<LogoutButton />} />
        <Route path="/auth/callback" element={<AuthCallback />} />
      </Routes>
    </Router>
  )
}

export default App
