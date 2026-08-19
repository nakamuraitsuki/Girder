import { Handle, Position, type NodeProps } from "@xyflow/react";
import { NodeShell } from "./NodeShell";
import type { RouterResource } from "../types/topology";

export type RouterNodeData = {
  resource: RouterResource;
  onPortClick: (routerId: string, portName: string) => void;
};

export function RouterNode({
  data,
  selected,
}: NodeProps & { data: RouterNodeData }) {
  const router = data.resource;
  return (
    <NodeShell kind="router" title={router.name} subtitle={router.id} selected={selected}>
      <Handle type="target" position={Position.Left} id="in" />
      <div className="node-row-list">
        {router.ports.map((port) => (
          <div
            className="node-row node-row--clickable"
            key={port.name}
            onClick={() => data.onPortClick(router.id, port.name)}
            title="クリックしてIPを設定"
          >
            <span className="node-row__label">{port.name}</span>
            <span className="node-row__meta">
              {port.networks?.join(", ") ?? "IP未設定"}
            </span>
          </div>
        ))}
      </div>
    </NodeShell>
  );
}
