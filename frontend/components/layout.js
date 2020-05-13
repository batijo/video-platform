import Head from 'next/head'
import { Container } from 'react-bootstrap'
import Navigation from './navigation'

export default ({ children, fluid }) => (
    <main>
      <Head>
        <title>Video Platform</title>
        <link rel="icon" href="/favicon.ico" />
        <link
          rel="stylesheet"
          href="https://maxcdn.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css"
          integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh"
          crossOrigin="anonymous"
        />
      </Head>

    <Navigation/>
    <Container className="mt-4" fluid={fluid}>
      {children}
    </Container>
  </main>
)