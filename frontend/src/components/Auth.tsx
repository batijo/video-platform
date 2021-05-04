import React from 'react'

import { login, register, resetLogin, resetRegister } from '../store/auth'
import { UserLogin, UserRegister } from '../types/user'
import { useAppSelector, useAppDispatch } from '../index'
import { useHistory } from 'react-router-dom'

const inputStyle = 'form-input mt-1 block w-full rounded-md bg-gray-100 border-transparent focus:border-gray-500 focus:bg-white focus:ring-0'

export const Login = () => {
  const dispatch = useAppDispatch()
  const history = useHistory()

  const [email, setEmail] = React.useState('')
  const [password, setPassword] = React.useState('')

  const handleEmail = (e: { target: HTMLInputElement }) => setEmail(e.target.value)
  const handlePassword = (e: { target: HTMLInputElement }) => setPassword(e.target.value)

  const loginResponse = useAppSelector(state => state.auth.login)
  const token = useAppSelector(state => state.auth.token)

  React.useEffect(() => {
    dispatch(resetLogin())
    if (token !== '') history.push('/')
  }, [token])

  const handleSubmit = (e: any) => {
    e.preventDefault()

    const credentials: UserLogin = {
      email,
      password
    }

    dispatch(login(credentials))

    setEmail('')
    setPassword('')
    dispatch(resetLogin())

    if (token !== '') history.push('/')
  }

  return (
    <div className="flex-grow">
      <div className="bg-white p-4 rounded-md mb-4">
        <p className="font-bold text-3xl text-gray-700">Login Form</p>
      </div>
      <div className="bg-white p-4 rounded-md mb-4">
        <div className="mb-4">
          <form onSubmit={handleSubmit}>
            <div className="p-6 grid grid-cols-2 gap-x-8 gap-y-4">
              {loginResponse.message !== '' &&
                <div className={`${loginResponse.status ? 'bg-blue-500' : 'bg-red-500'} col-span-2 p-4 rounded text-white`}>
                  {loginResponse.message}
                </div>
              }
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="email">Email</label>
                <input className={inputStyle} value={email} onChange={handleEmail} id="email" type="email" placeholder="example@example.com" />
              </div>
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password">Password</label>
                <input className={inputStyle} value={password} onChange={handlePassword} id="password" type="password" placeholder="hunter2" />
              </div>
              <div>
                <button className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded" type="submit">Sign in</button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
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

  const handleUsername = (e: { target: HTMLInputElement }) => setUsername(e.target.value)
  const handleEmail = (e: { target: HTMLInputElement }) => setEmail(e.target.value)
  const handlePassword = (e: { target: HTMLInputElement }) => setPassword(e.target.value)
  const handlePasswordRepeat = (e: { target: HTMLInputElement }) => setPasswordRepeat(e.target.value)
  const handleName = (e: { target: HTMLInputElement }) => setName(e.target.value)
  const handleLastname = (e: { target: HTMLInputElement }) => setLastname(e.target.value)

  const registerResponse = useAppSelector((state) => state.auth.register)

  React.useEffect(() => dispatch(resetRegister()), [])
  React.useEffect(() => {
    if (registerResponse.status === true) history.push('/')
  }, [registerResponse])

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
    <div className="flex-grow">
      <div className="bg-white p-4 rounded-md mb-4">
        <p className="font-bold text-3xl text-gray-700">Registration Form</p>
      </div>
      <div className="bg-white p-4 rounded-md mb-4">
        <div className="mb-4">
          <form onSubmit={handleSubmit}>
            <div className="p-6 grid grid-cols-2 gap-x-8 gap-y-4">
              {registerResponse.message !== '' &&
                <div className={`${registerResponse.status ? 'bg-blue-500' : 'bg-red-500'} col-span-2 p-4 rounded text-white`}>
                  {registerResponse.message}
                </div>
              }
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="name">Name</label>
                <input className={inputStyle} value={name} onChange={handleName} id="name" type="text" placeholder="John" />
              </div>
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="lastname">Lastname</label>
                <input className={inputStyle} value={lastname} onChange={handleLastname} id="lastname" type="text" placeholder="Doe" />
              </div>
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="username">Username*</label>
                {/* <input value={username} onChange={handleUsername} className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="Username" /> */}
                <input className={inputStyle} value={username} onChange={handleUsername} id="username" type="text" placeholder="Username" />
              </div>
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="email">Email*</label>
                <input className={inputStyle} value={email} onChange={handleEmail} id="email" type="email" placeholder="example@example.com" />
              </div>
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password">Password*</label>
                <input className={inputStyle} value={password} onChange={handlePassword} id="password" type="password" placeholder="hunter2" />
              </div>
              <div>
                <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password2">Repeat Password*</label>
                <input className={inputStyle} value={passwordRepeat} onChange={handlePasswordRepeat} id="password2" type="password" placeholder="hunter2" />
              </div>
              <div>
                <button className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded" type="submit">Sign up</button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}