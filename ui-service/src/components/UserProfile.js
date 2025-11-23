import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

function UserProfile() {
    const { user, logout } = useAuth();
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const dropdownRef = useRef(null);
    const navigate = useNavigate();

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
                setIsDropdownOpen(false);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleLogout = () => {
        logout();
        setIsDropdownOpen(false);
        navigate('/');
    };

    const getInitials = () => {
        if (!user) return '';
        const name = user.username || user.name || user.email || '';
        return name.charAt(0).toUpperCase();
    };

    if (!user) return null;

    return (
        <div className="user-profile" ref={dropdownRef}>
            <div 
                className="profile-avatar" 
                onClick={() => setIsDropdownOpen(!isDropdownOpen)}
            >
                {getInitials()}
            </div>
            
            {isDropdownOpen && (
                <div className="profile-dropdown">
                    <div className="profile-header">
                        <div className="profile-name">{user.username || user.name || 'User'}</div>
                        <div className="profile-email">{user.email}</div>
                    </div>
                    <button onClick={handleLogout} className="dropdown-item">
                        Logout
                    </button>
                </div>
            )}
        </div>
    );
}

export default UserProfile;
