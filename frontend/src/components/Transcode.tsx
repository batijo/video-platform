import React from 'react'

const selectStyle = 'block w-full mt-1 rounded-md bg-gray-100 border-transparent focus:border-gray-500 focus:bg-white focus:ring-0'

const Transcode = () => {
  const [manual, setManual] = React.useState(false)
  const handleManual = () => setManual(!manual)

  const presetUI = () => {
    if (manual) return (
      <>
        <div>
          <label htmlFor="codec-select">Video Codec</label>
          <select className={`form-select ${selectStyle}`} name="codec-select">
            <option>---</option>
            <option>---</option>
            <option>---</option>
            <option>---</option>
          </select>
        </div>
        <div>
          <label htmlFor="framerate-select">Audio Codec</label>
          <select className={`form-select ${selectStyle}`} name="framerate-select">
            <option>---</option>
            <option>---</option>
            <option>---</option>
            <option>---</option>
          </select>
        </div>
        <div>
          <label htmlFor="resolution-select">Resolution</label>
          <select className={`form-select ${selectStyle}`} name="resolution-select">
            <option>---</option>
            <option>---</option>
            <option>---</option>
            <option>---</option>
          </select>
        </div>
        <div>
          <label htmlFor="framerate-select">Framerate</label>
          <select className={`form-select ${selectStyle}`} name="framerate-select">
            <option>---</option>
            <option>---</option>
            <option>---</option>
            <option>---</option>
          </select>
        </div>
        {/* <div>
          <label htmlFor="framerate-select">Audio Channels</label>
          <select className={`form-select ${selectStyle}`} name="framerate-select">
            <option>---</option>
            <option>---</option>
            <option>---</option>
            <option>---</option>
          </select>
        </div> */}
      </>
    )
    else return (
      <>
        <div>
          <label htmlFor="video-presets">Video Presets</label>
          <select className={`form-multiselect ${selectStyle}`} multiple name="video-presets">
            <option>---</option>
            <option>---</option>
            <option>---</option>
            <option>---</option>
          </select>
        </div>
        <div>
          <label htmlFor="audio-presets">Audio Presets</label>
          <select className={`form-multiselect ${selectStyle}`} multiple name="audio-presets">
            <option>---</option>
            <option>---</option>
            <option>---</option>
            <option>---</option>
          </select>
        </div>
      </>
    )
  }

  return (
    <div className="flex flex-col flex-grow">
      <div className="flex flex-row bg-white p-4 rounded-md mb-4 items-center justify-between">
        <p className="font-bold text-3xl text-gray-700">Transcode Video</p>

        <div className="flex items-center">
          <label className="pr-4 text-lg" htmlFor="preset-mode">Manual Presets?</label>
          <input onClick={handleManual} className={`form-checkbox rounded bg-gray-200 border-transparent focus:border-transparent text-gray-700 focus:ring-1 focus:ring-offset-2 focus:ring-gray-500 p-3`} type="checkbox" name="preset-mode" />
        </div>
      </div>
      <div className=" bg-white p-6 rounded-md">
        <div className="grid grid-cols-2 gap-4 text-lg">
          <div className="col-span-2">
            <label htmlFor="video-select">Select Video</label>
            <select className={`form-select ${selectStyle}`} name="video-select">
              <option>---</option>
              <option>---</option>
              <option>---</option>
              <option>---</option>
            </select>
          </div>
          {presetUI()}
          <div>
            <label htmlFor="audio-tracks">Audio Tracks</label>
            <select className={`form-multiselect ${selectStyle}`} multiple name="audio-tracks">
              <option>---</option>
              <option>---</option>
              <option>---</option>
              <option>---</option>
            </select>
          </div>
          <div>
            <label htmlFor="audio-tracks">Subtitle Tracks</label>
            <select className={`form-multiselect ${selectStyle}`} multiple name="audio-tracks">
              <option>---</option>
              <option>---</option>
              <option>---</option>
              <option>---</option>
            </select>
          </div>
          <div>
            <button className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded" type="submit">Begin Transcoding</button>
          </div>
        </div>
      </div>
    </div >
  )
}

export default Transcode