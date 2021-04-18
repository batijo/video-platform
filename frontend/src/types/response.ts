interface APIResponse<Type> {
  status: boolean
  message: string
  error: string
  data: Type
}

export default APIResponse