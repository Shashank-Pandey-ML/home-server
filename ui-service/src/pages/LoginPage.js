import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import axios from 'axios';

function LoginPage() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const { login } = useAuth();
    const navigate = useNavigate();

    const capitalizeFirstLetter = (str) => {
        if (!str) return str;
        return str.charAt(0).toUpperCase() + str.slice(1);
    };

    const handleLogin = async (e) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            const response = await axios.post('http://localhost:8080/api/v1/auth/login', {
                email,
                password,
            });
            
            const { access_token } = response.data;
            await login(access_token);
            navigate('/');
        } catch (err) {
            const errorMsg = err.response?.data?.error || 'Login failed. Please try again.';
            setError(capitalizeFirstLetter(errorMsg));
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="login-page">
            <div className="login-container">
                <div className="login-box">
                    <h2 className="login-title">Login</h2>
                    
                    {error && <div className="error-message">{error}</div>}
                    
                    <form onSubmit={handleLogin} className="login-form">
                        <div className="form-group">
                            <label htmlFor="email">Email</label>
                            <input
                                type="email"
                                id="email"
                                placeholder="Enter your email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                                autoFocus
                            />
                        </div>
                        
                        <div className="form-group">
                            <label htmlFor="password">Password</label>
                            <input
                                type="password"
                                id="password"
                                placeholder="Enter your password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                            />
                        </div>
                        
                        <button type="submit" className="btn-login" disabled={loading}>
                            {loading ? 'Signing in...' : 'Sign In'}
                        </button>
                    </form>
                    
                    <div className="login-footer">
                        <a href="/" className="back-home">← Back to Home</a>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default LoginPage;
