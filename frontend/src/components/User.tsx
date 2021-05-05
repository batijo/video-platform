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
  const token = useAppSelector(state => state.auth.token)

  let tokenInfo = {
    user_id: '',
    email: '',
    admin: false,
    exp: 0
  }

  if (token !== '') tokenInfo = JSON.parse(atob(token.split('.')[1]))

  const handleEdit = () => {

  }

  const [name, setName] = React.useState(user.name ?? 'unspecified')
  const [lastname, setLastname] = React.useState(user.lastname ?? 'unspecified')
  const [isPublic, setIsPublic] = React.useState(user.public)
  const [isEdit, setIsEdit] = React.useState(false)

  const inputStyleConst = 'form-input mt-1 block w-full rounded-md border-transparent focus:border-gray-500 focus:bg-white focus:ring-0'
  const inputStyle = `${inputStyleConst} ${!isEdit ? 'bg-white' : 'bg-gray-100'}`

  return (
    <div className="flex-grow">
      <div className="bg-white p-4 rounded-md mb-4">
        <div className="flex justify-between pb-2 mb-4 border-b-2">
          <span className="font-bold text-3xl text-gray-700 block border-gray-200">User {user.username}</span>
          <button onClick={() => setIsEdit(!isEdit)} className="bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded">Edit</button>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 gap-x-8">
          <div className="flex flex-row justify-between items-center">
            <label className="pr-4" htmlFor="lastname">Username:</label>
            <input className={inputStyleConst} disabled={!isEdit} value={user.username}></input>
          </div>
          <div className="flex flex-row justify-between items-center">
            <label className="pr-4" htmlFor="lastname">Email:</label>
            <input className={inputStyleConst} disabled={!isEdit} value={user.email}></input>
          </div>
          <div className="flex flex-row justify-between items-center">
            <label className="pr-4" htmlFor="name" >Name:</label>
            <input className={inputStyle} disabled={!isEdit} value={name}></input>
          </div>
          <div className="flex flex-row justify-between items-center">
            <label className="pr-4" htmlFor="lastname">Lastname:</label>
            <input className={inputStyle} disabled={!isEdit} value={lastname}></input>
          </div>
          <div className="flex flex-row justify-between items-center">
            <label className="pr-4" htmlFor="lastname">Public:</label>
            <input onClick={() => setIsPublic(!isPublic)} checked={isPublic} disabled={!isEdit} className="form-checkbox rounded bg-gray-200 border-transparent focus:border-transparent text-gray-700 focus:ring-1 focus:ring-offset-2 focus:ring-gray-500 p-3" type="checkbox" name="preset-mode" />
          </div>
          <div className="flex flex-row justify-between items-center">
            <label htmlFor="lastname" className="pr-4">Admin:</label>
            <input checked={user.admin} className="form-checkbox rounded bg-gray-200 border-transparent focus:border-transparent text-gray-700 focus:ring-1 focus:ring-offset-2 focus:ring-gray-500 p-3" type="checkbox" name="preset-mode" />
          </div>
          {isEdit &&
            <div>
              <button onClick={handleEdit} className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded" type="submit">Save</button>
            </div>
          }
        </div>
      </div>
      <VideoList videos={videos} />
    </div>
  )
}

export const Settings = () => (
  <div></div>
)