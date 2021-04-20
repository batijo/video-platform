import { store } from '../index'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import axios from 'axios'

import APIResponse from '../types/response'
import { User, initialUser } from '../types/user'

const initialUserList: User[] = [];

export const userSlice = createSlice({
  name: 'users',
  initialState: {
    user: initialUser,
    userList: initialUserList
  },
  reducers: {
    userList: (state, action: PayloadAction<User[]>) => { state.userList = action.payload },
    userDetail: (state, action: PayloadAction<User>) => { return { ...action.payload, ...state } },
  }
})

export const getUser = (id: number) => {
  axios.get<APIResponse<User>>(`${window.origin}/api/auth/user/${id}`)
    .then(response => {
      store.dispatch(userSlice.actions.userDetail(response.data.data))
    })
}

export const getUsers = () => {
  axios.get<APIResponse<User[]>>(`${window.origin}/api/auth/user`)
    .then(response => {
      store.dispatch(userSlice.actions.userList(response.data.data))
    })
}

export default userSlice.reducer
