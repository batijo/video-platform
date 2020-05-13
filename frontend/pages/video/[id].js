import { useRouter } from 'next/router'
import dynamic from 'next/dynamic';
import Layout from '../../components/layout'
import { Jumbotron } from 'react-bootstrap'

const ShakaPlayer = dynamic(() => import('shaka-player-react'), { ssr: false });

export default () => {
  const router = useRouter()
  const { id } = router.query
  
  return (
    <Layout>
      <Jumbotron>
        <h1>Video {id}</h1>
        <ShakaPlayer 
          src="http://rdmedia.bbc.co.uk/dash/ondemand/bbb/2/client_manifest-high_profile-common_init.mpd"
          width={1920}
          autoPlay
        />
      </Jumbotron>
    </Layout>
  )
}