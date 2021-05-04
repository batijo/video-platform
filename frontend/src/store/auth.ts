import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import axios from 'axios'

import APIResponse from '../types/response'
import { UserLogin, UserRegister } from '../types/user'
import { AppDispatch, AppThunk } from '../index'

const initialRegister: APIResponse<{}> = {
  status: false,
  message: '',
  error: '',
  data: {}
}

const initialLogin: APIResponse<string> = {
  status: false,
  message: '',
  error: '',
  data: ''
}

export const authSlice = createSlice({
  name: 'auth',
  initialState: {
    register: initialRegister,
    login: initialLogin,
    token: ''
  },
  reducers: {
    register: (state, action: PayloadAction<APIResponse<{}>>) => { state.register = action.payload },
    resetRegister: (state) => { state.register = { ...initialRegister } },
    login: (state, action: PayloadAction<APIResponse<string>>) => {
      state.login = action.payload
      if (action.payload.data !== null && action.payload.data !== undefined) state.token = action.payload.data
    },
    logout: (state) => { state.token = '' },
    resetLogin: (state) => { state.login = { ...initialLogin } }
  }
})

export const resetRegister = (): AppThunk => async (dispatch: AppDispatch) => dispatch(authSlice.actions.resetRegister())
export const resetLogin = (): AppThunk => async (dispatch: AppDispatch) => dispatch(authSlice.actions.resetLogin())

export const register = (credentials: UserRegister): AppThunk => async (dispatch: AppDispatch) => {
  axios.post<APIResponse<{}>>(`${window.origin}/api/register`, credentials)
    .then(response => {
      dispatch(authSlice.actions.register(response.data))
    }).catch(error => dispatch(authSlice.actions.register(error.response.data)))
}

export const login = (credentials: UserLogin): AppThunk => async (dispatch: AppDispatch) => {
  axios.post<APIResponse<string>>(`${window.origin}/api/login`, credentials)
    .then(response => {
      dispatch(authSlice.actions.login(response.data))
    }).catch(error => dispatch(authSlice.actions.login(error.response.data)))
}

export const logout = (): AppThunk => async (dispatch: AppDispatch, getState) => {
  let token = getState().auth.token
  let headers = { 'Authorization': `Bearer ${token}` }

  axios.post<APIResponse<{}>>(`${window.origin}/api/auth/logout`, null, { headers })
    .then(response => {
      dispatch(authSlice.actions.logout())
    })
}

export default authSlice.reducer