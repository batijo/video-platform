import React from 'react'
import { Link, Route, Switch } from 'react-router-dom'

import 'tailwindcss/tailwind.css'
import '../index.css'

import thumbnail from 'url:../thumbnail.jpg'

import Login from './Login'
import Register from './Register'
import Video from './Video'

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

const Container: React.FC = props => (
  <div className="flex-1 min-w-full bg-gray-200">
    <div className="container mx-auto xl:max-w-screen-xl p-6">
      {props.children}
    </div>
  </div >
);

const LandingVideos: React.FC = () => {
  const videos = [];

  for (let i = 1; i < 10; i++) {
    videos.push(
      <div className=" bg-white rounded-md">
        <Link to={`video/${i}`}>
          <div className="relative">
            <div className="absolute antialiased font-bold text-white top-2 right-2 bg-gray-800 px-2 rounded-md opacity-80">+</div>
            <div className="absolute text-white bottom-2 right-2 bg-gray-800 px-1 rounded-md opacity-80">00:00</div>
            <img className="rounded-t-md" src={thumbnail} />
          </div>
        </Link>
        <div className="p-4 flex flex-row justify-between">
          <span className="font-bold text-xl" id="title">
            <Link to="/video">Video about {i}</Link>
          </span>
          <span className="text-md" id="title">2021-07-0{i}</span>
        </div>
      </div>
    )
  }

  return (
    <div className="">
      <div className="bg-white p-4 rounded-md mb-4">
        <p className="font-bold text-3xl text-gray-700">Latest Videos</p>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {videos}
      </div>
    </div>
  )
}

const Footer = () => (
  <footer className="bg-gray-800 p-4">
    <span className="text-md text-gray-100">Copyright TJB {new Date().getFullYear()}</span>
  </footer>
)

const App: React.FC = () => {
  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <Container>
        <Route exact path="/"><LandingVideos /></Route>
        <Route exact path="/login"><Login /></Route>
        <Route exact path="/register"><Register /></Route>

        <Route path="/video"><Video /></Route>
      </Container>
      <Footer />
    </div>
  )
}

export default App;