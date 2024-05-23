import { useState, useEffect } from "react";
import { Container, Button, Textarea, Card, Notification } from "@mantine/core";
import { jwtDecode } from "jwt-decode";

interface Message {
  GUID: string;
  Message: string;
  UserID: string;
}

function Dashboard() {
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState<Message[]>([]);
  const [notificationVisible, setNotificationVisible] = useState(false);

  useEffect(() => {
    fetchMessages();
  }, []);

  const handleLogout = () => {
    // Simulate logout
    console.log("Logged out");
  };

  const fetchMessages = () => {
    const idToken = sessionStorage.getItem("id_token");
    const apiUrl = `https://6xoa9t5ole.execute-api.eu-central-1.amazonaws.com/Prod/getMessages/${getUserId()}`;

    if (idToken) {
      fetch(apiUrl, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `${idToken}`,
        },
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.json();
        })
        .then((data) => {
          setMessages(data || []);
        })
        .catch((error) => {
          console.error("There was a problem fetching messages:", error);
        });
    }
  };

  const handleRefresh = () => {
    fetchMessages();
  };

  const handlePostMessage = () => {
    const idToken = sessionStorage.getItem("id_token");
    const apiUrl =
      "https://6xoa9t5ole.execute-api.eu-central-1.amazonaws.com/Prod/postMessage";
    if (idToken) {
      const decodedToken = jwtDecode(idToken);
      console.log(decodedToken.sub);
      const requestBody = {
        detail: {
          Message: message,
          UserID: decodedToken.sub,
        },
      };

      fetch(apiUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `${idToken}`,
        },
        body: JSON.stringify(requestBody),
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.json();
        })
        .then((data) => {
          console.log("Message posted successfully:", data);
          setMessage(""); // Clear message input after posting
          // After posting the message, fetch messages again to update the list
          fetchMessages();
          // Show notification
          setNotificationVisible(true);
          // Hide notification after 3 seconds
          setTimeout(() => {
            setNotificationVisible(false);
          }, 3000);
        })
        .catch((error) => {
          console.error("There was a problem posting the message:", error);
        });
    }
  };

  const getUserId = () => {
    const idToken = sessionStorage.getItem("id_token");
    if (idToken) {
      const decodedToken = jwtDecode(idToken);
      return decodedToken.sub;
    }
    return "";
  };

  return (
    <Container size="sm">
      <Button
        color="red"
        onClick={handleLogout}
        style={{ marginBottom: 16 }}
      >
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
    </Container>
  );
}

export default Dashboard;

