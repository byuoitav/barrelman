import { Component, OnInit} from '@angular/core';

import { ApiService } from "../../services/api.service";

@Component({
  selector: 'dashboard',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'dashboard';

  //public id: string;
  //private urlParams: URLSearchParams;

  constructor(public api: ApiService) {}

  public async ngOnInit() {
   // this.urlParams = new URLSearchParams(window.location.search);
    //this.id = await this.api.getDeviceID();
  }

  public wideSidenav(): boolean {
    //if (typeof this.urlParams === "undefined" || this.urlParams === null) {
      //return false;
    //}

    //if (this.urlParams.has("wide-sidenav")) {
     // return this.urlParams.get("wide-sidenav").toLowerCase() === "true";//}

    return true;
  }

  public setWideSidenav(val: boolean) {
    //this.urlParams.set("wide-sidenav", val.toString());

    window.history.replaceState(
      null,
      "System Health Dashboard",
      window.location.pathname + "?" //+ //this.urlParams.toString()
    );
  }

}
