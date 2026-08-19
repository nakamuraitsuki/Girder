import { http, HttpResponse, delay } from "msw";
import { initialTopology } from "./fixtures";
import type {
  TopologySnapshot,
  TopologyConnection,
  ConnectionKind,
} from "../types/topology";

// mswハンドラ内のインメモリ状態。実サーバーでは各Infrastructure(libvirt/OVN)が
// source of truth だが、モックではこのオブジェクトがその代わりを務める。
let db: TopologySnapshot = structuredClone(initialTopology);

let nextConnId = db.connections.length + 1;

/** kind ごとに接続元/先が既に他のノードと繋がっていないかを軽く検証する */
function validateConnection(
  kind: ConnectionKind,
  source: string,
  _target: string,
): string | null {
  if (kind === "nic-switch") {
    const [vmId, alias] = source.split(":");
    const vm = db.resources.find(
      (r) => r.kind === "vm" && r.id === vmId,
    );
    if (vm?.kind === "vm") {
      const nic = vm.nics.find((n) => n.alias === alias);
      if (nic?.connectedSwitchPort) {
        return `NIC "${alias}" は既に Switch に接続されています`;
      }
    }
  }
  if (kind === "switch-router") {
    const already = db.connections.find(
      (c) => c.kind === "switch-router" && c.source === source,
    );
    if (already) return "この Switch は既に Router に接続されています";
  }
  return null;
}

function applyConnection(conn: TopologyConnection) {
  if (conn.kind === "nic-switch") {
    const [vmId, alias] = conn.source.split(":");
    const vm = db.resources.find((r) => r.kind === "vm" && r.id === vmId);
    const sw = db.resources.find(
      (r) => r.kind === "switch" && r.id === conn.target,
    );
    if (vm?.kind === "vm" && sw?.kind === "switch") {
      const nic = vm.nics.find((n) => n.alias === alias);
      if (nic) nic.connectedSwitchPort = alias;
      if (!sw.ports.includes(alias)) sw.ports.push(alias);
    }
  }
  if (conn.kind === "gateway-router") {
    const gw = db.resources.find(
      (r) => r.kind === "gateway" && r.id === conn.source,
    );
    const router = db.resources.find(
      (r) => r.kind === "router" && r.id === conn.target,
    );
    if (gw?.kind === "gateway" && router?.kind === "router") {
      gw.connectedRouter = router.name;
    }
  }
  db.connections.push(conn);
}

function removeConnection(id: string) {
  const conn = db.connections.find((c) => c.id === id);
  if (!conn) return;
  if (conn.kind === "nic-switch") {
    const [vmId, alias] = conn.source.split(":");
    const vm = db.resources.find((r) => r.kind === "vm" && r.id === vmId);
    const sw = db.resources.find(
      (r) => r.kind === "switch" && r.id === conn.target,
    );
    if (vm?.kind === "vm") {
      const nic = vm.nics.find((n) => n.alias === alias);
      if (nic) nic.connectedSwitchPort = null;
    }
    if (sw?.kind === "switch") {
      sw.ports = sw.ports.filter((p) => p !== alias);
    }
  }
  if (conn.kind === "gateway-router") {
    const gw = db.resources.find(
      (r) => r.kind === "gateway" && r.id === conn.source,
    );
    if (gw?.kind === "gateway") gw.connectedRouter = null;
  }
  db.connections = db.connections.filter((c) => c.id !== id);
}

export const handlers = [
  // 関係操作: 現在のトポロジー全体を返す (source of truth のスナップショット)
  http.get("/api/topology", async () => {
    await delay(300);
    return HttpResponse.json(db);
  }),

  // 関係操作: Connect系APIをまとめて1エンドポイントに集約
  // (ConnectNICtoSwitch / ConnectSwitchToRouter / ConnectGatewayToRouter)
  http.post("/api/topology/connections", async ({ request }) => {
    await delay(500);
    const body = (await request.json()) as {
      kind: ConnectionKind;
      source: string;
      target: string;
    };

    const validationError = validateConnection(
      body.kind,
      body.source,
      body.target,
    );
    if (validationError) {
      return HttpResponse.json({ error: validationError }, { status: 409 });
    }

    const conn: TopologyConnection = {
      id: `conn-${nextConnId++}`,
      kind: body.kind,
      source: body.source,
      target: body.target,
    };
    applyConnection(conn);
    return HttpResponse.json(conn, { status: 201 });
  }),

  http.delete("/api/topology/connections/:id", async ({ params }) => {
    await delay(400);
    removeConnection(params.id as string);
    return new HttpResponse(null, { status: 204 });
  }),

  // 固有操作: Router Port への IP割り当て (IP/ACL Issue相当)
  http.patch(
    "/api/ovn/routers/:routerId/ports/:portName",
    async ({ params, request }) => {
      await delay(400);
      const body = (await request.json()) as { networks: string[] };
      const router = db.resources.find(
        (r) => r.kind === "router" && r.id === params.routerId,
      );
      if (router?.kind !== "router") {
        return HttpResponse.json(
          { error: "Router not found" },
          { status: 404 },
        );
      }
      const port = router.ports.find((p) => p.name === params.portName);
      if (!port) {
        return HttpResponse.json(
          { error: "Router port not found" },
          { status: 404 },
        );
      }
      port.networks = body.networks;
      return HttpResponse.json(port);
    },
  ),

  // デバッグ用: fixtureへリセット
  http.post("/api/topology/_reset", async () => {
    db = structuredClone(initialTopology);
    nextConnId = db.connections.length + 1;
    return new HttpResponse(null, { status: 204 });
  }),
];
