export type User = {
  id: number
  createdAt: Date
  updatedAt: Date
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

export default User