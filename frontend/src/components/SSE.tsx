import React from 'react'
import axios from 'axios'

const useEventSource = (url: string) => {
  const [data, updateData] = React.useState('')

  React.useEffect(() => {
    const source = new EventSource(url)

    source.onmessage = (event) => updateData(event.data)
  }, [])
}

const SSE = () => {

  // const updateSSE = async () => {
  //   let source = new EventSource(`${window.origin}/api/sse/dashboard`);
  // }

  // React.useEffect(() => { const runSSE = async () => await updateSSE() }, [])

  return (
    <div></div>
  )
}

export default SSE