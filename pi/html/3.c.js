(window.webpackJsonp=window.webpackJsonp||[]).push([[3],{"Q/Sb":function(e){e.exports={align:"align__3YvuJ",textField:"textField__3sCSg"}},heVN:function(e,t,n){(function(e){!function(e,t){"use strict";function n(e){delete s[e]}function r(e){if(f)setTimeout(r,0,e);else{var i=s[e];if(i){f=!0;try{!function(e){var n=e.callback,r=e.args;switch(r.length){case 0:n();break;case 1:n(r[0]);break;case 2:n(r[0],r[1]);break;case 3:n(r[0],r[1],r[2]);break;default:n.apply(t,r)}}(i)}finally{n(e),f=!1}}}}if(!e.setImmediate){var i,o,a,l,u,c=1,s={},f=!1,d=e.document,p=Object.getPrototypeOf&&Object.getPrototypeOf(e);p=p&&p.setTimeout?p:e,"[object process]"==={}.toString.call(e.process)?i=function(e){process.nextTick((function(){r(e)}))}:!function(){if(e.postMessage&&!e.importScripts){var t=!0,n=e.onmessage;return e.onmessage=function(){t=!1},e.postMessage("","*"),e.onmessage=n,t}}()?e.MessageChannel?((a=new MessageChannel).port1.onmessage=function(e){r(e.data)},i=function(e){a.port2.postMessage(e)}):d&&"onreadystatechange"in d.createElement("script")?(o=d.documentElement,i=function(e){var t=d.createElement("script");t.onreadystatechange=function(){r(e),t.onreadystatechange=null,o.removeChild(t),t=null},o.appendChild(t)}):i=function(e){setTimeout(r,0,e)}:(l="setImmediate$"+Math.random()+"$",u=function(t){t.source===e&&"string"==typeof t.data&&0===t.data.indexOf(l)&&r(+t.data.slice(l.length))},e.addEventListener?e.addEventListener("message",u,!1):e.attachEvent("onmessage",u),i=function(t){e.postMessage(l+t,"*")}),p.setImmediate=function(e){"function"!=typeof e&&(e=new Function(""+e));for(var t=new Array(arguments.length-1),n=0;n<t.length;n++)t[n]=arguments[n+1];return s[c]={callback:e,args:t},i(c),c++},p.clearImmediate=n}}("undefined"==typeof self?void 0===e?this:e:self)}).call(this,n("pCvA"))},oPUo:function(e,t,n){(function(e){function r(e,t){this._id=e,this._clearFn=t}var i=void 0!==e&&e||"undefined"!=typeof self&&self||window,o=Function.prototype.apply;t.setTimeout=function(){return new r(o.call(setTimeout,i,arguments),clearTimeout)},t.setInterval=function(){return new r(o.call(setInterval,i,arguments),clearInterval)},t.clearTimeout=t.clearInterval=function(e){e&&e.close()},r.prototype.unref=r.prototype.ref=function(){},r.prototype.close=function(){this._clearFn.call(i,this._id)},t.enroll=function(e,t){clearTimeout(e._idleTimeoutId),e._idleTimeout=t},t.unenroll=function(e){clearTimeout(e._idleTimeoutId),e._idleTimeout=-1},t._unrefActive=t.active=function(e){clearTimeout(e._idleTimeoutId);var t=e._idleTimeout;t>=0&&(e._idleTimeoutId=setTimeout((function(){e._onTimeout&&e._onTimeout()}),t))},n("heVN"),t.setImmediate="undefined"!=typeof self&&self.setImmediate||void 0!==e&&e.setImmediate||this&&this.setImmediate,t.clearImmediate="undefined"!=typeof self&&self.clearImmediate||void 0!==e&&e.clearImmediate||this&&this.clearImmediate}).call(this,n("pCvA"))},pCvA:function(e){var t;t=function(){return this}();try{t=t||new Function("return this")()}catch(e){"object"==typeof window&&(t=window)}e.exports=t},w7H0:function(e,t,n){"use strict";function r(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?r(n,!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):r(n).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}n.r(t),n.d(t,"default",(function(){return p}));var a=n("hosL"),l=n("Q/Sb"),u=n.n(l),c=n("QafI"),s=n.n(c),f=(n("DFTH"),n("xkhu")),d=n("oPUo"),p=function(e){function t(){for(var t,n=arguments.length,r=new Array(n),i=0;i<n;i++)r[i]=arguments[i];return(t=e.call.apply(e,[this].concat(r))||this).state={locationSettings:{latitude:null},alignStatus:{azAligned:null,altAligned:null,currentAz:null,currentAlt:null}},t.handleError=function(e){console.error("problem",e),t.setState({error:e})},t.statusLabel=function(e,t){return e&&t?Object(a.h)("h2",null,"Aligned"):Object(a.h)("h2",null,"Not Aligned")},t.azTarget=function(e){if(null!==e)return"AZ target: "+e},t.calculateAzTarget=function(e){return null===e?null:e<0?180:0},t.azValue=function(e,t,n){var r;if(null!==e&&null!==t&&null!==n)return r=n?"✔":180==e?t<e?"✘ >>":"✘ <<":t>180?"✘ >>":"✘ <<",t.toFixed(2)+" "+r},t.altValue=function(e,t,n){if(null!==e&&null!==t&&null!==n){var r,i=Math.abs(e);return r=n?"✔":t>i?"✘ ▼▼":"✘ ▲▲",t.toFixed(2)+" "+r}},t.altTarget=function(e){if(null!==e)return"Alt target: "+Math.abs(e)},t}var n,r;r=e,(n=t).prototype=Object.create(r.prototype),n.prototype.constructor=n,n.__proto__=r;var o=t.prototype;return o.componentDidMount=function(){return new Promise(function(e){var t=this;return Object(f.e)().then((function(e){t.setState({locationSettings:i({},e)}),console.log("Starting Interval"),t.timer=Object(d.setInterval)(t.refreshAlignmentStatus.bind(t),500)})).catch((function(e){return t.handleError(e)})),e()}.bind(this))},o.refreshAlignmentStatus=function(){return new Promise(function(e){var t=this;return Object(f.b)().then((function(e){return t.setState({alignStatus:i({},e)})})),e()}.bind(this))},o.componentWillUnmount=function(){console.log("Clearing timer"),clearInterval(this.timer._id)},o.render=function(e,t){!function(e){if(null==e)throw new TypeError("Cannot destructure undefined")}(e);var n=t.locationSettings,r=t.alignStatus,i=r.azAligned,o=r.altAligned,l=r.currentAz,c=r.currentAlt,f=n.latitude,d=this.calculateAzTarget(f);return Object(a.h)("div",{class:u.a.align},Object(a.h)("h1",null,"Align"),this.statusLabel(i,o),Object(a.h)("div",null,Object(a.h)("p",null,Object(a.h)(s.a,{label:"Azimuth",value:this.azValue(d,l,i),readOnly:!0}),Object(a.h)("br",null),this.azTarget(d)),Object(a.h)("p",null,Object(a.h)(s.a,{label:"Altitude",value:this.altValue(f,c,o),readOnly:!0}),Object(a.h)("br",null),this.altTarget(f))))},t}(a.Component)}}]);
//# sourceMappingURL=3.c.js.map