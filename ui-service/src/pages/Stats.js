import React, { useState, useEffect } from 'react';
import axios from 'axios';
import NavBar from '../components/NavBar';

function Stats() {
    const [stats, setStats] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetchStats();
        const interval = setInterval(fetchStats, 5000); // Refresh every 5 seconds
        return () => clearInterval(interval);
    }, []);

    const fetchStats = async () => {
        try {
            const token = localStorage.getItem('token');
            const response = await axios.get('http://localhost:8080/api/v1/stats', {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            setStats(response.data);
            setLoading(false);
            setError(null);
        } catch (err) {
            let errorMessage = 'Failed to fetch system statistics';
            
            if (err.response) {
                // Server responded with error status
                const statusCode = err.response.status;
                const errorData = err.response.data;
                
                if (statusCode === 401) {
                    errorMessage = 'Unauthorized - Please login again';
                } else if (statusCode === 403) {
                    errorMessage = 'Forbidden - You do not have access to view statistics';
                } else if (statusCode === 404) {
                    errorMessage = 'Stats endpoint not found';
                } else if (statusCode === 500) {
                    errorMessage = errorData?.error || errorData?.message || 'Internal server error';
                } else if (errorData?.error) {
                    errorMessage = `Error ${statusCode}: ${errorData.error}`;
                } else if (errorData?.message) {
                    errorMessage = `Error ${statusCode}: ${errorData.message}`;
                } else {
                    errorMessage = `Error ${statusCode}: Failed to fetch system statistics`;
                }
            } else if (err.request) {
                // Request made but no response received
                errorMessage = 'No response from server - Please check if the gateway is running';
            } else {
                // Error in setting up the request
                errorMessage = err.message || 'Failed to fetch system statistics';
            }
            
            setError(errorMessage);
            setLoading(false);
        }
    };

    const getStatusColor = (value) => {
        if (value < 50) return '#4caf50';
        if (value < 75) return '#ff9800';
        return '#f44336';
    };

    const formatBytes = (bytes) => {
        return (bytes / (1024 * 1024 * 1024)).toFixed(2);
    };

    if (loading) {
        return (
            <div className="stats-page">
                <NavBar />
                <div className="stats-content">
                    <div className="container">
                        <div className="loading">Loading system statistics...</div>
                    </div>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="stats-page">
                <NavBar />
                <div className="stats-content">
                    <div className="container">
                        <div className="error-message">{error}</div>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="stats-page">
            <NavBar />

            <div className="stats-content">
                <div className="container">
                    <h1 className="page-title">System Statistics</h1>
                    
                    {/* System Info Section */}
                    <div className="info-section">
                        <div className="info-grid">
                            <div className="info-item">
                                <span className="info-label">Hostname</span>
                                <span className="info-value">{stats.hostname}</span>
                            </div>
                            <div className="info-item">
                                <span className="info-label">Platform</span>
                                <span className="info-value">{stats.platform} ({stats.os})</span>
                            </div>
                            <div className="info-item">
                                <span className="info-label">Architecture</span>
                                <span className="info-value">{stats.architecture}</span>
                            </div>
                            <div className="info-item">
                                <span className="info-label">CPU Model</span>
                                <span className="info-value">{stats.cpu.model_name}</span>
                            </div>
                        </div>
                    </div>

                    {/* Main Stats Grid */}
                    <div className="stats-grid">
                        {/* CPU Usage */}
                        <div className="stat-card">
                            <div className="stat-icon" style={{ color: getStatusColor(stats.cpu.usage_percent) }}>
                                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <rect x="4" y="4" width="16" height="16" rx="2"/>
                                    <rect x="9" y="9" width="6" height="6"/>
                                    <line x1="9" y1="1" x2="9" y2="4"/>
                                    <line x1="15" y1="1" x2="15" y2="4"/>
                                    <line x1="9" y1="20" x2="9" y2="23"/>
                                    <line x1="15" y1="20" x2="15" y2="23"/>
                                    <line x1="20" y1="9" x2="23" y2="9"/>
                                    <line x1="20" y1="14" x2="23" y2="14"/>
                                    <line x1="1" y1="9" x2="4" y2="9"/>
                                    <line x1="1" y1="14" x2="4" y2="14"/>
                                </svg>
                            </div>
                            <h3>CPU Usage</h3>
                            <div className="stat-value">{stats.cpu.usage_percent.toFixed(1)}%</div>
                            <div className="stat-detail">{stats.cpu_count} Cores</div>
                            <div className="progress-bar">
                                <div 
                                    className="progress-fill" 
                                    style={{ 
                                        width: `${stats.cpu.usage_percent}%`,
                                        backgroundColor: getStatusColor(stats.cpu.usage_percent)
                                    }}
                                />
                            </div>
                        </div>

                        {/* Memory Usage */}
                        <div className="stat-card">
                            <div className="stat-icon" style={{ color: getStatusColor(stats.memory.used_percent) }}>
                                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <rect x="2" y="4" width="20" height="16" rx="2"/>
                                    <path d="M6 8h.01M10 8h.01M14 8h.01M18 8h.01M6 12h.01M10 12h.01M14 12h.01M18 12h.01M6 16h.01M10 16h.01M14 16h.01M18 16h.01"/>
                                </svg>
                            </div>
                            <h3>Memory Usage</h3>
                            <div className="stat-value">{stats.memory.used_percent.toFixed(1)}%</div>
                            <div className="stat-detail">{stats.memory.used_gb.toFixed(1)} / {stats.memory.total_gb.toFixed(1)} GB</div>
                            <div className="progress-bar">
                                <div 
                                    className="progress-fill" 
                                    style={{ 
                                        width: `${stats.memory.used_percent}%`,
                                        backgroundColor: getStatusColor(stats.memory.used_percent)
                                    }}
                                />
                            </div>
                        </div>

                        {/* Disk Usage */}
                        {stats.disk && stats.disk.length > 0 && (
                            <div className="stat-card">
                                <div className="stat-icon" style={{ color: getStatusColor(stats.disk[0].used_percent) }}>
                                    <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                        <circle cx="12" cy="12" r="10"/>
                                        <circle cx="12" cy="12" r="3"/>
                                        <line x1="12" y1="2" x2="12" y2="4"/>
                                        <line x1="12" y1="20" x2="12" y2="22"/>
                                    </svg>
                                </div>
                                <h3>Disk Usage</h3>
                                <div className="stat-value">{stats.disk[0].used_percent.toFixed(1)}%</div>
                                <div className="stat-detail">{stats.disk[0].used_gb.toFixed(1)} / {stats.disk[0].total_gb.toFixed(1)} GB</div>
                                <div className="progress-bar">
                                    <div 
                                        className="progress-fill" 
                                        style={{ 
                                            width: `${stats.disk[0].used_percent}%`,
                                            backgroundColor: getStatusColor(stats.disk[0].used_percent)
                                        }}
                                    />
                                </div>
                            </div>
                        )}

                        {/* System Uptime */}
                        <div className="stat-card">
                            <div className="stat-icon" style={{ color: '#2196f3' }}>
                                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <circle cx="12" cy="12" r="10"/>
                                    <polyline points="12 6 12 12 16 14"/>
                                </svg>
                            </div>
                            <h3>System Uptime</h3>
                            <div className="stat-value">{stats.uptime.uptime_formatted}</div>
                            <div className="stat-detail">{stats.uptime.uptime_days.toFixed(1)} days</div>
                        </div>

                        {/* Load Average */}
                        <div className="stat-card">
                            <div className="stat-icon" style={{ color: '#9c27b0' }}>
                                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <path d="M3 12h4l3 9 4-18 3 9h4"/>
                                </svg>
                            </div>
                            <h3>Load Average</h3>
                            <div className="stat-value">{stats.load.load1.toFixed(2)}</div>
                            <div className="stat-detail">1m: {stats.load.load1.toFixed(2)} | 5m: {stats.load.load5.toFixed(2)} | 15m: {stats.load.load15.toFixed(2)}</div>
                        </div>

                        {/* Network Stats */}
                        <div className="stat-card">
                            <div className="stat-icon" style={{ color: '#00bcd4' }}>
                                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                    <rect x="2" y="6" width="20" height="12" rx="2"/>
                                    <line x1="6" y1="10" x2="6" y2="10"/>
                                    <line x1="10" y1="10" x2="10" y2="10"/>
                                    <line x1="14" y1="10" x2="14" y2="10"/>
                                    <line x1="18" y1="10" x2="18" y2="10"/>
                                </svg>
                            </div>
                            <h3>Network Traffic</h3>
                            <div className="stat-value-small">
                                <div>↑ {formatBytes(stats.network.total_sent_bytes)} GB</div>
                                <div>↓ {formatBytes(stats.network.total_recv_bytes)} GB</div>
                            </div>
                        </div>
                    </div>

                    {/* Additional Disk Info */}
                    {stats.disk && stats.disk.length > 1 && (
                        <div className="disk-section">
                            <h2 className="section-title">Disk Partitions</h2>
                            <div className="disk-list">
                                {stats.disk.map((disk, index) => (
                                    <div key={index} className="disk-item">
                                        <div className="disk-header">
                                            <span className="disk-mount">{disk.mount_point}</span>
                                            <span className="disk-usage">{disk.used_gb.toFixed(1)} / {disk.total_gb.toFixed(1)} GB</span>
                                        </div>
                                        <div className="progress-bar">
                                            <div 
                                                className="progress-fill" 
                                                style={{ 
                                                    width: `${disk.used_percent}%`,
                                                    backgroundColor: getStatusColor(disk.used_percent)
                                                }}
                                            />
                                        </div>
                                        <div className="disk-details">
                                            <span>{disk.device}</span>
                                            <span>{disk.fs_type}</span>
                                            <span>{disk.used_percent.toFixed(1)}%</span>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

export default Stats;
