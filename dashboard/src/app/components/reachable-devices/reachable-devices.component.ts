import { Component, OnInit } from '@angular/core';

import { ApiService } from 'src/app/services/api.service';


@Component({
  selector: 'app-reachable-devices',
  templateUrl: './reachable-devices.component.html',
  styleUrls: ['./reachable-devices.component.scss']
})
export class ReachableDevicesComponent implements OnInit {
  //public pingResult: Map<string, PingResult>;
  //public roomHealth: Map<string, string>;

  constructor(private api: ApiService) {}

  async ngOnInit() {
    //this.pingResult = await this.api.getRoomPing();
    //console.log("ping result:", this.pingResult);

    //this.roomHealth = await this.api.getRoomHealth();
    //console.log("room health:", this.roomHealth);
  }

}
