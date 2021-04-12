import { store } from '../index'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import axios from 'axios'

type User = {
  username: string
}

const initialState: User[] = [];

export const userSlice = createSlice({
  name: 'users',
  initialState,
  reducers: {
    userList: (state, action: PayloadAction<User[]>) => action.payload,
    userDetail: (state, action: PayloadAction<User>) => { state.push(action.payload) },
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
