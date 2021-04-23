import React from 'react'
import { useParams } from 'react-router-dom'
import { VideoList } from './Video'
import { Video, initialVideo } from '../types/video'
import { User } from '../types/user'
import { useAppDispatch, useAppSelector } from '../index'
import { getUser } from '../store/user'
import { getUserVideoList } from '../store/video'


export const Profile = () => {
  const { id }: any = useParams()
  const dispatch = useAppDispatch()

  React.useEffect(() => {
    dispatch(getUser(id))
    dispatch(getUserVideoList(id))
  }, [])
  const user: User = useAppSelector(state => state.users.user)
  const videos: Video[] = useAppSelector(state => state.video.userVideoList)

  return (
    <div className="flex-grow">
      <div className="bg-white p-4 rounded-md mb-4">
        <span className="font-bold text-3xl text-gray-700 block border-gray-200 border-b-2 pb-2 mb-4">User {user.username} [{user.email}]</span>
        <p className="text-xl text-gray-600">TODO: User Profile</p>
      </div>
      <VideoList videos={videos} />
    </div>
  )
}

export const Settings = () => (
  <div></div>
)