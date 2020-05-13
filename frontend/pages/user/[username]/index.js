import { useRouter } from 'next/router'
import Error from 'next/error'
import useSWR from 'swr'
import Layout from '../../../components/layout'
import {
  Jumbotron,
  Row,
  Col
} from 'react-bootstrap'

export default () => {
  const router = useRouter()
  const { username } = router.query

  const { data, error } = useSWR(`https://${process.env.API_ROOT}/user/${username}`, fetch)

  
  // if (error || !data)
  //   return <Error statusCode={404} />

  return (
    <Layout>
      <h1>{username}'s Profile</h1>
      <Row>
        <Col xs={12}>
          <Jumbotron>
            <h2>Description</h2>
          </Jumbotron>
        </Col>

        <Col xs={6}>
          <Jumbotron>
            <h2>Stream</h2>
          </Jumbotron>
        </Col>

        <Col xs={6}>
          <Jumbotron>
            <h2>Latest Video</h2>
          </Jumbotron>
        </Col>
      </Row>
    </Layout>
  )
}