import React, { Component } from "react";
import { BrowserRouter } from "react-router-dom";
import "fontsource-roboto";

import { MuiPickersUtilsProvider } from "@material-ui/pickers";
import LuxonUtils from "@date-io/luxon";

import DateTime from "luxon/src/datetime";

import Button from "@material-ui/core/Button";
import { DateTimePicker } from "@material-ui/pickers";

import {
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
  Legend,
} from "recharts";

import "./App.css";

class App extends Component {
  constructor() {
    super();
    this.state = {
      csv: null,
      records: null,
      hasData: null,
      startDate: null,
      desiredStartDate: null,
      endDate: null,
      desiredEndDate: null,
      filteredRecords: null,
    };
  }

  componentDidMount() {
    const qp = window.location.search;
    const logFile = qp.split("file=")[1];

    if (logFile === undefined) {
      this.setState({ hasData: false });
      return;
    }

    fetch(`/logs/${logFile}`)
      .then((r) => r.text())
      .then((v) => {
        this.setState({ csv: v });
        this.parseCSV();
      });
  }

  rowHasValidData(row) {
    if (row === undefined) {
      return false;
    }

    if (row.length === 0) {
      return false;
    }

    const splittedRow = row.split(",");

    return splittedRow.length === 4;
  }

  parseCSV() {
    const splittedLines = this.state.csv.split("\n");
    const maxLen = splittedLines.length - 1;
    if (maxLen < 0) {
      this.setState({ hasData: false });
      return;
    }

    let maxRecordNo;

    for (let i = maxLen; i > 0; i--) {
      const thisRow = splittedLines[i];
      if (this.rowHasValidData(thisRow)) {
        maxRecordNo = i;
        break;
      }
    }

    if (maxRecordNo === undefined || maxRecordNo < 0) {
      console.log("No valid records!");
      this.setState({ hasData: false });
      return;
    }

    let records = [];
    let startDate, endDate;

    for (let i = 0; i < maxRecordNo; i++) {
      const splittedRecord = splittedLines[i].split(",");
      if (i === 0) {
        startDate = DateTime.fromISO(splittedRecord[0]);
      }

      if (i + 1 === maxRecordNo) {
        endDate = DateTime.fromISO(splittedRecord[0]);
      }

      const newRecord = {
        dateTime: DateTime.fromISO(splittedRecord[0]),
        setPoint: parseFloat(splittedRecord[1]),
        processValue: parseFloat(splittedRecord[2]),
      };
      records.push(newRecord);
    }

    this.setState({
      hasData: true,
      records: records,
      startDate: startDate,
      desiredStartDate: startDate,
      endDate: endDate,
      desiredEndDate: endDate,
    });
    this.filterRecords();
  }

  filterRecords() {
    let res = [];
    const { desiredStartDate, desiredEndDate } = this.state;

    this.state.records.forEach((e) => {
      if (e.dateTime >= desiredStartDate && e.dateTime <= desiredEndDate) {
        res.push(e);
      }
    });

    this.setState({ filteredRecords: res });
  }

  onSetStartDate = (dateTime) => {
    if (
      dateTime >= this.state.startDate &&
      dateTime <= this.state.desiredEndDate
    ) {
      this.setState({ desiredStartDate: dateTime });
      this.filterRecords();
    }
  };

  onSetEndDate = (dateTime) => {
    if (
      dateTime <= this.state.endDate &&
      dateTime >= this.state.desiredStartDate
    ) {
      this.setState({ desiredEndDate: dateTime });
      this.filterRecords();
    }
  };

  onResetClick() {
    const { startDate, endDate } = this.state;

    this.setState({
      desiredStartDate: startDate,
      desiredEndDate: endDate,
    });

    this.filterRecords();
  }

  datePickers() {
    if (this.state.hasData === true) {
      return (
        <div>
          <DateTimePicker
            label="Start"
            minDate={this.state.startDate.toJSDate()}
            maxDate={this.state.endDate.toJSDate()}
            value={this.state.desiredStartDate.toJSDate()}
            onChange={this.onSetStartDate.bind(this)}
          />
          <DateTimePicker
            label="End"
            minDate={this.state.startDate.toJSDate()}
            maxDate={this.state.endDate.toJSDate()}
            value={this.state.desiredEndDate.toJSDate()}
            onChange={this.onSetEndDate.bind(this)}
          />
          <br />
          <Button
            variant="contained"
            color="secondary"
            onClick={this.onResetClick.bind(this)}
          >
            Reset
          </Button>
        </div>
      );
    }
  }

  graph() {
    if (this.state.hasData === true) {
      return (
        <div>
          <ResponsiveContainer width="100%" height={400}>
            <LineChart
              width={1000}
              height={600}
              data={this.state.filteredRecords}
              margin={{
                top: 5,
                right: 30,
                left: 20,
                bottom: 5,
              }}
            >
              <XAxis dataKey="dateTime" />
              <CartesianGrid strokeDasharray="3 3" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Line type="monotone" dataKey="processValue" stroke="#8884d8" />
              <Line type="monotone" dataKey="setPoint" stroke="#82ca9d" />
            </LineChart>
          </ResponsiveContainer>
        </div>
      );
    }
  }

  noData() {
    if (this.state.hasData !== true) {
      return <div>No data!</div>;
    }
  }

  render() {
    return (
      <BrowserRouter basename={process.env.PUBLIC_URL}>
        <MuiPickersUtilsProvider utils={LuxonUtils}>
          <div className="App">
            {this.datePickers()}
            {this.graph()}
            {this.noData()}
          </div>
        </MuiPickersUtilsProvider>
      </BrowserRouter>
    );
  }
}

export default App;
