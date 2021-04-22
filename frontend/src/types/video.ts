export type Video = {
  id: number
  createdAt: string
  updatedAt: string
  title: string
  description: string
  userId: number
  public: boolean
  vstreamId: number
  strId: number
  fileName: string
  state: string
  videoCodec: string
  width: number
  height: number
  frameRate: number
  audioT: Array<Audio>
  subtitleT: Array<Subtitle>
  encData: Encode
}

export type Encode = {
  id: number
  createdAt: string
  videoId: number
  strId: number
  fileName: string
  videoCodec: string
  width: number
  height: number
  frameRate: number
  audioT: Array<Audio>
  subtitleT: Array<Subtitle>
}

export type Audio = {
  id: number
  videoId: number
  streamId: number
  codec: string
  language: string
  channels: number
}

export type Subtitle = {
  id: number
  videoId: number
  encId: number
  streamId: number
  language: string
}

export type Preset = {
  id: number
  name: string
  type: number
  resolution: string
  codec: string
  bitrate: string
}


export const initialEncode: Encode = {
  id: 0,
  createdAt: '',
  videoId: 0,
  strId: 0,
  fileName: '',
  videoCodec: '',
  width: 0,
  height: 0,
  frameRate: 0.0,
  audioT: [],
  subtitleT: []
}

export const initialVideo: Video = {
  id: 0,
  createdAt: '',
  updatedAt: '',
  title: '',
  description: '',
  userId: 0,
  public: false,
  vstreamId: 0,
  strId: 0,
  fileName: '',
  state: '',
  videoCodec: '',
  width: 0,
  height: 0,
  frameRate: 0.0,
  audioT: [],
  subtitleT: [],
  encData: initialEncode
}

export const initialPreset: Preset = {
  id: 0,
  name: '',
  type: 0,
  resolution: '',
  codec: '',
  bitrate: ''
}

export default Video