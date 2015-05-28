$( document ).ready(function() {

    var host = window.location.host
    var serversocket = new WebSocket("ws://" + host + "/managerws");

    serversocket.onopen = function() {
        serversocket.send("Connection init");
    }

    // Write message on receive
    serversocket.onmessage = function(e) {
        // console.log(e.data)
        obj = JSON.parse(e.data)
        // console.log("updating " + obj.id)
        // $("#s" + obj.id).html("");
        $("#s" + obj.id).jJsonViewer(e.data);
    };

});
