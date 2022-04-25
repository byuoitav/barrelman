import { HttpClient, HttpErrorResponse, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, of, throwError } from 'rxjs';
import { JsonConvert, OperationMode, ValueCheckingMode} from "json2typescript";
import {tap, map, catchError} from "rxjs/operators"
import { MatDialog } from '@angular/material/dialog';

export class DeviceInfo {
  hostname: string | undefined;
  id: string | undefined;
  internetConnectivity: boolean | undefined;
  ip: string | undefined;
  dhcp: DHCPInfo| undefined;
}

export class DHCPInfo {
  error: any | undefined;
  enabled: boolean| undefined;
  toogleable: boolean| undefined;
}


@Injectable({
  providedIn: 'root'
})
export class ApiService {
  public theme = "default";

  private jsonConvert: JsonConvert;
  private urlParams: URLSearchParams;

  constructor(private http: HttpClient, private dialog: MatDialog) { 
    this.jsonConvert = new JsonConvert();
    this.jsonConvert.ignorePrimitiveChecks = false;

    this.urlParams = new URLSearchParams(window.location.search);
    if (!this.urlParams.has("theme")) {
     // this.theme = this.urlParams.get("theme");
    }
  }

  public switchToUI() {
    window.location.pathname = "/ui"
  }

  public refresh() {
    window.location.reload();
  }

  public switchTheme(name: string) {
    console.log("switching theme to", name);

    this.theme = name;
    this.urlParams.set("theme", name);
    window.history.replaceState(
      null,
      "System Health Dashboard",
      window.location.pathname + "?" + this.urlParams.toString()
    );
  }

  public async reboot() {
    /*try {
      this.dialog.open(RebootComponent, { disableClose: true });
      const data = await this.http
        .put("device/reboot", {
          responseType: "text"
        })
        .toPromise();
    } catch (e) {
      // bug where responseType doesn't actually work
      if (e.status === 200) {
        console.log(e.error.text);
        return e.error.text;
      }

      throw new Error("error getting rebooting device: " + e);
    }*/
  }

  getDeviceInfo(): Observable<DeviceInfo>{
    return this.http.get<DeviceInfo>("/api/v1/device")
      .pipe(
        tap(_ => console.log('fetched DeviceInfo')),
        catchError(this.handleError<DeviceInfo>('getDeviceInfo',))
      );
  }




  public async getMaintenanceMode() {
    return false;
    /*
    try {
      const data = await this.http.get("maintenance").toPromise();

      return (<any>data) as boolean;
    } catch (e) {
      throw new Error("error getting maintenance mode: " + e);
    }
     */
  }

  public async toggleMaintenanceMode() {
    return false;
    /*
    try {
      const data = await this.http.put("maintenance", null).toPromise();

      return (<any>data) as boolean;
    } catch (e) {
      throw new Error("error toggling maintenance mode: " + e);
    }
     */
  }

  public async getSoftwareStati() {
    /*try {
      const data: any = await this.http.get("device/status").toPromise();
      const stati = this.jsonConvert.deserializeObject(data, Status);

      return stati;
    } catch (e) {
      throw new Error("error getting software status': " + e);
    }*/
  }
  
  getDeviceID() {
    return this.http.get<string>("/api/v1/device/id")
      .pipe(
        tap(_ => console.log('fetched DeviceID')),
        catchError(this.handleError<string>('getDeviceID',))
      );
  }

  public async getRoomPing() {
    /*try {
      const data = await this.http.get("room/ping").toPromise();

      // build the map
      const result = new Map<string, PingResult>();
      for (const key of Object.keys(data)) {
        if (key && data[key]) {
          const val = this.jsonConvert.deserializeObject(data[key], PingResult);
          result.set(key, val);
        }
      }

      return result;
    } catch (e) {
      throw new Error("error getting room ping info: " + e);
    }*/
  }

  public async getRoomHealth() {
    /*try {
      const data = await this.http.get("room/health").toPromise();

      // build the map
      const result = new Map<string, string>();
      for (const key of Object.keys(data)) {
        if (key && data[key]) {
          result.set(key, data[key]);
        }
      }

      return result;
    } catch (e) {
      throw new Error("error getting room ping info: " + e);
    }*/
  }

  public async getRunnerInfo() {
    /*try {
      const data: any = await this.http.get("device/runners").toPromise();
      const info = this.jsonConvert.deserializeArray(data, RunnerInfo);

      return info;
    } catch (e) {
      throw new Error("error getting device runner info: " + e);
    }*/
  }

  public async getViaInfo() {
    /*try {
      const data: any = await this.http.get("room/viainfo").toPromise();
      const info = this.jsonConvert.deserializeArray(data, ViaInfo);

      return info;
    } catch (e) {
      throw new Error("error getting via info: " + e);
    }*/
  }

  public async resetVia(address: string) {
    try {
      const data = await this.http
        .get("http://" + location.hostname + ":8014/via/" + address + "/reset")
        .toPromise();

      console.log("data", data);
    } catch (e) {
      throw new Error("error resetting via: " + e);
    }
  }

  public async rebootVia(address: string) {
    try {
      const data = await this.http
        .get("http://" + location.hostname + ":8014/via/" + address + "/reboot")
        .toPromise();

      console.log("data", data);
    } catch (e) {
      throw new Error("error rebooting via: " + e);
    }
  }

  public async getDividerSensorsStatus(address: string) {
    try {
      const data = await this.http
        .get("http://" + address + ":10000/divider/state")
        .toPromise();

      console.log("data", data);

      for (const [key] of Object.entries(data)) {
        if (key.includes("disconnected")) {
          return false;
        }
        if (key.includes("connected")) {
          return true;
        }
      }
    } catch (e) {
      throw new Error("error getting divider sensors connection status: " + e);
    }
    return true;
  }

  public async getHardwareInfo() {
    try {
      const data = await this.http.get("/device/hardwareinfo").toPromise();

      console.log("data", data);

      return data;
    } catch (e) {
      throw new Error("error getting hardware info: " + e)
    }
  }

  public async flushDNS() {
    this.http.get("/dns").subscribe((data: any) => {
      if (data == "success") {
        console.log("successfully flushed the dns cache");
      } else {
        console.log("failed to flush the dns cache");
      }
    });
  }

  private handleError<T>(operation = 'operation', result?:T){
    return(error: any): Observable<T> => {
      console.error( `error doing ${operation}`, error);
      return of(result as T);
    }
  }

}
