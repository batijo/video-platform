import { store } from '../index'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import axios from 'axios'

import { Video, Encode, initialVideo, initialEncode } from '../types/video'
import APIResponse from '../types/response'

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
  axios.get<APIResponse<Video>>(`${window.origin}/api/auth/video/${id}`)
    .then(response => {
      store.dispatch(videoSlice.actions.videoDetail(response.data.data))
    })
}

export default videoSlice.reducer