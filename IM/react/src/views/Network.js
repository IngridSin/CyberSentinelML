import React from "react";
import ChartistGraph from "react-chartist";
// react-bootstrap components
import {
  Badge,
  Button,
  Card,
  Navbar,
  Nav,
  Table,
  Container,
  Row,
  Col,
  Form,
  OverlayTrigger,
  Tooltip,
} from "react-bootstrap";

function Network() {
  return (
    <>
      <Container fluid>
        <Row>


          <Col lg="5.5" sm="6">
            <Card className="card-stats">
              <Card.Body>
                <Row>
                  <Col xs="3">
                    <div className="icon-big text-center icon-warning">
                      <i className="nc-icon nc-chart text-warning"></i>
                    </div>
                  </Col>
                  <Col xs="8">
                    <div className="numbers">
                      <p className="card-category" style={{ fontSize: '20px', color: 'black'}}>Total number of packet received </p>
                      <Card.Title as="h4">150GB</Card.Title>
                    </div>
                  </Col>
                </Row>
              </Card.Body>
              <Card.Footer>
              </Card.Footer>
            </Card>
          </Col>


          <Col lg="6.5" sm="6">
            <Card className="card-stats">
              <Card.Body>
                <Row>
                  <Col xs="2">
                    <div className="icon-big text-center icon-warning">
                      <i className="nc-icon nc-light-3 text-success"></i>
                    </div>
                  </Col>
                  <Col xs="10">
                    <div className="numbers">
                      <p className="card-category" style={{ fontSize: '20px', color: 'black'}}>Total number of suspected packet received</p>
                      <Card.Title as="h4">$ 1,345</Card.Title>
                    </div>
                  </Col>
                </Row>
              </Card.Body>
              <Card.Footer>
              </Card.Footer>
            </Card>
          </Col>

          </Row>


        <Row>

          <Col lg="6" sm="5">
            <Card className="card-stats">
              <Card.Body>
                <Row>
                  <Col xs="3">
                    <div className="icon-big text-center icon-warning">
                      <i className="nc-icon nc-vector text-danger"></i>
                    </div>
                  </Col>
                  <Col xs="8">
                    <div className="numbers">
                      <p className="card-category" style={{ fontSize: '20px', color: 'black'}}>Total number of packet received </p>
                      <Card.Title as="h4">150GB</Card.Title>
                    </div>
                  </Col>
                </Row>
              </Card.Body>
              <Card.Footer>
              </Card.Footer>
            </Card>
          </Col>


          <Col lg="6" sm="6">
            <Card className="card-stats">
              <Card.Body>
                <Row>
                  <Col xs="2">
                    <div className="icon-big text-center icon-warning">
                      <i className="nc-icon nc-light-3 text-success" style={{ fontSize: '45px' }}></i>
                    </div>
                  </Col>
                  <Col xs="10">
                    <div className="numbers">
                      <p className="card-category" style={{ fontSize: '20px', color: 'black'}}>Attack type detected last time</p>
                      <Card.Title as="h4">$ 1,345</Card.Title>
                    </div>
                  </Col>
                </Row>
              </Card.Body>
              <Card.Footer>
              </Card.Footer>
            </Card>
          </Col>

          </Row>

<Row>
          <Col md="4">
            <Card>
              <Card.Header>
                <Card.Title as="h4">Packet Statistics</Card.Title>
               {/*  <p className="card-category">Last Campaign Performance</p> */}
              </Card.Header>
              <Card.Body>
                <div
                  className="ct-chart ct-perfect-fourth"
                  id="chartPreferences"
                >
                  <ChartistGraph
                    data={{
                      labels: ["60%", "40%"],
                      series: [60,40],
                    }}
                    type="Pie"
                  />
                </div>
                <div className="legend">
                  <i className="fas fa-circle text-info"  style={{ marginLeft: '20px' }}></i>
                  Good Email
                  <i className="fas fa-circle text-danger" style={{ marginLeft: '30px' }}></i>
                  Suspected Email
                </div>
              </Card.Body>
            </Card>
          </Col>

    </Row>

      </Container>
    </>
  );
}

export default Network;
