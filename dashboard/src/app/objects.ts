import { toTypeScript } from "@angular/compiler";

export interface DHCPInfo {
    error: any;
    enabled: boolean;
    toogleable: boolean;
}

export interface DeviceInfo {
    hostname: string;
    id: string;
    internetConnectivity: boolean;
    ip: string;
    dhcp: DHCPInfo;
}

export interface PingResult {
    error: string;
    ip: string;
    packetsSent: Number;
    packetsReceived: Number;
    packetsLost: Number;
    averageRoundTrip: string;
}

export interface Trigger {
    tType: string;
    at: string;
    every: string;
    match: any;
}

export interface RunnerInfo {
    id: string;
    trigger: Trigger;
    context: any;
    lastRunTime: Date;
    lastRunDuration: string;
    lastRunError: string;
    currentlyRunning: boolean;
    runCount: number;
}

export interface ViaInfo {
    name: string
    address: string;
}