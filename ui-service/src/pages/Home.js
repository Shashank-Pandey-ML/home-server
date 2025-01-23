import React from 'react';
import Login from '../components/Login';
import NavBar from '../components/NavBar';
import ProfileHeader from '../components/ProfileHeader';
import About from '../components/About';

function Home() {
  return (
    <div>
      <NavBar />
      <ProfileHeader />
      <About />
      <Login />
    </div>
  );
}

export default Home;
