import React from 'react'
import { Route } from 'react-router-dom'

import { Login, Register } from './Auth'
import { VideoDetail, VideoList, Video } from './Video'
import { Profile, Settings } from './User'
import Navbar from './Navbar'

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