import React from 'react';
import './Overview.css'; 

const Overview = () => {
  return (
    <div className="overview">
      <h3>Overview</h3>
      <p>Last shift: 28/11 | 06:00 - 16:00</p>
      <p>Time in lieu: 165 hours</p>
      <p>Comp time: 165 hours</p>
      <p>Flex time: 165 hours</p>
      <p>Sick leave: 165 hours</p>
      <p>Leave: 165 hours</p>
      <div className="buttons">
        <button className="start-shift">Start shift</button>
        <button className="start-break">Start break</button>
        <button className="register-leave">Register leave</button>
      </div>
    </div>
  );
};

export default Overview;
