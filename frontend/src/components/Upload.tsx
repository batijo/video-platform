import React, { useCallback } from 'react'
import axios from 'axios'
import { useDropzone } from 'react-dropzone'
import { useAppSelector } from '../index'
import { Video, initialVideo } from '../types/video'
import { toCamelCaseObj } from '../utils'
import { Link } from 'react-router-dom'
import SSE from './SSE'

const dropzoneStyle = 'border-dashed border-2 border-gray-300 h-full w-full font-semibold text-blue-400 justify-center cursor-pointer'

const Upload = () => {

  const [isUploaded, setIsUploaded] = React.useState(false)
  const [isUploading, setIsUploading] = React.useState(false)
  const token = useAppSelector(state => state.auth.token)
  const [videoData, setVideoData] = React.useState(initialVideo)

  const onDrop = useCallback(files => {
    if (typeof files[0] === undefined) return

    let formData = new FormData()
    formData.append('file', files[0])

    let headers = {
      'Content-Type': 'multipart/form-data',
      'Authorization': `Bearer ${token}`
    }

    setIsUploading(true)

    axios.post(`${window.origin}/api/auth/upload`, formData, { headers: headers })
      .then(response => {
        // TODO: HANDLE UPLOADED
        setVideoData(toCamelCaseObj(response.data.data))
        setIsUploaded(true)
      })
  }, [])

  const { getRootProps, getInputProps, isDragActive } = useDropzone({ onDrop })

  const uploadUI = () => {
    if (isUploading) return <SSE />
    else return (
      <div className="flex-grow bg-white p-6 rounded-md">
        <div className={dropzoneStyle} {...getRootProps()}>
          <input {...getInputProps()} />
          <div className="tems-center justify-center text-2xl p-4 mx-auto">
            {isDragActive ?
              <p>Drop here</p> :
              <p>Drag and drop your video here, or click to select manually</p>
            }
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col flex-grow">
      <div className="flex flex-row justify-between items-center bg-white p-4 rounded-md mb-4">
        <span className="font-bold text-3xl text-gray-700">Upload Video</span>
        <Link className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded" to="/transcode" >Transcode</Link>
      </div>
      {uploadUI()}
    </div>
  )
}

export default Upload