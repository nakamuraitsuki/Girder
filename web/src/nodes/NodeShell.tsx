import type { ReactNode } from "react";
import "./nodes.css";

interface NodeShellProps {
  kind: "vm" | "switch" | "router" | "gateway";
  title: string;
  subtitle?: string;
  badge?: string;
  badgeVariant?: "running" | "shutoff" | "neutral";
  children?: ReactNode;
  selected?: boolean;
}

const KIND_LABEL: Record<NodeShellProps["kind"], string> = {
  vm: "VM",
  switch: "SWITCH",
  router: "ROUTER",
  gateway: "GATEWAY",
};

export function NodeShell({
  kind,
  title,
  subtitle,
  badge,
  badgeVariant = "neutral",
  children,
  selected,
}: NodeShellProps) {
  return (
    <div className={`node-shell node-shell--${kind} ${selected ? "is-selected" : ""}`}>
      <div className="node-shell__accent" />
      <div className="node-shell__header">
        <span className="node-shell__kind">{KIND_LABEL[kind]}</span>
        {badge && (
          <span className={`node-shell__badge node-shell__badge--${badgeVariant}`}>
            {badge}
          </span>
        )}
      </div>
      <div className="node-shell__title">{title}</div>
      {subtitle && <div className="node-shell__subtitle">{subtitle}</div>}
      {children}
    </div>
  );
}
