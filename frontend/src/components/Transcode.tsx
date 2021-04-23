import React from 'react'
import { Preset } from '../types/video'
import { useAppSelector, useAppDispatch } from '../index'
import { getUserVideoList } from '../store/video'
import { Video, Audio, Subtitle, Encode } from '../types/video'
import APIResponse from '../types/response'
import axios from 'axios'

const selectStyle = 'block w-full mt-1 rounded-md bg-gray-100 border-transparent focus:border-gray-500 focus:bg-white focus:ring-0'

const videoPresets: Preset[] = [
  {
    id: 1,
    name: "V_240p_H264_400",
    type: 0,
    resolution: "240p",
    codec: "h264",
    bitrate: "400"
  },
  {
    id: 3,
    name: "V_576p_HEVC_1200",
    type: 0,
    resolution: "576p",
    codec: "h265",
    bitrate: "1200"
  },
  {
    id: 4,
    name: "V_720p_HEVC_2800",
    type: 0,
    resolution: "720p",
    codec: "h265",
    bitrate: "2800"
  },
  {
    id: 5,
    name: "V_1080p_HEVC_3300",
    type: 0,
    resolution: "1080p",
    codec: "h265",
    bitrate: "3300"
  },
  {
    id: 7,
    name: "V_720p_H264_4500",
    type: 0,
    resolution: "720p",
    codec: "h264",
    bitrate: "4500"
  },
  {
    id: 8,
    name: "V_576p_H264_1800",
    type: 0,
    resolution: "576p",
    codec: "h264",
    bitrate: "1800"
  },
  {
    id: 9,
    name: "V_360p_HEVC_1000",
    type: 0,
    resolution: "360p",
    codec: "h265",
    bitrate: "800"
  },
  {
    id: 10,
    name: "V_ORIG_HEVC_1200",
    type: 0,
    resolution: "default",
    codec: "h265",
    bitrate: "1200"
  },
  {
    id: 11,
    name: "V_ORIG_HEVC_2000",
    type: 0,
    resolution: "default",
    codec: "h265",
    bitrate: "2000"
  },
  {
    id: 12,
    name: "V_1080p_HEVC_6600",
    type: 0,
    resolution: "1080p",
    codec: "h265",
    bitrate: "6600"
  }
]

const audioPresets: Preset[] = [
  {
    id: 2,
    name: "A_AAC_128",
    type: 1,
    resolution: "",
    codec: "aac",
    bitrate: "128"
  },
  {
    id: 6,
    name: "A_AAC_96",
    type: 1,
    resolution: "",
    codec: "aac",
    bitrate: "96"
  }
]

