let socket = new WebSocket("ws://localhost:12345/ws");
let username;

function sendMessage() {
    let form = document.forms.messageForm;
    let text = form.elements.messageInput.value;
    form.elements.messageInput.value = "";

    let date = new Date().toLocaleTimeString(
        'en-US', {
            hour12: false,
        hour: "numeric",
        minute: "numeric"
    });
    
    updateMessageContainer(`${date}: ${"You"}</br>${text}`)

    let message = {
        text: text,
        from_user: username,
        date: date
    }
    socket.send(JSON.stringify(message))
}

function updateMessageContainer(message) {
    let messageContainer = document.getElementById("messages");
    let newMessage = document.createElement("p");
    newMessage.innerHTML = message;
    messageContainer.appendChild(newMessage);
}

function initWebSocketConn() {
    socket.onopen = function (e) {
        alert("Join to chat");
    };
    socket.onmessage = function (e) {
        let msg = JSON.parse(e.data);
        console.log(msg);
        if ("username" in msg) {
            username = msg.username;
            return;
        }
        updateMessageContainer(`${msg.date}: ${msg.from_user}</br>${msg.text}`)
    };
}

window.onload = initWebSocketConn;