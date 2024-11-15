import React, { useContext, useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from '../components/AuthContext'; // Import the AuthContext
import { useNavigate } from 'react-router-dom'; // Import navigate hook from react-router-dom
import './WorkedTime.css';

const WorkedTime = () => {
  const [workedShifts, setWorkedShifts] = useState([]);
  const [errorMessage, setErrorMessage] = useState('');
  const { authToken } = useAuth(); // Get auth token from context
  const navigate = useNavigate(); // To navigate if unauthorized

  useEffect(() => {
    // If no authToken, navigate to login
    if (!authToken) {
      navigate('/login');
      return; // Stop execution if no authToken
    }

    // Function to fetch worked shifts
    const fetchWorkedShifts = async () => {
      try {
        const response = await axios.get('http://localhost:8080/v1/shifts/', {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        });
        const shifts = response.data?.data || [];
        setWorkedShifts(shifts);
      } catch (error) {
        setErrorMessage('Error: ' + (error.response?.data?.message || 'Failed to fetch worked shifts'));
      }
    };

    if (authToken) {
      fetchWorkedShifts();
    } else {
      setErrorMessage('Authorization token is missing.');
    }
  }, [authToken, navigate]);

  if (!authToken) {
    return <div>Loading...</div>;
  }

  if (errorMessage) {
    return <div>{errorMessage}</div>;
  }

  // Function to format date for display
  const formatShiftDate = (start) => {
    const startDate = new Date(start);
    return startDate.toLocaleDateString(); // Format the date as MM/DD/YYYY
  };

  // Function to format time for display (start-time to end-time)
  const formatShiftTime = (start, end) => {
    const startDate = new Date(start);
    const endDate = new Date(end);
    const startTime = startDate.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    const endTime = endDate.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    return `${startTime} - ${endTime}`;
  };

  return (
    <div className="worked-time">
      <h2>Worked Time</h2>
      <table>
        <thead>
          <tr>
            <th>Date</th>
            <th>Time</th>
          </tr>
        </thead>
        <tbody>
          {workedShifts.map((shift, index) => (
            <tr key={index}>
              <td>{formatShiftDate(shift.SignIn)}</td>
              <td>{formatShiftTime(shift.SignIn, shift.SignOut)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default WorkedTime;
