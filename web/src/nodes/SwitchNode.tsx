import { Handle, Position, type NodeProps } from "@xyflow/react";
import { NodeShell } from "./NodeShell";
import type { SwitchResource } from "../types/topology";

export type SwitchNodeData = { resource: SwitchResource };

export function SwitchNode({
  data,
  selected,
}: NodeProps & { data: SwitchNodeData }) {
  const sw = data.resource;
  return (
    <NodeShell kind="switch" title={sw.name} subtitle={sw.id} selected={selected}>
      <Handle type="target" position={Position.Left} id="in" />
      <Handle type="source" position={Position.Right} id="out" />
      <div className="node-row-list">
        {sw.ports.length === 0 && (
          <div className="node-row">
            <span className="node-row__meta">ポート未接続</span>
          </div>
        )}
        {sw.ports.map((port) => (
          <div className="node-row" key={port}>
            <span className="node-row__label">{port}</span>
          </div>
        ))}
      </div>
    </NodeShell>
  );
}
