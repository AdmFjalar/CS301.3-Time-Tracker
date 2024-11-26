import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { AuthProvider } from './components/AuthContext'; // Import AuthProvider
import HomePage from './components/HomePage';
import UserProfile from './components/UserProfile';
import ShiftDetails from './components/ShiftDetails';
import WorkedTimePage from './components/WorkedTimePage';
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
            <Route path="/worked-time" element={<WorkedTimePage />} />
            <Route path="/shifts/:id" element={<ShiftDetails />} />
            <Route path="/confirm/:token" element={<Activate />} />
            <Route path="/login" element={<Login />} />
            <Route path="/user-details" element={<UserProfile />} />
            <Route path="/users" element={<UserListPage />} />
            <Route path="/users/:id" element={<UserDetailPage />} />
            <Route path="/register" element={<Register />} />
            <Route path="/forgot-password" element={<ForgotPassword />} /> 
            <Route path="/reset-password/:token" element={<ResetPassword />} /> 
          </Route>
        </Routes>
      </Router>
    </AuthProvider>
  );
};

export default App;

