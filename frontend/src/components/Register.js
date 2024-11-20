import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import './Register.css';

const Register = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false); // State to track the loading status
  const navigate = useNavigate(); // Hook to navigate between pages

  const handleRegister = async (e) => {
    e.preventDefault();

    setMessage('');
    setLoading(true); // Disable the button during the request

    try {
      await axios.post('http://localhost:8080/v1/authentication/user', {
        email,
        password,
      });
      setMessage('Registration successful! Check your email for the activation link.');
    } catch (error) {
      setMessage('Error: ' + (error.response?.data?.message || 'Registration failed'));
    } finally {
      setLoading(false); // Re-enable the button after the request finishes
    }
  };

  return (
    <div className="register-container">
      <div className="register-form-container">
        <div className="register-form">
          {/* Back arrow */}
          <h2 className="register-title">Register</h2>
          <form onSubmit={handleRegister}>
            <input
              type="email"
              placeholder="Enter your email..."
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
            <input
              type="password"
              placeholder="Enter a password..."
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              minLength={12}
              maxLength={128}
              required
            />
            <button
              className="register-button"
              type="submit"
              disabled={loading} // Disable the button if loading is true
            >
              {loading ? 'Registering...' : 'Register'}
            </button>
          </form>
          <p className="register-message">{message}</p>
          <button className="back-button" onClick={() => navigate('/login')}>
            ← Back to Login
          </button>
        </div>
      </div>
    </div>
  );
};

export default Register;
