import React, { useCallback } from 'react'
import axios from 'axios'
import { useDropzone } from 'react-dropzone'
import { useAppSelector } from '../index'
import { Video, initialVideo } from '../types/video'
import { toCamelCaseObj } from '../utils'
import { Link } from 'react-router-dom'

const dropzoneStyle = 'border-dashed border-2 border-gray-300 h-full w-full font-semibold text-blue-400 justify-center cursor-pointer'

const Upload = () => {

  const [isUploaded, setIsUploaded] = React.useState(false)
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

    axios.post(`${window.origin}/api/auth/upload`, formData, { headers: headers })
      .then(response => {
        // TODO: HANDLE UPLOADED
        console.log(response.data)
        setVideoData(toCamelCaseObj(response.data.data))
        setIsUploaded(true)
      })
  }, [])

  const { getRootProps, getInputProps, isDragActive } = useDropzone({ onDrop })

  return (
    <div className="flex flex-col flex-grow">
      <div className="bg-white p-4 rounded-md mb-4">
        <p className="font-bold text-3xl text-gray-700">Upload Video</p>
        <Link to="/transcode">Transcode</Link>
      </div>
      <div className="flex-grow bg-white p-6 rounded-md">
        {isUploaded ?
          <div>
            <p>
              Uploaded<br />
              {videoData.title}<br />
              {videoData.videoCodec}<br />
              {videoData.width}x{videoData.height}
            </p>
          </div> :
          <div className={dropzoneStyle} {...getRootProps()}>
            <input {...getInputProps()} />
            <div className="tems-center justify-center text-2xl p-4 mx-auto">
              {isDragActive ?
                <p>Drop here</p> :
                <p>Drag and drop your video here, or click to select manually</p>
              }
            </div>
          </div>
        }


      </div>
    </div>
  )
}

export default Upload