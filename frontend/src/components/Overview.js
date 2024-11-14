import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from '../components/AuthContext';
import './Overview.css';

const Overview = () => {
  const { authToken } = useAuth();
  const [latestShift, setLatestShift] = useState(null); // To store the latest shift data
  const [latestStampType, setLatestStampType] = useState('sign-out'); // Default to 'sign-out' for initial visibility

  const validTransitions = {
    'sign-in': 'sign-out',          // 'sign-in' can only follow 'sign-out'
    'sign-out': 'sign-in,end-break', // 'sign-out' can follow 'sign-in' or 'start-break'
    'start-break': 'sign-in,end-break', // 'start-break' can only follow 'sign-in' or 'end-break'
    'end-break': 'start-break',       // 'end-break' can only follow 'start-break'
  };

  // Fetch the latest shift data (sign-in/sign-out)
  const fetchLatestShift = async () => {
    try {
      const response = await axios.get('http://localhost:8080/v1/shifts/', {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });

      console.log("Fetched shift data:", response.data);

      // Assuming you want the latest shift (the last object in the array)
      const shiftData = response.data?.data?.[response.data.data.length - 1];

      if (shiftData) {
        setLatestShift(shiftData);
      } else {
        console.error("No shift data found.");
      }
    } catch (error) {
      console.error('Failed to fetch the latest shift data:', error.response?.data || error.message);
    }
  };

  // Fetch the latest timestamp type (e.g., 'start-break', 'end-break', etc.)
  const fetchLatestTimestampType = async () => {
    try {
      const response = await axios.get('http://localhost:8080/v1/timestamps/latest', {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });

      const stampType = response.data?.data?.stamp_type;

      console.log(response);
      console.log(stampType);
      if (stampType) {
        setLatestStampType(stampType || 'sign-out'); // Set the latest timestamp type
      } else {
        console.error("No timestamp data found.");
      }
    } catch (error) {
      console.error('Failed to fetch the latest timestamp data:', error.response?.data || error.message);
    }
  };

  // Send POST request to record timestamp
  const sendTimestamp = async (stampType) => {
    const timestampData = {
      stamp_type: stampType,
      stamp_time: new Date().toISOString(),
    };

    try {
      await axios.post('http://localhost:8080/v1/timestamps/', timestampData, {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });
      console.log(`${stampType} recorded successfully`);

      // After a timestamp is sent, we re-fetch the latest shift and the latest timestamp type to update the state correctly
      await fetchLatestShift();
      await fetchLatestTimestampType();
    } catch (error) {
      console.error(`Failed to record ${stampType}:`, error.response?.data || error.message);
    }
  };

  // Format seconds into HH:MM format
  const formatTime = (seconds) => {
    if (isNaN(seconds)) return '00:00'; // Return a default time if input is not a valid number
    const hrs = Math.floor(seconds / 3600);
    const mins = Math.floor((seconds % 3600) / 60);
    return `${hrs.toString().padStart(2, '0')}:${mins.toString().padStart(2, '0')}`;
  };

  // Helper function to safely parse dates
  const parseDate = (dateStr) => {
    const date = new Date(dateStr);
    return isNaN(date) ? null : date;
  };

  // Determine if a button should be visible based on the latest stamp type
  const isButtonVisible = (stampType) => {
    if (!latestStampType) return stampType === 'sign-in'; // Show 'sign-in' if no stamp type exists
    const allowedTransitions = validTransitions[stampType]; // Get allowed transitions for the desired stamp type
    return allowedTransitions && allowedTransitions.split(',').includes(latestStampType); // Check if the latest stamp type is in the allowed transitions
  };

  useEffect(() => {
    fetchLatestShift();         // Fetch the latest shift data when the component mounts
    fetchLatestTimestampType(); // Fetch the latest timestamp type when the component mounts
  }, [authToken]);

  return (
    <div className="overview">
      <h3>Overview</h3>
      {latestShift ? (
        <>
          <p>
            Last shift: 
            {parseDate(latestShift.SignIn) ? `${parseDate(latestShift.SignIn).toLocaleDateString()} | ${parseDate(latestShift.SignIn).toLocaleTimeString()}` : 'Invalid Date'} - 
            {parseDate(latestShift.SignOut) ? parseDate(latestShift.SignOut).toLocaleTimeString() : 'Invalid Date'}
          </p>
          <p>Total shift time: {formatTime(latestShift.TotalShiftTime)}</p>
          <p>Total break time: {formatTime(latestShift.TotalBreakTime)}</p>
          <p>Net work time: {formatTime(latestShift.NetWorkTime)}</p>
        </>
      ) : (
        <p>Loading latest shift data...</p>
      )}
      <div className="buttons">
        {isButtonVisible('sign-in') && (
          <button className="start-shift" onClick={() => sendTimestamp('sign-in')}>
            Start shift
          </button>
        )}
        {isButtonVisible('sign-out') && (
          <button className="stop-shift" onClick={() => sendTimestamp('sign-out')}>
            Stop shift
          </button>
        )}
        {isButtonVisible('start-break') && (
          <button className="start-break" onClick={() => sendTimestamp('start-break')}>
            Start break
          </button>
        )}
        {isButtonVisible('end-break') && (
          <button className="end-break" onClick={() => sendTimestamp('end-break')}>
            End break
          </button>
        )}
        {isButtonVisible('register-leave') && (
          <button className="register-leave" onClick={() => sendTimestamp('register-leave')}>
            Register leave
          </button>
        )}
      </div>
    </div>
  );
};

export default Overview;
