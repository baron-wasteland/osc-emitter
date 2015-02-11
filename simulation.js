$( document ).ready(function() {

    var serversocket = new WebSocket("ws://localhost:8080/echo");
    var sensors = []
    var sensor0Interval = null

    // TODO query websocket server for sensor types ID map
    var SENSOR_CONTINUOUS = 0;
    var SENSOR_IMPULSE = 1;

    var sensorTypes = {
        "heat": {
            sensorType: SENSOR_CONTINUOUS,
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
                    s.value--;
                } else {
                    clearInterval(s.interval);
                }
            }
        },
        "whack": {
            sensorType: SENSOR_IMPULSE,
            minVal: 0,
            maxVal: 1024,
            initVal: 0,
            upIntervalMs: 20,
            downIntervalMs: 100,
            upFunc: function(s) {
                s.value = s.maxVal;
                clearInterval(s.interval);
            },
            downFunc: function(s) {
                if (s.value > s.minVal) {
                    s.value = Math.max(s.value-250, s.minVal);
                } else {
                    s.value = s.minVal;
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
                sensorType: sensorTypes[ sensor.attr("type")].sensorType,
                type: sensorTypes[ sensor.attr("type")]
            }
        }
        // console.log( sensors[ sensorId ] )
        return sensors[ sensorId ]
    }

});
