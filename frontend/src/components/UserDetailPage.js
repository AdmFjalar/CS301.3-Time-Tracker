import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate, useParams } from 'react-router-dom';
import { useAuth } from '../components/AuthContext';

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
        setUser(response.data);
      } catch (error) {
        setErrorMessage('Error: ' + (error.response?.data?.message || 'Failed to fetch user data'));
      }
    };

    fetchUserDetails();
  }, [authToken, id, navigate]);

  const handleEdit = async (e) => {
    e.preventDefault();

    try {
      await axios.put(`http://localhost:8080/v1/users/${id}`, user, {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });
      setErrorMessage('User updated successfully!');
    } catch (error) {
      setErrorMessage('Error: ' + (error.response?.data?.message || 'Failed to update user'));
    }
  };

  if (errorMessage) {
    return <div>{errorMessage}</div>;
  }

  if (!user) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h2>Edit User: {user.first_name} {user.last_name}</h2>
      <form onSubmit={handleEdit}>
        <label>
          First Name:
          <input
            type="text"
            value={user.first_name}
            onChange={(e) => setUser({ ...user, first_name: e.target.value })}
            required
          />
        </label>
        <br />
        <label>
          Last Name:
          <input
            type="text"
            value={user.last_name}
            onChange={(e) => setUser({ ...user, last_name: e.target.value })}
            required
          />
        </label>
        <br />
        <button type="submit">Save Changes</button>
      </form>
    </div>
  );
};

export default UserDetailPage;
