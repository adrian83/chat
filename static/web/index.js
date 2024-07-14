
const senderId = document.getElementById('session-id').value;
const senderName = document.getElementById('username').value;

console.log("senderId: " + senderId)
console.log("senderName: " + senderName)


const ID_ROOM_NAME_INPUT = "room-name-input";
const ID_ROOM_TABS_LIST = "room-tabs-list"
const ID_ROOM_NAMES_LIST = "ch-list"


const MSG_CREATE_ROOM = "CREATE_ROOM";
const MSG_USER_LEFT_ROOM = "USER_LEFT_ROOM"
const MSG_TEXT = "TEXT_MSG"




function escapeRoomName(name) {
    return name.replace(" ", "_");
}

function send(msgDict) {
    console.log("sending", msgDict);
    wsSocket.send(JSON.stringify(msgDict));
}


function createRoom() {
    console.log("createRoom");
    var roomName = document.getElementById(ID_ROOM_NAME_INPUT).value;
    var msgDict = {
        "msgType": MSG_CREATE_ROOM,
        "senderId": senderId,
        "room": roomName
    };
    send(msgDict)
}


function onConnect(event) {
    console.log(event);
    var connInfo = document.getElementById('connection-info');
    connInfo.style.display = 'none';
    var panels = document.getElementById('panels');
    panels.style.display = 'block';
    var logoutInfo = document.getElementById('logout-info');
    logoutInfo.style.display = 'block';
    var createRoomBtn = document.getElementById('ch-create');
    createRoomBtn.onclick = createRoom;
}


function setFocusOnTab(roomName) {
    var escaped = escapeRoomName(roomName);

    console.log("setFocusOnTab" + roomName);
    var tabs = document.getElementById(ID_ROOM_TABS_LIST).getElementsByTagName('li');
    for (let i = 0; i < tabs.length; i++) {
        tabs[i].classList.remove("active");
    }

    var tabsContent = document.getElementById('ch-contents');
    var children = tabsContent.children;
    for (var i = 0; i < children.length; i++) {
        children[i].style.display = "none";
    }

    document.getElementById('ch-' + escaped).classList.add("active");
    document.getElementById('content-' + escaped).style.display = "block";
    document.getElementById('msg-content-' + escaped).focus();

}


function isTabOpened(roomName) {
    var rooms = document.getElementById(ID_ROOM_TABS_LIST).getElementsByTagName('a');
    for (let i = 0; i < rooms.length; i++) {
        if(rooms[i].text == roomName) {
            return true;
        }
    }
    return false;
}


function generateSendMessageOnClickListener(room, msgInput) {
    return function() {
        var text = msgInput.value;
        if(!text){
            return;
        } 
        
        msg = {
            "senderId": senderId,
            "room": room,
            "senderName": senderName
        } 
        
        if(text == "exit") {
            msg["msgType"] = MSG_USER_LEFT_ROOM;
        } else {
            msg["msgType"] = MSG_TEXT;
            msg["content"] = text;
        }
        
        console.log("Sending message: '" + msg + "'");
    
        send(msg)

        msgInput.value = "";
    }
}


function addRoomTab(roomName) {
    var escaped = escapeRoomName(roomName);
    
    var roomLink = document.createElement("a");
    roomLink.text = roomName;
    roomLink.href = "#";
    roomLink.onclick = function() {
        console.log("clicked on " + roomName);
        setFocusOnTab(roomName);
    }

    var roomListElem = document.createElement("li");
    roomListElem.id = "ch-" + escaped;
    //withAttributes([strPair("role", "presentation")])
    roomListElem.appendChild(roomLink);

    var rooms = document.getElementById(ID_ROOM_TABS_LIST);
    rooms.appendChild(roomListElem);

    var conversationDiv = document.createElement("div");
    conversationDiv.id = "conversation-" + escaped;
    conversationDiv.style.dimaxHeightsplay = "400px";
    conversationDiv.style.overflowY = "scroll";

    var msgTextInput = document.createElement("input");
    msgTextInput.id = "msg-content-" + escaped;
    msgTextInput.classList.add("form-control");

    var sendMsgButton = document.createElement("button");
    sendMsgButton.type = "button";
    sendMsgButton.id = "msg-send-" + escaped;
    sendMsgButton.innerText = "Send";
    sendMsgButton.onclick = generateSendMessageOnClickListener(roomName, msgTextInput);
    sendMsgButton.classList.add("btn");
    sendMsgButton.classList.add("btn-default");

    var sendMsgSpan = document.createElement("span");
    sendMsgSpan.classList.add("input-group-btn");
    sendMsgSpan.appendChild(sendMsgButton);


    //.withOnKeyPressListener(handleEnter(_onSent))


    var inputGroupDiv = document.createElement("div");
    inputGroupDiv.classList.add("input-group");
    inputGroupDiv.appendChild(msgTextInput);
    inputGroupDiv.appendChild(sendMsgSpan);

    var contentDiv = document.createElement("div");
    contentDiv.id = "content-" + escaped;
    contentDiv.appendChild(document.createElement("br"));
    contentDiv.appendChild(inputGroupDiv);
    contentDiv.appendChild(document.createElement("br"));
    contentDiv.appendChild(conversationDiv);

    var content = document.getElementById('ch-contents');
    content.appendChild(contentDiv);
}


