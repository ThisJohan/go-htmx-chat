ws = new WebSocket("ws://localhost:5100/app/ws");

ws.onopen = function () {
  console.log("Connected");
};

ws.onmessage = function (evt) {
  console.log({ evt });
};

setInterval(function () {
  ws.send("Hello, Server!");
}, 2000);
