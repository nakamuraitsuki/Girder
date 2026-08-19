import type { TopologySnapshot } from "../types/topology";

/**
 * 実際のGirder API (lsp-list等で確認済み) のデータ形状に準拠したfixture。
 * シナリオ: Context.md 14章の3層Webアプリ構成を模している。
 */
export const initialTopology: TopologySnapshot = {
  resources: [
    {
      kind: "vm",
      id: "vm-web-01",
      name: "web-server",
      state: "running",
      nics: [
        {
          alias: "web-eth0",
          target: "tap3f2a",
          macAddress: "52:54:00:12:34:01",
          connectedSwitchPort: "web-eth0",
        },
      ],
    },
    {
      kind: "vm",
      id: "vm-app-01",
      name: "app-server",
      state: "running",
      nics: [
        {
          alias: "app-eth0",
          target: "tap7c91",
          macAddress: "52:54:00:12:34:02",
          connectedSwitchPort: null,
        },
      ],
    },
    {
      kind: "vm",
      id: "vm-db-01",
      name: "db-server",
      state: "shutoff",
      nics: [
        {
          alias: "db-eth0",
          target: "tap1e88",
          macAddress: "52:54:00:12:34:03",
          connectedSwitchPort: null,
        },
      ],
    },
    {
      kind: "switch",
      id: "sw-ext-01",
      name: "sw-external",
      ports: ["web-eth0", "sw-external-router"],
    },
    {
      kind: "switch",
      id: "sw-int-01",
      name: "sw-internal",
      ports: ["sw-internal-router"],
    },
    {
      kind: "router",
      id: "lr-main-01",
      name: "lr-main",
      ports: [
        { name: "sw-external", networks: ["10.0.0.1/24"] },
        { name: "sw-internal", networks: null },
      ],
    },
    {
      kind: "gateway",
      id: "gw-01",
      name: "gw-uplink",
      physicalIface: "enp3s0",
      connectedRouter: "lr-main",
    },
  ],
  connections: [
    {
      id: "conn-1",
      kind: "nic-switch",
      source: "vm-web-01:web-eth0",
      target: "sw-ext-01",
    },
    {
      id: "conn-2",
      kind: "switch-router",
      source: "sw-ext-01",
      target: "lr-main-01",
    },
    {
      id: "conn-3",
      kind: "gateway-router",
      source: "gw-01",
      target: "lr-main-01",
    },
  ],
};
