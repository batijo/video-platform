import Layout from '../../components/layout'
import {
  Jumbotron,
  Form,
  Button
} from 'react-bootstrap'

export default () => (
  <Layout>
    <Jumbotron>
      <h1>Video Upload</h1>
      <Form>
        <Form.File 
            id="video-file"
            label="Upload your video"
            custom
        />
        <Button variant="primary" type="submit">
          Upload
        </Button>
      </Form>
    </Jumbotron>
  </Layout>
)