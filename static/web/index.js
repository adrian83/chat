
const senderId = document.getElementById('session-id').value;
const senderName = document.getElementById('username').value;

console.log("senderId: " + senderId);
console.log("senderName: " + senderName);

const MAIN_ROOM_NAME = "main";

const ID_ROOM_NAME_INPUT = "ch-name";
const ID_ROOM_TABS_LIST = "ch-tabs";
const ID_ROOM_NAMES_LIST = "ch-list";
const ID_PREFIX_ROOM_TAB = "room-";
const ID_PREFIX_CONVERSATION_PANEL = "conversation-";
const ID_PREFIX_MSG_CONTENT_PREFIX = "msg-content-";
const ID_PREFIX_CONTENT_PANEL = "content-";

const MSG_CREATE_ROOM = "CREATE_ROOM";
const MSG_USER_LEFT_ROOM = "USER_LEFT_ROOM";
const MSG_TEXT = "TEXT_MSG";
const MSG_ROOMS_LIST = "ROOMS_LIST";
const MSG_REMOVE_ROOM = "REMOVE_ROOM";
const MSG_USER_JOINED_ROOM = "USER_JOINED_ROOM";
const MSG_LOGOUT = "LOGOUT_USER";




function escapeText(text) {
    return text.replace(" ", "_");
}

function createMessageInputId(roomName) {
    return ID_PREFIX_MSG_CONTENT_PREFIX + escapeText(roomName);
}

function createRoomTabId(roomName) {
    return ID_PREFIX_ROOM_TAB + escapeText(roomName);
}

function createConversationPanelId(roomName) {
    return ID_PREFIX_CONVERSATION_PANEL + escapeText(roomName);
}

function createContentPanelId(roomName) {
    return ID_PREFIX_CONTENT_PANEL + escapeText(roomName);
}

function send(msgDict) {
    console.log("sending", msgDict);
    wsSocket.send(JSON.stringify(msgDict));
}


function sendCreateRoomMessage() {
    console.log("sendCreateRoomMessage");
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
    document.getElementById('connection-info').style.display = 'none';
    document.getElementById('panels').style.display = 'block';
    document.getElementById('logout-info').style.display = 'block';
    document.getElementById('ch-create').onclick = sendCreateRoomMessage;
}


function setFocusOnTab(roomName) {
    console.log("setFocusOnTab: " + roomName);

    var tabs = document.getElementById(ID_ROOM_TABS_LIST).getElementsByTagName('li');
    for (let i = 0; i < tabs.length; i++) {
        tabs[i].classList.remove("active");
    }

    var tabsContent = document.getElementById('ch-contents');
    var children = tabsContent.children;
    for (var i = 0; i < children.length; i++) {
        children[i].style.display = "none";
    }

    document.getElementById(createRoomTabId(roomName)).classList.add("active");
    document.getElementById(createContentPanelId(roomName)).style.display = "block";
    document.getElementById(createMessageInputId(roomName)).focus();
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
            if(room == MAIN_ROOM_NAME) {
                return;
            }
            msg["msgType"] = MSG_USER_LEFT_ROOM;
        } else if(text == "logout") {
            msg["msgType"] = MSG_LOGOUT;
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
    var roomLink = document.createElement("a");
    roomLink.text = roomName;
    roomLink.href = "#";
    roomLink.onclick = () => setFocusOnTab(roomName);
    

    var roomListElem = document.createElement("li");
    roomListElem.id = createRoomTabId(roomName);
    roomListElem.appendChild(roomLink);

    var rooms = document.getElementById(ID_ROOM_TABS_LIST);
    rooms.appendChild(roomListElem);

    var conversationDiv = document.createElement("div");
    conversationDiv.id = createConversationPanelId(roomName);
    conversationDiv.style.dimaxHeightsplay = "400px";
    conversationDiv.style.overflowY = "scroll";

    var msgTextInput = document.createElement("input");
    msgTextInput.id = createMessageInputId(roomName);
    msgTextInput.classList.add("form-control");

    var sendMsgButton = document.createElement("button");
    sendMsgButton.type = "button";
    sendMsgButton.innerText = "Send";
    sendMsgButton.onclick = generateSendMessageOnClickListener(roomName, msgTextInput);
    sendMsgButton.classList.add("btn");
    sendMsgButton.classList.add("btn-default");

    var sendMsgSpan = document.createElement("span");
    sendMsgSpan.classList.add("input-group-btn");
    sendMsgSpan.appendChild(sendMsgButton);

    var inputGroupDiv = document.createElement("div");
    inputGroupDiv.classList.add("input-group");
    inputGroupDiv.appendChild(msgTextInput);
    inputGroupDiv.appendChild(sendMsgSpan);

    var contentDiv = document.createElement("div");
    contentDiv.id = createContentPanelId(roomName);
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
        
        if(!opened) {
            addRoomTab(roomName);
        }

        setFocusOnTab(roomName);
    }
}


