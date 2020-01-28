import { h, Component } from "preact";
import { Router, route } from "preact-router";
import { getFlags } from "../lib/settings";

import Header from "./header";

// Code-splitting is automated for routes
import Debug from "../routes/debug";
import Align from "../routes/align";
import Track from "../routes/track";
import APSettings from "../routes/ap_settings";
import LocationSettings from "../routes/location_settings";

export default class App extends Component {
  state = {
    flags: {
      needsAPSettings: null,
      needsLocationSettings: null
    }
  };
  /** Gets fired when the route changes.
   *	@param {Object} event		"change" event from [preact-router](http://git.io/preact-router)
   *	@param {string} event.url	The newly routed URL
   */
  flagsAreTrueOrNull() {
    const { needsAPSettings, needsLocationSettings } = this.state.flags;

    return (
      needsAPSettings === null ||
      needsAPSettings === true ||
      needsLocationSettings === null || needsAPSettings === true
    );
  }

  async handleRoute(e) {
    var flags;
    var currentUrl;

    if (this.flagsAreTrueOrNull) {
      flags = await getFlags();
      this.setState({ flags: { ...flags } });
    } else {
      flags = this.state.flags;
    }

    if (e.url == "/debug") {
      currentUrl = "/debug";
    } else if (flags.needsAPSettings) {
      currentUrl = "/ap_settings";
    } else if (
      flags.needsLocationSettings &&
      !["/ap_settings", "/location_settings"].includes(e.url)
    ) {
      currentUrl = "/location_settings";
    } else {
      currentUrl = e.url;
    }

    if (currentUrl != e.url) {
      route(currentUrl);
    }

    this.setState({
      currentUrl: currentUrl
    });
  }

  render() {
    return (
      <div id="app">
        <Header selectedRoute={this.state.currentUrl} />
        <Router onChange={this.handleRoute.bind(this)}>
          <Debug path="/debug" />
          <Align path="/align" />
          <Track path="/" />
          <APSettings path="/ap_settings" />
          <LocationSettings path="/location_settings" />
        </Router>
      </div>
    );
  }
}
