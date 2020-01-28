import { h, Component } from "preact";
import linkState from "linkstate";

import Button from "preact-material-components/Button";
import "preact-material-components/Button/style.css";
import TextField from "preact-material-components/TextField";
import "preact-material-components/TextField/style.css";

import style from "./style";

import { getLocationSettings, setLocationSettings } from "../../lib/settings";

export default class LocationSettings extends Component {
  state = {
    locationSettings: {
      latitude: null,
      magDeclination: null,
      azError: null,
      altError: null,
      xOffset: null,
      yOffset: null,
      zOffset: null
    },
    error: null,
    info: null
  };

  async componentDidMount() {
    getLocationSettings()
      .then(r => this.setState({ locationSettings: { ...r } }))
      .catch(e => this.handleError(e));
  }

  handleError = e => {
    console.error("problem", e);
    this.setState({ error: e });
  };

  onSubmit = e => {
    e.preventDefault();
    this.setState({ error: null, info: null });
    const {
      latitude,
      magDeclination,
      azError,
      altError,
      xOffset,
      yOffset,
      zOffset
    } = this.state.locationSettings;
    setLocationSettings(
      latitude,
      magDeclination,
      azError,
      altError,
      xOffset,
      yOffset,
      zOffset
    )
      .then(r => this.setState({ info: "Location settings updated." }))
      .catch(e => this.handleError(e));
  };

  errorToast() {
    if (this.state.error != null) {
      return <p>{this.state.error.toString()}</p>;
    }
  }

  infoToast() {
    if (this.state.info != null) {
      return <p>{this.state.info.toString()}</p>;
    }
  }

  render({}, { locationSettings }) {
    return (
      <div class={style.ap}>
        <h1>Location Settings</h1>
        {this.infoToast()}
        {this.errorToast()}
        <form onSubmit={this.onSubmit.bind(this)}>
          <h2>Basic Settings</h2>
          <p>
            <TextField
              label="Latitude (Decimal degrees)"
              value={locationSettings.latitude}
              onInput={linkState(this, "locationSettings.latitude")}
            ></TextField>
          </p>
          <p>
            <TextField
              label="Magnetic Declination"
              value={locationSettings.magDeclination}
              onInput={linkState(this, "locationSettings.magDeclination")}
            ></TextField>
          </p>

          <h2>Advanced Settings</h2>
          <p>How many degrees over or under is permissable for alignment.</p>
          <p>
            <TextField
              label="Azimuth Error"
              value={locationSettings.azError}
              onInput={linkState(this, "locationSettings.azError")}
            ></TextField>
          </p>
          <p>
            <TextField
              label="Altitude Error"
              value={locationSettings.altError}
              onInput={linkState(this, "locationSettings.altError")}
            ></TextField>
          </p>
          <p>
            Change these settings only if the device does not show zero pitch
            when laying flat and level and zero degrees pointed in the direction
            of the magnetic north or south pole.
          </p>
          <p>
            <TextField
              label="X offset"
              value={locationSettings.xOffset}
              onInput={linkState(this, "locationSettings.xOffset")}
            ></TextField>
          </p>
          <p>
            <TextField
              label="Y offset"
              value={locationSettings.yOffset}
              onInput={linkState(this, "locationSettings.yOffset")}
            ></TextField>
          </p>
          <p>
            <TextField
              label="Z offset"
              value={locationSettings.zOffset}
              onInput={linkState(this, "locationSettings.zOffset")}
            ></TextField>
          </p>
          <Button raised ripple onClick={this.onSubmit.bind(this)}>
            Update
          </Button>
        </form>
      </div>
    );
  }
}
