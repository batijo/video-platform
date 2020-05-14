import React, { createRef } from 'react'
import axios from 'axios'
import Layout from '../../components/layout'
import {
  Jumbotron,
  Form,
  Button,
  Row,
  Col
} from 'react-bootstrap'

const PresetMenu = () => (
  <div>
    <Form.Group controlId="video-presets-1">
      <Form.Label>Video presets</Form.Label>
      <Form.Control as="select" custom>
        <option selected value="nochange">None</option>
      </Form.Control>
    </Form.Group>

    <Form.Group controlId="audio-presets-1">
      <Form.Label>Audio presets</Form.Label>
      <Form.Control as="select" custom>
        <option selected value="nochange">None</option>
      </Form.Control>
    </Form.Group>

    <Form.Group controlId="audio-select-1">
      <Form.Label>Audio Tracks</Form.Label>
      <Form.Control as="select" custom multiple={true}>
        <option selected value="nochange">None</option>
      </Form.Control>
    </Form.Group>


    <Form.Group controlId="subtitle-select-1">
      <Form.Label>Subtitles</Form.Label>
      <Form.Control as="select" custom multiple={true}>
        <option selected value="nochange">None</option>
      </Form.Control>
    </Form.Group>
  </div>
)

export default class Upload extends React.Component {
  constructor(props) {
    super(props)

    this.handleUpload = this.handleUpload.bind(this);
    this.handleCodec = this.handleCodec.bind(this);
    this.handleResolution = this.handleResolution.bind(this);
    this.handleFramerate = this.handleFramerate.bind(this);
    this.handleAudio = this.handleAudio.bind(this);
    this.handleAudioChannels = this.handleAudioChannels.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);

    this.state = {
      video: '',
      audio: '',
      audioChannels: '',
      uploaded: false,
      form: {
        FileName: '',
        StrId: '',
        VideoCodec: '',
        FrameRate: '',
        Width: '',
        Height: '',
        AudioT: [],
        SubtitleT: [] 
      }
    }
  }

  handleUpload(event) {
    let file = event.target.files[0]
    let name = event.target.value.lastIndexOf('\\')

    if (!file) return

    let formData = new FormData()
    formData.append('file', file)

    axios.post('/upload', formData)
    .then((res) => {
      let form = {...this.state.form }

      form.FileName = name
      form.VtId = res.data['videotrack'][0]['Index']

      this.setState({
        video: res.data,
        uploaded: true,
        form
      })
    })
    .catch((err) => {
      console.log(err)
    })
  }

  handleCodec(event) {
    let form = {...this.state.form}
    form.VideoCodec = event.target.value
    this.setState({ form })
  }

  handleResolution(event) {
    let form = {...this.state.form}
    let resolution = event.target.value.split(':')
    form.Width = resolution[0]
    form.Height = resolution[1]
    this.setState({ form })
  }

  handleFramerate(event) {
    let form = {...this.state.form}
    form.FrameRate = event.target.value
    this.setState({ form })
  }

  handleAudio(event) {
    this.setState({ audio: event.target.value })
  }

  handleAudioChannels(event) {
    this.setState({ audioChannels: event.target.value })
  }

  handleSubmit() {
    let video = this.state.video
    let form = {...this.state.form}
    let audioTrack = video['audiotrack'].findIndex(i => i['index'] == this.state.audio)
    let audio = {
      AtId: '',
      AtCodec: '',
      Language: '',
      Channels: ''
    }

    if (!form.VtCodec)
      form.VtCodec = video['videotrack'][0]['codecName']

    if (!form.Width)
      form.Width = video['videotrack'][0]['width']

    if (!form.Height)
      form.Height = video['videotrack'][0]['height']

    if (!form.FrameRate)
      form.FrameRate = video['videotrack'][0]['frameRate']

    event.preventDefault();
  }

  render() {
    return (
      <Layout>
        <Jumbotron>
          <h1>Video Upload</h1>

          <Form method="post" encType="multipart/form-data">
            <Form.Group>
              <Form.Label>File Upload</Form.Label>
              <Form.File
                  onChange={this.handleUpload}
                  disabled={this.state.uploaded}
                  id="input-file"
                  label={this.state.form.FileName}
                  custom
              />
            </Form.Group>
          </Form>

          <Form onSubmit={this.handleSubmit}>
            <Row>
              <Col xs={12} md={6}>
                <Form.Group>
                  <Form.Label>Codec</Form.Label>
                  <Form.Control as="select" custom onChange={this.handleCodec}>
                    <option selected value="">Keep</option>
                    <option value="h265">H265</option>
                    <option value="h264">H264</option>
                  </Form.Control>
                </Form.Group>
              </Col>
  
              <Col xs={12} md={6}>
                <Form.Group controlId="resolution-select">
                  <Form.Label>Resolution</Form.Label>
                  <Form.Control as="select" custom onChange={this.handleResolution}>
                    <option selected value="">Keep</option>
                    <option value="1920:1080">1080p</option>
                    <option value="1280:720">720p</option>
                    <option value="858:480">480p</option>
                    <option value="480:360">360p</option>
                  </Form.Control>
                </Form.Group>
              </Col>
  
  
              <Col xs={12} md={6}>
                <Form.Group controlId="fr-select">
                  <Form.Label>Framerate</Form.Label>
                  <Form.Control as="select" custom onChange={this.handleFramerate}>
                    <option selected value="">Keep</option>
                    <option value="25.0">25</option>
                  </Form.Control>
                </Form.Group>
              </Col>
  
              <Col xs={12} md={6}>
                <Form.Group controlId="audio-select">
                  <Form.Label>Audio Codec</Form.Label>
                  <Form.Control as="select" custom onChange={this.handleAudio}>
                    <option selected value="">Keep</option>
                    <option value="aac">AAC</option>
                  </Form.Control>
                </Form.Group>
              </Col>
  
              <Col xs={12} md={6}>
                <Form.Group controlId="channels-select">
                  <Form.Label>Audio Channels</Form.Label>
                  <Form.Control as="select" custom onChange={this.handleAudioChannels}>
                    <option selected value="">Keep</option>
                    <option value="1">Mono</option>
                    <option value="2">Stereo</option>
                  </Form.Control>
                </Form.Group>
              </Col>
  
              <Col xs={12}>
                <Button
                  onChange={this.handleSubmit}
                  disabled={!this.state.uploaded}
                  variant="primary"
                  type="submit">
                  Upload
                </Button>
              </Col>
            </Row>
          </Form>
        </Jumbotron>
      </Layout>
    )
  }
}