import React from 'react'
import axios from 'axios'
import Layout from '../components/layout'
import {
  Form,
  Button
} from 'react-bootstrap'

// const JWT = React.createContext('jwt')

export default class Login extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      email: '',
      password: ''
    }

    this.handleInput = this.handleInput.bind(this)
    this.handleSubmit = this.handleSubmit.bind(this)
  }

  handleInput(event) {
    this.setState({ [event.target.name]: event.target.value })
  }

  handleSubmit(event) {
    event.preventDefault()

    axios.post(`${process.env.API_ROOT}/login`, this.state)
    .then((res) => {
      console.log(res.data)
    })
  }

  render() {
    return (
      <Layout>
        <h1>Login</h1>
        <Form onSubmit={this.handleSubmit}>

          <Form.Group>
            <Form.Label>Email address</Form.Label>
            <Form.Control type="email" name="email" onChange={this.handleInput}></Form.Control>
          </Form.Group>
    
          <Form.Group>
            <Form.Label>Password</Form.Label>
            <Form.Control type="password" name="password" onChange={this.handleInput}></Form.Control>
          </Form.Group>
    
          <Button variant="primary" type="submit">
            Log in
          </Button>

        </Form>
      </Layout>
    )
  }
}