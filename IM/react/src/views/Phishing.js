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




function Phishing() {
  return (
    <>
      <Container fluid>
        <Row>

          <Col lg="4" sm="6">
            <Card className="card-stats">
              <Card.Body>
                <Row>
                  <Col xs="3">
                    <div className="icon-big text-center icon-warning">
                      <i className="nc-icon nc-email-83 text-warning"></i>
                    </div>
                  </Col>
                  <Col xs="8">
                    <div className="numbers">
                      <p className="card-category" style={{ fontSize: '20px', color: 'black'}}>Total number of email received </p>
                      <Card.Title as="h4">150GB</Card.Title>
                    </div>
                  </Col>
                </Row>
              </Card.Body>
              <Card.Footer>
              </Card.Footer>
            </Card>
          </Col>


          <Col lg="4" sm="6">
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
                      <p className="card-category" style={{ fontSize: '20px', color: 'black'}}>Total number of suspected email received</p>
                      <Card.Title as="h4">$ 1,355</Card.Title>
                    </div>
                  </Col>
                </Row>
              </Card.Body>
              <Card.Footer>
              </Card.Footer>
            </Card>
          </Col>



          <Col lg="4" sm="6">
            <Card className="card-stats">
              <Card.Body>
                <Row>
                  <Col xs="2">
                    <div className="icon-big text-center icon-warning">
                      <i className="nc-icon nc-watch-time text-danger"></i>
                    </div>
                  </Col>
                  <Col xs="10">
                    <div className="numbers">
                      <p className="card-category" style={{ fontSize: '20px', color: 'black'}}>Last time a suspected email detected</p>
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
                <Card.Title as="h4">Email Statistics</Card.Title>
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


      <Col md="8">
        <Card>
          <Card.Body>
            <Form.Group>
              <label style={{ textTransform: 'none', fontSize: '20px', color: 'black' }}> Content of the suspected email detected last time</label>
              <Form.Control
                cols="80"
                defaultValue="CONTENT"
                placeholder="content"
                rows="14"
                as="textarea"
              />
            </Form.Group>
          </Card.Body>
        </Card>
      </Col>
    </Row>

      </Container>
    </>
  );
}

export default Phishing;
