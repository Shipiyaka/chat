const socket = new WebSocket("ws://localhost:12345/ws");
const messageInput = document.getElementById("messageInput");

messageInput.addEventListener("keyup", ({ key }) => {
    if (key === "Enter") {
        sendText();
    }
});

let user;

function sendText() {
    let text = messageInput.value;
    if (text === "") {
        return;
    }
    messageInput.value = "";

    message = prepareMessageObject("text", text)
    
    updateMessageContainer(message);

    socket.send(JSON.stringify(message))
}

function sendImg(src) {
    message = prepareMessageObject("img", src)
    
    updateMessageContainer(message);

    socket.send(JSON.stringify(message))
}

function prepareMessageObject(type, content) {
    let date = new Date().toLocaleTimeString(
        'en-US', {
        hour12: false,
        hour: "numeric",
        minute: "numeric"
    }
    );

    let message = {
        from_user: user.username,
        username_color: user.username_color,
        date: date,
    };

    if (type === "text") {
        message.text = content;
    } else if (type === "img") {
        message.img = content;
    }

    return message
}

function updateMessageContainer(message) {
    let messageContainer = document.getElementById("messages");
    let kindOfMessage;
    if (user === undefined) {
        kindOfMessage = "incomingMessage"
    } else {
        kindOfMessage = (message.from_user === user.username ? "outgoingMessage" : "incomingMessage");
    }

    let photo = (
        kindOfMessage === "outgoingMessage" ?
            "../static/images/profile_pic/photo_2021-03-12 00.39.30.jpeg" :
            "../static/images/profile_pic/photo_2021-03-13 14.44.40.jpeg"
    )
    
    essenceOfMessage = ("img" in message ? `<img class="messagePic" src=${message.img}>` : message.text)

    messageContainer.innerHTML += `<div class="${kindOfMessage}">
    <!-- ${kindOfMessage} -->
    <div class="${kindOfMessage}Item">
        <img class="${kindOfMessage}ProfilePic"
            src="${photo}">
        <div class="${kindOfMessage}Content">
            <div style="color: ${message.username_color}" class="${kindOfMessage}Login">${message.from_user}</div>
            ${essenceOfMessage}
            <div class="messageTime">${message.date}</div>
        </div>
    </div>
</div>`

    document.body.scrollTop = document.body.scrollHeight - document.body.clientHeight;
}

function createWebSocketEvents() {
    socket.onopen = function (e) {
        alert("Join to chat");
    };
    socket.onmessage = function (e) {
        let message = JSON.parse(e.data);

        if ("username" in message) {
            user = { username: message.username, username_color: message.username_color }
            console.log(user)
            return;
        }
        updateMessageContainer(message)
    };
}

function uploadFile() {
    let filesSelected = document.getElementById("getFile").files;
    if (filesSelected.length > 0) {
        let fileToLoad = filesSelected[0];

        let fileReader = new FileReader();

        fileReader.onload = function (fileLoadedEvent) {
            let srcData = fileLoadedEvent.target.result;
            sendImg(srcData);
        }
        fileReader.readAsDataURL(fileToLoad);
    }
}

window.onload = createWebSocketEvents;