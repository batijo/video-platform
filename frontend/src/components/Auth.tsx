import React from 'react'

import { login, register, resetLogin, resetRegister } from '../store/auth'
import { UserLogin, UserRegister } from '../types/user'
import APIResponse from '../types/response'
import { useAppSelector, useAppDispatch } from '../index'
import { useHistory } from 'react-router-dom'

export const Login = () => {
  const dispatch = useAppDispatch()
  const history = useHistory()

  const [email, setEmail] = React.useState('')
  const [password, setPassword] = React.useState('')
  const [message, setMessage] = React.useState('')
  const [messageColor, setMessageColor] = React.useState('bg-blue-500')

  const handleEmail = (e: { target: HTMLInputElement; }) => setEmail(e.target.value)
  const handlePassword = (e: { target: HTMLInputElement; }) => setPassword(e.target.value)

  const loginResponse = useAppSelector((state) => state.auth.login)

  const handleSubmit = (e: any) => {
    e.preventDefault()

    const credentials: UserLogin = {
      email: email,
      password: password
    }

    dispatch(login(credentials))

    setMessage(loginResponse.message)

    if (loginResponse.status == true) {
      setMessageColor('bg-blue-500')
    } else {
      setMessageColor('bg-red-500')
    }

    setEmail('')
    setPassword('')
    dispatch(resetLogin())

    if (loginResponse.status) history.push("/")
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
                <button className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded" type="submit">Submit</button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </>
  )
}



export const Register = () => (
  <>
    <div className="bg-white p-4 rounded-md mb-4">
      <p className="font-bold text-3xl text-gray-700">Registration Form</p>
    </div>
    <div className="bg-white p-4 rounded-md mb-4">
      <div className="mb-4">
        <div className="p-6 grid grid-cols-2 gap-x-8 gap-y-4">
          <div>
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="username">Username*</label>
            <input className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="Username" />
          </div>
          <div>
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="email">Email*</label>
            <input className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="email" type="email" placeholder="example@example.com" />
          </div>
          <div>
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password">Password*</label>
            <input className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="password" type="password" placeholder="hunter2" />
          </div>
          <div>
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password2">Password*</label>
            <input className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="password2" type="password" placeholder="hunter2" />
          </div>
          <div>
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="name">Name</label>
            <input className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="name" type="text" placeholder="John" />
          </div>
          <div>
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="lastname">Lastname</label>
            <input className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="lastname" type="text" placeholder="Doe" />
          </div>
        </div>
      </div>
    </div>
  </>
)