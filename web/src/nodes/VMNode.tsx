import { Handle, Position, type NodeProps } from "@xyflow/react";
import { NodeShell } from "./NodeShell";
import type { VMResource } from "../types/topology";

export type VMNodeData = { resource: VMResource };

export function VMNode({ data, selected }: NodeProps & { data: VMNodeData }) {
  const vm = data.resource;
  return (
    <NodeShell
      kind="vm"
      title={vm.name}
      subtitle={vm.id}
      badge={vm.state}
      badgeVariant={
        vm.state === "running"
          ? "running"
          : vm.state === "shutoff"
            ? "shutoff"
            : "neutral"
      }
      selected={selected}
    >
      <div className="node-row-list">
        {vm.nics.map((nic) => (
          <div className="node-row" key={nic.alias}>
            <span
              className={`node-row__dot ${
                nic.connectedSwitchPort
                  ? "node-row__dot--connected"
                  : "node-row__dot--disconnected"
              }`}
            />
            <span className="node-row__label" style={{ flex: 1 }}>
              {nic.alias}
            </span>
            <span className="node-row__meta">{nic.macAddress}</span>
            <Handle
              type="source"
              position={Position.Right}
              id={nic.alias}
              style={{ position: "absolute", right: -11, top: "50%" }}
            />
          </div>
        ))}
      </div>
    </NodeShell>
  );
}
