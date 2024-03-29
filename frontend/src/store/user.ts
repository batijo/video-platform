import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import axios from 'axios'

import APIResponse from '../types/response'
import { User, initialUser } from '../types/user'
import { AppDispatch, AppThunk } from '../index'
import { toCamelCaseObj } from '../utils'

const initialUserList: User[] = [];

export const userSlice = createSlice({
  name: 'users',
  initialState: {
    user: initialUser,
    userList: initialUserList
  },
  reducers: {
    userList: (state, action: PayloadAction<User[]>) => { state.userList = action.payload },
    userDetail: (state, action: PayloadAction<User>) => { state.user = action.payload },
  }
})

export const getUser = (id: number): AppThunk => async (dispatch: AppDispatch, getState) => {
  let headers = { 'Authorization': `Bearer ${getState().auth.token}` }

  axios.get<APIResponse<User>>(`${window.origin}/api/auth/user/${id}`, { headers })
    .then(response => {
      dispatch(userSlice.actions.userDetail(toCamelCaseObj(response.data.data)))
    })
}

export const getUsers = (): AppThunk => async (dispatch: AppDispatch, getState) => {
  let headers = { 'Authorization': `Bearer ${getState().auth.token}` }

  axios.get<APIResponse<User[]>>(`${window.origin}/api/auth/user`, { headers })
    .then(response => {
      dispatch(userSlice.actions.userList(toCamelCaseObj(response.data.data)))
    })
}

export const editProfile = (id: number, user: any): AppThunk => async (dispatch: AppDispatch, getState) => {
  let headers = { 'Authorization': `Bearer ${getState().auth.token}`, 'Content-Type': 'application/json' }

  axios({
    method: 'post',
    url: `${window.origin}/api/auth/user/update/${id}`,
    headers,
    data: user
  })
}

export default userSlice.reducer
