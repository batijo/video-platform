import React from 'react'
import ReactDOM from 'react-dom'
import { BrowserRouter } from 'react-router-dom'
import { combineReducers, configureStore } from '@reduxjs/toolkit'
import App from './components/App'
import "tailwindcss/tailwind.css"
import './index.css'

import user from './store/user'

const rootReducer = combineReducers({
  user
})

export const store = configureStore({
  reducer: rootReducer
})

store.subscribe(() => console.log(store.getState()))

ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>,
  document.getElementById('root')
);