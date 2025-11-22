import React from 'react';
import NavBar from '../components/NavBar';
import ProfileHeader from '../components/ProfileHeader';
import About from '../components/About';
import Skills from '../components/Skills';
import Portfolio from '../components/Portfolio';
import Blog from '../components/Blog';
import ContactForm from '../components/ContactForm';
import Contact from '../components/Contact';

function Home() {
  return (
    <div>
      <NavBar />
      <ProfileHeader />
      <About />
      <Skills />
      <Portfolio />
      <Blog />
      <ContactForm />
      <Contact />
    </div>
  );
}

export default Home;
