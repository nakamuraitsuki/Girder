import { Handle, Position, type NodeProps } from "@xyflow/react";
import { NodeShell } from "./NodeShell";
import type { GatewayResource } from "../types/topology";

export type GatewayNodeData = { resource: GatewayResource };

export function GatewayNode({
  data,
  selected,
}: NodeProps & { data: GatewayNodeData }) {
  const gw = data.resource;
  return (
    <NodeShell kind="gateway" title={gw.name} subtitle={gw.id} selected={selected}>
      <Handle type="source" position={Position.Right} id="out" />
      <div className="gateway-iface">物理NIC: {gw.physicalIface}</div>
    </NodeShell>
  );
}
