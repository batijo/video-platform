type User = {
  id: number
  createdAt: Date
  updatedAt: Date
  DeletedAt: Date
  username: string
  name: string
  lastname: string
  email: string
  admin: boolean
  public: boolean
  token: string
}

export default User