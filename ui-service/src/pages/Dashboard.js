import React from 'react';
import { Link } from 'react-router-dom';
import NavBar from '../components/NavBar';

function Dashboard() {
    return (
        <div className="dashboard">
            <NavBar />

            <div className="dashboard-content">
                <div className="container">
                    <h1 className="dashboard-title">Admin Console</h1>
                    <p className="dashboard-subtitle">Manage your home server</p>

                    <div className="dashboard-grid">
                        <Link to="/dashboard/stats" className="dashboard-card">
                            <div className="card-icon">
                                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <line x1="12" y1="20" x2="12" y2="10"></line>
                                    <line x1="18" y1="20" x2="18" y2="4"></line>
                                    <line x1="6" y1="20" x2="6" y2="16"></line>
                                </svg>
                            </div>
                            <h3>Stats</h3>
                            <p>Hardware metrics and service status</p>
                        </Link>

                        <div className="dashboard-card coming-soon">
                            <div className="card-icon">
                                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"></path>
                                    <circle cx="12" cy="13" r="4"></circle>
                                </svg>
                            </div>
                            <h3>Camera</h3>
                            <p>Live camera streams from your home</p>
                            <span className="badge">Coming Soon</span>
                        </div>

                        <div className="dashboard-card coming-soon">
                            <div className="card-icon">
                                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path>
                                    <polyline points="13 2 13 9 20 9"></polyline>
                                </svg>
                            </div>
                            <h3>File Server</h3>
                            <p>Browse and manage your files</p>
                            <span className="badge">Coming Soon</span>
                        </div>

                        <div className="dashboard-card coming-soon">
                            <div className="card-icon">
                                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <circle cx="12" cy="12" r="3"></circle>
                                    <path d="M12 1v6m0 6v6m5.66-13.66l-4.24 4.24m-2.82 2.82l-4.24 4.24m13.66-4.24l-4.24-4.24m-2.82-2.82l-4.24-4.24"></path>
                                </svg>
                            </div>
                            <h3>Settings</h3>
                            <p>Configure system preferences</p>
                            <span className="badge">Coming Soon</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Dashboard;
