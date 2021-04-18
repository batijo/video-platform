import { store } from '../index'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import axios from 'axios'

import User from '../types/user'
// import Video from '../types/video'

const initialUser: User = {
  id: 0,
  createdAt: new Date(),
  updatedAt: new Date(),
  DeletedAt: new Date(),
  username: '',
  name: '',
  lastname: '',
  email: '',
  admin: false,
  public: false,
  token: '',
}

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

export const getUsers = () => {
  axios.get<User[]>('https://localhost/api/auth/user')
    .then(response => {
      store.dispatch(userSlice.actions.userList(response.data))
    })
}

export const getUser = (id: number) => {
  axios.get<User>(`https://localhost/api/auth/user/${id}`)
    .then(response => {
      store.dispatch(userSlice.actions.userDetail(response.data))
    })
}

export default userSlice.reducer
