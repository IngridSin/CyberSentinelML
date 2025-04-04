import React, { useEffect, useState } from "react";
import {
  Card,
  Table,
  Container,
  Row,
  Col,
  Pagination,
  Button,
  Collapse,
} from "react-bootstrap";

function EmailTable() {
  const [emails, setEmails] = useState([]);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [expandedRow, setExpandedRow] = useState(null);

  const fetchEmails = async (pageNum) => {
    try {
      const res = await fetch(`http://localhost:8080/api/emails?page=${pageNum}`);
      const data = await res.json();

      const safeTotal = Number(data.total ?? 0);
      const safePageSize = Number(data.pageSize ?? 10);
      const safeEmails = Array.isArray(data.emails) ? data.emails : [];

      setEmails(safeEmails);
      setTotalPages(Math.max(1, Math.ceil(safeTotal / safePageSize)));
    } catch (err) {
      console.error("Failed to load emails:", err);
    }
  };

  useEffect(() => {
    fetchEmails(page);
  }, [page]);

  const toggleRow = (id) => {
    setExpandedRow(expandedRow === id ? null : id);
  };

  const getTopTriggerWords = (email) => {
    const allWords = [];

    Object.keys(email).forEach((key) => {
      if (key.startsWith("top_5_words_from_")) {
        const words = Object.keys(email[key] || {});
        allWords.push(...words);
      }
    });

    const wordCounts = allWords.reduce((acc, word) => {
      acc[word] = (acc[word] || 0) + 1;
      return acc;
    }, {});

    const sorted = Object.entries(wordCounts).sort((a, b) => b[1] - a[1]);
    return sorted.slice(0, 5).map(([word]) => word).join(", ");
  };

  return (
    <Container fluid>
      <Row>
        <Col md="12">
          <Card className="strpied-tabled-with-hover">
            <Card.Header>
              <Card.Title as="h4">All Emails</Card.Title>
              <p className="card-category">Paginated email records</p>
            </Card.Header>
            <Card.Body className="table-full-width table-responsive px-0">
              <Table className="table-hover table-striped mb-0">
                <thead>
                <tr>
                  <th>ID</th>
                  <th>Subject</th>
                  <th>Sender</th>
                  <th>Prediction</th>
                  <th>Date</th>
                  <th>Body Score</th>
                  <th>Header Score</th>
                  <th>Risk Score</th>
                  <th>Top Words</th>
                  <th>Details</th>
                </tr>
                </thead>
                <tbody>
                {emails.length > 0 ? (
                  emails.map((email) => (
                    <React.Fragment key={email.id}>
                      <tr>
                        <td>{email.id}</td>
                        <td>{email.subject}</td>
                        <td>{email.sender}</td>
                        <td>
                          {email.prediction === 1 ? (
                            <span className="text-danger">Phishing</span>
                          ) : (
                            <span className="text-success">Safe</span>
                          )}
                        </td>
                        <td>{new Date(email.date).toLocaleString()}</td>
                        <td>{(email.winner_probability * 100).toFixed(2)}%</td>
                        <td>{email.header_valid ? "100" : "0"}</td>
                        <td>{email.risk_score}</td>
                        <td>{getTopTriggerWords(email)}</td>
                        <td>
                          <Button
                            variant="info"
                            size="sm"
                            onClick={() => toggleRow(email.id)}
                          >
                            {expandedRow === email.id ? "Hide" : "Show"}
                          </Button>
                        </td>
                      </tr>
                      <tr>
                        <td colSpan="10" className="p-0 border-0">
                          <Collapse in={expandedRow === email.id}>
                            <div className="p-3 bg-light text-dark border-top">
                              <strong>Email Body:</strong>
                              <pre
                                style={{
                                  whiteSpace: "pre-wrap",
                                  wordBreak: "break-word",
                                  marginTop: "10px",
                                }}
                              >
                                  {email.body || "No body available."}
                                </pre>
                            </div>
                          </Collapse>
                        </td>
                      </tr>
                    </React.Fragment>
                  ))
                ) : (
                  <tr>
                    <td colSpan="10" className="text-center">
                      No emails found.
                    </td>
                  </tr>
                )}
                </tbody>
              </Table>

              {totalPages > 1 && (
                <Pagination className="justify-content-center mt-3">
                  {Array.from({ length: totalPages }, (_, i) => (
                    <Pagination.Item
                      key={i + 1}
                      active={i + 1 === page}
                      onClick={() => setPage(i + 1)}
                    >
                      {i + 1}
                    </Pagination.Item>
                  ))}
                </Pagination>
              )}
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Container>
  );
}

export default EmailTable;
