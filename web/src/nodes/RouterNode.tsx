import { Handle, Position, type NodeProps } from "@xyflow/react";
import { NodeShell } from "./NodeShell";
import type { RouterResource } from "../types/topology";

export type RouterNodeData = {
  resource: RouterResource;
  connectedGatewayName: string | null;
  onPortClick: (routerId: string, portName: string) => void;
};

export function RouterNode({ data, selected }: NodeProps & { data: RouterNodeData }) {
  const router = data.resource;

  return (
    <NodeShell kind="router" title={router.name} subtitle={router.id} selected={selected}>
      {/* 左: Gateway = Default Static Route側 */}
      <div className="node-row-list">
        <div className="node-row" style={{ position: "relative" }}>
          <Handle
            type="target"
            position={Position.Left}
            id="uplink"
            style={{ position: "absolute", left: -11, top: "50%" }}
          />
          <span className="node-row__label">
            {data.connectedGatewayName ? `← ${data.connectedGatewayName}` : "(Gateway未接続)"}
          </span>
          <span className="node-row__meta">default route</span>
        </div>
      </div>

      {/* 右: 配下Switch = Network境界Port */}
      <div className="node-row-list" style={{ marginTop: 6 }}>
        {router.ports.map((port) => (
          <div
            className="node-row node-row--clickable"
            style={{ position: "relative" }}
            key={port.name}
            onClick={() => data.onPortClick(router.id, port.name)}
            title="クリックしてIPを設定"
          >
            <span className="node-row__label">{port.name}</span>
            <span className="node-row__meta">{port.networks?.join(", ") ?? "IP未設定"}</span>
            <Handle
              type="target"
              position={Position.Right}
              id={port.name}
              style={{ position: "absolute", right: -11, top: "50%" }}
            />
          </div>
        ))}
        <div className="node-row node-row--ghost" style={{ position: "relative" }}>
          <span className="node-row__meta">+ 新規Switch</span>
          <Handle
            type="target"
            position={Position.Right}
            id="new-switch"
            style={{ position: "absolute", right: -11, top: "50%" }}
          />
        </div>
      </div>
    </NodeShell>
  );
}