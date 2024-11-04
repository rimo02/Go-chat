var ws;
let currentUser = null;

async function login() {
  const username = document.getElementById("username").value;
  const password = document.getElementById("password").value;

  // Send login request to the backend
  try {
    const response = await fetch("/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, password }),
    });
    const result = await response.json();

    if (result.success) {
      document.getElementById("login").style.display = "none";
      document.getElementById("chatApp").style.display = "block";

      ws = new WebSocket("ws://localhost:7777/ws"); // Change to your WebSocket URL as needed

      ws.onopen = function () {
        console.log("Connected to WebSocket server.");
        currentUser = username; // Set currentUser to the logged-in username
        ws.send(JSON.stringify({ type: "login", username }));
      };

      ws.onmessage = function (event) {
        const data = JSON.parse(event.data);
        displayMessage(data);
      };

      ws.onclose = function () {
        console.log("Disconnected from WebSocket server.");
      };
    } else {
      alert(result.message); // Show error message from the server
    }
  } catch (error) {
    console.error("Error logging in:", error);
    alert("Login failed. Please try again.");
  }
}

// Display the message in the chat
function displayMessage(data) {
  const chat = document.getElementById("chat");

  const messageContainer = document.createElement("div");
  messageContainer.classList.add("message-container");

  const usernameElement = document.createElement("p");
  usernameElement.classList.add("username");

  if (data.username === currentUser) {
    usernameElement.textContent = "You";
    messageContainer.classList.add("my-message");
  } else {
    usernameElement.textContent = data.username;
    messageContainer.classList.add("other-message");
  }

  messageContainer.appendChild(usernameElement);

  // Add the message text
  const message = document.createElement("p");
  message.classList.add("message-text");
  message.textContent = data.text;
  messageContainer.appendChild(message);

  // Append the message container to the chat
  chat.appendChild(messageContainer);

  // Auto-scroll to the bottom when a new message is added
  chat.scrollTop = chat.scrollHeight;
}

// Sends the message to backend server
function sendMessage() {
  const input = document.getElementById("message");

  const messageData = JSON.stringify({
    username: currentUser,
    text: input.value,
  });
  ws.send(messageData);

  input.value = "";
}
