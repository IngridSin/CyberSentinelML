import React from "react";
import ChartistGraph from "react-chartist";
import { Card, Col, Container, Form, Row } from "react-bootstrap";
import useEmailDashboardStore from "../store/emailDashboardStore";

const getPieChartData = (total, phishing) => {
  const safe = total - phishing;
  const safePercent = total > 0 ? Math.round((safe / total) * 100) : 0;
  const phishingPercent = 100 - safePercent;
  return {
    labels: [`${safePercent}%`, `${phishingPercent}%`],
    series: [safePercent, phishingPercent],
  };
};

const formatDateTime = (timestamp) =>
  timestamp ? new Date(timestamp).toLocaleString() : "N/A";

function Phishing() {
  const {
    totalEmails,
    phishingEmails,
    lastPhishingTime,
    lastPhishingEmail,
  } = useEmailDashboardStore();

  const chartData = getPieChartData(totalEmails, phishingEmails);

  return (
    <Container fluid>
      <Row>
        {/* Total Emails */}
        <Col lg="4" sm="6">
          <Card className="card-stats">
            <Card.Body>
              <Row>
                <Col xs="3">
                  <div className="icon-big text-center icon-warning">
                    <i className="nc-icon nc-chart text-warning" />
                  </div>
                </Col>
                <Col xs="8">
                  <div className="numbers">
                    <p className="card-category text-dark fw-bold">
                      Total Emails Received
                    </p>
                    <Card.Title as="h4">{totalEmails}</Card.Title>
                  </div>
                </Col>
              </Row>
            </Card.Body>
          </Card>
        </Col>

        {/* Suspected Phishing Emails */}
        <Col lg="4" sm="6">
          <Card className="card-stats">
            <Card.Body>
              <Row>
                <Col xs="2">
                  <div className="icon-big text-center icon-success">
                    <i className="nc-icon nc-light-3 text-success" />
                  </div>
                </Col>
                <Col xs="10">
                  <div className="numbers">
                    <p className="card-category text-dark fw-bold">
                      Suspected Phishing Emails
                    </p>
                    <Card.Title as="h4">{phishingEmails}</Card.Title>
                  </div>
                </Col>
              </Row>
            </Card.Body>
          </Card>
        </Col>

        {/* Last Detection Time */}
        <Col lg="4" sm="6">
          <Card className="card-stats">
            <Card.Body>
              <Row>
                <Col xs="2">
                  <div className="icon-big text-center icon-danger">
                    <i className="nc-icon nc-vector text-danger" />
                  </div>
                </Col>
                <Col xs="10">
                  <div className="numbers">
                    <p className="card-category text-dark fw-bold">
                      Last Detected Phishing Attempt
                    </p>
                    <Card.Title as="h4">{formatDateTime(lastPhishingTime)}</Card.Title>
                  </div>
                </Col>
              </Row>
            </Card.Body>
          </Card>
        </Col>
      </Row>

      <Row>
        {/* Pie Chart */}
        <Col md="4">
          <Card>
            <Card.Header>
              <Card.Title as="h4">Email Classification Summary</Card.Title>
            </Card.Header>
            <Card.Body>
              <div className="ct-chart ct-perfect-fourth" id="chartPreferences">
                <ChartistGraph data={chartData} type="Pie" />
              </div>
              <div className="legend">
                <i className="fas fa-circle text-info ms-3" />
                Safe Emails
                <i className="fas fa-circle text-danger ms-4" />
                Suspected Phishing
              </div>
            </Card.Body>
          </Card>
        </Col>

        {/* Last Phishing Email Body */}
        <Col md="8">
          <Card>
            <Card.Header>
              <Card.Title as="h4">Last Detected Phishing Email</Card.Title>
            </Card.Header>
            <Card.Body>
              {lastPhishingEmail?.body ? (
                <>
                  <p><strong>Subject:</strong> {lastPhishingEmail.subject}</p>
                  <p><strong>Sender:</strong> {lastPhishingEmail.sender}</p>
                  <p><strong>Recipient:</strong> {lastPhishingEmail.recipient}</p>
                  <p><strong>Return Path:</strong> {lastPhishingEmail.return_path}</p>
                  <p><strong>DKIM:</strong> {lastPhishingEmail.dkim}</p>
                  <p><strong>SPF:</strong> {lastPhishingEmail.spf}</p>
                  <Form.Group>
                    <Form.Label className="text-dark fw-bold" style={{ fontSize: "18px" }}>
                      Email Body
                    </Form.Label>
                    <Form.Control
                      as="textarea"
                      rows="10"
                      readOnly
                      value={lastPhishingEmail.body}
                    />
                  </Form.Group>
                </>
              ) : (
                <p>No phishing email detected yet.</p>
              )}
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Container>
  );
}

export default Phishing;