const Transcode = () => {
  const dispatch = useAppDispatch()
  const token = useAppSelector(state => state.auth.token)

  let user = {
    user_id: 0,
    email: '',
    admin: false,
    exp: 0
  }

  if (token !== '') user = JSON.parse(atob(token.split('.')[1]))

  const [manual, setManual] = React.useState(true)
  const handleManual = () => setManual(!manual)
  const videoList: Video[] = useAppSelector(state => state.video.userVideoList)
  React.useEffect(() => dispatch(getUserVideoList(user.user_id)), [])

  const [video, setVideo] = React.useState('none')
  const [videoCodec, setVideoCodec] = React.useState('nochange')
  const [audioCodec, setAudioCodec] = React.useState('nochange')
  const [resolution, setResolution] = React.useState('nochange')
  const [framerate, setFramerate] = React.useState('nochange')
  const [audioTracks, setAudioTracks] = React.useState<Audio[]>([])
  const [subtitleTracks, setSubtitleTracks] = React.useState<Subtitle[]>([])

  const handleSubmit = (e: any) => {
    e.preventDefault()

    let headers = { 'Authorization': `Bearer ${token}` }

    if (manual) {
      let videoDetail: Video = videoList.filter(v => String(v.id) === video)[0]
      let encode: Encode = {
        videoId: videoDetail.id,
        fileName: videoDetail.fileName,
        videoCodec: videoCodec,
        width: Number(resolution.split(':')[0]),
        height: Number(resolution.split(':')[1]),
        frameRate: Number(framerate),
        audioT: audioTracks,
        subtitleT: subtitleTracks
      }

      encode.videoId = videoDetail.id
      encode.fileName = videoDetail.fileName
      encode.videoCodec = videoCodec
      encode.width = Number(resolution.split(':')[0])
      encode.height = Number(resolution.split(':')[1])
      encode.frameRate = Number(framerate)
      encode.audioT = audioTracks
      encode.subtitleT = subtitleTracks

      console.log(encode)

      axios.post<APIResponse<{}>>(`${window.origin}/api/auth/transcode/${video}`, encode, { headers })
        .then(response => {
          console.log(response.data)
        })
    } else {
      console.log('no')
    }
  }

  const presetUI = () => {
    if (manual) return (
      <>
        <div>
          <label htmlFor="video-codec-select">Video Codec</label>
          <select
            className={`form-select ${selectStyle}`} name="video-codec-select" defaultValue="nochange"
            onChange={e => setVideoCodec(e.target.value)} value={videoCodec}
          >
            <option value="nochange">Keep current</option>
            <option value="h264">H264</option>
            <option value="h265">H265</option>
            <option value="vp9">VP9</option>
          </select>
        </div>
        <div>
          <label htmlFor="audio-codec-select">Audio Codec</label>
          <select
            className={`form-select ${selectStyle}`} name="audio-codec-select" defaultValue="nochange"
            onChange={e => setAudioCodec(e.target.value)} value={audioCodec}
          >
            <option value="nochange">Keep current</option>
            <option value="aac">AAC</option>
          </select>
        </div>
        <div>
          <label htmlFor="resolution-select">Resolution</label>
          <select
            className={`form-select ${selectStyle}`} name="resolution-select" defaultValue="nochange"
            onChange={e => setResolution(e.target.value)} value={resolution}
          >
            <option value="nochange">Keep current</option>
            <option value="1920:1080">1080p</option>
            <option value="1280:720">720p</option>
            <option value="858:480">480p</option>
            <option value="480:360">360p</option>
          </select>
        </div>
        <div>
          <label htmlFor="framerate-select">Framerate</label>
          <select
            className={`form-select ${selectStyle}`} name="framerate-select" defaultValue="nochange"
            onChange={e => setFramerate(e.target.value)} value={framerate}
          >
            <option value="nochange">Keep current</option>
            <option value="50">50</option>
            <option value="25">25</option>
            <option value="60">60</option>
            <option value="24">24</option>
          </select>
        </div>
      </>
    )
    else return (
      <>
        <div>
          <label htmlFor="video-presets">Video Presets</label>
          <select className={`form-multiselect ${selectStyle}`} multiple name="video-presets">
            {videoPresets.map(preset =>
              <option key={preset.id} value={preset.id}>{preset.name}</option>
            )}
          </select>
        </div>
        <div>
          <label htmlFor="audio-presets">Audio Presets</label>
          <select className={`form-multiselect ${selectStyle}`} multiple name="audio-presets">
            {audioPresets.map(preset =>
              <option key={preset.id} value={preset.id}>{preset.name}</option>
            )}
          </select>
        </div>
      </>
    )
  }

  return (
    <div className="flex flex-col flex-grow">
      <div className="flex flex-row bg-white p-4 rounded-md mb-4 items-center justify-between">
        <p className="font-bold text-3xl text-gray-700">Transcode Video</p>

        <div className="flex items-center">
          <label className="pr-4 text-lg" htmlFor="preset-mode">Manual Presets?</label>
          <input onClick={handleManual} checked={manual} className={`form-checkbox rounded bg-gray-200 border-transparent focus:border-transparent text-gray-700 focus:ring-1 focus:ring-offset-2 focus:ring-gray-500 p-3`} type="checkbox" name="preset-mode" />
        </div>
      </div>
      <div className=" bg-white p-6 rounded-md">
        <form onSubmit={handleSubmit}>
          <div className="grid grid-cols-2 gap-4 text-lg">
            <div className="col-span-2">
              <label htmlFor="video-select">Select Video</label>
              <select
                className={`form-select ${selectStyle}`} name="video-select" defaultValue="none"
                onChange={e => setVideo(e.target.value)} value={video}
              >
                <option value="none">---</option>
                {videoList.map(video =>
                  <option key={video.id} value={video.id}>{video.fileName}</option>
                )}
              </select>
            </div>
            {presetUI()}
            <div>
              <label htmlFor="audio-tracks">Audio Tracks</label>
              <select className={`form-multiselect ${selectStyle}`} multiple name="audio-tracks">
                <option>---</option>
                <option>---</option>
                <option>---</option>
                <option>---</option>
              </select>
            </div>
            <div>
              <label htmlFor="audio-tracks">Subtitle Tracks</label>
              <select className={`form-multiselect ${selectStyle}`} multiple name="audio-tracks">
                <option>---</option>
                <option>---</option>
                <option>---</option>
                <option>---</option>
              </select>
            </div>
            <div>
              <button
                className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded"
                type="submit" disabled={video === 'none' ? true : false}
              >
                Begin Transcoding
              </button>
            </div>
          </div>
        </form>
      </div>
    </div >
  )
}

export default Transcode