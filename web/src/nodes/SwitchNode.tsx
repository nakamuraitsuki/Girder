import { Handle, Position, type NodeProps } from "@xyflow/react";
import { NodeShell } from "./NodeShell";
import type { SwitchResource } from "../types/topology";

export type SwitchNodeData = {
  resource: SwitchResource;
  /** このSwitchが接続されうるRouter名の集合。ports内のどれがRouter向けポートかの判定に使う */
  routerNames: Set<string>;
};

export function SwitchNode({ data, selected }: NodeProps & { data: SwitchNodeData }) {
  const sw = data.resource;
  const routerPort = sw.ports.find((p) => data.routerNames.has(p));
  const nicPorts = sw.ports.filter((p) => !data.routerNames.has(p));

  return (
    <NodeShell kind="switch" title={sw.name} subtitle={sw.id} selected={selected}>
      {/* 左: Router方向。1本のみなので固定id、既に接続済みならその行にハンドルをぶら下げる */}
      <div className="node-row-list">
        <div className="node-row" style={{ position: "relative" }}>
          <Handle
            type="source"
            position={Position.Left}
            id="to-router"
            style={{ position: "absolute", left: -11, top: "50%" }}
          />
          <span className="node-row__label">
            {routerPort ? `→ ${routerPort}` : "(Router未接続)"}
          </span>
        </div>
      </div>

      {/* 右: NICごとに独立したハンドル */}
      <div className="node-row-list" style={{ marginTop: 6 }}>
        {nicPorts.length === 0 && (
          <div className="node-row"><span className="node-row__meta">NIC未接続</span></div>
        )}
        {nicPorts.map((port) => (
          <div className="node-row" style={{ position: "relative" }} key={port}>
            <span className="node-row__label">{port}</span>
            <Handle
              type="target"
              position={Position.Right}
              id={port}
              style={{ position: "absolute", right: -11, top: "50%" }}
            />
          </div>
        ))}
        {/* 新規NIC接続受け入れ用の汎用ハンドル。繋がると上のnicPorts行として独立表示される */}
        <div className="node-row node-row--ghost" style={{ position: "relative" }}>
          <span className="node-row__meta">+ 新規NIC</span>
          <Handle
            type="target"
            position={Position.Right}
            id="new-nic"
            style={{ position: "absolute", right: -11, top: "50%" }}
          />
        </div>
      </div>
    </NodeShell>
  );
}