import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';
import './ResetPassword.css';

const ResetPassword = () => {
  const { token } = useParams(); // Token from URL
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    if (!token) {
      setMessage('Invalid or expired reset token.');
    }
  }, [token]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setMessage('');

    if (newPassword !== confirmPassword) {
      setMessage('Passwords do not match!');
      setLoading(false);
      return;
    }

    try {
      await axios.post('http://localhost:8080/v1/auth/reset-password', { token, newPassword });
      setMessage('Password reset successful!');
      navigate('/login'); // Redirect to login after successful reset
    } catch (error) {
      setMessage('Error: ' + (error.response?.data?.message || 'Failed to reset password'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h2>Reset Password</h2>
      <form onSubmit={handleSubmit}>
        <input
          type="password"
          placeholder="New password"
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Confirm new password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          required
        />
        <button type="submit" disabled={loading}>
          {loading ? 'Resetting...' : 'Reset Password'}
        </button>
      </form>
      <p>{message}</p>
    </div>
  );
};

export default ResetPassword;
