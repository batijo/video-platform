import { store } from '../index'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import axios from 'axios'

import { Video, Encode, initialVideo, initialEncode, Queue } from '../types/video'
import APIResponse from '../types/response'
import { AppDispatch, AppThunk } from '../index'
import { toCamelCaseObj } from '../utils'

const initialVideoList: Video[] = []
const initialVideoQueue: Queue[] = []

export const videoSlice = createSlice({
  name: 'videos',
  initialState: {
    video: initialVideo,
    videoList: initialVideoList,
    userVideoList: initialVideoList,
    videoQueue: initialVideoQueue
  },
  reducers: {
    videoDetail: (state, action: PayloadAction<Video>) => { state.video = action.payload },
    videoList: (state, action: PayloadAction<Video[]>) => { state.videoList = action.payload },
    userVideoList: (state, action: PayloadAction<Video[]>) => { state.userVideoList = action.payload },
    videoQueue: (state, action: PayloadAction<any>) => { state.videoQueue = action.payload.elements ?? [] }
  }
})

export const getVideo = (id: number): AppThunk => async (dispatch: AppDispatch, getState) => {
  let token = getState().auth.token
  let headers = { 'Authorization': `Bearer ${token}` }

  axios.get<APIResponse<Video>>(`${window.origin}/api/auth/video/${id}`, { headers })
    .then(response => {
      dispatch(videoSlice.actions.videoDetail(toCamelCaseObj(response.data.data)))
    })
}

export const getVideoList = (): AppThunk => async (dispatch: AppDispatch, getState) => {
  let token = getState().auth.token
  let headers = { 'Authorization': `Bearer ${token}` }

  axios.get<APIResponse<Video[]>>(`${window.origin}/api/auth/video`, { headers })
    .then(response => {
      dispatch(videoSlice.actions.videoList(toCamelCaseObj(response.data.data)))
    })
}

export const getUserVideoList = (id: number): AppThunk => async (dispatch: AppDispatch, getState) => {
  let token = getState().auth.token
  let headers = { 'Authorization': `Bearer ${token}` }

  axios.get<APIResponse<Video[]>>(`${window.origin}/api/auth/video/user/${id}`, { headers })
    .then(response => {
      dispatch(videoSlice.actions.userVideoList(toCamelCaseObj(response.data.data)))
    })
}

export const getVideoQueue = (): AppThunk => async (dispatch: AppDispatch, getState) => {
  let token = getState().auth.token
  let headers = { 'Authorization': `Bearer ${token}` }

  axios.get<APIResponse<any>>(`${window.origin}/api/auth/queue`, { headers })
    .then(response => {
      dispatch(videoSlice.actions.videoQueue(toCamelCaseObj(response.data.data)))
    })
}

export default videoSlice.reducer