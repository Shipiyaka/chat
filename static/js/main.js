const socket = new WebSocket("ws://localhost:12345/ws");
const messageInput = document.getElementById("messageInput");
messageInput.addEventListener("keyup", ({ key }) => {
    if (key === "Enter") {
        sendText();
    }
});

let username;

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
        from_user: username,
        date: date
    };

    if (type === "text") {
        message.text = content;
    } else  if (type === "img") {
        message.img = content;
    }

    return message
}

function updateMessageContainer(message) {
    let messageContainer = document.getElementById("messages");
    let kindOfMessage = (message.from_user === username ? "outgoingMessage" : "incomingMessage");
    let photo = (
        message.from_user === username ?
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
            <div class="${kindOfMessage}Login">${message.from_user}</div>
            ${essenceOfMessage}
            <div class="messageTime">${message.date}</div>
        </div>
    </div>
</div>`

    document.body.scrollTop = document.body.scrollHeight - document.body.clientHeight;
}

function createWebSocketEvents() {
    socket.onopen = function(e) {
        alert("Join to chat");
    };
    socket.onmessage = function(e) {
        let message = JSON.parse(e.data);

        if ("username" in message) {
            username = message.username;
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

      fileReader.onload = function(fileLoadedEvent) {
        let srcData = fileLoadedEvent.target.result;
        sendImg(srcData);
      }
      fileReader.readAsDataURL(fileToLoad);
    }
  }

window.onload = createWebSocketEvents;
/*
$(function(){
    $('.minimized').click(function(event) {
        var i_path = $(this).attr('src');
        $('body').append('<div id="overlay"></div><div id="magnify"><img src="'+i_path+'"><div id="close-popup"><i></i></div></div>');
        $('#magnify').css({
            left: ($(document).width() - $('#magnify').outerWidth())/2,   
            top: ($(window).height() - $('#magnify').outerHeight())/2
        });
        $('#overlay, #magnify').fadeIn('fast');
    });
    
    $('body').on('click', '#close-popup, #overlay', function(event) {
      event.preventDefault();
  
      $('#overlay, #magnify').fadeOut('fast', function() {
        $('#close-popup, #magnify, #overlay').remove();
      });
    });
  });
  */