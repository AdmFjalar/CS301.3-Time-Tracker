import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useParams } from 'react-router-dom';
import './Activate.css';

const Activate = () => {
  const { token: urlToken } = useParams(); // Extract token from URL parameters
  const [token, setToken] = useState(urlToken || ''); // Initialize with URL token or empty string
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false); // Track whether the request is in progress
  const [activated, setActivated] = useState(false); // Track whether activation was successful

  useEffect(() => {
    // If there's a token in the URL, set it in state
    if (urlToken) {
      setToken(urlToken);
    }
  }, [urlToken]);

  const handleActivate = async (e) => {
    e.preventDefault();
    setLoading(true); // Disable the button and show processing
    setMessage('');

    try {
      await axios.put(`http://localhost:8080/v1/users/activate/${token}`);
      setMessage('Account activated successfully!');
      setActivated(true); // Mark activation as successful
    } catch (error) {
      setMessage('Error: ' + (error.response?.data?.message || 'Activation failed'));
    } finally {
      setLoading(false); // Re-enable button only if activation failed
    }
  };

  return (
    <div>
      <h2>Account Activation</h2>
      {!activated && ( // Hide form if activation is successful
        <form onSubmit={handleActivate}>
          <input
            className="token-input"
            type="text"
            placeholder="Activation Token"
            value={token}
            onChange={(e) => setToken(e.target.value)}
            required
          />
          <button type="submit" disabled={loading}>
            {loading ? 'Processing...' : 'Activate'}
          </button>
        </form>
      )}
      <p>{message}</p>
    </div>
  );
};

export default Activate;
