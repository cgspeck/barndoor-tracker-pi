import { h, Component } from "preact";
import { route } from "preact-router";
import TopAppBar from "preact-material-components/TopAppBar";
import Drawer from "preact-material-components/Drawer";
import "preact-material-components/Drawer/style.css";
import "preact-material-components/Dialog/style.css";
import "preact-material-components/List/style.css";
import "preact-material-components/TopAppBar/style.css";

export default class Header extends Component {
  closeDrawer() {
    this.drawer.MDComponent.open = false;
  }

  openDrawer = () => (this.drawer.MDComponent.open = true);

  drawerRef = drawer => (this.drawer = drawer);

  linkTo = path => () => {
    route(path);
    this.closeDrawer();
  };

  goAlign = this.linkTo("/align");
  goAPSettings = this.linkTo("/ap_settings");
  goLocationSettings = this.linkTo("/location_settings");
  goIntervalometerSettings = this.linkTo("/intervalometer_settings");
  goDebug = this.linkTo("/debug");
  goTrack = this.linkTo("/");
  goDewController = this.linkTo("/dew_controller");


  render({ selectedRoute }, {}) {
    return (
      <div>
        <TopAppBar className="topappbar">
          <TopAppBar.Row>
            <TopAppBar.Section align-start>
              <TopAppBar.Icon menu onClick={this.openDrawer}>
                menu
              </TopAppBar.Icon>
              <TopAppBar.Title>Barndoor Tracker</TopAppBar.Title>
            </TopAppBar.Section>
          </TopAppBar.Row>
        </TopAppBar>

        <Drawer modal ref={this.drawerRef}>
          <Drawer.DrawerContent>
            <Drawer.DrawerItem
              selected={selectedRoute == "/align"}
              onClick={this.goAlign}
            >
              Align
            </Drawer.DrawerItem>

            <Drawer.DrawerItem
              selected={selectedRoute == "/"}
              onClick={this.goTrack}
            >
              {/*<List.ItemGraphic>home</List.ItemGraphic> */}
              Track
            </Drawer.DrawerItem>

            <Drawer.DrawerItem
              selected={selectedRoute == "/intervalometer_settings"}
              onClick={this.goIntervalometerSettings}
            >
              Intervalometer Settings
            </Drawer.DrawerItem>

            <Drawer.DrawerItem
              selected={selectedRoute == "/dew_controller"}
              onClick={this.goDewController}
            >
              Dew Controller
            </Drawer.DrawerItem>

            <Drawer.DrawerItem
              selected={selectedRoute == "/location_settings"}
              onClick={this.goLocationSettings}
            >
              Location Settings
            </Drawer.DrawerItem>

            <Drawer.DrawerItem
              selected={selectedRoute == "/ap_settings"}
              onClick={this.goAPSettings}
            >
              Access Point Settings
            </Drawer.DrawerItem>

            <Drawer.DrawerItem
              selected={selectedRoute == "/debug"}
              onClick={this.goDebug}
            >
              Debug
            </Drawer.DrawerItem>
          </Drawer.DrawerContent>
        </Drawer>
      </div>
    );
  }
}
