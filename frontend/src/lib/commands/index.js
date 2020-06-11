import axios from 'axios';
import config from '../../config';

async function toggleIgnoreAz(enabled) {
  return axios
    .post(`${config.endpoint}/toggle/ignoreAz`, {
      enabled,
    })
    .then((r) => {
      return r.data.enabled;
    });
}

async function toggleIgnoreAlt(enabled) {
  return axios
    .post(`${config.endpoint}/toggle/ignoreAlt`, {
      enabled,
    })
    .then((r) => {
      return r.data.enabled;
    });
}

async function toggleIntervalometer(enabled) {
  return axios
    .post(`${config.endpoint}/toggle/intervalometer`, {
      enabled,
    })
    .then((r) => {
      return r.data.enabled;
    });
}

async function toggleDewController(enabled) {
  return axios
    .post(`${config.endpoint}/toggle/dewcontroller`, {
      enabled,
    })
    .then((r) => {
      return r.data.enabled;
    });
}

async function startHoming() {
  return axios
    .post(`${config.endpoint}/track`, {
      command: 'home',
    })
    .then((r) => r.data.state);
}

async function startTracking() {
  return axios
    .post(`${config.endpoint}/track`, {
      command: 'track',
    })
    .then((r) => r.data.state);
}

async function stopTracking() {
  return axios
    .post(`${config.endpoint}/track`, {
      command: 'stop',
    })
    .then((r) => r.data.state);
}

export {
  toggleIgnoreAz,
  toggleIgnoreAlt,
  toggleIntervalometer,
  toggleDewController,
  startHoming,
  startTracking,
  stopTracking,
};
