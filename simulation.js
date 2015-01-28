var serversocket = new WebSocket("ws://localhost:8080/echo");

obj = { Data: [] };

serversocket.onopen = function() {
    serversocket.send("Connection init");
}

// Write message on receive
serversocket.onmessage = function(e) {
    document.getElementById('comms').innerHTML += "Received: " + e.data + "<br>";
};

function senddata() {
    setInterval( function(){
        for ( i=0; i<6; i++ ) {
            obj["Data"][i] = Math.ceil( Math.random() * 283);
        }
        data = JSON.stringify(obj)
        serversocket.send(data);
        document.getElementById('comms').innerHTML += "Sent: " + data + "<br>";
    }, 20);
}