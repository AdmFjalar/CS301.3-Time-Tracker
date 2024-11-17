import React, { useState } from 'react';
import axios from 'axios';
import { useAuth } from '../components/AuthContext';
import './Login.css';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { signIn } = useAuth();
  

  const handleLogin = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      const response = await axios.post('http://localhost:8080/v1/authentication/token', {
        email,
        password,
      });
      const token = response.data.data;
      signIn(token);
      setMessage('Login successful!');
    
    } catch (error) {
      setMessage('Error: ' + (error.response?.data?.message || 'Login failed'));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-sidebar">
        <h1>Thyme Flies</h1>
        <p>Time Tracker</p>
      </div>
      <div className="login-form-container">
        <form onSubmit={handleLogin} className="login-form">
          <h2>User ID or Email</h2>
          <input
            type="text"
            placeholder="Enter username..."
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <h2>Password</h2>
          <input
            type="password"
            placeholder="Enter password..."
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <button type="submit" disabled={isLoading}>
            {isLoading ? 'Loading...' : 'Log in'}
          </button>
          <p className="login-message">{message}</p>
          <div className="login-links">
            <a href="/register">No account?</a>
            <a href="/forgot-password">Forgotten password?</a>
          </div>
        </form>
      </div>
    </div>
  );
};

export default Login;
