import React, { useEffect, useState } from "react";
import {
  Card,
  Table,
  Container,
  Row,
  Col,
  Pagination,
  Badge,
  Collapse,
  Button,
  Form,
  InputGroup,
} from "react-bootstrap";

function NetworkPacketTable() {
  const [packets, setPackets] = useState([]);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [expandedRow, setExpandedRow] = useState(null);
  const [showOnlyMalicious, setShowOnlyMalicious] = useState(false);
  const [inputPage, setInputPage] = useState("");

  const fetchPackets = async (pageNum) => {
    try {
      const endpoint = showOnlyMalicious
        ? `http://localhost:8080/api/malicious-packets?page=${pageNum}`
        : `http://localhost:8080/api/network-packets?page=${pageNum}`;

      const res = await fetch(endpoint);
      const data = await res.json();

      const safeTotal = Number(data.total ?? 0);
      const safePageSize = Number(data.pageSize ?? 15);
      const safePackets = Array.isArray(data.packets) ? data.packets : [];

      setPackets(safePackets);
      setTotalPages(Math.max(1, Math.ceil(safeTotal / safePageSize)));
    } catch (err) {
      console.error("Failed to load packets:", err);
    }
  };

  // Refetch when page or filter changes
  useEffect(() => {
    fetchPackets(page);
  }, [page, showOnlyMalicious]);

  const toggleRow = (flowID) => {
    setExpandedRow(expandedRow === flowID ? null : flowID);
  };

  const goToPage = () => {
    const parsed = parseInt(inputPage);
    if (!isNaN(parsed) && parsed >= 1 && parsed <= totalPages) {
      setPage(parsed);
    }
    setInputPage("");
  };

  const toggleMaliciousFilter = () => {
    setPage(1); // reset to first page
    setShowOnlyMalicious((prev) => !prev);
  };

  return (
    <Container fluid>
      <Row>
        <Col md="12">
          <Card className="strpied-tabled-with-hover">
            <Card.Header>
              <div className="d-flex justify-content-between align-items-center">
                <div>
                  <Card.Title as="h4">Captured Network Flows</Card.Title>
                  <p className="card-category">Paginated flow records with risk info</p>
                </div>
                <Button
                  variant={showOnlyMalicious ? "danger" : "secondary"}
                  onClick={toggleMaliciousFilter}
                >
                  {showOnlyMalicious ? "Show All" : "Show Only Malicious"}
                </Button>
              </div>
            </Card.Header>

            <Card.Body className="table-full-width table-responsive px-0">
              <Table className="table-hover table-striped mb-0">
                <thead>
                <tr>
                  <th>Flow ID</th>
                  <th>Source IP</th>
                  <th>Destination IP</th>
                  <th>Protocol</th>
                  <th>Timestamp</th>
                  <th>Risk Score</th>
                  <th>Prediction</th>
                  <th>Details</th>
                </tr>
                </thead>
                <tbody>
                {packets.length > 0 ? (
                  packets.map((pkt) => (
                    <React.Fragment key={pkt.flow_id}>
                      <tr>
                        <td>{pkt.flow_id}</td>
                        <td>{pkt.src_ip}</td>
                        <td>{pkt.dst_ip}</td>
                        <td>{pkt.protocol}</td>
                        <td>{new Date(pkt.timestamp).toLocaleString()}</td>
                        <td>{pkt.risk_score}</td>
                        <td>
                          {pkt.prediction === 1 ? (
                            <Badge bg="danger">Malicious</Badge>
                          ) : (
                            <Badge bg="success">Benign</Badge>
                          )}
                        </td>
                        <td>
                          <Button
                            variant="info"
                            size="sm"
                            onClick={() => toggleRow(pkt.flow_id)}
                          >
                            {expandedRow === pkt.flow_id ? "Hide" : "Show"}
                          </Button>
                        </td>
                      </tr>
                      <tr>
                        <td colSpan="8" className="p-0 border-0">
                          <Collapse in={expandedRow === pkt.flow_id}>
                            <div className="p-3 bg-light text-dark border-top">
                              <strong>More Info:</strong>
                              <pre
                                style={{
                                  whiteSpace: "pre-wrap",
                                  wordBreak: "break-word",
                                  marginTop: "10px",
                                }}
                              >
                                  {JSON.stringify(pkt, null, 2)}
                                </pre>
                            </div>
                          </Collapse>
                        </td>
                      </tr>
                    </React.Fragment>
                  ))
                ) : (
                  <tr>
                    <td colSpan="8" className="text-center">
                      No packet data available.
                    </td>
                  </tr>
                )}
                </tbody>
              </Table>

              {totalPages > 1 && (
                <div className="d-flex justify-content-center align-items-center mt-3">
                  <Pagination>
                    <Pagination.Prev
                      onClick={() => setPage((p) => Math.max(1, p - 1))}
                      disabled={page === 1}
                    />
                    <Pagination.Item active disabled>
                      Page {page} / {totalPages}
                    </Pagination.Item>
                    <Pagination.Next
                      onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                      disabled={page === totalPages}
                    />
                  </Pagination>
                  <InputGroup className="ms-3" style={{ maxWidth: 120 }}>
                    <Form.Control
                      type="number"
                      placeholder="Page #"
                      value={inputPage}
                      onChange={(e) => setInputPage(e.target.value)}
                      onKeyDown={(e) => e.key === "Enter" && goToPage()}
                    />
                    <Button variant="primary" onClick={goToPage}>
                      Go
                    </Button>
                  </InputGroup>
                </div>
              )}
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Container>
  );
}

export default NetworkPacketTable;
