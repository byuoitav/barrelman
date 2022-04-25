import { Component, OnInit } from '@angular/core';

import { ApiService } from 'src/app/services/api.service';
import { DeviceInfo, PingResult } from '../../objects';
import { map } from 'rxjs/operators';

@Component({
  selector: 'app-overview',
  templateUrl: './overview.component.html',
  styleUrls: ['./overview.component.scss']
})
export class OverviewComponent implements OnInit {

  public hasDividerSensors: boolean | void | undefined;
  deviceInfo : DeviceInfo | undefined;
  public pingResult:  Map<string, PingResult> | undefined | void;
  public dividerSensorStatus: string | undefined;
  public dividerSensorAddr: string | undefined;
  // public maintenanceMode: boolean;

  constructor(public api: ApiService) {}

  async ngOnInit() {
    this.getDeviceInfo();
   
    this.pingResult = await this.api.getRoomPing();
    console.log("ping result", this.pingResult);
    this.hasDividerSensors = await this.getDividerSensors();
    this.connected();
    setInterval(() => {
      this.connected();
    }, 2000);

    /*
    this.maintenanceMode = await this.api.getMaintenanceMode();
    console.log("maintenanceMode", this.maintenanceMode);
     */
  }

  getDeviceInfo(): void {
    this.api.getDeviceInfo().subscribe(deviceInfo => this.deviceInfo = deviceInfo);
    console.log("device info", this.deviceInfo);
  }

  public isDefined(test: any): boolean {
    return typeof test !== "undefined" && test !== null;
  }

  public async toggleMaintenanceMode() {
    /*
    console.log("toggling maintenance mode");

    this.maintenanceMode = await this.api.toggleMaintenanceMode();
    console.log("maintenanceMode", this.maintenanceMode);
     */
  }

  public reachable(): number {
    if (!this.pingResult) {
      return 0;
    }

    return Array.from(this.pingResult.values()).filter(r => r.packetsLost === 0)
    .length;
  }

  public unreachable(): number {
    if (!this.pingResult) {
      return 0;
    }

    return Array.from(this.pingResult.values()).filter(r => r.packetsLost > 0)
      .length;
    
  }

  public getDividerSensors() {
    if (!this.pingResult) {
      return console.error("no divider sensors");
    }

    for (const k of Array.from(this.pingResult.keys())) {
      if (k.includes("DS1")) {
        this.dividerSensorAddr = k + ".byu.edu";
        return true;
      }
    }
    return false;
  }

  public async connected() {
    if (
      (await this.api.getDividerSensorsStatus(this.dividerSensorAddr || "")) == true
    ) {
      this.dividerSensorStatus = "Connected";
    } else {
      this.dividerSensorStatus = "Disconnected";
    }
  }
 
}
