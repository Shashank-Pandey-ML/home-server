import React, { createContext, useState, useContext, useEffect } from 'react';
import axios from 'axios';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    const fetchUserProfile = async (token) => {
        try {
            const response = await axios.get('http://localhost:8080/api/v1/auth/users/profile', {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            return response.data;
        } catch (error) {
            console.error('Failed to fetch user profile:', error);
            return null;
        }
    };

    useEffect(() => {
        // Check if user is already logged in (token in localStorage)
        const token = localStorage.getItem('token');
        if (token) {
            fetchUserProfile(token).then(profile => {
                if (profile) {
                    setUser({ ...profile, token });
                } else {
                    // Token invalid, clear it
                    localStorage.removeItem('token');
                }
                setLoading(false);
            });
        } else {
            setLoading(false);
        }
    }, []);

    const login = async (token) => {
        localStorage.setItem('token', token);
        const profile = await fetchUserProfile(token);
        if (profile) {
            setUser({ ...profile, token });
        }
    };

    const logout = () => {
        localStorage.removeItem('token');
        setUser(null);
    };

    return (
        <AuthContext.Provider value={{ user, login, logout, loading }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};
