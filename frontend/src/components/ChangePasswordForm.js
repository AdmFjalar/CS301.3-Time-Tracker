import React, { useState } from 'react';
import './ChangePasswordForm.css';

const ChangePasswordForm = () => {
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [error, setError] = useState(null);
  const [successMessage, setSuccessMessage] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();

    const payload = {
      old_password: oldPassword,
      new_password: newPassword,
    };

    try {
      const response = await fetch('http://localhost:8080/v1/change-password', {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        // Handle error if the response is not OK
        const errorData = await response.json();
        setError(errorData.message || 'Failed to change password');
        setSuccessMessage(null);
      } else {
        // Success response
        setSuccessMessage('Password updated successfully!');
        setError(null);
      }
    } catch (error) {
      setError('An error occurred while changing the password.');
      setSuccessMessage(null);
    }
  };

  return (
    <div className="change-password-container">
      <h3>Change Password</h3>
      <form onSubmit={handleSubmit} className="change-password-form">
        <div>
          <label htmlFor="oldPassword">Old Password: </label>
          <input
            type="password"
            id="oldPassword"
            value={oldPassword}
            onChange={(e) => setOldPassword(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="newPassword">New Password: </label>
          <input
            type="password"
            id="newPassword"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit">Update Password</button>
      </form>

      {error && <p className="error-message">{error}</p>}
      {successMessage && <p className="success-message">{successMessage}</p>}
    </div>
  );
};

export default ChangePasswordForm;
