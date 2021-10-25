import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { RouterModule, Routes } from "@angular/router";
import { HttpClientModule } from "@angular/common/http";

import { MatSidenavModule } from "@angular/material/sidenav"; 
import { MatButtonModule } from "@angular/material/button";
import { MatToolbar, MatToolbarModule } from "@angular/material/toolbar";
import { MatCardModule } from "@angular/material/card";
import { MatDividerModule } from "@angular/material/divider";
import { MatListModule } from "@angular/material/list";
import { MatExpansionModule } from "@angular/material/expansion"
import { MatIconModule } from "@angular/material/icon";
import { MatProgressBarModule } from "@angular/material/progress-bar";
import { MatDialogModule } from "@angular/material/dialog";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import "hammerjs"; 

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './components/app/app.component';
import { OverviewComponent } from './components/overview/overview.component';
import { ReachableDevicesComponent } from './components/reachable-devices/reachable-devices.component';
import { ViaControlComponent } from './components/via-control/via-control.component';
import { APP_BASE_HREF } from '@angular/common';

const routes: Routes = [
  {
    path: "",
    redirectTo: "/overview",
    pathMatch: "full"
  },
  {
    path: "overview",
    component: OverviewComponent
  },
  {
    path: "reachable-devices",
    component: ReachableDevicesComponent
  },
  {
    path: "via-control",
    component: ViaControlComponent
  }
]

@NgModule({
  declarations: [
    AppComponent,
    OverviewComponent,
    ReachableDevicesComponent,
    ViaControlComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    MatSidenavModule,
    MatButtonModule,
    MatToolbarModule,
    MatCardModule,
    MatDividerModule,
    MatListModule,
    MatExpansionModule,
    MatIconModule,
    MatProgressBarModule,
    MatDialogModule,
    MatProgressSpinnerModule,
    AppRoutingModule,
    RouterModule.forRoot(routes)
  ],
  providers: [
    
    {
    provide: APP_BASE_HREF,
    useValue: "/dashboard"
    }

  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
