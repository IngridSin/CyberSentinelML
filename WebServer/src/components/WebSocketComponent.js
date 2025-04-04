import React, { useState, useEffect } from "react";

const WebSocketComponent = () => {
  const [recordCount, setRecordCount] = useState(0);
  const [groupedData, setGroupedData] = useState([]);
  const [lastData, setLastData] = useState(null);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080");

    socket.onmessage = (event) => {
      const receivedData = JSON.parse(event.data);

      // Update the state with the received data
      setRecordCount(receivedData.totalRecords);
      setGroupedData(receivedData.groupedData);
      setLastData(receivedData.lastData); // Set last data (latest record)
    };

    return () => {
      socket.close(); // Clean up the socket when the component unmounts
    };
  }, []);

  return (
    <div className="p-4">
      <h1 className="text-xl font-bold">id</h1>
      {recordCount > 0 ? (
        <div className="mt-2">
          <p><strong>Total Records:</strong> {recordCount}</p>

          <h2 className="mt-4">Grouped Data</h2>
          {groupedData.length > 0 ? (
            groupedData.map((group, index) => (
              <div key={index} className="mt-2">
                <p><strong>Category:</strong> {group.category}</p>
                <p><strong>Record Count:</strong> {group.record_count}</p>
              </div>
            ))
          ) : (
            <p>No grouped data available</p>
          )}

          <h2 className="mt-4">Last Data</h2>
          {lastData ? (
            <div className="mt-2">
              {/* Display each field of the last data */}
              {Object.entries(lastData).map(([key, value]) => (
                <p key={key}><strong>{key}:</strong> {value}</p>
              ))}
            </div>
          ) : (
            <p>No last data available</p>
          )}
        </div>
      ) : (
        <p>Loading data...</p>
      )}
    </div>
  );
};

export default WebSocketComponent;
