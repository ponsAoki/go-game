const express = require("express");
const app = express();

app.use(express.static("."));

app.get("/", (_, res) => {
  res.redirect("/game/main.html");
});

app.listen(3000, () => {
  console.log("Listening on port 3000");
});
