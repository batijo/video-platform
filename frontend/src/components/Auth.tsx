import React from 'react'

import { login, register, resetLogin, resetRegister } from '../store/auth'
import { UserLogin, UserRegister } from '../types/user'
import { useAppSelector, useAppDispatch } from '../index'
import { useHistory } from 'react-router-dom'

export const Login = () => {
  const dispatch = useAppDispatch()
  const history = useHistory()

  const [email, setEmail] = React.useState('')
  const [password, setPassword] = React.useState('')

  const [message, setMessage] = React.useState('')
  const [messageColor, setMessageColor] = React.useState('bg-blue-500')

  const handleEmail = (e: { target: HTMLInputElement }) => setEmail(e.target.value)
  const handlePassword = (e: { target: HTMLInputElement }) => setPassword(e.target.value)

  const loginResponse = useAppSelector(state => state.auth.login)

  const handleSubmit = (e: any) => {
    e.preventDefault()

    const credentials: UserLogin = {
      email,
      password
    }

    dispatch(login(credentials))
    setMessage(loginResponse.message)

    if (loginResponse.status === true) setMessageColor('bg-blue-500')
    else setMessageColor('bg-red-500')

    setEmail('')
    setPassword('')
    dispatch(resetLogin())

    if (loginResponse.status === true) history.push("/")
  }

  return (
    <>
      <div className="bg-white p-4 rounded-md mb-4">
        <p className="font-bold text-3xl text-gray-700">Login Form</p>
      </div>
      <div className="bg-white p-4 rounded-md mb-4">
        <div className="mb-4">
          <form onSubmit={handleSubmit}>
            <div className="p-6 grid grid-cols-2 gap-x-8 gap-y-4">
              {message !== '' &&
                <div className={`${messageColor} col-span-2 p-4 rounded text-white`}>
                  {message}
                </div>
              }
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="email">Email</label>
                <input value={email} onChange={handleEmail} className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="email" type="email" placeholder="example@example.com" />
              </div>
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password">Password</label>
                <input value={password} onChange={handlePassword} className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="password" type="password" placeholder="hunter2" />
              </div>
              <div>
                <button className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded" type="submit">Sign in</button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </>
  )
}



export const Register = () => {
  const dispatch = useAppDispatch()
  const history = useHistory()

  const [username, setUsername] = React.useState('')
  const [email, setEmail] = React.useState('')
  const [password, setPassword] = React.useState('')
  const [passwordRepeat, setPasswordRepeat] = React.useState('')
  const [name, setName] = React.useState('')
  const [lastname, setLastname] = React.useState('')

  const [message, setMessage] = React.useState('')
  const [messageColor, setMessageColor] = React.useState('bg-blue-500')

  const handleUsername = (e: { target: HTMLInputElement }) => setUsername(e.target.value)
  const handleEmail = (e: { target: HTMLInputElement }) => setEmail(e.target.value)
  const handlePassword = (e: { target: HTMLInputElement }) => setPassword(e.target.value)
  const handlePasswordRepeat = (e: { target: HTMLInputElement }) => setPasswordRepeat(e.target.value)
  const handleName = (e: { target: HTMLInputElement }) => setName(e.target.value)
  const handleLastname = (e: { target: HTMLInputElement }) => setLastname(e.target.value)

  const registerResponse = useAppSelector((state) => state.auth.register)

  const handleSubmit = (e: any) => {
    e.preventDefault()

    if (password !== passwordRepeat) {
      history.push('/register')
      return
    }

    const credentials: UserRegister = {
      username,
      email,
      password,
      name,
      lastname,
    }

    dispatch(register(credentials))
    setMessage(registerResponse.message)

    if (registerResponse.status === true) setMessageColor('bg-blue-500')
    else setMessageColor('bg-red-500')

    setUsername('')
    setEmail('')
    setPassword('')
    setPasswordRepeat('')
    setName('')
    setLastname('')

    dispatch(resetRegister())

    if (registerResponse.status === true) history.push('/login')
  }

  return (
    <>
      <div className="bg-white p-4 rounded-md mb-4">
        <p className="font-bold text-3xl text-gray-700">Registration Form</p>
      </div>
      <div className="bg-white p-4 rounded-md mb-4">
        <div className="mb-4">
          <form onSubmit={handleSubmit}>
            <div className="p-6 grid grid-cols-2 gap-x-8 gap-y-4">
              {message !== '' &&
                <div className={`${messageColor} col-span-2 p-4 rounded text-white`}>
                  {message}
                </div>
              }
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="username">Username*</label>
                <input value={username} onChange={handleUsername} className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="Username" />
              </div>
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="email">Email*</label>
                <input value={email} onChange={handleEmail} className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="email" type="email" placeholder="example@example.com" />
              </div>
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password">Password*</label>
                <input value={password} onChange={handlePassword} className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="password" type="password" placeholder="hunter2" />
              </div>
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password2">Password*</label>
                <input value={passwordRepeat} onChange={handlePasswordRepeat} className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="password2" type="password" placeholder="hunter2" />
              </div>
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="name">Name</label>
                <input value={name} onChange={handleName} className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="name" type="text" placeholder="John" />
              </div>
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="lastname">Lastname</label>
                <input value={lastname} onChange={handleLastname} className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="lastname" type="text" placeholder="Doe" />
              </div>
              <div>
                <button className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded" type="submit">Sign up</button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </>
  )
}