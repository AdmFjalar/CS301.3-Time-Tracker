import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { AuthProvider } from './components/AuthContext'; // Import AuthProvider
import HomePage from './components/HomePage';
import WorkedTime from './components/WorkedTime';
import Login from './components/Login';
import Register from './components/Register';
import Activate from './components/Activate';
import AppLayout from './AppLayout';
import UserListPage from './components/UserListPage'; 
import UserDetailPage from './components/UserDetailPage';  
import ForgotPassword from './components/ForgotPassword'; 
import ResetPassword from './components/ResetPassword'; 
import './App.css';

const App = () => {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/" element={<AppLayout />}>
            <Route path="/" element={<HomePage />} />
            <Route path="/worked-time" element={<WorkedTime />} />
            <Route path="/confirm/:token" element={<Activate />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/users" element={<UserListPage />} /> 
            <Route path="/user/:id" element={<UserDetailPage />} /> 
            <Route path="/forgot-password" element={<ForgotPassword />} /> 
            <Route path="/reset-password/:token" element={<ResetPassword />} /> 
          </Route>
        </Routes>
      </Router>
    </AuthProvider>
  );
};

export default App;

