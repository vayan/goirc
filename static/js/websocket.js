var ws;

if(ws != null) {
  ws.close();
  ws = null;
}

var host = window.location.hostname;
ws = new WebSocket("ws://" + host + ":1112/ws");


ws.onopen = function() {
  console.log("open ws");
  if ($("#yuid").val() != "") {
    ws.send("co]"+$("#yuid").val());
  }
};

ws.onmessage = function(e) {
  console.log("receive : " + e.data);
  parse_irc(e.data);

};

ws.onclose = function(e) {
  console.log("close ws");
};
