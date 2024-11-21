import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from '../components/AuthContext'; // Assuming you have an AuthContext to provide the token
import './UserProfile.css'; // Assuming you have a separate CSS file for styling
import ChangePasswordForm from './ChangePasswordForm';

const UserProfile = () => {
  const { authToken } = useAuth(); // Get the authToken from the AuthContext
  const [userDetails, setUserDetails] = useState(null); // State to store user details

  // Fetch user profile data
  const fetchUserProfile = async () => {
    try {
      const response = await axios.get('http://localhost:8080/v1/', {
        headers: {
          Authorization: `Bearer ${authToken}`, // Use the token in the Authorization header
        },
      });
      setUserDetails(response.data.data); // Access the 'data' property in the response
    } catch (error) {
      console.error('Failed to fetch user details:', error.response?.data || error.message);
    }
  };

  // Fetch user details on component mount (when authToken is available)
  useEffect(() => {
    if (authToken) {
      fetchUserProfile();
    }
  }, [authToken]); // Re-run when the authToken changes

  // Handle case where user details are not yet loaded
  if (!userDetails) {
    return <div>Loading user profile...</div>;
  }

  return (
    <div className="user-profile">
      <h2>User Profile</h2>
      <div className="profile-details">
        <p><strong>First Name:</strong> {userDetails.first_name}</p>
        <p><strong>Last Name:</strong> {userDetails.last_name}</p>
        <p><strong>Email:</strong> {userDetails.email}</p>
        <p><strong>Role:</strong> {userDetails.role.name}</p>
      </div>
      <ChangePasswordForm />
    </div>
  );
};

export default UserProfile;
