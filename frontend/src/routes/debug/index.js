import { h, Component } from "preact";
import style from "./style";
import { getAllSettings } from "../../lib/settings";

export default class Debug extends Component {
  state = {
    appContext: null
  };

  async componentDidMount() {
    getAllSettings().then(r => {
      this.setState({ appContext: r });
    });
  }

  appContext() {
    if (this.state.appContext === null) {
      return <p>LOADING</p>;
    }
    let jsonStr = JSON.stringify(this.state.appContext, null, 2);

    jsonStr = jsonStr.replace(/\n/g, "<br/>");
    jsonStr = jsonStr.replace(/ /g, "&nbsp");

    return <div dangerouslySetInnerHTML={{ __html: jsonStr }} />;
  }

  render({}, {}) {
    return (
      <div class={style.debug}>
        <h1>Debug</h1>
        <a href="/config.json">View/Download config.json</a>
        <p>NODE_ENV: {JSON.stringify(process.env.NODE_ENV)}</p>
        <h2>Backend State:</h2>
        {this.appContext()}
      </div>
    );
  }
}
