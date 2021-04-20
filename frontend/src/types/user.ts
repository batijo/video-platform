export type User = {
  id: number
  createdAt: string
  updatedAt: string
  username: string
  name: string
  lastname: string
  email: string
  admin: boolean
  public: boolean
}

export type UserLogin = {
  email: string
  password: string
}

export type UserRegister = {
  username: string
  email: string
  password: string
  name?: string
  lastname?: string
}

export const initialUser: User = {
  id: 0,
  createdAt: '',
  updatedAt: '',
  username: '',
  name: '',
  lastname: '',
  email: '',
  admin: false,
  public: false,
}

export default User