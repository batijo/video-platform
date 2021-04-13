import React from 'react'
import thumbnail from 'url:../thumbnail.jpg'


export const Stream: React.FC = () => (
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