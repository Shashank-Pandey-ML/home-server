import React from "react";
import man1 from '../assets/images/man1.png';

function About() {
    return (
        <section class="section pt-0" id="about">
        <div class="container text-center">
            <div class="about" id="about">
                <div class="about-img-holder">
                    <img src={man1} class="about-img" alt="Failed to render" />
                </div>
                <div class="about-caption">
                    <p class="section-subtitle">Who Am I ?</p>
                    <h2 class="section-title mb-3">About Me</h2>
                    <p>
                        I am a backend developer with 5+ years of professional experience in Python, Golang, Flask, Django and Docker with a passion for problem solving. 
                        <br/>Throughout my career at Xoriant/Maplelabs, I have led significant projects that streamlined processes, enhanced system efficiency, and maintained high standards of code quality. My contributions include spearheading design decisions, conducting thorough code reviews, and developing robust solutions that address complex technical challenges. Despite having limited experience with scalable products, I am a fast learner with a strong passion for innovation and continuous improvement.              
                    </p>
                    <a href="assets/misc/Shashank Resume.pdf" target="_blank" download="Shashank Pandey Resume.pdf">
                        <button class="btn-rounded btn btn-outline-primary mt-4">Download CV</button>
                    </a>
                </div>              
            </div>
        </div>
    </section>
    );
}

export default About;
