import React from 'react';
import './WorkedTime.css';

const WorkedTime = () => {
  // Sample worked time data (replace with actual data or dynamic content later)
  const workedShifts = [
    { date: '28/11', time: '06:00 - 16:00' },
    { date: '27/11', time: '06:00 - 16:00' },
    { date: '26/11', time: '06:00 - 16:00' },
    { date: '25/11', time: '06:00 - 16:00' },
  ];

  return (
    <div className="worked-time">
      <h2>Worked Time</h2>
      <table>
        <thead>
          <tr>
            <th>Date</th>
            <th>Shift</th>
          </tr>
        </thead>
        <tbody>
          {workedShifts.map((shift, index) => (
            <tr key={index}>
              <td>{shift.date}</td>
              <td>{shift.time}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default WorkedTime;
