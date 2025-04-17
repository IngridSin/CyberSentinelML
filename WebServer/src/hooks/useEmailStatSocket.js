import { useEffect } from "react";
import useEmailDashboardStore from "../store/emailDashboardStore";
import useNetworkStore from "../store/networkStore";

const useStatsSocket = () => {
  useEffect(() => {
    const fetchInitialStats = async () => {
      try {
        const [emailRes, networkRes] = await Promise.all([
          fetch("http://localhost:8080/api/email-stats"),
          fetch("http://localhost:8080/api/network-stats"),
        ]);
        const emailData = await emailRes.json();
        const networkData = await networkRes.json();

        useEmailDashboardStore.getState().setDashboardStats(emailData);
        useNetworkStore.getState().setNetworkStats(networkData);
      } catch (err) {
        console.error("Failed to fetch initial stats:", err);
      }
    };

    fetchInitialStats();

    let ws;
    let reconnectTimeout;

    const connectWebSocket = () => {
      ws = new WebSocket("ws://localhost:8080/ws");

      ws.onopen = () => console.log("🔗 Unified WebSocket connected");

      ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        console.log("🛰️ WebSocket message:", message);

        switch (message.type) {
          case "EMAIL_DASHBOARD_STATS":
            useEmailDashboardStore.getState().setDashboardStats(message.payload);
            break;
          case "NETWORK_DASHBOARD_STATS":
            useNetworkStore.getState().setNetworkStats(message.payload);
            break;
          default:
            console.warn("Unknown WebSocket message type:", message.type);
        }
      };

      ws.onclose = () => {
        console.warn("🔌 WebSocket disconnected. Reconnecting in 10s...");
        reconnectTimeout = setTimeout(connectWebSocket, 10000);
      };

      ws.onerror = (err) => {
        console.error("WebSocket error:", err);
        ws.close();
      };
    };

    connectWebSocket();

    return () => {
      if (ws) ws.close();
      if (reconnectTimeout) clearTimeout(reconnectTimeout);
    };
  }, []);
};

export default useStatsSocket;
