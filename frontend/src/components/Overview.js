import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from '../components/AuthContext';
import './Overview.css';

const Overview = () => {
  const { authToken } = useAuth();
  const [shiftStats, setShiftStats] = useState({
    totalShiftTime: 0,
    totalBreakTime: 0,
    netWorkTime: 0,
  });
  const [latestShift, setLatestShift] = useState(null); // To store the latest shift
  const [latestStampType, setLatestStampType] = useState(null); // To store the latest timestamp type

  const validTransitions = {
    'sign-in': 'sign-out',
    'sign-out': 'sign-in,end-break',
    'start-break': 'sign-in,end-break',
    'end-break': 'start-break',
  };

  // Fetch the shifts and sum the stats from the last 7 days
  const fetchLatestShift = async () => {
    try {
      const response = await axios.get('http://localhost:8080/v1/shifts/', {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });

      const shifts = response.data?.data || [];
      const sevenDaysAgo = new Date();
      sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7);

      let totalShiftTime = 0;
      let totalBreakTime = 0;
      let netWorkTime = 0;

      // Filter and aggregate the stats from the last 7 days
      shifts.forEach(shift => {
        const signInDate = new Date(shift.SignIn);
        if (signInDate >= sevenDaysAgo) {
          totalShiftTime += shift.TotalShiftTime || 0;
          totalBreakTime += shift.TotalBreakTime || 0;
          netWorkTime += shift.NetWorkTime || 0;
        }
      });

      setShiftStats({
        totalShiftTime,
        totalBreakTime,
        netWorkTime,
      });

      // Set the latest shift (this is always the most recent shift in the array)
      if (shifts.length > 0) {
        setLatestShift(shifts[shifts.length - 1]);
      }

    } catch (error) {
      console.error('Failed to fetch the latest shift data:', error.response?.data || error.message);
    }
  };

  // Fetch the latest timestamp type
  const fetchLatestTimestampType = async () => {
    try {
      const response = await axios.get('http://localhost:8080/v1/timestamps/latest', {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });

      const stampType = response.data?.data?.stamp_type;

      if (stampType) {
        setLatestStampType(stampType || 'sign-out');
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
    if (!latestStampType) return stampType === 'sign-in'; 
    const allowedTransitions = validTransitions[stampType]; 
    return allowedTransitions && allowedTransitions.split(',').includes(latestStampType); 
  };

  useEffect(() => {
    fetchLatestShift();         
    fetchLatestTimestampType(); 
  }, [authToken]);

  return (
    <div className="overview">
      <h3>Overview</h3>
      
      {/* Display the latest shift's info */}
      {latestShift ? (
  <>
      <h4>Latest shift:</h4> {parseDate(latestShift.SignIn) ? (
        `${parseDate(latestShift.SignIn).toLocaleDateString()} | ${parseDate(latestShift.SignIn).toLocaleTimeString()}`
      ) : ' Invalid Date'} - {parseDate(latestShift.SignOut) ? parseDate(latestShift.SignOut).toLocaleTimeString() : ' Invalid Date'}
  </>
) : null}

      {/* Display aggregated stats for the last 7 days */}
        <br></br>
        <h4>Total shift time (last 7 days):</h4>    {formatTime(shiftStats.totalShiftTime)}<br />

        <h4>Total break time (last 7 days):</h4>    {formatTime(shiftStats.totalBreakTime)}<br />

        <h4>Effective work time (last 7 days):</h4> {formatTime(shiftStats.netWorkTime)}

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
