import React from "react";
import image from "./07.png";
import ReactDOM from "react-dom";
import { render } from "react-dom";
import "./App.css";
import $ from "jquery";
import auth0 from "auth0-js";
import axios from "axios";

const AUTH0_CLIENT_ID = "DxKXzYxf9B12bDwME3Og6YN6p8RoqvFA";
const AUTH0_DOMAIN = "dev-l0ini8h1.us.auth0.com";
const AUTH0_CALLBACK_URL = "https://hunt.zozimus.in";
const AUTH0_API_AUDIENCE = "https://dev-l0ini8h1.us.auth0.com/api/v2/";

// <div className="navbar" id="mainNavBar">
// 	<div className="container">
// 	<div className="row">
// 	<div className="two columns"><a><div className="nav-links">Rules</div></a></div>
//   		<div className="two columns"><a><div className="nav-links">Cryptex 2019</div></a></div>
// 	<div className="two columns"><a><div id="nav-links-main">C R Y P T E X</div></a></div>
//   		<div className="two columns"><a><div className="nav-links">Sponsors</div></a></div>
//   		<div className="three columns"><a><div className="nav-links">About Us</div></a></div>
//   		</div>
//   		</div>
// </div>
class Navbar extends React.Component {
  render() {
    return (
      <nav class="animated fadeInDown">
        <ul>
          <a href="/rules">
            <li>Guidelines</li>
          </a>
          <a href="/" className="main animated flipInX">
            <li>H U N T</li>
          </a>
          <a href="/leaderboardtable" id="leaderboard-nav-link">
            <li>Leaderboard</li>
          </a>
        </ul>
      </nav>
    );
  }
}

class App extends React.Component {
  parseHash() {
    this.auth0 = new auth0.WebAuth({
      domain: AUTH0_DOMAIN,
      clientID: AUTH0_CLIENT_ID
    });
    this.auth0.parseHash((err, authResult) => {
      if (err) {
        return console.log(err);
      }
      if (
        authResult !== null &&
        authResult.accessToken !== null &&
        authResult.idToken !== null
      ) {
        localStorage.setItem("access_token", authResult.accessToken);
        localStorage.setItem("id_token", authResult.idToken);
        localStorage.setItem(
          "email",
          JSON.stringify(authResult.idTokenPayload)
        );
        window.location = window.location.href.substr(
          0,
          window.location.href.indexOf("")
        );
      }
    });
  }

  setup() {
    $.ajaxSetup({
      beforeSend: function(xhr) {
        if (localStorage.getItem("access_token")) {
          xhr.setRequestHeader(
            "Authorization",
            "Bearer " + localStorage.getItem("access_token")
          );
        }
      }
    });
  }

  setState() {
    let idToken = localStorage.getItem("id_token");
    if (idToken) {
      this.loggedIn = true;
    } else {
      this.loggedIn = false;
    }
  }
  logout() {
    localStorage.removeItem("id_token");
    localStorage.removeItem("access_token");
    localStorage.removeItem("profile");
    window.location.reload();
  }
  componentWillMount() {
    this.setup();
    this.parseHash();
    this.setState();
  }
  renderBody() {
    if (this.loggedIn)
      return (
        <div>
          {" "}
          <Navbar /> <LoggedIn />
        </div>
      );
    else
      return (
        <div>
          <Navbar />
          <Home />
        </div>
      );
  }
  render() {
    return this.loggedIn == undefined ? (
      <div className="loader"></div>
    ) : (
      this.renderBody()
    );
  }
}
class LoggedIn extends React.Component {
  constructor(props) {
    super(props);
    this.state = { value: "", level: "", client: {} };
    this.handleChange = this.handleChange.bind(this);
    this.fetchLevel = this.fetchLevel.bind(this);
  }
  handleChange(event) {
    this.setState({ value: event.target.value });
  }

  fetchLevel() {
    let url =
      "https://hunt.zozimus.in/whichlevel/" +
      JSON.parse(localStorage.getItem("email")).email;
    fetch(url)
      .then(response => response.json())
      .then(result => {
        this.setState({ level: result.message });
      });
  }

  componentDidMount() {
    this.fetchLevel();
  }
  render() {
    const level = this.state.level;
    if (level) {
      switch (level) {
        case "-2":
          return <LevelUsername />;
          break;
        case "-1":
          return <LevelRules />;
          break;
        default:
          return <LevelText />;
      }
    } else {
      return <div className="loader"></div>;
    }
  }
}
class LevelWon extends React.Component {
  logout() {
    localStorage.removeItem("id_token");
    localStorage.removeItem("access_token");
    localStorage.removeItem("profile");
    window.location.reload();
  }
  render() {
    return (
      <div className="won congrats">
        You have won. <br /> Congrats.
        <br />
        <p class="mobile">Credits:</p>
        <p class="mobile">
          Questions by: Rajnish Gupta, Saurav Madhusoodanan, Rishika Rao, Riddhi
          Shah, Rachit Keerti Das,Angad Narula
        </p>
        <p class="mobile">
          Website by: Vishnu VS, Lambda Coordinator,IIT Hyderabad
        </p>
        <p class="mobile">
          Please give us your feedback{" "}
          <a href="https://forms.gle/cXDErCBHFpQva38k7">here</a>.
        </p>
        <button onClick={this.logout}>Logout</button>
      </div>
    );
  }
}

