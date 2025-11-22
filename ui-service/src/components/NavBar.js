import React, { useState, useEffect } from 'react';

function NavBar() {
    const [isMenuOpen, setIsMenuOpen] = useState(false);
    const [isScrolled, setIsScrolled] = useState(false);

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
        const element = document.querySelector(target);
        if (element) {
            const offset = 70; // Height of fixed navbar
            const elementPosition = element.getBoundingClientRect().top;
            const offsetPosition = elementPosition + window.pageYOffset - offset;

            window.scrollTo({
                top: offsetPosition,
                behavior: 'smooth'
            });
        }
        setIsMenuOpen(false);
    };

    return (
        <nav className={`custom-navbar ${isScrolled ? 'scrolled' : ''}`}>
            <div className="container">
                <a href="/" className="logo" onClick={(e) => handleNavClick(e, '#home')}>
                    Shashank Pandey
                </a>         
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
                </ul>
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
