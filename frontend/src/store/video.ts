import { store } from '../index'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import axios from 'axios'

import { Video, Encode } from '../types/video'

const initialEncode: Encode = {
  id: 0,
  createdAt: new Date(),
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

const initialVideo: Video = {
  id: 0,
  createdAt: new Date(),
  updatedAt: new Date(),
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

const initialVideoList: Video[] = [];

export const videoSlice = createSlice({
  name: 'videos',
  initialState: {
    video: initialVideo,
    videoList: initialVideoList
  },
  reducers: {
    videoDetail: (state, action: PayloadAction<Video>) => { state.video = action.payload }
    // TODO (?): videoList:
  }
})

export const getVideo = (id: number) => {
  axios.get<Video>(`https://localhost/api/auth/video/${id}`)
    .then(response => {
      store.dispatch(videoSlice.actions.videoDetail(response.data))
    })
}

export default videoSlice.reducer