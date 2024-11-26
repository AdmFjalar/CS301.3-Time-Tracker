import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../components/AuthContext';
import './Login.css';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { signIn } = useAuth();
  const { authToken } = useAuth(); // Get the authToken and signOut function from the context
  const navigate = useNavigate();

useEffect(() => {
  if (authToken) {
    navigate('/'); // Redirect to the home page if already logged in
  }
}, [authToken]);

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
      navigate('/'); // Navigate to the home page on successful login
    } catch (error) {
      setMessage('Error: ' + (error.response?.data?.message || 'Login failed'));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-form-container">
        <div className="login-form">
          <h2>Log in</h2>
          <form onSubmit={handleLogin}>
            <input
              type="email"
              placeholder="Enter email..."
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
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
          </form>
          <p className="login-message">{message}</p>
          <div className="login-links">
            <button
              className="link-button"
              onClick={() => navigate('/register')}
            >
              No account?
            </button>
            <button
              className="link-button"
              onClick={() => navigate('/forgot-password')}
            >
              Forgotten password?
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Login;
