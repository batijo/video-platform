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

export default Video