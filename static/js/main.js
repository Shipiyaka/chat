function sendMessage() {
    let form = document.forms.messageForm;
    let message = form.elements.messageInput.value;
    form.elements.messageInput.value = "";

    let messageContainer = document.getElementById("messages");
    let newText = document.createElement("p");
    newText.innerHTML = message;
    messageContainer.appendChild(newText);
}

function initWebSocketConn() {
    let socket = new WebSocket("ws://localhost:12345/ws");
    socket.onopen = function(e) {
        alert("Join to chat");
        socket.send("Ping");
    };
    socket.onmessage = function(event) {
        alert(`${event.data}`);
    };
}

window.onload = initWebSocketConn;