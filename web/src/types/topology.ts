/**
 * Girder API のレスポンス実データ形状に準拠した型定義。
 *
 * 方針 (Context.md 原則1,3,4 に対応):
 * - Girder独自のドメインモデルを作らない。ここにあるのは各Infrastructure
 *   (libvirt / OVN) のリソースを「GET /api/topology」がまとめて返す形そのもの。
 * - status/pending 等のUI専用フィールドは別レイヤー(store)で持たせ、
 *   ここには混ぜない。これはGirder Core側の永続状態ではなくReact側の
 *   楽観的更新のためのローカルUI状態であるため。
 */

export type ResourceKind = "vm" | "switch" | "router" | "gateway";

/** libvirt Domain (VM) */
export interface VMResource {
  kind: "vm";
  /** libvirt domain UUID */
  id: string;
  name: string;
  state: "running" | "shutoff" | "paused" | "unknown";
  /** このVMが持つNIC一覧 (libvirtxml.DomainInterface 由来) */
  nics: NIC[];
}

/** libvirtxml.DomainInterface を最小限で表現。Girder独自のNIC型は作らない。 */
export interface NIC {
  /** libvirt <alias> element に保存された論理名 */
  alias: string;
  /** TAPデバイス名 (crypto/rand 4byte hex から生成) */
  target: string;
  macAddress: string;
  /** 接続先の OVN Logical Switch Port 名。未接続なら null */
  connectedSwitchPort: string | null;
}

/** OVN Logical Switch */
export interface SwitchResource {
  kind: "switch";
  /** OVN Logical Switch UUID */
  id: string;
  name: string;
  /** このSwitchに存在するPort名一覧 (NIC側のalias、またはRouter接続用ポート名) */
  ports: string[];
}

/** OVN Logical Router */
export interface RouterResource {
  kind: "router";
  /** OVN Logical Router UUID */
  id: string;
  name: string;
  ports: RouterPort[];
}

export interface RouterPort {
  /** Router Port名 = 接続先SwitchName (命名規則) */
  name: string;
  /** 割り当て済みIP (CIDR)。IP割当Issueで設定されるまではnull */
  networks: string[] | null;
}

/** Gateway (物理NIC -> OVS Bridge -> External Logical Switch の固定操作結果) */
export interface GatewayResource {
  kind: "gateway";
  id: string;
  name: string;
  /** ユーザが割り当てた物理NIC名 */
  physicalIface: string;
  /** ConnectGatewayToRouter で接続されたRouter名。未接続ならnull */
  connectedRouter: string | null;
}

export type TopologyResource =
  | VMResource
  | SwitchResource
  | RouterResource
  | GatewayResource;

/** 関係操作の種別。Edge作成時にどのAPIを叩くかを決める */
export type ConnectionKind =
  | "nic-switch" // ConnectNICtoSwitch
  | "switch-router" // ConnectSwitchToRouter
  | "gateway-router"; // ConnectGatewayToRouter

/** ノード間の関係 (Edge)。関係操作の結果としてのみ存在する */
export interface TopologyConnection {
  id: string;
  kind: ConnectionKind;
  /** 送信元ノードID (NIC の場合は "vmId:nicAlias" 形式) */
  source: string;
  target: string;
}

/** GET /api/topology のレスポンス全体。Girderのsource of truthスナップショット */
export interface TopologySnapshot {
  resources: TopologyResource[];
  connections: TopologyConnection[];
}
