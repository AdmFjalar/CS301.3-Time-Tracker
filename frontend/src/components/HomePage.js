import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../components/AuthContext'; // Use the AuthContext to get the token
import UpcomingShifts from '../components/UpcomingShifts'; // Default import
import Overview from '../components/Overview'; // Default import

const HomePage = () => {
  const [userData, setUserData] = useState(null);
  const [errorMessage, setErrorMessage] = useState('');
  const { authToken } = useAuth(); // Get auth token from context
  const navigate = useNavigate(); // To navigate if unauthorized

  useEffect(() => {
    // If no authToken, navigate to login
    if (!authToken) {
      navigate('/login');
      return; // Stop execution if no authToken
    }

    // Function to fetch user data from the API
    const fetchUserData = async () => {
      try {
        const response = await axios.get('http://localhost:8080/v1/', {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        });
        setUserData(response.data.data); // Set user data
      } catch (error) {
        setErrorMessage('Error: ' + (error.response?.data?.message || 'Failed to fetch user data'));
      }
    };

    // Fetch data only if userData isn't already loaded
    if (!userData) {
      fetchUserData();
    }
  }, [authToken, userData, navigate]); // Add userData to dependencies to avoid re-fetching

  // Handle error or loading states
  if (errorMessage) {
    return <div>{errorMessage}</div>; // Show error if there's an issue
  }

  if (!userData) {
    return <div>Loading...</div>; // Show loading message while data is fetched
  }

  // Personalized greeting message
  const welcomeMessage = userData.first_name ? `Welcome back, ${userData.first_name}` : 'Welcome back!';

  return (
    <div>
      <h2>{welcomeMessage}</h2>
      {/* Shifts and overview section */}
      <div className="shifts-and-overview">
        <UpcomingShifts />
        <Overview />
      </div>
    </div>
  );
};

export default HomePage;