class LevelText extends React.Component {
  constructor(props) {
    super(props);
    this.state = { value: "", url: "", level: -3 };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }
  logout() {
    localStorage.removeItem("id_token");
    localStorage.removeItem("access_token");
    localStorage.removeItem("profile");
    window.location.reload();
  }
  handleChange(event) {
    this.setState({ value: event.target.value });
  }
  handleSubmit(event) {
    event.preventDefault();
    const headers = {
      headers: {
        "Content-Type": "application/json",
        authorization: "Bearer " + localStorage.getItem("access_token")
      }
    };
    let url =
      "https://hunt.zozimus.in/answer/" +
      this.state.level.toString() +
      "/" +
      this.state.value +
      "?id_token=" +
      localStorage.getItem("id_token");
    axios
      .get(url, headers)
      .then(result => {
        window.location.reload();
      })
      .catch(error => {
        localStorage.clear();
        window.location.reload();
      });
  }
  componentWillMount() {
    const headers = {
      headers: {
        "Content-Type": "application/json",
        authorization: "Bearer " + localStorage.getItem("access_token")
      }
    };
    let url =
      "https://hunt.zozimus.in/level?id_token=" +
      localStorage.getItem("id_token");
    axios
      .get(url, headers)
      .then(result => {
        this.setState({ url: result.data.URL, level: result.data.Level });
        console.log(result);
      })
      .catch(error => {
        localStorage.clear();
        window.location.reload();
      });
  }
  render() {
    if (this.state.url == "") {
      return <div className="loader"></div>;
    } else {
      return (
        <div className="won congrats">
          <p
            className="mobile"
            dangerouslySetInnerHTML={{ __html: this.state.url }}
          ></p>
          <p className="mobile">
            <a
              href="https://chat.whatsapp.com/D1AEiLDkCbs5GYplh5KK3U"
              style={{ color: "#61A8D6" }}
            >
              Discourse Forum
            </a>
          </p>
          <form onSubmit={this.handleSubmit}>
            <input
              type="name"
              className="answerTextbox"
              value={this.state.value}
              onChange={this.handleChange}
            />
            <br />
            <br />
            <input type="submit" className="answer-button" value="Submit" />
          </form>
          <br />
          <button onClick={this.logout}>Logout</button>
        </div>
      );
    }
  }
}

class LevelUsername extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      value: "Team Name",
      name1: "Participant Name",
      name2: "Participant Name",
      name3: "Participant Name",
      name4: "Participant Name",
      name5: "Participant Name"
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleChange1 = this.handleChange1.bind(this);
    this.handleChange2 = this.handleChange2.bind(this);
    this.handleChange3 = this.handleChange3.bind(this);
    this.handleChange4 = this.handleChange4.bind(this);
    this.handleChange5 = this.handleChange5.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }
  logout() {
    localStorage.removeItem("id_token");
    localStorage.removeItem("access_token");
    localStorage.removeItem("profile");
    window.location.reload();
  }
  handleChange(event) {
    this.setState({ value: event.target.value });
  }
  handleChange1(event) {
    this.setState({ name1: event.target.value });
  }
  handleChange2(event) {
    this.setState({ name2: event.target.value });
  }
  handleChange3(event) {
    this.setState({ name3: event.target.value });
  }
  handleChange4(event) {
    this.setState({ name4: event.target.value });
  }
  handleChange5(event) {
    this.setState({ name5: event.target.value });
  }
  handleSubmit(event) {
    event.preventDefault();
    let url = "https://hunt.zozimus.in/doesUsernameExist/" + this.state.value;
    fetch(url)
      .then(response => response.json())
      .then(result => {
        if (result.message === "true") {
          alert("This username already exists");
        } else {
          const headers = {
            headers: {
              "Content-Type": "application/json",
              authorization: "Bearer " + localStorage.getItem("access_token")
            }
          };
          var loginUrl =
            "https://hunt.zozimus.in" +
            "/adduser/" +
            JSON.parse(localStorage.getItem("email")).email +
            "/" +
            this.state.value +
            "/" +
            this.state.name1 +
            "/" +
            this.state.name2 +
            "/" +
            this.state.name3 +
            "/" +
            this.state.name4 +
            "/" +
            this.state.name5;
          console.log(loginUrl);
          axios
            .get(loginUrl, headers)
            .then(() => {
              //window.location.reload();
            })
            .catch(error => {
              localStorage.clear();
              console.log(error);
              //window.location.reload();
            });
        }
      });
  }
  render() {
    return (
      <div className="username-form won">
        <p>
          You are logged in, {JSON.parse(localStorage.getItem("email")).email}.{" "}
        </p>
        <p>Give us a username.</p>
        <form onSubmit={this.handleSubmit}>
          <input
            type="name"
            className="username"
            value={this.state.value}
            onChange={this.handleChange}
            placeholder="Team Name"
          />
          <br />
          <input
            type="name"
            className="username"
            value={this.state.name1}
            onChange={this.handleChange1}
            placeholder="Participant Name"
          />
          <br />
          <input
            type="name"
            className="username"
            value={this.state.name2}
            onChange={this.handleChange2}
            placeholder="Participant Name"
          />
          <br />
          <input
            type="name"
            className="username"
            value={this.state.name3}
            onChange={this.handleChange3}
            placeholder="Participant Name"
          />
          <br />
          <input
            type="name"
            className="username"
            value={this.state.name4}
            onChange={this.handleChange4}
            placeholder="Participant Name"
          />
          <br />
          <input
            type="name"
            className="username"
            value={this.state.name5}
            onChange={this.handleChange5}
            placeholder="Participant Name"
          />
          <br />
          <br />
          <input type="submit" className="dive" value="Submit" />
        </form>
        <br />
        <button onClick={this.logout}>Logout</button>
      </div>
    );
  }
}

