const socket = new WebSocket("ws://localhost:12345/ws");
const messageInput = document.getElementById("messageInput");
messageInput.addEventListener("keyup", ({key}) => {
    if (key === "Enter") {
        sendMessage();
    }
});

let username;

function sendMessage() {
    let text = messageInput.value;
    if (text === "")  {
        return;
    }
    messageInput.value = "";

    let date = new Date().toLocaleTimeString(
        'en-US', {
        hour12: false,
        hour: "numeric",
        minute: "numeric"
    });

    let message = {
        text: text,
        from_user: username,
        date: date
    }

    updateMessageContainer(message);

    socket.send(JSON.stringify(message))
}

function updateMessageContainer(message) {
    let messageContainer = document.getElementById("messages");
    let kindOfMessage = (message.from_user == username ? "outgoingMessage" : "incomingMessage");
    let photo = (
        message.from_user == username 
        ? "../static/images/profile_pic/photo_2021-03-12 00.39.30.jpeg" 
        : "../static/images/profile_pic/photo_2021-03-13 14.44.40.jpeg"
    );

    messageContainer.innerHTML += `<div class="${kindOfMessage}">
    <!-- ${kindOfMessage} -->
    <div class="${kindOfMessage}Item">
        <img class="${kindOfMessage}ProfilePic"
            src="${photo}">
        <div class="${kindOfMessage}Text">
            <div class="${kindOfMessage}Login">${message.from_user}</div>
            ${message.text}
        </div>
        <div class="messageTime">${message.date}</div>
    </div>
</div>`

    messageContainer.scrollTop = messageContainer.scrollHeight - messageContainer.clientHeight;
}

function createWebSocketEvents() {
    socket.onopen = function (e) {
        alert("Join to chat");
    };
    socket.onmessage = function (e) {
        let message = JSON.parse(e.data);

        if ("username" in message) {
            username = message.username;
            return;
        }
        updateMessageContainer(message)
    };
}

window.onload = createWebSocketEvents;
