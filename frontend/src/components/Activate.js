import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useParams } from 'react-router-dom';

const Activate = () => {
  const { token: urlToken } = useParams(); // Extract token from URL parameters
  const [token, setToken] = useState(urlToken || ''); // Initialize with URL token or empty string
  const [message, setMessage] = useState('');

  useEffect(() => {
    // If there's a token in the URL, set it in state
    if (urlToken) {
      setToken(urlToken);
    }
  }, [urlToken]); // Dependency array ensures this only runs when urlToken changes

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