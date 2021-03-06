import MockAdapter from 'axios-mock-adapter';
import {
  getAllSettings,
  getAlignStatus,
  getAPSettings,
  getFlags,
  getLocationSettings,
  setAPSettings,
  setLocationSettings,
  setPID,
} from '../../src/lib/settings';
import Axios from 'axios';

import config from '../../src/config';

const mock = new MockAdapter(Axios);

afterEach(() => {
  mock.reset();
});

describe('getAllSettings', () => {
  beforeEach(() => {
    mock.onGet(`${config.endpoint}/backend/status/debug`).reply(200, {
      apSettings: {
        ssid: 'foo',
        key: 'bar',
      },
    });
  });

  it('returns expected value', async () => {
    const res = await getAllSettings();
    expect(res).toEqual({
      apSettings: {
        ssid: 'foo',
        key: 'bar',
      },
    });
  });
});

describe('getAPSettings', () => {
  beforeEach(() => {
    mock.onGet(`${config.endpoint}/backend/settings/network/ap`).reply(200, {
      ssid: 'foo',
      key: 'bar',
    });
  });

  it('returns expected value', async () => {
    const res = await getAPSettings();
    expect(res).toEqual({
      ssid: 'foo',
      key: 'bar',
    });
  });
});

describe('setAPSettings', () => {
  beforeEach(() => {
    mock.onPost(`${config.endpoint}/backend/settings/network`).reply(200, {
      accessPointMode: true,
    });
    mock.onPost(`${config.endpoint}/backend/settings/network/ap`).reply(200, {
      ssid: 'foo2',
      key: 'bar2',
    });
  });

  it('returns posted value', async () => {
    const res = await setAPSettings('foo2', 'bar2');
    expect(res).toEqual({
      ssid: 'foo2',
      key: 'bar2',
    });
  });
});

describe('getLocationSettings', () => {
  beforeEach(() => {
    mock.onGet(`${config.endpoint}/backend/settings/location`).reply(200, {
      latitude: -37.74,
      magDeclination: 11.64,
      azError: 1.3,
      altError: 2.4,
      xOffset: 0,
      yOffset: 0,
      zOffset: 0,
    });
  });

  it('returns expected value', async () => {
    const res = await getLocationSettings();
    expect(res).toEqual({
      latitude: -37.74,
      magDeclination: 11.64,
      azError: 1.3,
      altError: 2.4,
      xOffset: 0,
      yOffset: 0,
      zOffset: 0,
    });
  });
});

describe('setLocationSettings', () => {
  beforeEach(() => {
    mock
      .onPost(`${config.endpoint}/backend/settings/location`, {
        latitude: -37.74,
        magDeclination: 11.64,
        azError: 4,
        altError: 5,
        xOffset: 1,
        yOffset: 2,
        zOffset: 3,
      })
      .reply(200, {
        latitude: -37.74,
        magDeclination: 11.64,
        azError: 4,
        altError: 5,
        xOffset: 1,
        yOffset: 2,
        zOffset: 3,
      });
  });

  it('returns expected value', async () => {
    const res = await setLocationSettings(
      '-37.74',
      '11.64',
      '4',
      '5',
      '1',
      '2',
      '3',
    );
    expect(res).toEqual({
      latitude: -37.74,
      magDeclination: 11.64,
      azError: 4,
      altError: 5,
      xOffset: 1,
      yOffset: 2,
      zOffset: 3,
    });
  });
});

describe('getFlags', () => {
  beforeEach(() => {
    mock.onGet(`${config.endpoint}/backend/status/flags`).reply(200, {
      needsAPSettings: false,
      needsLocationSettings: true,
    });
  });

  it('returns expected value', async () => {
    const res = await getFlags();
    expect(res).toEqual({
      needsAPSettings: false,
      needsLocationSettings: true,
    });
  });
});

describe('getAlignStatus', () => {
  beforeEach(() => {
    mock.onGet(`${config.endpoint}/backend/status/align`).reply(200, {
      azAligned: false,
      altAligned: true,
      currentAz: 110.3,
      currentAlt: 37.3,
    });
  });

  it('returns expected value', async () => {
    const res = await getAlignStatus();
    expect(res).toEqual({
      azAligned: false,
      altAligned: true,
      currentAz: 110.3,
      currentAlt: 37.3,
    });
  });
});

describe('setPID', () => {
  beforeEach(() => {
    mock
      .onPost(`${config.endpoint}/backend/settings/pid`, {
        p: 1.23,
        i: 4.56,
        d: 7.89,
      })
      .reply(200, {
        p: 1.23,
        i: 4.56,
        d: 7.89,
      });
  });

  it('returns expected value', async () => {
    const res = await setPID('1.23', '4.56', '7.89');
    expect(res).toEqual({
      p: 1.23,
      i: 4.56,
      d: 7.89,
    });
  });
});
