import './App.css'
import '@mantine/core/styles.css';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { MantineProvider } from '@mantine/core';
import LoginButton from './components/LoginButton';
import AuthCallback from './components/AuthCallback';
import Dashboard from './components/Dashboard';
import PrivacyPolicy from './components/PrivacyPolicy';


function App() {

  return (
    <MantineProvider>
      <Router>
        <Routes>
          <Route path="/" element={<LoginButton />} />
          <Route path="/auth/callback" element={<AuthCallback />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/privacy" element={<PrivacyPolicy />} />
        </Routes>
      </Router>
    </MantineProvider>
  )
}

export default App
