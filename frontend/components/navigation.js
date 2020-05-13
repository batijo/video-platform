import Link from 'next/link'

import { 
    Navbar,
    Nav
} from 'react-bootstrap'
  
export default () => (
  <Navbar bg="light" expand="lg">
    <Navbar.Brand>
      <Link href="/">
        <a className="navbar-brand">Video Platform</a>
      </Link>
    </Navbar.Brand>
    <Navbar.Toggle aria-controls="navbar-nav" />
    <Navbar.Collapse id="navbar-nav">
      <Nav className="ml-auto">
        <Navbar.Text>
          <Link href="/login">
            <a className="nav-link">Log in</a>
          </Link>
        </Navbar.Text>
        <Navbar.Text>
          <Link href="/register">
            <a className="nav-link">Sign up</a>
          </Link>
        </Navbar.Text>
      </Nav>
    </Navbar.Collapse>
  </Navbar>
)