import React from "react";
import ChartistGraph from "react-chartist";
import {
  Card,
  Container,
  Row,
  Col,
} from "react-bootstrap";

import useNetworkStore from "../store/networkStore";
import useStatsSocket from "../hooks/useEmailStatSocket";

function Network() {
  useStatsSocket();

  const {
    totalFlows,
    maliciousFlows,
    lastMaliciousFlow,
  } = useNetworkStore();

  const safeTotal = totalFlows > 0 ? totalFlows : 1; // avoid divide-by-zero
  const goodFlows = totalFlows - maliciousFlows;

  return (
    <Container fluid>
      <Row className="mb-4">
        <Col lg="6" sm="6">
          <Card className="card-stats">
            <Card.Body>
              <Row>
                <Col xs="3">
                  <div className="icon-big text-center icon-warning">
                    <i className="nc-icon nc-chart text-primary" />
                  </div>
                </Col>
                <Col xs="9">
                  <div className="numbers">
                    <p className="card-category" style={{ fontSize: '20px', color: 'black' }}>
                      Total Flows Captured
                    </p>
                    <Card.Title as="h4">{totalFlows}</Card.Title>
                  </div>
                </Col>
              </Row>
            </Card.Body>
          </Card>
        </Col>

        <Col lg="6" sm="6">
          <Card className="card-stats">
            <Card.Body>
              <Row>
                <Col xs="3">
                  <div className="icon-big text-center icon-danger">
                    <i className="nc-icon nc-lock-circle-open text-danger" />
                  </div>
                </Col>
                <Col xs="9">
                  <div className="numbers">
                    <p className="card-category" style={{ fontSize: '20px', color: 'black' }}>
                      Malicious Flows Detected
                    </p>
                    <Card.Title as="h4">{maliciousFlows}</Card.Title>
                  </div>
                </Col>
              </Row>
            </Card.Body>
          </Card>
        </Col>
      </Row>

      <Row className="mb-4">
        <Col lg="6" sm="6">
          <Card className="card-stats">
            <Card.Body>
              <Row>
                <Col xs="3">
                  <div className="icon-big text-center icon-warning">
                    <i className="nc-icon nc-vector text-warning" />
                  </div>
                </Col>
                <Col xs="9">
                  <div className="numbers">
                    <p className="card-category" style={{ fontSize: '20px', color: 'black' }}>
                      Last Malicious Flow Detected
                    </p>
                    <Card.Title as="h5">
                      {lastMaliciousFlow?.flow_id || "N/A"}
                    </Card.Title>
                    <p style={{ fontSize: "14px" }}>
                      Protocol: {lastMaliciousFlow?.protocol || "N/A"}<br />
                      Risk Score: {lastMaliciousFlow?.risk_score?.toFixed(2) || "N/A"}
                    </p>
                  </div>
                </Col>
              </Row>
            </Card.Body>
          </Card>
        </Col>

        <Col md="6">
          <Card>
            <Card.Header>
              <Card.Title as="h4">Packet Flow Pie Chart</Card.Title>
            </Card.Header>
            <Card.Body>
              <div className="ct-chart ct-perfect-fourth">
                <ChartistGraph
                  data={{
                    labels: [
                      `${((goodFlows / safeTotal) * 100).toFixed(0)}%`,
                      `${((maliciousFlows / safeTotal) * 100).toFixed(0)}%`,
                    ],
                    series: [goodFlows, maliciousFlows],
                  }}
                  type="Pie"
                />
              </div>
              <div className="legend" style={{ paddingLeft: "20px" }}>
                <i className="fas fa-circle text-info" /> Safe Packets
                <span style={{ marginLeft: "20px" }} />
                <i className="fas fa-circle text-danger" /> Malicious Packets
              </div>
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Container>
  );
}

export default Network;
