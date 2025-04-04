import { useEffect } from "react";
import useEmailDashboardStore from "../store/emailDashboardStore";

const useEmailStatsSocket = () => {
  useEffect(() => {
    const fetchInitialStats = async () => {
      try {
        const res = await fetch("http://localhost:8080/api/email-stats");
        const data = await res.json();
        useEmailDashboardStore.getState().setDashboardStats(data);
      } catch (err) {
        console.error("Failed to fetch initial email stats:", err);
      }
    };

    fetchInitialStats();

    let ws;
    let reconnectTimeout;

    const connectWebSocket = () => {
      ws = new WebSocket("ws://localhost:8080/ws");

      ws.onopen = () => console.log("WebSocket connected");

      ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        console.log("WS DATA:", message);

        if (message.type === "EMAIL_DASHBOARD_STATS") {
          useEmailDashboardStore.getState().setDashboardStats(message.payload); // Use .payload
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

export default useEmailStatsSocket;
