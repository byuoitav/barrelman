import { Component, OnInit } from '@angular/core';

import { ApiService } from 'src/app/services/api.service';

@Component({
  selector: 'app-via-control',
  templateUrl: './via-control.component.html',
  styleUrls: ['./via-control.component.scss']
})
export class ViaControlComponent implements OnInit {
  //public viainfo: ViaInfo[] = [];

  constructor(private api: ApiService) {}

  async ngOnInit() {
    //this.viainfo = await this.api.getViaInfo();
    console.log("via info");
  }

  reset() {
    //this.api.resetVia(via.address);
  }

  reboot() {
    //this.api.rebootVia(via.address);
  }

}
