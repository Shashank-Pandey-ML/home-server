import React from "react";

function About() {
    return (
        <section className="section pt-0" id="about">
            <div className="container text-center">
                <div className="about">
                    <div className="about-img-holder">
                        <img src="/an/ShashankProfileImg.jpeg" className="about-img" alt="Shashank Pandey" style={{ borderRadius: '15px' }} />
                    </div>
                    <div className="about-caption">
                        <p className="section-subtitle">Who Am I ?</p>
                        <h2 className="section-title mb-3">About Me</h2>
                        <p>
                            I am a backend developer with 5+ years of professional experience in Python, Golang, Flask, Django and Docker with a passion for problem solving. 
                            <br/>Throughout my career at Xoriant/Maplelabs, I have led significant projects that streamlined processes, enhanced system efficiency, and maintained high standards of code quality. My contributions include spearheading design decisions, conducting thorough code reviews, and developing robust solutions that address complex technical challenges. Despite having limited experience with scalable products, I am a fast learner with a strong passion for innovation and continuous improvement.              
                        </p>
                        <a href="/Shashank-Resume-v1.1.pdf" target="_blank" rel="noopener noreferrer" download="Shashank-Pandey-Resume.pdf">
                            <button className="btn-rounded btn btn-outline-primary mt-4">Download CV</button>
                        </a>
                    </div>              
                </div>
            </div>
        </section>
    );
}

export default About;
