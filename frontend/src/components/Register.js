import React, { useState } from 'react';
import axios from 'axios';
import './Register.css'; 

const Register = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');

  const handleRegister = async (e) => {
    e.preventDefault();
    setMessage('');
    try {
      await axios.post('http://localhost:8080/v1/authentication/user', {
        email,
        password,
      });
      setMessage('Registration successful! Check your email for activation link.');
    } catch (error) {
      setMessage('Error: ' + (error.response?.data?.message || 'Registration failed'));
    }
  };

  return (
    
    <div className="register-container">
  {/* <div className="register-sidebar">
    <h1>Thyme Flies</h1>
    <p>Time Tracker</p>
  </div> */}
  <div className="register-form-container">
    <div className="register-form">
      <h2 className="register-title">Register</h2>
      <form onSubmit={handleRegister}>
        <input
          type="email"
          placeholder="Email"
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
        <button type="submit">Register</button>
      </form>
      <p className="register-message">{message}</p>
    </div>
  </div>
</div>

  );
};

export default Register;

