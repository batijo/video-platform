import { useRouter } from 'next/router'
import Error from 'next/error'
import useSWR from 'swr'
import dynamic from 'next/dynamic';
import Layout from '../../../components/layout'
import {
  Jumbotron,
  Row,
  Col
} from 'react-bootstrap'

const ShakaPlayer = dynamic(() => import('shaka-player-react'), { ssr: false });

export default () => {
  const router = useRouter()
  const { username } = router.query

  const { data, error } = useSWR(`https://${process.env.API_ROOT}/user/${username}/stream`, fetch)

  
  // if (error || !data)
  //   return <Error statusCode={404} />

  return (
    <Layout fluid>
      <Row>
        <Col xs={8}>
        <Jumbotron>
          <h1>{username}'s Stream</h1>
          <ShakaPlayer 
            src="http://rdmedia.bbc.co.uk/dash/ondemand/bbb/2/client_manifest-high_profile-common_init.mpd"
            width={1920}
            autoPlay
          />
        </Jumbotron>
        </Col>
        <Col xs={4}>
          <Jumbotron>
            <h1>Chat</h1>
          </Jumbotron>
        </Col>
      </Row>
    </Layout>
  )
}