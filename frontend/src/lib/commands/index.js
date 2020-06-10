import axios from "axios";
import config from "../../config";


async function toggleIntervalometer(enabled) {
  return axios
  .post(`${config.endpoint}/track/toggle/intervalometer`, {
    enabled: enabled
  })
  .then(r => r.data);
}

async function toggleDewController(enabled) {
  return axios
  .post(`${config.endpoint}/track/toggle/dewcontroller`, {
    enabled: enabled
  })
  .then(r => r.data);
}

async function startHoming() {
  return axios
  .post(`${config.endpoint}/track`, {
    "command": "home"
  }).then(r => r.data.state);
}

async function startTracking() {
  return axios
  .post(`${config.endpoint}/track`, {
    "command": "track"
  }).then(r => r.data.state);
}

async function stopTracking() {
  return axios
  .post(`${config.endpoint}/track`, {
    "command": "stop"
  }).then(r => r.data.state);
}

export {
  toggleIntervalometer,
  toggleDewController,
  startHoming,
  startTracking,
  stopTracking
};
