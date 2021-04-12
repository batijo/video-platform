import React from 'react'
import { VideoList } from './Video'

export const Profile: React.FC = () => (
  <>
    <div className="bg-white p-4 rounded-md mb-4">
      <span className="font-bold text-3xl text-gray-700 block border-gray-200 border-b-2 pb-2 mb-4">User 123456</span>
      <p className="text-xl text-gray-600">blah blah</p>
    </div>
    <VideoList />
  </>
)

export const Settings: React.FC = () => (
  <div></div>
)