import React from 'react'
import axios from 'axios'
import { useAppSelector } from '../index'

// const useEventSource = (url: string) => {
//   const [data, updateData] = React.useState('')

//   React.useEffect(() => {
//     const source = new EventSource(url)

//     source.onmessage = (event) => {
//       console.log(event.data)
//       updateData(data + event.data)
//     }
//   }, [])

//   return data
// }

const SSE = () => {
  const token = useAppSelector(state => state.auth.token)
  const [data, updateData] = React.useState('')

  console.log(data)

  return (
    <div>{data}</div>
  )
}

export default SSE