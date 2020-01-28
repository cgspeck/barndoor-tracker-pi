import { h, Component } from "preact";
import style from "./style";

import TextField from "preact-material-components/TextField";
import "preact-material-components/TextField/style.css";

import { getAlignStatus, getLocationSettings } from "../../lib/settings";
import { setInterval } from "timers";

export default class Align extends Component {
  state = {
    locationSettings: {
      latitude: null
    },
    alignStatus: {
      azAligned: null,
      altAligned: null,
      currentAz: null,
      currentAlt: null
    }
  };

  async componentDidMount() {
    getLocationSettings()
      .then(r => {
        this.setState({ locationSettings: { ...r } });
        console.log("Starting Interval");
        this.timer = setInterval(this.refreshAlignmentStatus.bind(this), 500);
      })
      .catch(e => this.handleError(e));
  }

  handleError = e => {
    console.error("problem", e);
    this.setState({ error: e });
  };

  async refreshAlignmentStatus() {
    getAlignStatus().then(r => this.setState({ alignStatus: { ...r } }));
  }

  componentWillUnmount() {
    console.log("Clearing timer");
    clearInterval(this.timer._id);
  }

  statusLabel = (azAligned, altAligned) => {
    if (azAligned && altAligned) {
      return <h2>Aligned</h2>;
    }
    return <h2>Not Aligned</h2>;
  };

  azTarget = azTarget => {
    if (azTarget !== null) {
      return `AZ target: ${azTarget}`;
    }
  };

  calculateAzTarget = latitude => {
    if (latitude === null) return null;

    if (latitude < 0) {
      return 180;
    } else {
      return 0;
    }
  };

  azValue = (azTarget, currentAz, azAligned) => {
    if (azTarget !== null && currentAz !== null && azAligned !== null) {
      let indicator;

      if (azAligned) {
        indicator = "✔";
      } else if (azTarget == 180) {
        indicator = currentAz < azTarget ? "✘ >>" : "✘ <<";
      } else {
        indicator = currentAz > 180 ? "✘ >>" : "✘ <<";
      }

      return `${currentAz.toFixed(2)} ${indicator}`;
    }
  };

  altValue = (latitude, currentAlt, altAligned) => {
    if (latitude !== null && currentAlt !== null && altAligned !== null) {
      const altTarget = Math.abs(latitude);
      let indicator;

      if (altAligned) {
        indicator = "✔";
      } else {
        indicator = currentAlt > altTarget ? "✘ ▼▼" : "✘ ▲▲";
      }

      return `${currentAlt.toFixed(2)} ${indicator}`;
    }
  };

  altTarget = latitude => {
    if (latitude !== null) {
      const targetVal = Math.abs(latitude);

      return `Alt target: ${targetVal}`;
    }
  };

  render({}, { locationSettings, alignStatus }) {
    const { azAligned, altAligned, currentAz, currentAlt } = alignStatus;
    const { latitude, azError } = locationSettings;

    const azTarget = this.calculateAzTarget(latitude);

    return (
      <div class={style.align}>
        <h1>Align</h1>
        {this.statusLabel(azAligned, altAligned)}
        <div>
          <p>
            <TextField
              label="Azimuth"
              value={this.azValue(azTarget, currentAz, azAligned)}
              readOnly
            ></TextField>
            <br />
            {this.azTarget(azTarget)}
          </p>
          <p>
            <TextField
              label="Altitude"
              value={this.altValue(latitude, currentAlt, altAligned)}
              readOnly
            ></TextField>
            <br />
            {this.altTarget(latitude)}
          </p>
        </div>
      </div>
    );
  }
}
