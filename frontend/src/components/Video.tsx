import React from 'react'
import thumbnail from 'url:../thumbnail.jpg'
import { Link } from 'react-router-dom'
import ReactPlayer from 'react-player'

import Video from '../types/video'
import video from '../store/video'

export const VideoDetail = ({ video }: { video: Video }) => {
  return (
    <div className="flex-grow">
      <div className="bg-white p-4 rounded-md mb-4">
        <p className="font-bold text-3xl text-gray-700">{video.title}</p>
      </div>
      <div className="aspect-w-16 aspect-h-9">
        <ReactPlayer
          width='100%'
          height='100%'
          light={`http://localhost/thumb/${video.fileName}720p.mp4/thumb-5000.jpg`}
          url={`http://localhost/hls/${video.fileName},360p.mp4,480p.mp4,720p.mp4,.en_US.vtt,.urlset/master.m3u8`}
          controls
        />
      </div>
    </div>
  )
}

export const VideoList = ({ videos }: { videos: Video[] }) => {
  const fVideos: Video[] = [...videos].filter(v => v.state === 'transcoded')

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-4">
      {fVideos.map(v => (
        <div className="bg-white rounded-md" key={v.id}>
          <Link to={`/video/${v.id}`}>
            <div className="relative">
              <div className="absolute antialiased font-bold text-white top-2 right-2 bg-gray-800 px-2 rounded-md opacity-80">+</div>
              <div className="absolute text-white bottom-2 right-2 bg-gray-800 px-1 rounded-md opacity-80">00:00</div>
              <img className="rounded-t-md" src={thumbnail} />
            </div>
          </Link>
          <div className="p-4 flex flex-row justify-between">
            <span className="font-bold text-xl" id="title">
              <Link to={`/video/${v.id}`}>{v.fileName}</Link>
            </span>
            <span className="text-md" id="title">2021-07-01</span>
          </div>
        </div>
      ))}
    </div>
  )
}