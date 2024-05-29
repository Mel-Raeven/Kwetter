// src/postMessage.ts
import { jwtDecode } from "jwt-decode";

interface DecodedToken {
  sub: string;
}

interface PostMessageResponse {
  success: boolean;
  data?: any;
  error?: string;
}

export async function postMessage(
  idToken: string,
  message: string
): Promise<PostMessageResponse> {
  const apiUrl =
    "https://0a43x0s4q4.execute-api.eu-central-1.amazonaws.com/Prod/postMessage";

  const decodedToken: DecodedToken = jwtDecode(idToken);
  const requestBody = {
    detail: {
      Message: message,
      UserID: decodedToken.sub,
    },
  };

  try {
    const response = await fetch(apiUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `${idToken}`,
      },
      body: JSON.stringify(requestBody),
    });

    if (!response.ok) {
      throw new Error("Network response was not ok");
    }

    const data = await response.json();
    return { success: true, data };
  } catch (error) {
    console.error("There was a problem posting the message:", error);
    return { success: false };
  }
}

