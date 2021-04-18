import React from 'react'
import ReactDOM from 'react-dom'
import { BrowserRouter } from 'react-router-dom'

import { Provider } from 'react-redux'
import { combineReducers, configureStore } from '@reduxjs/toolkit'

import App from './components/App'
import "tailwindcss/tailwind.css"
import './index.css'

import user from './store/user'
// import video from './store/video'

const rootReducer = combineReducers({
  user
})

export const store = configureStore({
  reducer: rootReducer
})

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