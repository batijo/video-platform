import Layout from '../components/layout'
import { Jumbotron, Row, Col } from 'react-bootstrap';

export default () => (
  <Layout>
    <h1>Featured Users</h1>

    <Row>
      <Col xs={6}>
        <Jumbotron>
          <h2>User</h2>
        </Jumbotron>
      </Col>
      <Col xs={6}>
        <Jumbotron>
          <h2>User</h2>
        </Jumbotron>
      </Col>
      <Col xs={6}>
        <Jumbotron>
          <h2>User</h2>
        </Jumbotron>
      </Col>
      <Col xs={6}>
        <Jumbotron>
          <h2>User</h2>
        </Jumbotron>
      </Col>
    </Row>

  </Layout>
)