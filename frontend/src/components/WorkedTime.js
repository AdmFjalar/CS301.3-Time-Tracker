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
        shifts.reverse();
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
    const day = startDate.getDate().toString().padStart(2, '0');
    const month = (startDate.getMonth() + 1).toString().padStart(2, '0');
    const year = startDate.getFullYear();
    return `${day}/${month}-${year}`; // Format the date as DD/MM-YYYY
  };

  // Function to format time for display (start-time to end-time)
  const formatShiftTime = (time) => {
    const timeDate = new Date(time);
    const timeTime = timeDate.toLocaleTimeString( 'en-NZ', { hour: '2-digit', minute: '2-digit', hour12: false });    
    return `${timeTime}`;
  };

  const secondsToTimeSpan = (seconds) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${hours}h ${minutes}m`;
  };

  return (
    <div className="worked-time">
      <h2>Worked Time</h2>
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
          {workedShifts.map((shift, index) => (
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

export default WorkedTime;
