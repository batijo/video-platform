type Video = {
  title: string
  description: string
  userID: number
  public: boolean
  strID: number
  filename: string
  state: string
  codec: string
  width: number
  height: number
  framerate: number
}

type Audio = {
  id: number
  videoID: number
  streamID: number
  codec: string,
  language: string,
  channels: number
}



export default Video