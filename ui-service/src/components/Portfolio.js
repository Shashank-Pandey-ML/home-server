import React from 'react';
import folio1 from '../assets/images/folio1.png';
import folio2 from '../assets/images/folio2.png';
import folio3 from '../assets/images/folio3.png';

function Portfolio() {
    const projects = [
        { 
            img: folio1, 
            title: 'Home Server Infrastructure',
            description: 'Microservices architecture with Docker, Go, React, and PostgreSQL. Features JWT authentication, API gateway with routing, and service orchestration.',
            tags: ['Go', 'Docker', 'React', 'PostgreSQL', 'Microservices'],
            link: 'https://github.com/Shashank-Pandey-ML/home-server'
        },
        { 
            img: folio2, 
            title: 'Open Source Contributions',
            description: 'Identified and resolved critical multithreading deadlock issues in Cisco UCS Python SDKs (ucsmsdk, ucscsdk, imcsdk) by enhancing synchronization logic. Also fixed critical data reading issue in boltdb Python repository caused by incorrect free-ids storage.',
            tags: ['Python', 'Open Source', 'Multithreading', 'Bug Fixes', 'SDK'],
            links: [
                { name: 'ucsmsdk', url: 'https://github.com/Shashank-Pandey-ML/ucsmsdk' },
                { name: 'ucscsdk', url: 'https://github.com/Shashank-Pandey-ML/ucscsdk' },
                { name: 'imcsdk', url: 'https://github.com/Shashank-Pandey-ML/imcsdk' },
                { name: 'boltdb', url: 'https://github.com/Shashank-Pandey-ML/boltdb' }
            ]
        },
        { 
            img: folio3, 
            title: 'ClockWrk Mobile App',
            description: 'Flutter-based mobile application for employee clock-in/clock-out management with authentication and time tracking features.',
            tags: ['Flutter', 'Dart', 'Mobile', 'Authentication'],
            link: 'https://github.com/Shashank-Pandey-ML/clockwrk_login'
        }
    ];

    return (
        <section className="section portfolio-section" id="portfolio">
            <div className="container text-center">
                <p className="section-subtitle">What I Built</p>
                <h6 className="section-title">Portfolio</h6>
                <div className="portfolio-grid">
                    {projects.map((project, index) => (
                        <div key={index}>
                            {project.link ? (
                                <a 
                                    href={project.link} 
                                    target="_blank" 
                                    rel="noopener noreferrer" 
                                    className="portfolio-item-link"
                                >
                                    <div className="portfolio-item">
                                        <img 
                                            src={project.img} 
                                            className="portfolio-item-img" 
                                            alt={project.title}
                                        />
                                        <div className="portfolio-item-content">
                                            <h4 className="portfolio-item-title">{project.title}</h4>
                                            <p className="portfolio-item-description">{project.description}</p>
                                            <div className="portfolio-item-tags">
                                                {project.tags.map((tag, idx) => (
                                                    <span className="portfolio-tag" key={idx}>{tag}</span>
                                                ))}
                                            </div>
                                        </div>
                                    </div>
                                </a>
                            ) : (
                                <div className="portfolio-item">
                                    <img 
                                        src={project.img} 
                                        className="portfolio-item-img" 
                                        alt={project.title}
                                    />
                                    <div className="portfolio-item-content">
                                        <h4 className="portfolio-item-title">{project.title}</h4>
                                        <p className="portfolio-item-description">{project.description}</p>
                                        <div className="portfolio-item-tags">
                                            {project.tags.map((tag, idx) => (
                                                <span className="portfolio-tag" key={idx}>{tag}</span>
                                            ))}
                                        </div>
                                        {project.links && (
                                            <div className="portfolio-item-links">
                                                {project.links.map((link, idx) => (
                                                    <a 
                                                        href={link.url}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        className="portfolio-link-btn"
                                                        key={idx}
                                                    >
                                                        {link.name}
                                                    </a>
                                                ))}
                                            </div>
                                        )}
                                    </div>
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
}

export default Portfolio;