class LevelRules extends React.Component {
  constructor(props) {
    super(props);
    this.handleAccepted = this.handleAccepted.bind(this);
  }
  logout() {
    localStorage.removeItem("id_token");
    localStorage.removeItem("access_token");
    localStorage.removeItem("profile");
    window.location.reload();
  }
  handleAccepted(event) {
    event.preventDefault();
    const headers = {
      headers: {
        "Content-Type": "application/json",
        authorization: "Bearer " + localStorage.getItem("access_token")
      }
    };
    let url =
      "https://hunt.zozimus.in/acceptedrules?id_token=" +
      localStorage.getItem("id_token");
    axios
      .get(url, headers)
      .then(() => {
        window.location.reload();
      })
      .catch(error => {
        localStorage.clear();
        window.location.reload();
      });
  }
  render() {
    return (
      <div className="rules-container won">
        <br />
        <br />
        <br />
        <br />
        <br />
        <br />
        <br />
        <br />
        <br />
        <br />
        <br />
        <br />
        <br />
        <br />
        <br />
        <h1 className="rules">Rules</h1>
        <div class="rules" style={{ textAlign: "left" }}>
          <div class="rules-content">
            <ol>
              <li>
                Check out the rules{" "}
                <a href="https://hunt.zozimus.in/rules">here</a>.
              </li>
            </ol>
          </div>
        </div>
        <form onSubmit={this.handleAccepted}>
          <input type="submit" className="username-button" value="I accept" />
        </form>
        <br />
        <button onClick={this.logout}>Logout</button>
      </div>
    );
  }
}

class Home extends React.Component {
  constructor(props) {
    super(props);
    this.authenticate = this.authenticate.bind(this);
  }
  authenticate() {
    this.WebAuth = new auth0.WebAuth({
      domain: AUTH0_DOMAIN,
      clientID: AUTH0_CLIENT_ID,
      scope: "openid email",
      audience: AUTH0_API_AUDIENCE,
      responseType: "token id_token",
      redirectUri: AUTH0_CALLBACK_URL
    });
    this.WebAuth.authorize();
  }

  render() {
    return (
      <div>
        <br />
        <div class="jumbotron animated fadeIn">
          <img src={image} class="main-image" />
          <p class="jumbotron-heading animated fadeIn">The Hunt</p>
          <p class="jumbotron-subtitle">Zozimus 2022</p>
          <p class="jumbotron-subtitle">
            Online till 1830 hours, 10th March.{" "}
          </p>
          <button className="DiveInButton" onClick={this.authenticate}>
            <div className="transform">D I V E &nbsp; I N</div>
          </button>
        </div>
      </div>
    );
  }
}

class Callback extends React.Component {
  render() {
    return <h1>Loading</h1>;
  }
}
class User {
  constructor(level, clientSecret, username) {
    this.level = level;
    this.clientSecret = clientSecret;
    this.username = username;
  }
  getLevel() {
    return this.level;
  }
  getClientSecret() {
    return this.clientSecret;
  }
  getUsername() {
    return this.username;
  }
}

export default App;
