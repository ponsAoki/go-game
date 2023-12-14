const express = require("express");
const path = require("path");

const app = express();

app.use(express.static("./game"));

app.get("/", (_, res) => {
  res.sendFile(path.join(__dirname, "/game/main.html"));
});

app.listen(3000, () => {
  console.log("Listening on port 3000");
});
