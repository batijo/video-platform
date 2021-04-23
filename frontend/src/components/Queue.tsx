import React from 'react'
import { useAppDispatch, useAppSelector } from '../index'
import { getVideoQueue } from '../store/video'
import { Queue } from '../types/video'

const VideoQueue = () => {
  const dispatch = useAppDispatch()
  React.useEffect(() => dispatch(getVideoQueue()), [])
  const queue: Queue[] = useAppSelector(state => state.video.videoQueue) ?? []

  return (
    <div className="flex flex-col flex-grow">
      <div className="flex flex-row justify-between items-center bg-white p-4 rounded-md mb-4">
        <span className="font-bold text-3xl text-gray-700">Queue</span>
      </div>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 text-lg">
        {queue.length < 1 &&
          <div className="lg:col-span-2 justify-between items-center bg-white rounded-md p-4">The Queue is currently empty.</div>
        }
        {queue.map((v, i) =>
          <div className="flex flex-row justify-between items-center bg-white rounded-md p-4">
            <span className="bg-gray-500 text-white px-3 py-1 rounded-md font-semibold">{v.position}</span>
            <span>{v.videoTitle ?? ''}</span>
          </div>
        )}
      </div>
    </div>
  )
}

export default VideoQueue