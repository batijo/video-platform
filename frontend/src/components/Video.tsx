import React from 'react'
import thumbnail from 'url:../thumbnail.jpg'

const Video: React.FC = () => (
  <>
    <div className="bg-white p-4 rounded-md mb-4">
      <p className="font-bold text-3xl text-gray-700">Video about ??</p>
    </div>
    <div className="bg-white rounded-md mb-4">
      <img className="rounded-md md:rounded-r-none" src={thumbnail} />
    </div>
  </>
)

export default Video