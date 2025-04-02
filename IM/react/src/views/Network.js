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
                      <i className="nc-icon nc-cloud-download-93 text-warning"></i>
                    </div>
                  </Col>
                  <Col xs="8">
                    <div className="numbers">
                      <p className="card-category" style={{ fontSize: '20px', color: 'black'}}>Total number of packet received </p>
                      <Card.Title as="h4">15GB</Card.Title>
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
                      <i className="nc-icon nc-alien-33 text-success"></i>
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
                      <i className="nc-icon nc-delivery-fast text-danger"></i>
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
                      <i className="nc-icon nc-notes text-success" style={{ fontSize: '45px' }}></i>
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
      </Card.Header>
      <Card.Body>
        <div className="ct-chart ct-perfect-fourth" id="chartPreferences">
          <ChartistGraph
            data={{
              labels: ["50%", "10%", "20%", "20%"],
              series: [50, 10, 20, 20],
            }}
            type="Pie"
            options={{
              // Customizing the colors of the slices
              series: [
                { value: 50, className: 'ct-series-a' },
                { value: 10, className: 'ct-series-b' },
                { value: 20, className: 'ct-series-c' },
                { value: 20, className: 'ct-series-d' },
              ],
              chartPadding: 30,
              donut: true, // Optional: To make it a donut chart
              labelOffset: 50,
            }}
          />
        </div>
        <div className="legend">
          <i className="fas fa-circle" style={{ color: '#dc3545', marginLeft: '20px' }}></i> Dos
          <i className="fas fa-circle success" style={{ color: '#17a2b8', marginLeft: '30px' }}></i> SQL Injection
          <i className="fas fa-circle text-danger" style={{ color: '#ffc107', marginLeft: '30px' }}></i> Ping
          <i className="fas fa-circle text-warning" style={{ color: '#28a745', marginLeft: '30px' }}></i> XSS Attack
        </div>
      </Card.Body>
    </Card>
  </Col>


<Col md="6">
            <Card>
              <Card.Header>
                <Card.Title as="h4">Packet number</Card.Title>
              </Card.Header>
              <Card.Body>
                <div className="ct-chart" id="chartActivity">
                  <ChartistGraph
                    data={{
                      labels: [
                        "TCP",
                        "CDP",
                        "ICMP",
                      ],
                      series: [
                        [
                          542,
                          443,
                          320,
                        ],
                      ],
                    }}
                    type="Bar"
                    options={{
                      seriesBarDistance: 10,
                      axisX: {
                        showGrid: false,
                      },
                      height: "245px",
                    }}
                    responsiveOptions={[
                      [
                        "screen and (max-width: 640px)",
                        {
                          seriesBarDistance: 5,
                          axisX: {
                            labelInterpolationFnc: function (value) {
                              return value[0];
                            },
                          },
                        },
                      ],
                    ]}
                  />
                </div>
              </Card.Body>
              <Card.Footer>
                <div className="legend">
                  <i className="fas fa-circle text-info"></i>
                  Number of packets
                </div>
              </Card.Footer>
            </Card>
          </Col>
</Row>
      </Container>
    </>
  );
}

export default Network;
