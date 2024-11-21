import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../components/AuthContext'; // Use AuthContext for the token

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

    // Fetch the list of users
    const fetchUsers = async () => {
      try {
        const response = await axios.get('http://localhost:8080/v1/users', {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        });
        setUsers(response.data); // Assuming users data is returned as an array
      } catch (error) {
        setErrorMessage('Error: ' + (error.response?.data?.message || 'Failed to fetch users'));
      }
    };

    fetchUsers();
  }, [authToken, navigate]);

  if (errorMessage) {
    return <div>{errorMessage}</div>;
  }

  if (!users.length) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h2>All Users</h2>
      <ul>
        {users.map((user) => (
          <li key={user.id}>
            <a href={`/user/${user.id}`}>{user.first_name} {user.last_name}</a>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default UserListPage;
