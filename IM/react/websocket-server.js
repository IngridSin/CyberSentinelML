const WebSocket = require("ws");
const mysql = require("mysql2/promise");

const server = new WebSocket.Server({ port: 8080 });

const db = mysql.createPool({
  host: "localhost",
  user: "user",
  password: "user",
  database: "mydb",
});

// Fetch total record count
async function getRecordCount() {
  try {
    const [rows] = await db.query("SELECT COUNT(*) AS total_records FROM test");
    return rows[0].total_records;
  } catch (error) {
    console.error("Database Error:", error);
    return 0;
  }
}

// Fetch grouped data by category
async function getGroupedData() {
  try {
    const [rows] = await db.query("SELECT type, COUNT(*) AS type_count FROM test GROUP BY type");
    return rows;
  } catch (error) {
    console.error("Database Error:", error);
    return [];
  }
}

// Fetch the most recent data (latest row)
async function getLastData() {
  try {
    const [rows] = await db.query("SELECT * FROM your_table ORDER BY id DESC LIMIT 1");
    return rows.length ? rows[0] : null;
  } catch (error) {
    console.error("Database Error:", error);
    return null;
  }
}

server.on("connection", async (ws) => {
  console.log("Client connected");

  // Fetch record count, grouped data, and last data
  const recordCount = await getRecordCount();
  const groupedData = await getGroupedData();
  const lastData = await getLastData();

  const dataToSend = {
    totalRecords: recordCount,
    groupedData: groupedData,
    lastData: lastData, // Add last data to the response
  };

  // Send aggregated data to the client
  ws.send(JSON.stringify(dataToSend));

  ws.on("close", () => console.log("Client disconnected"));
});
