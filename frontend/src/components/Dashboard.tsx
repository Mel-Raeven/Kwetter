import { useState, useEffect } from "react";
import { Container, Button, Textarea, Card, Notification } from "@mantine/core";
import { postMessage } from "../functions/postMessage.ts"; // Import the postMessage function
import { jwtDecode } from "jwt-decode";

interface Message {
  GUID: string;
  Message: string;
  UserID: string;
  ts: string;
  lastEvaluatedKey?: string;
}

function Dashboard() {
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState<Message[]>([]);
  const [notificationVisible, setNotificationVisible] = useState(false);
  const [lastEvaluatedKey, setLastEvaluatedKey] = useState<string | null>(null);

  useEffect(() => {
    fetchMessages();
  }, []);

  const handleLogout = () => {
    // Simulate logout
    console.log("Logged out");
  };

  const fetchMessages = (lastEvaluatedKey: string | null = null) => {
    const idToken = sessionStorage.getItem("id_token");
    let apiUrl = `https://0a43x0s4q4.execute-api.eu-central-1.amazonaws.com/Prod/getMessages/${getUserId(idToken)}`;

    if (idToken) {
      const headers: Record<string, string> = {
        "Content-Type": "application/json",
        Authorization: idToken,
      };

      // Include the last evaluated key in the request headers if available
      if (lastEvaluatedKey) {
        headers["Last-Evaluated-Key"] = lastEvaluatedKey;
      }

      fetch(apiUrl, {
        method: "GET",
        headers: headers,
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          // Retrieve Last-Evaluated-Key from response headers
          const lastEvaluatedKeyHeader = response.headers.get("Last-Evaluated-Key");
          return response.json().then((data) => ({ data, lastEvaluatedKeyHeader }));
        })
        .then(({ data, lastEvaluatedKeyHeader }) => {
          // Filter out duplicate messages by comparing GUIDs
          const uniqueMessages = data.filter((newMessage: Message) => {
            return !messages.some(
              (existingMessage) => existingMessage.GUID === newMessage.GUID
            );
          });
          setMessages((prevMessages) => [...prevMessages, ...uniqueMessages]);
          // Set the last evaluated key from the response headers
          if (lastEvaluatedKeyHeader) {
            setLastEvaluatedKey(lastEvaluatedKeyHeader);
          } else {
            setLastEvaluatedKey(null);
          }
        })
        .catch((error) => {
          console.error("There was a problem fetching messages:", error);
        });
    }
  };

  const handleRefresh = () => {
    fetchMessages();
  };

  const handlePostMessage = async () => {
    const idToken = sessionStorage.getItem("id_token");
    if (idToken) {
      const response = await postMessage(idToken, message);
      if (response.success) {
        setMessage(""); // Clear message input after posting
        fetchMessages(); // Fetch messages again to update the list
        // Show notification
        setNotificationVisible(true);
        // Hide notification after 3 seconds
        setTimeout(() => {
          setNotificationVisible(false);
        }, 3000);
      } else {
        console.error("Failed to post message:", response.error);
      }
    }
  };

  const getUserId = (idToken: string | null) => {
    if (idToken) {
      const decodedToken = jwtDecode(idToken);
      return decodedToken.sub;
    }
    return "";
  };

  return (
    <Container size="sm">
      <Button color="red" onClick={handleLogout} style={{ marginBottom: 16 }}>
        Logout
      </Button>
      <Textarea
        placeholder="Enter your message..."
        value={message}
        onChange={(event) => setMessage(event.target.value)}
      />
      <Button
        color="blue"
        variant="outline"
        onClick={handlePostMessage}
        style={{ marginTop: 16, marginRight: 8 }}
      >
        Post Message
      </Button>
      <Button
        color="gray"
        variant="outline"
        onClick={handleRefresh}
        style={{ marginTop: 16 }}
      >
        Refresh
      </Button>
      <div style={{ marginTop: 100 }}>
        {messages.map((msg, index) => (
          <Card key={index} shadow="sm" padding="md" style={{ marginBottom: 12 }}>
            <div>{msg.Message}</div>
            <div>Posted by: {msg.UserID}</div>
          </Card>
        ))}
      </div>
      {notificationVisible && (
        <Notification
          title="Post Received"
          color="teal"
          onClose={() => setNotificationVisible(false)}
          style={{ position: "fixed", bottom: 20, right: 20, zIndex: 9999 }}
        >
          Your message has been received successfully and is being processed!
        </Notification>
      )}
      {lastEvaluatedKey && (
        <Button
          color="cyan"
          onClick={() => fetchMessages(lastEvaluatedKey)}
          style={{ marginTop: 16 }}
        >
          Load More
        </Button>
      )}
    </Container>
  );
}

export default Dashboard;

