import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import {
  ReactFlow,
  Background,
  Controls,
  MiniMap,
  useNodesState,
  useEdgesState,
  addEdge,
  BackgroundVariant,
  type Connection,
  type Edge,
  type Node,
  type EdgeMouseHandler,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";

import { useTopologyStore } from "./store/topologyStore";
import type { ConnectionKind, TopologyResource } from "./types/topology";
import { VMNode } from "./nodes/VMNode";
import { SwitchNode } from "./nodes/SwitchNode";
import { RouterNode } from "./nodes/RouterNode";
import { GatewayNode } from "./nodes/GatewayNode";
import { RouterPortForm } from "./forms/RouterPortForm";

const nodeTypes = {
  vm: VMNode,
  switch: SwitchNode,
  router: RouterNode,
  gateway: GatewayNode,
};

const POSITIONS_STORAGE_KEY = "girder-web:node-positions";

function loadStoredPositions(): Record<string, { x: number; y: number }> {
  try {
    const raw = localStorage.getItem(POSITIONS_STORAGE_KEY);
    return raw ? JSON.parse(raw) : {};
  } catch {
    return {};
  }
}

function saveStoredPositions(positions: Record<string, { x: number; y: number }>) {
  localStorage.setItem(POSITIONS_STORAGE_KEY, JSON.stringify(positions));
}

/** kind別に列を分け、簡易的な自動レイアウトを行う (座標はGirderのDriver状態には含まれないため、UI側でのみ保持) */
function layoutResources(
  resources: TopologyResource[],
  stored: Record<string, { x: number; y: number }>,
): Record<string, { x: number; y: number }> {
  const columnX: Record<TopologyResource["kind"], number> = {
    gateway: 0,
    vm: 0,
    switch: 340,
    router: 680,
  };
  const columnCursor: Record<string, number> = {};
  const result: Record<string, { x: number; y: number }> = {};

  // gatewayを最上段、VMをその下に積む都合上、種別ごとにy開始位置をずらす
  const rowStartY: Record<TopologyResource["kind"], number> = {
    gateway: 0,
    vm: 140,
    switch: 40,
    router: 40,
  };

  for (const r of resources) {
    if (stored[r.id]) {
      result[r.id] = stored[r.id];
      continue;
    }
    const x = columnX[r.kind];
    const key = `${r.kind}`;
    const cursor = columnCursor[key] ?? rowStartY[r.kind];
    const height = r.kind === "vm" ? 140 : 120;
    result[r.id] = { x, y: cursor };
    columnCursor[key] = cursor + height;
  }
  return result;
}

function connectionKindFor(
  sourceType: string,
  targetType: string,
): ConnectionKind | null {
  if (sourceType === "vm" && targetType === "switch") return "nic-switch";
  if (sourceType === "switch" && targetType === "router") return "switch-router";
  if (sourceType === "gateway" && targetType === "router") return "gateway-router";
  return null;
}

const STATUS_COLOR: Record<string, string> = {
  confirmed: "#3a4354",
  pending: "#f0a93c",
  error: "#ef5d5d",
};

export function Canvas() {
  const {
    snapshot,
    connectionEntries,
    loading,
    loadError,
    loadTopology,
    connect,
    disconnect,
  } = useTopologyStore();

  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);
  const [toast, setToast] = useState<string | null>(null);
  const [portForm, setPortForm] = useState<{
    routerId: string;
    portName: string;
    currentNetworks: string[] | null;
  } | null>(null);

  const positionsRef = useRef(loadStoredPositions());

  useEffect(() => {
    loadTopology();
  }, [loadTopology]);

  const handlePortClick = useCallback((routerId: string, portName: string) => {
    const router = snapshot?.resources.find(
      (r) => r.kind === "router" && r.id === routerId,
    );
    const port =
      router?.kind === "router"
        ? router.ports.find((p) => p.name === portName)
        : undefined;
    setPortForm({
      routerId,
      portName,
      currentNetworks: port?.networks ?? null,
    });
  }, [snapshot]);

  // snapshot -> React Flow nodes への変換。既存の位置(ドラッグ済み)は保持する。
  useEffect(() => {
    if (!snapshot) return;
    const positions = layoutResources(snapshot.resources, positionsRef.current);
    positionsRef.current = { ...positionsRef.current, ...positions };
    saveStoredPositions(positionsRef.current);

    setNodes(
      snapshot.resources.map((r): Node => ({
        id: r.id,
        type: r.kind,
        position: positions[r.id],
        data:
          r.kind === "router"
            ? { resource: r, onPortClick: handlePortClick }
            : { resource: r },
      })),
    );
  }, [snapshot, setNodes, handlePortClick]);

  // connectionEntries -> React Flow edges への変換 (status別に色分け)
  useEffect(() => {
    const nextEdges: Edge[] = Object.values(connectionEntries).map((entry) => {
      const { conn, status } = entry;
      let source = conn.source;
      let sourceHandle: string | undefined;
      if (conn.kind === "nic-switch") {
        const [vmId, alias] = conn.source.split(":");
        source = vmId;
        sourceHandle = alias;
      }
      return {
        id: conn.id,
        source,
        sourceHandle,
        target: conn.target,
        targetHandle: conn.kind === "nic-switch" ? "in" : "in",
        animated: status === "pending",
        style: {
          stroke: STATUS_COLOR[status],
          strokeDasharray: status === "pending" ? "4 3" : undefined,
        },
        label: status === "error" ? entry.errorMessage : undefined,
        labelStyle: { fill: "#ef5d5d", fontSize: 10 },
        data: { status },
      };
    });
    setEdges(nextEdges);
  }, [connectionEntries, setEdges]);

  const handleConnect = useCallback(
    (params: Connection) => {
      const sourceNode = nodes.find((n) => n.id === params.source);
      const targetNode = nodes.find((n) => n.id === params.target);
      if (!sourceNode || !targetNode) return;

      const kind = connectionKindFor(
        sourceNode.type as string,
        targetNode.type as string,
      );
      if (!kind) {
        setToast(
          `${sourceNode.type} → ${targetNode.type} の接続はサポートされていません`,
        );
        setTimeout(() => setToast(null), 3000);
        return;
      }

      const source =
        kind === "nic-switch"
          ? `${params.source}:${params.sourceHandle}`
          : params.source!;
      connect(kind, source, params.target!);
    },
    [nodes, connect],
  );

  const handleEdgeDoubleClick: EdgeMouseHandler = useCallback(
    (_, edge) => {
      if (edge.data?.status !== "confirmed") return;
      if (confirm("この接続を切断しますか?")) {
        disconnect(edge.id);
      }
    },
    [disconnect],
  );

  const onConnectHandler = useCallback(
    (params: Connection) => {
      handleConnect(params);
    },
    [handleConnect],
  );

  const onEdgesConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges],
  );

  const handleNodeDragStop = useCallback((_: unknown, node: Node) => {
    positionsRef.current = { ...positionsRef.current, [node.id]: node.position };
    saveStoredPositions(positionsRef.current);
  }, []);

  const stats = useMemo(() => {
    if (!snapshot) return null;
    const counts: Record<string, number> = {};
    for (const r of snapshot.resources) counts[r.kind] = (counts[r.kind] ?? 0) + 1;
    return counts;
  }, [snapshot]);

  return (
    <div style={{ width: "100%", height: "100%", position: "relative" }}>
      <div className="topbar">
        <div className="topbar__title">Girder — Topology</div>
        {stats && (
          <div className="topbar__stats">
            <span>VM {stats.vm ?? 0}</span>
            <span>Switch {stats.switch ?? 0}</span>
            <span>Router {stats.router ?? 0}</span>
            <span>Gateway {stats.gateway ?? 0}</span>
          </div>
        )}
        {loading && <span className="topbar__loading">読み込み中...</span>}
      </div>

      {loadError && <div className="banner banner--error">{loadError}</div>}

      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={(params) => {
          onEdgesConnect(params);
          onConnectHandler(params);
        }}
        onEdgeDoubleClick={handleEdgeDoubleClick}
        onNodeDragStop={handleNodeDragStop}
        nodeTypes={nodeTypes}
        fitView
        proOptions={{ hideAttribution: true }}
        defaultEdgeOptions={{ type: "smoothstep" }}
      >
        <Background variant={BackgroundVariant.Dots} gap={24} size={1} color="#1a2029" />
        <Controls showInteractive={false} />
        <MiniMap
          pannable
          zoomable
          style={{ background: "#12161d" }}
          maskColor="rgba(10,13,18,0.75)"
          nodeColor={(n) =>
            n.type === "vm"
              ? "#5b9dfb"
              : n.type === "switch"
                ? "#35c98b"
                : n.type === "router"
                  ? "#f0a93c"
                  : "#ef5d5d"
          }
        />
      </ReactFlow>

      {toast && <div className="toast">{toast}</div>}

      {portForm && (
        <RouterPortForm
          routerId={portForm.routerId}
          portName={portForm.portName}
          currentNetworks={portForm.currentNetworks}
          onClose={() => setPortForm(null)}
          onSaved={loadTopology}
        />
      )}
    </div>
  );
}
