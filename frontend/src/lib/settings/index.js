import axios from 'axios';
import config from '../../config';

async function getAllSettings() {
  return axios
    .get(`${config.endpoint}/backend/status/debug`)
    .then((r) => r.data);
}

async function getAPSettings() {
  return axios
    .get(`${config.endpoint}/backend/settings/network/ap`)
    .then((r) => r.data);
}

async function setAPSettings(ssid, key) {
  // remove when frontend actually has wireless management features
  await axios.post(`${config.endpoint}/backend/settings/network`, {
    accessPointMode: true,
  });
  return axios
    .post(`${config.endpoint}/backend/settings/network/ap`, {
      ssid,
      key,
    })
    .then((r) => r.data);
}

async function getIntervalometerSettings() {
  return axios
    .get(`${config.endpoint}/backend/settings/intervalometer`)
    .then((r) => r.data);
}

async function setIntervalometerSettings(bulbInterval, restInterval) {
  return axios
    .post(
      `${config.endpoint}/backend/settings/intervalometer`,
      {
        bulbInterval,
        restInterval,
      },
      {
        transformRequest: [
          function (data, _) {
            const res = {
              bulbInterval: parseInt(data.bulbInterval, 10),
              restInterval: parseInt(data.restInterval, 10),
            };
            return JSON.stringify(res);
          },
        ],
      },
    )
    .then((r) => r.data);
}

async function getLocationSettings() {
  return axios
    .get(`${config.endpoint}/backend/settings/location`)
    .then((r) => r.data);
}

async function setLocationSettings(
  latitude,
  magDeclination,
  azError,
  altError,
  xOffset,
  yOffset,
  zOffset,
) {
  return axios
    .post(
      `${config.endpoint}/backend/settings/location`,
      {
        latitude,
        magDeclination,
        azError,
        altError,
        xOffset,
        yOffset,
        zOffset,
      },
      {
        transformRequest: [
          function (data, _) {
            const res = {
              latitude: parseFloat(data.latitude),
              magDeclination: parseFloat(data.magDeclination),
              azError: parseFloat(data.azError),
              altError: parseFloat(data.altError),
              xOffset: parseInt(data.xOffset, 10),
              yOffset: parseInt(data.yOffset, 10),
              zOffset: parseInt(data.zOffset, 10),
            };
            return JSON.stringify(res);
          },
        ],
      },
    )
    .then((r) => r.data);
}

async function getFlags() {
  return axios
    .get(`${config.endpoint}/backend/status/flags`)
    .then((r) => r.data);
}

async function getAlignStatus() {
  return axios
    .get(`${config.endpoint}/backend/status/align`)
    .then((r) => r.data);
}

async function getTrackState() {
  return axios.get(`${config.endpoint}/backend/track`).then((r) => r.data);
}

async function getDewControllerStatus() {
  return axios
    .get(`${config.endpoint}/backend/status/dew_controller`)
    .then((r) => r.data);
}

async function setTargetTemperature(targetTemp) {
  return axios
    .post(
      `${config.endpoint}/backend/settings/dew_controller`,
      {
        targetTemp,
      },
      {
        transformRequest: [
          function (data) {
            const res = {
              targetTemp: parseFloat(data.targetTemp),
            };
            return JSON.stringify(res);
          },
        ],
      },
    )
    .then((r) => r.data);
}

async function setDutyCycle(dutyCycle) {
  return axios
    .post(
      `${config.endpoint}/backend/settings/dew_controller/duty_cycle`,
      {
        dutyCycle,
      },
      {
        transformRequest: [
          function (data) {
            const res = {
              dutyCycle: parseInt(data.dutyCycle, 10),
            };
            return JSON.stringify(res);
          },
        ],
      },
    )
    .then((r) => {
      return r.data.enabled;
    });
}

async function setPID(p, i, d) {
  return axios
    .post(
      `${config.endpoint}/backend/settings/pid`,
      {
        p,
        i,
        d,
      },
      {
        transformRequest: [
          function (data) {
            const res = {
              p: parseFloat(data.p),
              i: parseFloat(data.i),
              d: parseFloat(data.d),
            };
            return JSON.stringify(res);
          },
        ],
      },
    )
    .then((r) => r.data);
}

export {
  getAllSettings,
  getAlignStatus,
  getAPSettings,
  getFlags,
  getLocationSettings,
  getTrackState,
  setAPSettings,
  setLocationSettings,
  getIntervalometerSettings,
  setIntervalometerSettings,
  getDewControllerStatus,
  setDutyCycle,
  setTargetTemperature,
  setPID,
};
