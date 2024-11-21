import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useAuth } from '../components/AuthContext'; // Assuming you have an AuthContext to provide the token
import { useNavigate, useParams } from 'react-router-dom'; // Import useParams for URL params and navigate for redirect
import './WorkedTime.css';

const ShiftDetails = () => {
  const { authToken } = useAuth(); // Get auth token from context
  const [shiftDetails, setShiftDetails] = useState([]);
  const [errorMessage, setErrorMessage] = useState('');
  const { id } = useParams(); // Extract the 'id' param from the URL
  const navigate = useNavigate(); // To navigate if unauthorized

  useEffect(() => {
    // If no authToken, navigate to login
    if (!authToken) {
      navigate('/login');
      return; // Stop execution if no authToken
    }

    // Function to fetch specific shift details by ID
    const fetchShiftDetails = async () => {
      try {
        const response = await axios.get(`http://localhost:8080/v1/shifts/${id}`, {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        });
        setShiftDetails(response.data?.data || []); // Set fetched shift details
      } catch (error) {
        setErrorMessage('Error: ' + (error.response?.data?.message || 'Failed to fetch shift details'));
      }
    };

    if (authToken && id) {
      fetchShiftDetails();
    } else {
      setErrorMessage('Authorization token is missing or invalid shift ID.');
    }
  }, [authToken, id, navigate]);

  if (!authToken) {
    return <div>Loading...</div>;
  }

  if (errorMessage) {
    return <div>{errorMessage}</div>;
  }

  // Function to format date for display
  const formatShiftDate = (start) => {
    const startDate = new Date(start);
    const day = startDate.getDate().toString().padStart(2, '0');
    const month = (startDate.getMonth() + 1).toString().padStart(2, '0');
    const year = startDate.getFullYear();
    return `${day}/${month}-${year}`; // Format the date as DD/MM-YYYY
  };

  // Function to format time for display (start-time to end-time)
  const formatShiftTime = (time) => {
    const timeDate = new Date(time);
    const timeTime = timeDate.toLocaleTimeString('en-NZ', { hour: '2-digit', minute: '2-digit', hour12: false });
    return `${timeTime}`;
  };

  const secondsToTimeSpan = (seconds) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${hours}h ${minutes}m`;
  };

  return (
    <div className="worked-time">
      <h2>Shift Details</h2>
      <table>
        <thead>
          <tr>
            <th>Date</th>
            <th>Sign-in</th>
            <th>Break time</th>
            <th>Shift time</th>
            <th>Sign-out</th>
          </tr>
        </thead>
        <tbody>
          {shiftDetails.map((shift, index) => (
            <tr key={index}>
              <td>{formatShiftDate(shift.SignIn)}</td>
              <td>{formatShiftTime(shift.SignIn)}</td>
              <td>{secondsToTimeSpan(shift.TotalBreakTime)}</td>
              <td>{secondsToTimeSpan(shift.TotalShiftTime)}</td>
              <td>{formatShiftTime(shift.SignOut)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default ShiftDetails;
