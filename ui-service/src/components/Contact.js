import React from 'react';

function Contact() {
    return (
        <div className="container">
            <footer className="footer">       
                <p className="mb-0">Copyright © {new Date().getFullYear()} <a href="/">Shashank Pandey</a></p>
                <div className="social-links text-right m-auto ml-sm-auto">
                    <a href="https://www.linkedin.com/in/iam-shashank-pandey/" className="link" target="_blank" rel="noopener noreferrer">
                        <i className="ti-linkedin"></i>
                    </a>
                    <a href="https://github.com/Shashank-Pandey-ML" className="link" target="_blank" rel="noopener noreferrer">
                        <i className="ti-github"></i>
                    </a>
                    <a href="https://x.com/iam_shashankp" className="link" target="_blank" rel="noopener noreferrer">
                        <i className="ti-twitter-alt"></i>
                    </a>
                    <a href="mailto:shashankp2022@gmail.com" className="link" target="_blank" rel="noopener noreferrer">
                        <i className="ti-email"></i>
                    </a>
                    <a href="https://instagram.com/iam_shashank16" className="link" target="_blank" rel="noopener noreferrer">
                        <i className="ti-instagram"></i>
                    </a>
                </div>
            </footer>
        </div>
    );
}

export default Contact;
