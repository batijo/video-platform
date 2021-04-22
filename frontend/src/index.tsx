import React from 'react'
import ReactDOM from 'react-dom'
import { BrowserRouter } from 'react-router-dom'
import { Provider, TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux'
import { combineReducers, configureStore, Action } from '@reduxjs/toolkit'
import { ThunkAction } from 'redux-thunk'

import App from './components/App'
import "tailwindcss/tailwind.css"
import './index.css'

import auth from './store/auth'
import users from './store/user'
import video from './store/video'

export const store = configureStore({
  reducer: combineReducers({
    auth,
    users,
    video
  }),
  middleware: (getDefaultMiddleware) => getDefaultMiddleware()
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch


export const useAppDispatch = () => useDispatch<AppDispatch>()
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector
export type AppThunk = ThunkAction<void, RootState, unknown, Action>

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </Provider>
  </React.StrictMode>,
  document.getElementById('root')
);