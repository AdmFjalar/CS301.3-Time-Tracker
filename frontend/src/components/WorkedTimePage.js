import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../components/AuthContext'; // Use the AuthContext to get the token
import WorkedTime from '../components/WorkedTime'; // Default import
import Overview from '../components/Overview'; // Default import
import './WorkedTimePage.css';

const WorkedTimePage = () => {
    const { authToken } = useAuth(); // Get auth token from context
    const navigate = useNavigate(); // To navigate if unauthorized

    useEffect(() => {
        // If no authToken, navigate to login
        if (!authToken) {
            navigate('/login');
            return; // Stop execution if no authToken
        }
    })

    return (
        <>
            <WorkedTime />
            <Overview />

        </>
    );
}

export default WorkedTimePage;