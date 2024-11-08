import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useLocation } from 'react-router-dom';

const Activate = () => {
  const [token, setToken] = useState('');
  const [message, setMessage] = useState('');
  const location = useLocation();

  useEffect(() => {
    // Extract the token from the URL query
    const params = new URLSearchParams(location.search);
    const urlToken = params.get('token');
    if (urlToken) setToken(urlToken);
  }, [location.search]);

  const handleActivate = async (e) => {
    e.preventDefault();
    try {
      await axios.put(`http://localhost:8080/v1/users/activate/${token}`);
      setMessage('Account activated successfully!');
    } catch (error) {
      setMessage('Error: ' + (error.response?.data?.message || 'Activation failed'));
    }
  };

  return (
    <div>
      <h2>Account Activation</h2>
      <form onSubmit={handleActivate}>
        <input
          type="text"
          placeholder="Activation Token"
          value={token}
          onChange={(e) => setToken(e.target.value)}
          required
        />
        <button type="submit">Activate</button>
      </form>
      <p>{message}</p>
    </div>
  );
};

export default Activate;
