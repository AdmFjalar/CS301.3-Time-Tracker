import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate, useParams } from 'react-router-dom';
import { useAuth } from '../components/AuthContext';
import './UserDetailPage.css';

const UserDetailPage = () => {
  const { id } = useParams(); // Get user id from URL
  const [user, setUser] = useState(null);
  const [errorMessage, setErrorMessage] = useState('');
  const { authToken } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!authToken) {
      navigate('/login');
      return;
    }

    // Fetch user data by ID
    const fetchUserDetails = async () => {
      try {
        const response = await axios.get(`http://localhost:8080/v1/users/${id}`, {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        });

        // Set user data from the response
        setUser(response.data.data); // Access the `data` property
      } catch (error) {
        setErrorMessage('Error: ' + (error.response?.data?.message || 'Failed to fetch user data'));
      }
    };

    fetchUserDetails();
  }, [authToken, id, navigate]);

  const handleEdit = async (e) => {
    e.preventDefault();

    // Construct the payload based on the expected backend structure
    const payload = {
      first_name: user.first_name || '',
      last_name: user.last_name || '',
      email: user.email || '',
      manager_id: user.manager_id || 0,
      role_id: user.role?.id || 0,
    };

    try {
      await axios.patch(`http://localhost:8080/v1/users/${id}`, payload, {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });
      setErrorMessage('User updated successfully!');
    } catch (error) {
      setErrorMessage('Error: ' + (error.response?.data?.message || 'Failed to update user'));
    }
    setTimeout(() => {
      navigate('/users');
    }, 1000);
  };

  const handleDelete = async () => {
    if (!window.confirm('Are you sure you want to delete this user?')) {
      return;
    }

    try {
      await axios.delete(`http://localhost:8080/v1/users/${id}`, {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });
      alert('User deleted successfully!');
      setTimeout(() => {
        navigate('/users');
      }, 1000);
    } catch (error) {
      setErrorMessage('Error: ' + (error.response?.data?.message || 'Failed to delete user'));
      setTimeout(() => {
        navigate('/users');
      }, 1000);
    }
  };

  if (errorMessage) {
    return <div className="error-message">{errorMessage}</div>;
  }

  if (!user) {
    return <div className="loading-message">Loading...</div>;
  }

  return (
    <div className="user-detail-container">
      <h2 className="user-detail-title">Edit User: {user.first_name} {user.last_name}</h2>
      <form className="user-detail-form" onSubmit={handleEdit}>
        <label>
          First Name:
          <input
            type="text"
            value={user.first_name || ''}
            onChange={(e) => setUser({ ...user, first_name: e.target.value })}
            required
          />
        </label>
        <label>
          Last Name:
          <input
            type="text"
            value={user.last_name || ''}
            onChange={(e) => setUser({ ...user, last_name: e.target.value })}
            required
          />
        </label>
        <label>
          Email:
          <input
            type="email"
            value={user.email || ''}
            onChange={(e) => setUser({ ...user, email: e.target.value })}
            required
          />
        </label>
        <label>
          Manager ID:
          <input
            type="number"
            value={user.manager_id || 0}
            onChange={(e) => setUser({ ...user, manager_id: parseInt(e.target.value, 10) || 0 })}
          />
        </label>
        <label>
          Role ID:
          <input
            type="number"
            min={1}
            max={3}
            value={user.role?.id || 0}
            onChange={(e) =>
              setUser({
                ...user,
                role: { ...user.role, id: parseInt(e.target.value, 10) || 0 },
              })
            }
          />
        </label>
        <button type="submit" className="save-button">Save Changes</button>
      </form>

      <button onClick={() => navigate(`/shifts/${id}`)} className="get-timestamps-button">Get users timestamps</button>
      <button onClick={handleDelete} className="delete-button">Delete User</button>
    </div>
  );
};

export default UserDetailPage;
