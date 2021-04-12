import React from 'react'
import thumbnail from 'url:../thumbnail.jpg'
import { Link, Route, Switch } from 'react-router-dom'


export const Video: React.FC = () => (
  <>
    <div className="bg-white p-4 rounded-md mb-4">
      <p className="font-bold text-3xl text-gray-700">Video about ??</p>
    </div>
    <div className="flex flex-row">
      <div className="rounded-md flex-1">
        <img className="rounded-md" src={thumbnail} />
      </div>
      <div className="bg-white rounded-md p-4 md:ml-4 flex-grow">
        <p className="font-bold text-xl text-gray-700">info</p>
      </div>
    </div>
  </>
)

export const VideoList: React.FC = () => {
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
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
      {videos}
    </div>
  )
}