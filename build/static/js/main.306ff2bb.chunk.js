(this["webpackJsonphall-of-fame"]=this["webpackJsonphall-of-fame"]||[]).push([[0],{22:function(e,t,a){e.exports=a.p+"static/media/07.d6c382f2.png"},24:function(e,t,a){e.exports=a(49)},29:function(e,t,a){},30:function(e,t,a){},49:function(e,t,a){"use strict";a.r(t);var n=a(0),l=a.n(n),o=a(9),r=a.n(o),c=(a(29),a(1)),i=a(3),u=a(4),s=a(6),m=a(5),h=a(7),d=a(22),b=a.n(d),g=(a(30),a(23)),v=a.n(g),p=a(11),f=a(8),E=a.n(f),k=function(e){function t(){return Object(i.a)(this,t),Object(s.a)(this,Object(m.a)(t).apply(this,arguments))}return Object(h.a)(t,e),Object(u.a)(t,[{key:"render",value:function(){return l.a.createElement("nav",{class:"animated fadeInDown"},l.a.createElement("ul",null,l.a.createElement("a",{href:"/rules"},l.a.createElement("li",null,"Guidelines")),l.a.createElement("a",{href:"/",className:"main animated flipInX"},l.a.createElement("li",null,"H U N T")),l.a.createElement("a",{href:"/leaderboardtable",id:"leaderboard-nav-link"},l.a.createElement("li",null,"Leaderboard"))))}}]),t}(l.a.Component),S=function(e){function t(){return Object(i.a)(this,t),Object(s.a)(this,Object(m.a)(t).apply(this,arguments))}return Object(h.a)(t,e),Object(u.a)(t,[{key:"parseHash",value:function(){this.auth0=new p.a.WebAuth({domain:"zozimus-hunt.auth0.com",clientID:"n6UsLego812KDbGksvSi5KQdZ8okPBQ2"}),this.auth0.parseHash((function(e,t){if(e)return console.log(e);null!==t&&null!==t.accessToken&&null!==t.idToken&&(localStorage.setItem("access_token",t.accessToken),localStorage.setItem("id_token",t.idToken),localStorage.setItem("email",JSON.stringify(t.idTokenPayload)),window.location=window.location.href.substr(0,window.location.href.indexOf("")))}))}},{key:"setup",value:function(){v.a.ajaxSetup({beforeSend:function(e){localStorage.getItem("access_token")&&e.setRequestHeader("Authorization","Bearer "+localStorage.getItem("access_token"))}})}},{key:"setState",value:function(){var e=localStorage.getItem("id_token");this.loggedIn=!!e}},{key:"logout",value:function(){localStorage.removeItem("id_token"),localStorage.removeItem("access_token"),localStorage.removeItem("profile"),window.location.reload()}},{key:"componentWillMount",value:function(){this.setup(),this.parseHash(),this.setState()}},{key:"renderBody",value:function(){return this.loggedIn?l.a.createElement("div",null," ",l.a.createElement(k,null)," ",l.a.createElement(j,null)):l.a.createElement("div",null,l.a.createElement(k,null),l.a.createElement(w,null))}},{key:"render",value:function(){return void 0==this.loggedIn?l.a.createElement("div",{className:"loader"}):this.renderBody()}}]),t}(l.a.Component),j=function(e){function t(e){var a;return Object(i.a)(this,t),(a=Object(s.a)(this,Object(m.a)(t).call(this,e))).state={value:"",level:"",client:{}},a.handleChange=a.handleChange.bind(Object(c.a)(a)),a.fetchLevel=a.fetchLevel.bind(Object(c.a)(a)),a}return Object(h.a)(t,e),Object(u.a)(t,[{key:"handleChange",value:function(e){this.setState({value:e.target.value})}},{key:"fetchLevel",value:function(){var e=this,t="https://hunt.zozimus.in/whichlevel/"+JSON.parse(localStorage.getItem("email")).email;fetch(t).then((function(e){return e.json()})).then((function(t){e.setState({level:t.message})}))}},{key:"componentDidMount",value:function(){this.fetchLevel()}},{key:"render",value:function(){var e=this.state.level;if(!e)return l.a.createElement("div",{className:"loader"});switch(e){case"-2":return l.a.createElement(C,null);case"-1":return l.a.createElement(O,null);default:return l.a.createElement(y,null)}}}]),t}(l.a.Component),y=(l.a.Component,function(e){function t(e){var a;return Object(i.a)(this,t),(a=Object(s.a)(this,Object(m.a)(t).call(this,e))).state={value:"",url:"",level:-3},a.handleChange=a.handleChange.bind(Object(c.a)(a)),a.handleSubmit=a.handleSubmit.bind(Object(c.a)(a)),a}return Object(h.a)(t,e),Object(u.a)(t,[{key:"logout",value:function(){localStorage.removeItem("id_token"),localStorage.removeItem("access_token"),localStorage.removeItem("profile"),window.location.reload()}},{key:"handleChange",value:function(e){this.setState({value:e.target.value})}},{key:"handleSubmit",value:function(e){e.preventDefault();var t={headers:{"Content-Type":"application/json",authorization:"Bearer "+localStorage.getItem("access_token")}},a="https://hunt.zozimus.in/answer/"+this.state.level.toString()+"/"+this.state.value+"?id_token="+localStorage.getItem("id_token");E.a.get(a,t).then((function(e){window.location.reload()})).catch((function(e){localStorage.clear(),window.location.reload()}))}},{key:"componentWillMount",value:function(){var e=this,t={headers:{"Content-Type":"application/json",authorization:"Bearer "+localStorage.getItem("access_token")}},a="https://hunt.zozimus.in/level?id_token="+localStorage.getItem("id_token");E.a.get(a,t).then((function(t){e.setState({url:t.data.URL,level:t.data.Level}),console.log(t)})).catch((function(e){localStorage.clear(),window.location.reload()}))}},{key:"render",value:function(){return""==this.state.url?l.a.createElement("div",{className:"loader"}):l.a.createElement("div",{className:"won congrats"},l.a.createElement("p",{className:"mobile",dangerouslySetInnerHTML:{__html:this.state.url}}),l.a.createElement("form",{onSubmit:this.handleSubmit},l.a.createElement("input",{type:"name",className:"answerTextbox",value:this.state.value,onChange:this.handleChange}),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("input",{type:"submit",className:"answer-button",value:"Submit"})),l.a.createElement("br",null),l.a.createElement("button",{onClick:this.logout},"Logout"))}}]),t}(l.a.Component)),C=function(e){function t(e){var a;return Object(i.a)(this,t),(a=Object(s.a)(this,Object(m.a)(t).call(this,e))).state={value:"",name1:"",name2:"",name3:"",name4:"",name5:""},a.handleChange=a.handleChange.bind(Object(c.a)(a)),a.handleChange1=a.handleChange1.bind(Object(c.a)(a)),a.handleChange2=a.handleChange2.bind(Object(c.a)(a)),a.handleChange3=a.handleChange3.bind(Object(c.a)(a)),a.handleChange4=a.handleChange4.bind(Object(c.a)(a)),a.handleChange5=a.handleChange5.bind(Object(c.a)(a)),a.handleSubmit=a.handleSubmit.bind(Object(c.a)(a)),a}return Object(h.a)(t,e),Object(u.a)(t,[{key:"logout",value:function(){localStorage.removeItem("id_token"),localStorage.removeItem("access_token"),localStorage.removeItem("profile"),window.location.reload()}},{key:"handleChange",value:function(e){this.setState({value:e.target.value})}},{key:"handleChange1",value:function(e){this.setState({name1:e.target.value})}},{key:"handleChange2",value:function(e){this.setState({name2:e.target.value})}},{key:"handleChange3",value:function(e){this.setState({name3:e.target.value})}},{key:"handleChange4",value:function(e){this.setState({name4:e.target.value})}},{key:"handleChange5",value:function(e){this.setState({name5:e.target.value})}},{key:"handleSubmit",value:function(e){var t=this;e.preventDefault();var a="https://hunt.zozimus.in/doesUsernameExist/"+this.state.value;fetch(a).then((function(e){return e.json()})).then((function(e){if("true"==e.message)alert("That username exists");else{var a={headers:{"Content-Type":"application/json",authorization:"Bearer "+localStorage.getItem("access_token")}},n="https://hunt.zozimus.in/adduser/"+JSON.parse(localStorage.getItem("email")).email+"/"+t.state.value+"/"+t.state.name1+"/"+t.state.name2+"/"+t.state.name3+"/"+t.state.name4+"/"+t.state.name5;console.log(n),E.a.get(n,a).then((function(){window.location.reload()})).catch((function(e){localStorage.clear(),window.location.reload()}))}}))}},{key:"render",value:function(){return l.a.createElement("div",{className:"username-form won"},l.a.createElement("p",null,"You are logged in, ",JSON.parse(localStorage.getItem("email")).email,"."," "),l.a.createElement("p",null,"Give us a username."),l.a.createElement("form",{onSubmit:this.handleSubmit},l.a.createElement("input",{type:"name",className:"username",value:this.state.value,onChange:this.handleChange,placeholder:"Team Name"}),l.a.createElement("br",null),l.a.createElement("input",{type:"name",className:"username",value:this.state.name1,onChange:this.handleChange1,placeholder:"Participant Name"}),l.a.createElement("br",null),l.a.createElement("input",{type:"name",className:"username",value:this.state.name2,onChange:this.handleChange2,placeholder:"Participant Name"}),l.a.createElement("br",null),l.a.createElement("input",{type:"name",className:"username",value:this.state.name3,onChange:this.handleChange3,placeholder:"Participant Name"}),l.a.createElement("br",null),l.a.createElement("input",{type:"name",className:"username",value:this.state.name4,onChange:this.handleChange4,placeholder:"Participant Name"}),l.a.createElement("br",null),l.a.createElement("input",{type:"name",className:"username",value:this.state.name5,onChange:this.handleChange5,placeholder:"Participant Name"}),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("input",{type:"submit",className:"dive",value:"Submit"})),l.a.createElement("br",null),l.a.createElement("button",{onClick:this.logout},"Logout"))}}]),t}(l.a.Component),O=function(e){function t(e){var a;return Object(i.a)(this,t),(a=Object(s.a)(this,Object(m.a)(t).call(this,e))).handleAccepted=a.handleAccepted.bind(Object(c.a)(a)),a}return Object(h.a)(t,e),Object(u.a)(t,[{key:"logout",value:function(){localStorage.removeItem("id_token"),localStorage.removeItem("access_token"),localStorage.removeItem("profile"),window.location.reload()}},{key:"handleAccepted",value:function(e){e.preventDefault();var t={headers:{"Content-Type":"application/json",authorization:"Bearer "+localStorage.getItem("access_token")}},a="https://hunt.zozimus.in/acceptedrules?id_token="+localStorage.getItem("id_token");E.a.get(a,t).then((function(){window.location.reload()})).catch((function(e){localStorage.clear(),window.location.reload()}))}},{key:"render",value:function(){return l.a.createElement("div",{className:"rules-container won"},l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("br",null),l.a.createElement("h1",{className:"rules"},"Rules"),l.a.createElement("div",{class:"rules",style:{textAlign:"left"}},l.a.createElement("div",{class:"rules-content"},l.a.createElement("ol",null,l.a.createElement("li",null,"Check out the rules"," ",l.a.createElement("a",{href:"https://hunt.zozimus.in/rules"},"here"),".")))),l.a.createElement("form",{onSubmit:this.handleAccepted},l.a.createElement("input",{type:"submit",className:"username-button",value:"I accept"})),l.a.createElement("br",null),l.a.createElement("button",{onClick:this.logout},"Logout"))}}]),t}(l.a.Component),w=function(e){function t(e){var a;return Object(i.a)(this,t),(a=Object(s.a)(this,Object(m.a)(t).call(this,e))).authenticate=a.authenticate.bind(Object(c.a)(a)),a}return Object(h.a)(t,e),Object(u.a)(t,[{key:"authenticate",value:function(){this.WebAuth=new p.a.WebAuth({domain:"zozimus-hunt.auth0.com",clientID:"n6UsLego812KDbGksvSi5KQdZ8okPBQ2",scope:"openid email",audience:"https://zozimus-hunt.auth0.com/api/v2/",responseType:"token id_token",redirectUri:"https://hunt.zozimus.in"}),this.WebAuth.authorize()}},{key:"render",value:function(){return l.a.createElement("div",null,l.a.createElement("br",null),l.a.createElement("div",{class:"jumbotron animated fadeIn"},l.a.createElement("img",{src:b.a,class:"main-image"}),l.a.createElement("p",{class:"jumbotron-heading animated fadeIn"},"The Hunt"),l.a.createElement("p",{class:"jumbotron-subtitle"},"Zozimus 2019"),l.a.createElement("p",{class:"jumbotron-subtitle"},"Online till 1900 hours, 10th November."," "),l.a.createElement("button",{className:"DiveInButton",onClick:this.authenticate},l.a.createElement("div",{className:"transform"},"D I V E \xa0 I N"))))}}]),t}(l.a.Component),I=(l.a.Component,S);Boolean("localhost"===window.location.hostname||"[::1]"===window.location.hostname||window.location.hostname.match(/^127(?:\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$/));r.a.render(l.a.createElement(I,null),document.getElementById("root")),"serviceWorker"in navigator&&navigator.serviceWorker.ready.then((function(e){e.unregister()}))}},[[24,1,2]]]);
//# sourceMappingURL=main.306ff2bb.chunk.js.map