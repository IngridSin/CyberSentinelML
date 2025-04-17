/*!

=========================================================
* Light Bootstrap Dashboard React - v2.0.1
=========================================================

* Product Page: https://www.creative-tim.com/product/light-bootstrap-dashboard-react
* Copyright 2022 Creative Tim (https://www.creative-tim.com)
* Licensed under MIT (https://github.com/creativetimofficial/light-bootstrap-dashboard-react/blob/master/LICENSE.md)

* Coded by Creative Tim

=========================================================

* The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

*/
import Phishing from "./views/Phishing.js";
import Network from "./views/Network.js";
import EmailTable from "./views/EmailTable";
import NetworkPacketTable from "./views/NetworkPacketTable";

const dashboardRoutes = [


  {
    path: "/Phishing",
    name: "Emails",
    icon: "nc-icon nc-chart-pie-35",
    component: Phishing,
    layout: "/admin"
  },

    {
    path: "/Network",
    name: "Network",
    icon: "nc-icon nc-notes",
    component: Network,
    layout: "/admin"
  },
  {
    path: "/emails",
    name: "All Emails",
    icon: "nc-icon nc-notes",
    component: EmailTable,
    layout: "/admin"
  },
  {
    path: "/all-packets",
    name: "All Packets",
    icon: "nc-icon nc-notes",
    component: NetworkPacketTable,
    layout: "/admin"
  },


];

export default dashboardRoutes;
