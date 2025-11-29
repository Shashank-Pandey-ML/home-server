import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import UserProfile from './UserProfile';

function NavBar() {
    const [isMenuOpen, setIsMenuOpen] = useState(false);
    const [isScrolled, setIsScrolled] = useState(false);
    const { user } = useAuth();
    const navigate = useNavigate();
    const location = useLocation();
    
    const isHomePage = location.pathname === '/' || location.pathname === '';

    const toggleMenu = () => {
        setIsMenuOpen(!isMenuOpen);
    };

    useEffect(() => {
        const handleScroll = () => {
            setIsScrolled(window.scrollY > 50);
        };

        window.addEventListener('scroll', handleScroll);
        return () => window.removeEventListener('scroll', handleScroll);
    }, []);

    const handleNavClick = (e, target) => {
        e.preventDefault();
        
        // If not on home page, navigate to home first
        if (!isHomePage) {
            navigate('/');
            setTimeout(() => {
                const element = document.querySelector(target);
                if (element) {
                    const offset = 70;
                    const elementPosition = element.getBoundingClientRect().top;
                    const offsetPosition = elementPosition + window.pageYOffset - offset;
                    window.scrollTo({
                        top: offsetPosition,
                        behavior: 'smooth'
                    });
                }
            }, 100);
        } else {
            const element = document.querySelector(target);
            if (element) {
                const offset = 70;
                const elementPosition = element.getBoundingClientRect().top;
                const offsetPosition = elementPosition + window.pageYOffset - offset;
                window.scrollTo({
                    top: offsetPosition,
                    behavior: 'smooth'
                });
            }
        }
        setIsMenuOpen(false);
    };

    return (
        <nav className={`custom-navbar ${isScrolled ? 'scrolled' : ''}`}>
            <div className="container">
                <a href="/" className="logo" onClick={(e) => handleNavClick(e, '#home')}>
                    Shashank Pandey
                </a>
                <div className="nav-group">
                    <ul className={`nav ${isMenuOpen ? 'active' : ''}`}>
                    <li className="item">
                        <a className="link" href="#home" onClick={(e) => handleNavClick(e, '#home')}>
                            Home
                        </a>
                    </li>
                    <li className="item">
                        <a className="link" href="#about" onClick={(e) => handleNavClick(e, '#about')}>
                            About
                        </a>
                    </li>
                    <li className="item">
                        <a className="link" href="#service" onClick={(e) => handleNavClick(e, '#service')}>
                            Skills
                        </a>
                    </li>
                    <li className="item">
                        <a className="link" href="#portfolio" onClick={(e) => handleNavClick(e, '#portfolio')}>
                            Portfolio
                        </a>
                    </li>
                    <li className="item">
                        <a className="link" href="#blog" onClick={(e) => handleNavClick(e, '#blog')}>
                            Blog
                        </a>
                    </li>
                    <li className="item">
                        <a className="link" href="#contact" onClick={(e) => handleNavClick(e, '#contact')}>
                            Contact
                        </a>
                    </li>
                    {user && (
                        <li className="item">
                            <a className="link" href="/an/dashboard">
                                Admin Console
                            </a>
                        </li>
                    )}
                    {!user && (
                        <li className="item login-item">
                            <a className="link" href="/an/login">
                                Login
                            </a>
                        </li>
                    )}
                </ul>
                {user && <UserProfile />}
                </div>
                <button
                    type="button"
                    id="nav-toggle"
                    className={`hamburger hamburger--elastic ${isMenuOpen ? 'is-active' : ''}`}
                    onClick={toggleMenu}
                >
                    <div className="hamburger-box">
                        <div className="hamburger-inner"></div>
                    </div>
                </button>
            </div>          
        </nav>
    );
}

export default NavBar;
