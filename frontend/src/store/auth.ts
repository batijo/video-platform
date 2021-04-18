import { store } from '../index'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import axios from 'axios'

import APIResponse from '../types/response'
import { UserLogin, UserRegister } from '../types/user'

export const authSlice = createSlice({
  name: 'auth',
  initialState: {
    register: {},
    login: {},
    token: ''
  },
  reducers: {
    register: (state, action: PayloadAction<APIResponse<{}>>) => { state.register = action.payload },
    login: (state, action: PayloadAction<APIResponse<string>>) => { state.login = action.payload; state.token = action.payload.data },
  }
})

export const register = (credentials: UserLogin) => {
  axios.post<APIResponse<{}>>(`https://localhost/api/register`)
    .then(response => {
      store.dispatch(authSlice.actions.register(response.data))
    })
}

export const login = (credentials: UserLogin) => {
  axios.post<APIResponse<string>>(`https://localhost/api/login`)
    .then(response => {
      store.dispatch(authSlice.actions.login(response.data))
    })
}

export default authSlice.reducer