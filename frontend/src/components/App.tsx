import React from 'react'
import { Link, Route, Switch } from 'react-router-dom'

import 'tailwindcss/tailwind.css'
import '../index.css'

//import thumbnail from 'url:../thumbnail.jpg'

import { Login, Register } from './Auth'
import { VideoDetail, VideoList, Video } from './Video'
import { Profile, Settings } from './User'
import { TypeFlags } from 'typescript'

const MenuButton = () => (
  <div className="absolute inset-y-0 left-0 flex items-center sm:hidden">
    <button type="button" className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white" aria-controls="mobile-menu" aria-expanded="false">
      <svg className="block h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
      </svg>
      <svg className="hidden h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>
  </div>
)

const Logo = () => (
  <div className="flex-shrink-0 flex items-center">
    <Link to="/" className="text-blue-500 text-2xl font-extrabold">Video Platform</Link>
  </div>
)

const Navbar = () => (
  <nav className="bg-gray-800">
    <div className="max-w-7xl mx-auto px-2 sm:px-6 lg:px-8">
      <div className="relative flex items-center justify-between h-16">
        <MenuButton />
        <div className="flex-1 flex items-center justify-center sm:items-stretch sm:justify-between">
          <Logo />
          <div className="hidden sm:block sm:ml-6">
            <div className="flex space-x-4">
              <Link to="/login" className="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium">Login</Link>
              <Link to="/register" className="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium">Register</Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </nav>
)

const Container = ({ children }: React.PropsWithChildren<{}>) => (
  <div className="flex-1 min-w-full bg-gray-200">
    <div className="container mx-auto xl:max-w-screen-xl 2xl:max-w-screen-2xl p-6">
      {children}
    </div>
  </div >
);

const LandingVideos = () => {
  const videos: Video[] = [];

  for (let i = 1; i < 13; i++) {
    videos.push({
      title: `Video about ${i}`,
      description: 'test_description'
    })
  }

  return (
    <>
      <div className="bg-white p-4 rounded-md mb-4">
        <p className="font-bold text-3xl text-gray-700">Latest Videos</p>
      </div>
      <VideoList videos={videos} />
    </>
  )
}

const Footer = () => (
  <footer className="bg-gray-800 p-4">
    <span className="text-md text-gray-100">Copyright TJB {new Date().getFullYear()}</span>
  </footer>
)

const App = () => (
  <div className="flex flex-col min-h-screen">
    <Navbar />
    <Container>
      <Route exact path="/"><LandingVideos /></Route>
      <Route exact path="/login"><Login /></Route>
      <Route exact path="/register"><Register /></Route>
      <Route exact path="/settings"><Settings /></Route>
      <Route path="/video/:id"><VideoDetail filename="idk" title="test_title" description="test_description" /></Route>
      <Route path="/user/:id"><Profile /></Route>
    </Container>
    <Footer />
  </div>
)

export default App;