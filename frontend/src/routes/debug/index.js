import { h, Component } from "preact";
import style from "./style";
import { getAllSettings } from "../../lib/settings";

export default class Debug extends Component {
  state = {
    apSettings: {
      ssid: null,
      key: null
    },
    debug: {
      currentMillis: null
    },
    location: {
      latitude: null,
      magDeclination: null,
      altError: null,
      azError: null,
      xOffset: null,
      yOffset: null,
      zOffset: null,
      locationSet: null
    }
  };

  async componentDidMount() {
    getAllSettings().then(r => {
      console.log("got debug response", r);
      this.setState({ ...r });
    });
  }

  render({}, { apSettings, debug, location }) {
    return (
      <div class={style.debug}>
        <h1>Debug</h1>
        <a href="/config.json">View/Download config.json</a>
        <p>NODE_ENV: {JSON.stringify(process.env.NODE_ENV)}</p>
        <p>currentMillis: {debug.currentMillis}</p>
        <h2>AP Settings</h2>
        <p>SSID: {apSettings.ssid}</p>
        <p>Key: {apSettings.key}</p>
        <h2>Location Settings</h2>
        <p>Latitude: {location.latitude}</p>
        <p>magDeclination: {location.magDeclination}</p>
        <p>azError: {location.azError}</p>
        <p>altError: {location.altError}</p>
        <p>xOffset: {location.xOffset}</p>
        <p>yOffset: {location.yOffset}</p>
        <p>zOffset: {location.zOffset}</p>
        <p>locationSet: {location.locationSet ? "true" : "false"}</p>
      </div>
    );
  }
}
