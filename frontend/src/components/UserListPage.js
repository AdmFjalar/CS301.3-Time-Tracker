import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../components/AuthContext';
import './UserListPage.css';

const UserListPage = () => {
  const [users, setUsers] = useState([]);
  const [errorMessage, setErrorMessage] = useState('');
  const { authToken } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!authToken) {
      navigate('/login');
      return;
    }

    const fetchUsers = async () => {
      try {
        const response = await axios.get('http://localhost:8080/v1/users/', {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        });

        setUsers(response.data.data);
      } catch (error) {
        setErrorMessage('Error: ' + (error.response?.data?.message || 'Failed to fetch users'));
      }
    };

    fetchUsers();
  }, [authToken, navigate]);

  if (errorMessage) {
    return <div className="error-message">{errorMessage}</div>;
  }

  if (!users.length) {
    return <div className="loading-message">Loading...</div>;
  }

  return (
    <div className="user-list-container">
      <h2 className="user-list-title">All Users</h2>
      <table className="user-list-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Email</th>
            <th>First Name</th>
            <th>Last Name</th>
          </tr>
        </thead>
        <tbody>
          {users.map((user) => (
            <tr
              key={user.id}
              className="clickable-row"
              onClick={() => navigate(`/users/${user.id}`)}
            >
              <td>{user.id}</td>
              <td>{user.email}</td>
              <td>{user.first_name || '-'}</td>
              <td>{user.last_name || '-'}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default UserListPage;
