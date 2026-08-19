import type {
  TopologySnapshot,
  TopologyConnection,
  ConnectionKind,
} from "../types/topology";

// Phase2で実Goサーバーに繋ぐ際は、ここ(あるいはビルド時の環境変数)だけ変更すればよい。
// ビルド成果物をGoサーバーが配信する構成なら相対パスのままで動く。
const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "";

export class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
  if (!res.ok) {
    let message = res.statusText;
    try {
      const body = await res.json();
      if (body?.error) message = body.error;
    } catch {
      // レスポンスボディがJSONでない場合は statusText をそのまま使う
    }
    throw new ApiError(message, res.status);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export const api = {
  getTopology: () => request<TopologySnapshot>("/api/topology"),

  createConnection: (kind: ConnectionKind, source: string, target: string) =>
    request<TopologyConnection>("/api/topology/connections", {
      method: "POST",
      body: JSON.stringify({ kind, source, target }),
    }),

  deleteConnection: (id: string) =>
    request<void>(`/api/topology/connections/${id}`, { method: "DELETE" }),

  assignRouterPortIp: (
    routerId: string,
    portName: string,
    networks: string[],
  ) =>
    request<{ name: string; networks: string[] }>(
      `/api/ovn/routers/${routerId}/ports/${portName}`,
      { method: "PATCH", body: JSON.stringify({ networks }) },
    ),
};
