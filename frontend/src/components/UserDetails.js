// UserDetails.js
import React from 'react';

const UserDetails = ({ firstName, lastName, email }) => {
  return (
    <div>
      <h3>User Details</h3>
      <p><strong>First Name:</strong> {firstName}</p>
      <p><strong>Last Name:</strong> {lastName}</p>
      <p><strong>Email:</strong> {email}</p>
    </div>
  );
};

export default UserDetails;
