import React from 'react'
import axios from 'axios'
import Layout from '../components/layout'
import {
  Form,
  Button
} from 'react-bootstrap'

export default class Register extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      name: '',
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

    axios.post(`${process.env.API_ROOT}/register`, this.state)
    .then((res) => {
      console.log(res.data)
    })
  }

  render() {
    return (
      <Layout>
        <h1>Register</h1>
        <Form onSubmit={this.handleSubmit}>

          <Form.Group>
            <Form.Label>Name</Form.Label>
            <Form.Control name="name" onChange={this.handleInput}></Form.Control>
          </Form.Group>

          <Form.Group>
            <Form.Label>Email address</Form.Label>
            <Form.Control type="email" name="email" onChange={this.handleInput}></Form.Control>
          </Form.Group>
    
          <Form.Group>
            <Form.Label>Password</Form.Label>
            <Form.Control type="password" name="password" onChange={this.handleInput}></Form.Control>
          </Form.Group>
    
          <Button variant="primary" type="submit">
            Sign up
          </Button>

        </Form>
      </Layout>
    )
  }
}