function createSelectRoomOnClickListener(roomName) {
    return function() { 
        var opened = isTabOpened(roomName);
        console.log('clicked ' + roomName + " opened: " + opened);
        
        if(!opened) {
            addRoomTab(roomName);
        }
        setFocusOnTab(roomName);
    }
}


function existingRooms() {
    var rooms = document.getElementById(ID_ROOM_NAMES_LIST).getElementsByTagName('a');
    var names = [];
    for (let i = 0; i < rooms.length; i++) {
        var room = rooms[i].text;
        names.push(room);
    }
    console.log("existing rooms: " + names);
    return names;
}


function addRoomToRoomsList(room) {
    console.log("addRoomToRoomsList: " + room);

    var eRooms = existingRooms();
    if(eRooms.includes(room)){
        return;
    }

    var rooms = document.getElementById(ID_ROOM_NAMES_LIST);

    var room = document.createElement("a");
    room.id = "ch-list-name-" + room;
    room.text = room;
    room.classList.add("list-group-item");
    room.href = "#";
    room.onclick = createSelectRoomOnClickListener(room);
    rooms.appendChild(room);
}


function removeRoomTab(room) {
    console.log("removeRoomTab: " + room);

    var escaped = escapeRoomName(room);

    const tabsPanel = document.getElementById(ID_ROOM_TABS_LIST);
    const tab = document.getElementById("ch-" + escaped);
    tabsPanel.removeChild(tab);

    setFocusOnTab("main");
}


function refreshRoomsList(roomList) {
    console.log("refreshRoomsList" + roomList);
    var rooms = document.getElementById(ID_ROOM_NAMES_LIST);
    while (rooms.firstChild) {
        rooms.firstChild.remove();
    }

    roomList.forEach((name) => {
        addRoomToRoomsList(name)
    });
}


function displayMessage(room, senderName, content) {
    var escaped = escapeRoomName(room);
    var conversationDiv = document.getElementById("conversation-" + escaped);

    var textParagraph = document.createElement("p");
    textParagraph.innerText = senderName + ": " + content;

    //console.log(escaped, conversationDiv, textParagraph);

    conversationDiv.appendChild(textParagraph);
}


function handleMessage(message) {
    var stringMsg = message['data'];
    console.log("Received: " + stringMsg);
    var jsonMsg = JSON.parse(stringMsg);
    var msgType = jsonMsg['msgType'];
    switch (msgType) {
        case 'ROOMS_LIST':
            var rooms = jsonMsg['rooms'];
            refreshRoomsList(rooms);
            break;
        case 'CREATE_ROOM':
            var room = jsonMsg['room'];
            addRoomToRoomsList(room);
            break;
        case 'REMOVE_ROOM':
            var room = jsonMsg['room'];
            removeRoomToRoomsList(room);
            break;
        case 'USER_JOINED_ROOM':
            var room = jsonMsg['room'];
            addRoomTab(room);
            setFocusOnTab(room);
            break;
        case MSG_USER_LEFT_ROOM:
            var room = jsonMsg['room'];
            removeRoomTab(room);
            break;
        case MSG_TEXT:
            displayMessage(jsonMsg['room'], jsonMsg['senderName'], jsonMsg['content'])
            break;
        case 'LOGOUT_USER':
            break;
        default:
          console.log(`Unknown message type ${msgType}.`);
      }
}


var host = window.location.hostname + (window.location.port != null ? ':' + window.location.port : '');
var wsSocket = new WebSocket("ws://" + host + "/talk");

wsSocket.onopen = onConnect;
wsSocket.onmessage = handleMessage;
