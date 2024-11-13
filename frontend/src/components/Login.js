import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom'; // Import the useNavigate hook
import { useAuth } from '../components/AuthContext';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false); // Track loading state
  const { signIn } = useAuth(); // Use signIn method from context
  const navigate = useNavigate(); // Initialize useNavigate

  const handleLogin = async (e) => {
    e.preventDefault();
    setIsLoading(true); // Set loading state to true during request
    try {
      const response = await axios.post('http://localhost:8080/v1/authentication/token', {
        email,
        password,
      });
      const token = response.data.data; // Extract token from response
      signIn(token); // Store the token in context via signIn
      setMessage('Login successful!');
      navigate('/'); // Redirect to home page after successful login
    } catch (error) {
      setMessage('Error: ' + (error.response?.data?.message || 'Login failed'));
    } finally {
      setIsLoading(false); // Set loading state back to false
    }
  };

  return (
    <div>
      <h2>Login</h2>
      <form onSubmit={handleLogin}>
        <input
          type="text"
          placeholder="ID or Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit" disabled={isLoading}>Login</button> {/* Disable button while loading */}
      </form>
      <p>{message}</p>
      {isLoading && <p>Loading...</p>} {/* Show loading state */}
    </div>
  );
};

export default Login;
