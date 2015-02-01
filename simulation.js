$( document ).ready(function() {

    var serversocket = new WebSocket("ws://localhost:8080/echo");
    var sensors = []
    var sensor0Interval = null

    var sensorTypes = {
        "heat": {
            minVal: 0,
            maxVal: 2500,
            initVal: 0,
            upIntervalMs: 20,
            downIntervalMs: 60,
            upFunc: function( s ) {
                if ( s.value < s.maxVal ) { 
                    s.value++ 
                } else {
                    clearInterval(s.interval);
                }
            },
            downFunc:function( s ) { 
                if ( s.value > s.minVal ) {
                    s.value-- 
                } else {
                    clearInterval(s.interval);
                }
            }
        }
    }

    serversocket.onopen = function() {
        serversocket.send("Connection init");
    }

    // Write message on receive
    serversocket.onmessage = function(e) {
        document.getElementById('comms').innerHTML += "Received: " + e.data + "<br>";
    };

    $(".sensor").mousedown(function() {
        sensor = getSensor( $(this) );
        updateSensor( 
            sensor, 
            sensor.type.upIntervalMs,
            sensor.type.upFunc
        );
    });

    $(".sensor").mouseup(function() {
        sensor = getSensor( $(this) );
        updateSensor( 
            sensor, 
            sensor.type.downIntervalMs,
            sensor.type.downFunc  
        );
    });

    function updateSensor( sensor, intervalMs, updateFunc ) {
        clearInterval(sensor.interval);
        sensor.interval = setInterval( function(){
            updateFunc(sensor);
            data = JSON.stringify(sensor)
            serversocket.send(data);
            $( "#sensor-" + sensor.id ).html("<span>" + sensor.value + "</span>")
            // $('#comms').append( "Sent: " + data + "<br>");
        }, intervalMs );
    }

    function getSensor( sensor ) {
        sensorId = sensor.attr('id')
        if ( sensors[ sensorId ] === undefined ) {
            console.log( "Creating new sensor " + sensorId)
            sensors[ sensorId ] = {
                interval: null,
                value: sensorTypes[ sensor.attr("type")].initVal,
                minVal: sensorTypes[ sensor.attr("type")].minVal,
                maxVal: sensorTypes[ sensor.attr("type")].maxVal,
                id: parseInt( sensorId.split("-")[1] ),
                type: sensorTypes[ sensor.attr("type")]
            }
        }
        // console.log( sensors[ sensorId ] )
        return sensors[ sensorId ]
    }

});
