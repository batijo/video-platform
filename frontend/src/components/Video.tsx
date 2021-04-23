import React from 'react'
import { Link, useParams } from 'react-router-dom'
import ReactPlayer from 'react-player'

import Video from '../types/video'
import { getVideo } from '../store/video'
import { useAppDispatch, useAppSelector } from '../index'

const thumbnail = (video: Video) => {
  if (video.resolutions !== null && video.resolutions !== undefined)
    return `http://localhost/thumb/${video.fileName}${video.resolutions[0]}.mp4/thumb-5000.jpg`
  else
    return `http://localhost/thumb/${video.fileName}/thumb-5000.jpg`
}

const url = (video: Video) => {
  if (video.resolutions !== null && video.resolutions !== undefined)
    return `http://localhost/hls/${video.fileName},${video.resolutions.join('.mp4,')}.mp4,.urlset/master.m3u8`
  else
    return `http://localhost/hls/${video.fileName},.urlset/master.m3u8`
}

const duration = (video: Video) => {
  return `${Math.floor(Math.floor(video.duration) / 60)}:${Math.floor(video.duration) % 60}`
}

export const VideoDetail = () => {
  const dispatch = useAppDispatch()
  const { id }: any = useParams()

  React.useEffect(() => dispatch(getVideo(id)), [])
  const video = useAppSelector<Video>(state => state.video.video)

  return (
    <div className="flex-grow">
      <div className="bg-white p-4 rounded-md mb-4">
        <p className="font-bold text-3xl text-gray-700">{video.fileName}</p>
      </div>
      <div className="aspect-w-16 aspect-h-9">
        <ReactPlayer
          width='100%'
          height='100%'
          // light={thumbnail(video)}
          url={url(video)}
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
              <div className="relative aspect-w-16 aspect-h-9">
                <img className="rounded-t-md" src={thumbnail(v)} />
              </div>
              <div className="absolute antialiased font-bold text-white top-2 right-2 bg-gray-800 px-2 rounded-md opacity-80">+</div>
              <div className="absolute text-white bottom-2 right-2 bg-gray-800 px-1 rounded-md opacity-80">{duration(v)}</div>
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