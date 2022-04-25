import { toTypeScript } from "@angular/compiler";

export interface DHCPInfo {
    error: any | undefined;
    enabled: boolean | undefined;
    toogleable: boolean | undefined;
}

export interface DeviceInfo {
    hostname: string | undefined;
    id: string | undefined;
    internetConnectivity: boolean | undefined;
    ip: string | undefined;
    dhcp: DHCPInfo | undefined;
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