function existingRooms() {
    var roomsList = document.getElementById(ID_ROOM_NAMES_LIST);
    var rooms = roomsList.getElementsByTagName('a');
    var names = [];
    for (let i = 0; i < rooms.length; i++) {
        var roomName = rooms[i].text;
        names.push(roomName);
    }
    console.log("existing rooms: " + names);
    return names;
}


function addRoomToRoomsList(roomName) {
    console.log("addRoomToRoomsList: " + roomName);

    var existingRoomsList = existingRooms();
    if(existingRoomsList.includes(roomName)){
        return;
    }

    var roomLinkElem = document.createElement("a");
    roomLinkElem.text = roomName;
    roomLinkElem.classList.add("list-group-item");
    roomLinkElem.href = "#";
    roomLinkElem.onclick = createSelectRoomOnClickListener(roomName);

    var roomLinksList = document.getElementById(ID_ROOM_NAMES_LIST);
    roomLinksList.appendChild(roomLinkElem);
}


function removeRoomTab(roomName) {
    console.log("removeRoomTab: " + roomName);

    const tabsPanel = document.getElementById(ID_ROOM_TABS_LIST);
    const tab = document.getElementById(createRoomTabId(roomName));
    tabsPanel.removeChild(tab);

    setFocusOnTab(MAIN_ROOM_NAME);
}

function removeRoomFromRoomsList(roomName) {
    console.log("removeRoomFromRoomsList: " + roomName);
    var ul = document.getElementById(ID_ROOM_NAMES_LIST);
    var items = ul.getElementsByTagName("li");

    var toDel = null;
    for (var i = 0; i < items.length; ++i) {
        if(items[i].text == roomName){
            toDel = items[i];
            break;
        }
    }
    if(toDel != null) {
        toDel.remove();
    }
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


function displayMessage(roomName, senderName, content) {
    var textParagraph = document.createElement("p");
    textParagraph.innerText = senderName + ": " + content;

    var conversationDiv = document.getElementById(createConversationPanelId(roomName));
    conversationDiv.appendChild(textParagraph);
}


function logout() {
    // _closeClient();
    // changeLocation('/logout');
    wsSocket.close();
    window.location.href = "/logout";
}


function handleMessage(message) {
    var stringMsg = message['data'];
    console.log("Received: " + stringMsg);
    var jsonMsg = JSON.parse(stringMsg);
    var msgType = jsonMsg['msgType'];
    switch (msgType) {
        case MSG_ROOMS_LIST:
            var rooms = jsonMsg['rooms'];
            refreshRoomsList(rooms);
            break;
        case MSG_CREATE_ROOM:
            var room = jsonMsg['room'];
            addRoomToRoomsList(room);
            break;
        case MSG_REMOVE_ROOM:
            var room = jsonMsg['room'];
            removeRoomFromRoomsList(room);
            break;
        case MSG_USER_JOINED_ROOM:
            var room = jsonMsg['room'];
            addRoomTab(room);
            setFocusOnTab(room);
            break;
        case MSG_USER_LEFT_ROOM:
            var room = jsonMsg['room'];
            removeRoomTab(room);
            removeRoomFromRoomsList(room);
            break;
        case MSG_TEXT:
            displayMessage(jsonMsg['room'], jsonMsg['senderName'], jsonMsg['content'])
            break;
        case 'LOGOUT_USER':
            logout();
            break;
        default:
          console.log(`Unknown message type ${msgType}.`);
      }
}


var host = window.location.hostname + (window.location.port != null ? ':' + window.location.port : '');
var wsSocket = new WebSocket("ws://" + host + "/talk");

wsSocket.onopen = onConnect;
wsSocket.onmessage = handleMessage;
