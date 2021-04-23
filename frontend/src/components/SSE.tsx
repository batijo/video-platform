import React from 'react'
import axios from 'axios'
import { useAppSelector } from '../index'

const SSE = () => {
  const token = useAppSelector(state => state.auth.token)
  const [data, updateData] = React.useState<string[]>([])

  const dataRef = React.useRef<string[]>([])
  dataRef.current = data

  React.useEffect(() => {
    const source = new EventSource(`${window.origin}/api/sse/dashboard/${token}`)
    source.onmessage = (event) => updateData(arr => [...arr, event.data])
  }, [])

  return (
    <div className="flex-grow max-h-96 bg-gray-700 text-white p-4 rounded-md font-mono overflow-x-auto">
      {data.map(line =>
        <span className="block">{line}</span>
      )}
    </div>
  )
}

export default SSE