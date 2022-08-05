const LEFT = "left";
const RIGHT = "right";

const EVENT_MESSAGE = "message";
const EVENT_OTHER = "other";

const userPhotots = [
    "/assets/img/1.png",
    "/assets/img/2.png",
    "/assets/img/3.png",
    "/assets/img/4.png",
    "/assets/img/5.png",
    "/assets/img/6.png",
    "/assets/img/7.png",
    "/assets/img/8.png",
    "/assets/img/9.png",
    "/assets/img/10.png"
];

var PERSON_IMG = userPhotots[getRandomNum(0, userPhotots.length)];
var PERSON_NAME = "Guset" + Math.floor(Math.random() * 1000);

var url = "ws://" + window.location.host + "/ws?id=" + PERSON_NAME;
var ws = new WebSocket(url);
var chatroom = document.getElementsByClassName("msger-chat");
var text = document.getElementById("msg");
var send = document.getElementById("send");

function  getRandomNum(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

send.onclick = function (e) {
    handleMessageEvent()
}

text.onkeydown = function (e) {
    if (e.keyCode === 13 && text.value !== "") {
        handleMessageEvent()
    }
}

function handleMessageEvent() {
    ws.send(JSON.stringify({
        "event": "message",
        "photo": PERSON_IMG,
        "name": PERSON_NAME,
        "content": text.value,
    }));
    text.value = "";
}

ws.onmessage = function (e) {
    var m = JSON.parse(e.data)
    var msg = ""
    switch (m.event) {
        case EVENT_MESSAGE:
            if (m.name == PERSON_NAME) {
                msg = getMessage(m.name, m.photo, RIGHT, m.content);
            } else {
                msg = getMessage(m.name, m.photo, LEFT, m.content);
            }
            break;
        case EVENT_OTHER:
            if (m.name != PERSON_NAME) {
                msg = getEventMessage(m.name + " " + m.content);
            } else {
                msg = getEventMessage("您已" + m.content);
            }
            break;
        default:
            break;
    }
    insertMsg(msg, chatroom[0]);
}

function getEventMessage(msg) {
    var msg = `<div class="msg-left">${msg}</div>`
    return msg;
}

function getMessage(name, img, side, text) {
    const d = new Date()
    var msg = `
        <div class="msg ${side}-msg">
            <div class="msg-img" style="background-image:url(${img})"></div>
            <div class="msg-bubble">
                <div class="msg-info">
                    <div class="msg-info-name">${name}</div>
                    <div class="msg-info-time">${d.getFullYear()}/${d.getMonth()}/${d.getDay()} ${d.getHours()}:${d.getMinutes()}</div>
                </div>
                <div class="msg-text">${text}</div>
            </div>
        </div>
    `;
    return msg;
}

function insertMsg(msg, domObj) {
    domObj.insertAdjacentHTML("beforeend", msg);
    domObj.scrollTop += 500;
}