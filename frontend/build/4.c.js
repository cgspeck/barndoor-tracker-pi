(window.webpackJsonp=window.webpackJsonp||[]).push([[4],{N1tf:function(t,n,e){"use strict";function r(t){return new Promise((function(n){return n(g.a.post(f.a.endpoint+"/backend/toggle/ignoreAz",{enabled:t}).then((function(t){return t.data.enabled})))}))}function o(t){return new Promise((function(n){return n(g.a.post(f.a.endpoint+"/backend/toggle/ignoreAlt",{enabled:t}).then((function(t){return t.data.enabled})))}))}function u(t){return new Promise((function(n){return n(g.a.post(f.a.endpoint+"/backend/toggle/intervalometer",{enabled:t}).then((function(t){return t.data.enabled})))}))}function a(t){return new Promise((function(n){return n(g.a.post(f.a.endpoint+"/backend/toggle/dewcontroller",{enabled:t}).then((function(t){return t.data.enabled})))}))}function i(t){return new Promise((function(n){return n(g.a.post(f.a.endpoint+"/backend/toggle/dewcontroller/logging",{enabled:t}).then((function(t){return t.data.enabled})))}))}function l(){return new Promise((function(t){return t(g.a.post(f.a.endpoint+"/backend/track",{command:"home"}).then((function(t){return t.data.state})))}))}function c(){return new Promise((function(t){return t(g.a.post(f.a.endpoint+"/backend/track",{command:"track"}).then((function(t){return t.data.state})))}))}function d(){return new Promise((function(t){return t(g.a.post(f.a.endpoint+"/backend/track",{command:"stop"}).then((function(t){return t.data.state})))}))}e.d(n,"g",(function(){return r})),e.d(n,"f",(function(){return o})),e.d(n,"h",(function(){return u})),e.d(n,"d",(function(){return a})),e.d(n,"e",(function(){return i})),e.d(n,"a",(function(){return l})),e.d(n,"b",(function(){return c})),e.d(n,"c",(function(){return d}));var s=e("czhI"),g=e.n(s),f=e("rOcY")},"Q/Sb":function(t){t.exports={align:"align__3YvuJ",textField:"textField__3sCSg"}},w7H0:function(t,n,e){"use strict";function r(t,n){var e=Object.keys(t);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(t);n&&(r=r.filter((function(n){return Object.getOwnPropertyDescriptor(t,n).enumerable}))),e.push.apply(e,r)}return e}function o(t){for(var n=1;n<arguments.length;n++){var e=null!=arguments[n]?arguments[n]:{};n%2?r(Object(e),!0).forEach((function(n){u(t,n,e[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(t,Object.getOwnPropertyDescriptors(e)):r(Object(e)).forEach((function(n){Object.defineProperty(t,n,Object.getOwnPropertyDescriptor(e,n))}))}return t}function u(t,n,e){return n in t?Object.defineProperty(t,n,{value:e,enumerable:!0,configurable:!0,writable:!0}):t[n]=e,t}e.r(n),e.d(n,"default",(function(){return p}));var a=e("hosL"),i=e("Q/Sb"),l=e.n(i),c=e("QafI"),d=e.n(c),s=(e("DFTH"),e("Qe1s")),g=e.n(s),f=(e("nU75"),e("xkhu")),b=e("N1tf"),h=e("oPUo"),p=function(t){function n(){for(var n,e=arguments.length,r=new Array(e),o=0;o<e;o++)r[o]=arguments[o];return(n=t.call.apply(t,[this].concat(r))||this).state={locationSettings:{latitude:null},ignoreAlt:!1,ignoreAz:!1,alignStatus:{azAligned:null,altAligned:null,currentAz:null,currentAlt:null}},n.handleError=function(t){console.error("problem",t),n.setState({error:t})},n.statusLabel=function(t,n){return t&&n?Object(a.h)("h2",null,"Aligned"):Object(a.h)("h2",null,"Not Aligned")},n.azTarget=function(t){if(null!==t)return"AZ target: "+t},n.calculateAzTarget=function(t){return null===t?null:t<0?180:0},n.azValue=function(t,n,e){var r;if(null!==t&&null!==n&&null!==e)return r=e?"✔":180==t?n<t?"✘ >>":"✘ <<":n>180?"✘ >>":"✘ <<",n.toFixed(2)+" "+r},n.altValue=function(t,n,e){if(null!==t&&null!==n&&null!==e){var r,o=Math.abs(t);return r=e?"✔":n>o?"✘ ▼▼":"✘ ▲▲",n.toFixed(2)+" "+r}},n.altTarget=function(t){if(null!==t)return"Alt target: "+Math.abs(t)},n.onIgnoreAzToggled=function(t){var e=t.target.checked;console.log("IgnoreAz toggled to: "+(e?"enabled":"disabled")),n.setState({ignoreAz:e}),Object(b.g)(e).then((function(t){return n.setState({ignoreAz:t})})).catch((function(t){return n.handleError(t)}))},n.onIgnoreAltToggled=function(t){var e=t.target.checked;console.log("IgnoreAlt toggled to: "+(e?"enabled":"disabled")),n.setState({ignoreAlt:e}),Object(b.f)(e).then((function(t){return n.setState({ignoreAlt:t})})).catch((function(t){return n.handleError(t)}))},n}var e,r;r=t,(e=n).prototype=Object.create(r.prototype),e.prototype.constructor=e,e.__proto__=r;var u=n.prototype;return u.componentDidMount=function(){return new Promise(function(t){var n=this;return Object(f.g)().then((function(t){n.setState({locationSettings:{latitude:t.latitude}}),n.setState({ignoreAlt:t.ignoreAlt}),n.setState({ignoreAz:t.ignoreAz}),console.log("Starting Interval"),n.timer=Object(h.setInterval)(n.refreshAlignmentStatus.bind(n),500)})).catch((function(t){return n.handleError(t)})),t()}.bind(this))},u.refreshAlignmentStatus=function(){return new Promise(function(t){var n=this;return Object(f.b)().then((function(t){return n.setState({alignStatus:o({},t)})})),t()}.bind(this))},u.componentWillUnmount=function(){console.log("Clearing timer"),clearInterval(this.timer._id)},u.render=function(t,n){!function(t){if(null==t)throw new TypeError("Cannot destructure undefined")}(t);var e=n.alignStatus,r=n.ignoreAz,o=n.ignoreAlt,u=e.azAligned,i=e.altAligned,c=e.currentAz,s=e.currentAlt,f=n.locationSettings.latitude,b=this.calculateAzTarget(f);return Object(a.h)("div",{class:l.a.align},Object(a.h)("h1",null,"Align"),this.statusLabel(u,i),Object(a.h)("div",null,Object(a.h)("p",null,Object(a.h)(d.a,{label:"Azimuth",value:this.azValue(b,c,u),readOnly:!0}),Object(a.h)("br",null),this.azTarget(b)),Object(a.h)("p",null,"Ignore Azimuth:"," ",Object(a.h)(g.a,{onChange:this.onIgnoreAzToggled.bind(this),checked:!0===r})),Object(a.h)("p",null,Object(a.h)(d.a,{label:"Altitude",value:this.altValue(f,s,i),readOnly:!0}),Object(a.h)("br",null),this.altTarget(f)),Object(a.h)("p",null,"Ignore Altitude:"," ",Object(a.h)(g.a,{onChange:this.onIgnoreAltToggled.bind(this),checked:!0===o}))))},n}(a.Component)}}]);
//# sourceMappingURL=4.c.js.map