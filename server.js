const express = require("express");
const bodyParser = require("body-parser");
const app = express();
const port = 8080;
const jwt = require("express-jwt");
const jwtAuthz = require("express-jwt-authz");
const jwksRsa = require("jwks-rsa");
const { Pool } = require("pg");
const pool = new Pool();
app.use(logger("dev"));
app.use(express.json());
app.use(
  express.urlencoded({
    extended: false
  })
);
const checkJwt = jwt({
  // Dynamically provide a signing key
  // based on the kid in the header and
  // the signing keys provided by the JWKS endpoint.
  secret: jwksRsa.expressJwtSecret({
    cache: true,
    rateLimit: true,
    jwksRequestsPerMinute: 5,
    jwksUri: `https://zozimus-hunt.auth0.com/.well-known/jwks.json`
  }),

  // Validate the audience and the issuer.
  audience: "https://dev-l0ini8h1.us.auth0.com/api/v2/",
  issuer: `https://dev-l0ini8h1.us.auth0.com/`,
  algorithms: ["RS256"]
});

app.use(express.json());

app.get("/whichlevel/", (req, res) => {
  console.log(req);
  const query = {
    text: "SELECT * from users WHERE email = $1",
    values: [req.body.email]
  };
  pool.query(query, (err, res) => {
    if (err) {
      console.log(err.stack);
    } else {
      res.send(res.rows[0].level);
    }
  });
});

app.post("/api/private/create", checkJwt, function(req, res) {
  const query = {
    text:
      "INSERT INTO users(teamname, name1, name2, name3, name4, name5, email, level) VALUES($1, $2, $3, $4, $5)",
    values: [
      req.body.teamname,
      req.body.name1,
      req.body.name2,
      req.body.name3,
      req.body.name4,
      req.body.name5,
      req.body.email,
      -1
    ]
  };
  pool.query(query, (err, res) => {
    if (err) {
      console.log(err.stack);
    } else {
      console.log(res.rows[0]);
    }
  });
});

app.listen(port, () => {
  console.log(`App running on port ${port}.`);
});
