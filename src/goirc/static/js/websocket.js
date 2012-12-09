var ws;

if(ws != null) {
  ws.close();
  ws = null;
}

var host = window.location.hostname;
ws = new WebSocket("ws://" + host + ":1111/ws");


ws.onopen = function() {
  console.log("open ws");

};

ws.onmessage = function(e) {
  console.log("receive :" + e.data);

};

ws.onclose = function(e) {
  console.log("close ws");
};
