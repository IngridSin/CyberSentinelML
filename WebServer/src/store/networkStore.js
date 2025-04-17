import { create } from 'zustand';
import axios from 'axios';

const useNetworkStore = create((set) => ({
  totalFlows: 0,
  maliciousFlows: 0,
  lastMaliciousTime: null,
  lastMaliciousFlow: {
    flow_id: '',
    src_ip: '',
    dst_ip: '',
    protocol: '',
    risk_score: 0,
    timestamp: '',
  },

  // Paginated Flows
  packets: [],
  totalPackets: 0,
  currentPage: 1,
  pageSize: 10,
  loading: false,

  setNetworkStats: (data) =>
    set(() => ({
      totalFlows: data.total_flows || 0,
      maliciousFlows: data.malicious_flows || 0,
      lastMaliciousTime: data.last_malicious_time || null,
      lastMaliciousFlow: data.last_malicious_flow || {},
    })),


  fetchNetworkPackets: async (page = 1, pageSize = 10) => {
    set({ loading: true });
    try {
      const res = await axios.get(`http://localhost:8080/api/network-packets?page=${page}&pageSize=${pageSize}`);
      const { packets, total } = res.data;

      set({
        packets,
        totalPackets: total,
        currentPage: page,
        pageSize,
        loading: false,
      });
    } catch (err) {
      console.error("Failed to fetch network packets:", err);
      set({ loading: false });
    }
  },
}));

export default useNetworkStore;
