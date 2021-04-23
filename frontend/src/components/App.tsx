import React from 'react'
import { Route } from 'react-router-dom'

import { Login, Register } from './Auth'
import { VideoDetail, VideoList } from './Video'
import { Profile, Settings } from './User'
import Navbar from './Navbar'
import Upload from './Upload'
import Transcode from './Transcode'

import { Video, initialVideo, initialEncode } from '../types/video'
import { useAppDispatch, useAppSelector } from '../index'
import { getVideoList } from '../store/video'

const Container = ({ children }: React.PropsWithChildren<{}>) => (
  <div className="flex flex-1 min-w-full bg-gray-200">
    <div className="flex flex-grow container mx-auto xl:max-w-screen-xl py-4 px-1">
      {children}
    </div>
  </div>
);

const LandingVideos = () => {
  const dispatch = useAppDispatch()

  dispatch(getVideoList())

  const videos: Video[] = useAppSelector(state => state.video.videoList)

  // for (let i = 1; i < 13; i++) {
  //   let v = { ...initialVideo }
  //   v.title = `Video about ${i}`
  //   v.description = 'test_description'
  //   videos.push(v)
  // }

  return (
    <div className="flex-grow">
      <div className="bg-white p-4 rounded-md mb-4">
        <p className="font-bold text-3xl text-gray-700">Latest Videos</p>
      </div>
      <VideoList videos={videos} />
    </div>
  )
}

const Footer = () => (
  <footer className="bg-gray-800 p-4">
    <span className="text-md text-gray-100">Copyright TJB {new Date().getFullYear()}</span>
  </footer>
)

const App = () => {
  const devitoVideo: Video = {
    id: 0,
    createdAt: '',
    updatedAt: '',
    title: 'How Danny DeVito Eats An Egg on Broadway | Acting Class',
    description: 'Danny DeVito gives an acting class on how he eats an egg while still clearly delivering lines in the play \'The Price\' by Arthur Miller, now on Broadway.',
    userId: 0,
    queueId: 0,
    public: false,
    vstreamId: 0,
    strId: 0,
    fileName: 'devito',
    state: '',
    videoCodec: 'H264',
    width: 0,
    height: 0,
    frameRate: 0.0,
    audioT: [],
    subtitleT: [],
    encData: initialEncode
  }

  return (
    <div className="flex flex-col min-h-screen">
      <Navbar />
      <Container>
        <Route exact path="/"><LandingVideos /></Route>
        <Route exact path="/login"><Login /></Route>
        <Route exact path="/register"><Register /></Route>
        <Route exact path="/settings"><Settings /></Route>
        <Route exact path="/upload"><Upload /></Route>
        <Route exact path="/transcode"><Transcode /></Route>
        <Route path="/video/:id"><VideoDetail video={devitoVideo} /></Route>
        <Route path="/user/:id"><Profile /></Route>
      </Container>
      <Footer />
    </div>
  )
}

export default App;