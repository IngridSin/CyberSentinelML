// components/PhishingAlertModal.js
import React from "react";
import { Modal, Button } from "react-bootstrap";
import useEmailDashboardStore from "../store/emailDashboardStore";

const PhishingAlertModal = () => {
  const { modalVisible, hideModal, dashboardStats } = useEmailDashboardStore();

  return (
    <Modal show={modalVisible} onHide={hideModal} centered>
      <Modal.Header closeButton className="bg-danger text-white">
        <Modal.Title>⚠️ Phishing Alert</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <p>A new phishing email was just detected.</p>
        <p><strong>Subject:</strong> {dashboardStats?.last_phishing_email?.subject}</p>
        <p><strong>From:</strong> {dashboardStats?.last_phishing_email?.sender}</p>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={hideModal}>
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

export default PhishingAlertModal;
