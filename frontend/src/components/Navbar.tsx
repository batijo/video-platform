import React from 'react'

import { Link } from 'react-router-dom'
import { useAppSelector, useAppDispatch } from '../index'
import { logout } from '../store/auth'
import { useHistory } from 'react-router-dom'

const MenuButton = () => (
  <div className="absolute inset-y-0 left-0 flex items-center sm:hidden">
    <button type="button" className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white" aria-controls="mobile-menu" aria-expanded="false">
      <svg className="block h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16M4 18h16" />
      </svg>
      <svg className="hidden h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>
  </div>
)

const Logo = () => (
  <div className="flex-shrink-0 flex items-center">
    <Link to="/" className="text-blue-400 text-2xl font-extrabold">Video Platform</Link>
  </div>
)

const Navbar = () => {
  const dispatch = useAppDispatch()
  const history = useHistory()
  const token = useAppSelector(state => state.auth.token)

  let tokenInfo = {
    user_id: '',
    email: '',
    admin: false,
    exp: 0
  }

  if (token !== '') tokenInfo = JSON.parse(atob(token.split('.')[1]))

  const authMenu = () => {
    if (token === '') {
      return (
        <>
          <Link to="/login" className="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium">Sign in</Link>
          <Link to="/register" className="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium border border-white">Sign up</Link>
        </>
      )
    } else {
      return (
        <>
          <Link to="/upload" className="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium border border-white">Upload</Link>
          <Link to={`/user/${tokenInfo.user_id}`} className="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium">{tokenInfo.email}</Link>
          <button onClick={handleLogout} className="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium">Log out</button>
        </>
      )
    }
  }

  const handleLogout = () => {
    dispatch(logout())
    history.push('/')
  }

  return (
    <nav className="bg-gray-800">
      <div className="max-w-7xl mx-auto px-2 sm:px-6 lg:px-8">
        <div className="relative flex items-center justify-between h-16">
          <MenuButton />
          <div className="flex-1 flex items-center justify-center sm:items-stretch sm:justify-between">
            <Logo />
            <div className="hidden sm:block sm:ml-6">
              <div className="flex space-x-4">
                {authMenu()}
              </div>
            </div>
          </div>
        </div>
      </div>
    </nav>
  )
}

export default Navbar