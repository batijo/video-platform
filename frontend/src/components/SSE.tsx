import React from 'react'
import axios from 'axios'
import { useAppSelector } from '../index'

const useEventSource = (url: string) => {
  const [data, updateData] = React.useState('')

  React.useEffect(() => {
    const source = new EventSource(url)

    // source.onmessage = (event) => {
    //   console.log(event.data)
    //   updateData(event.data)
    // }

    source.onmessage = function logEvents(event) {
      updateData(event.data);
    }
  }, [])

  return data
}

const SSE = () => {
  const token = useAppSelector(state => state.auth.token)
  const data = useEventSource(`${window.origin}/api/sse/dashboard/?token=${token}`)

  console.log(data)

  return (
    <div>{data}</div>
  )
}

export default SSE