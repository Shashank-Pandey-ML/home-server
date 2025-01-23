import React, { useState } from 'react';
import { Link } from "react-router-dom";

function NavBar() {
    const [isMenuOpen, setIsMenuOpen] = useState(false);

    const toggleMenu = () => {
        setIsMenuOpen(!isMenuOpen);
    };

    return (
        <nav class="custom-navbar" data-spy="affix" data-offset-top="20">
            <div class="container">
                <p class="logo">Shashank Pandey</p>         
                <ul class="nav">
                    <li class="item">
                        {/* <a class="link" href="#home">Home</a> */}
                        <Link className="link" to="/home">Home</Link>
                    </li>
                    <li class="item">
                        {/* <a class="link" href="#about">About</a> */}
                        <Link className="link" to="/about">About</Link>
                    </li>
                    {/* <li class="item">
                        <a class="link" href="#portfolio">Portfolio</a>
                    </li>
                    <li class="item">
                        <a class="link" href="#blog">Blog</a>
                    </li>
                    <li class="item">
                        <a class="link" href="#contact">Contact</a>
                    </li>
                    <li class="item">
                        <a class="link" href="#login">Login</a>
                    </li> */}
                </ul>
                {/* <a href="javascript:void(0)" id="nav-toggle" class="hamburger hamburger--elastic">
                    <div class="hamburger-box">
                    <div class="hamburger-inner"></div>
                    </div>
                </a> */}
                <a
                    href="#"
                    id="nav-toggle"
                    class={`hamburger hamburger--elastic ${isMenuOpen ? 'is-active' : ''}`}
                    onClick={toggleMenu}
                    >
                    <div class="hamburger-box">
                        <div class="hamburger-inner"></div>
                    </div>
                </a>
            </div>          
        </nav>
    );
}

export default NavBar;
