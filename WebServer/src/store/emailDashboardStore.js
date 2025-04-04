import { create } from 'zustand';

const useEmailDashboardStore = create((set) => ({
  totalEmails: 0,
  phishingEmails: 0,
  lastPhishingTime: null,
  lastPhishingEmail: {
    subject: '',
    body: '',
    message_id: '',
    timestamp: '',
  },

  setDashboardStats: (data) =>
    set(() => ({
      totalEmails: data.total_emails || 0,
      phishingEmails: data.phishing_emails || 0,
      lastPhishingTime: data.last_phishing_time || null,
      lastPhishingEmail: {
        subject: data.last_phishing_email?.subject || "",
        body: data.last_phishing_email?.body || "",
        sender: data.last_phishing_email?.sender || "",
        recipient: data.last_phishing_email?.recipient || "",
        return_path: data.last_phishing_email?.return_path || "",
        dkim: data.last_phishing_email?.dkim || "",
        spf: data.last_phishing_email?.spf || "",
        timestamp: data.last_phishing_email?.timestamp || "",
      },
    })),
}));

export default useEmailDashboardStore;
