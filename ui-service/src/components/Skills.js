import React from 'react';
import pythonIcon from '../assets/images/python.svg';
import djangoIcon from '../assets/images/django.png';
import flaskIcon from '../assets/images/flask.svg';
import goIcon from '../assets/images/go.png';
import oopsIcon from '../assets/images/object-oriented-programming.png';
import ansibleIcon from '../assets/images/ansible.svg';
import dockerIcon from '../assets/images/docker.svg';
import problemSolvingIcon from '../assets/images/problem-solving.png';
import postgresIcon from '../assets/images/postgres.svg';

function Skills() {
    const skills = [
        { name: 'Python', icon: pythonIcon },
        { name: 'Django', icon: djangoIcon },
        { name: 'Flask', icon: flaskIcon },
        { name: 'Go', icon: goIcon },
        { name: 'OOP Design', icon: oopsIcon },
        { name: 'Ansible', icon: ansibleIcon },
        { name: 'Docker', icon: dockerIcon },
        { name: 'Problem Solving', icon: problemSolvingIcon },
        { name: 'PostgreSQL', icon: postgresIcon }
    ];

    return (
        <section className="section" id="service">
            <div className="container text-center">
                <p className="section-subtitle">What I Do</p>
                <h6 className="section-title">Technical Skills</h6>
                <div className="row">
                    {skills.map((skill, index) => (
                        <div className="col-sm-6 col-md-4 col-lg-3 mb-4" key={index}>
                            <div className="service-card">
                                <div className="body">
                                    <img 
                                        src={skill.icon} 
                                        alt={skill.name}
                                        className="icon"
                                    />
                                    <h6 className="title">{skill.name}</h6>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
}

export default Skills;
