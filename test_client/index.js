let ws;

const messagesDiv = document.getElementById("messages");
const msgInput = document.getElementById("msgInput");
const tokenInput = document.getElementById("token");
const connectBtn = document.getElementById("connectBtn");
const sendBtn = document.getElementById("sendBtn");

function log(msg) {
    const p = document.createElement("p");
    p.textContent = msg;
    messagesDiv.appendChild(p);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}


connectBtn.onclick = () => {
    const access_token = tokenInput.value;
    const connectionUrl = 'ws://' + window.location.host + '/ws'
    ws = new WebSocket(connectionUrl);

    ws.onopen = () => {
        log("Connected to server");
        ws.send(JSON.stringify({
            access_token,
            type: 'auth'
        }));
    };

    ws.onmessage = (event) => {
        log("Received: " + event.data);
    };

    ws.onclose = (event) => {
        log(`Connection closed: ${event.code}, ${event.reason}`);
    };

    ws.onerror = (event) => {
        log("Error: " + event);
    };
};

sendBtn.onclick = () => {
    const msg = msgInput.value;
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(msg);
        log("Sent: " + msg);
        msgInput.value = "";
    } else {
        log("WebSocket not connected");
    }
};
