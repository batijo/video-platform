import Layout from '../components/layout'
import {
  Form,
  Button
} from 'react-bootstrap'

export default () => (
  <Layout>
    <h1>Registration Form</h1>

    <Form>
      <Form.Group>
        <Form.Label>Email address</Form.Label>
        <Form.Control type="email"></Form.Control>
      </Form.Group>

      <Form.Group>
        <Form.Label>Password</Form.Label>
        <Form.Control type="password"></Form.Control>
      </Form.Group>

      <Button variant="primary" type="submit">
        Sign up
      </Button>
    </Form>
  </Layout>
)