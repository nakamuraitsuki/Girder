import { create } from "zustand";
import { api, ApiError } from "../api/client";
import type {
  TopologySnapshot,
  TopologyConnection,
  ConnectionKind,
} from "../types/topology";

/**
 * 各Edgeの表示状態。Girder Core側の永続状態ではなく、
 * APIレスポンス確定を待つ間だけ存在するReact側のローカルUI状態。
 * (design doc 「Node/Edgeの表示状態とDriverの真実のズレ」への対応)
 */
export type ConnectionStatus = "confirmed" | "pending" | "error";

interface ConnectionEntry {
  conn: TopologyConnection;
  status: ConnectionStatus;
  errorMessage?: string;
}

interface TopologyState {
  snapshot: TopologySnapshot | null;
  connectionEntries: Record<string, ConnectionEntry>;
  loading: boolean;
  loadError: string | null;

  loadTopology: () => Promise<void>;
  connect: (
    kind: ConnectionKind,
    source: string,
    target: string,
  ) => Promise<void>;
  disconnect: (connectionId: string) => Promise<void>;
}

let tempIdCounter = 0;

export const useTopologyStore = create<TopologyState>((set, get) => ({
  snapshot: null,
  connectionEntries: {},
  loading: false,
  loadError: null,

  loadTopology: async () => {
    set({ loading: true, loadError: null });
    try {
      const snapshot = await api.getTopology();
      const entries: Record<string, ConnectionEntry> = {};
      for (const conn of snapshot.connections) {
        entries[conn.id] = { conn, status: "confirmed" };
      }
      set({ snapshot, connectionEntries: entries, loading: false });
    } catch (err) {
      set({
        loading: false,
        loadError:
          err instanceof ApiError ? err.message : "トポロジーの取得に失敗しました",
      });
    }
  },

  connect: async (kind, source, target) => {
    const tempId = `pending-${tempIdCounter++}`;
    const optimisticConn: TopologyConnection = { id: tempId, kind, source, target };

    // 楽観的にEdgeをpending状態で即座に表示する
    set((state) => ({
      connectionEntries: {
        ...state.connectionEntries,
        [tempId]: { conn: optimisticConn, status: "pending" },
      },
    }));

    try {
      const created = await api.createConnection(kind, source, target);
      set((state) => {
        const { [tempId]: _discard, ...rest } = state.connectionEntries;
        return {
          connectionEntries: {
            ...rest,
            [created.id]: { conn: created, status: "confirmed" },
          },
        };
      });
      // NIC等の接続先フィールド(connectedSwitchPort等)も更新が必要なので再取得する。
      // Phase2でPOSTレスポンスに更新後リソースを含める形に変えられれば省略可。
      await get().loadTopology();
    } catch (err) {
      const message =
        err instanceof ApiError ? err.message : "接続に失敗しました";
      set((state) => ({
        connectionEntries: {
          ...state.connectionEntries,
          [tempId]: { conn: optimisticConn, status: "error", errorMessage: message },
        },
      }));
      // ロールバック表示: 数秒後にerror Edgeを消す
      setTimeout(() => {
        set((state) => {
          const { [tempId]: _discard, ...rest } = state.connectionEntries;
          return { connectionEntries: rest };
        });
      }, 2500);
    }
  },

  disconnect: async (connectionId) => {
    const prevEntry = get().connectionEntries[connectionId];
    if (!prevEntry) return;

    set((state) => ({
      connectionEntries: {
        ...state.connectionEntries,
        [connectionId]: { ...prevEntry, status: "pending" },
      },
    }));

    try {
      await api.deleteConnection(connectionId);
      set((state) => {
        const { [connectionId]: _discard, ...rest } = state.connectionEntries;
        return { connectionEntries: rest };
      });
      await get().loadTopology();
    } catch (err) {
      const message =
        err instanceof ApiError ? err.message : "切断に失敗しました";
      // ロールバック: confirmed状態に戻す
      set((state) => ({
        connectionEntries: {
          ...state.connectionEntries,
          [connectionId]: { ...prevEntry, status: "error", errorMessage: message },
        },
      }));
      setTimeout(() => {
        set((state) => ({
          connectionEntries: {
            ...state.connectionEntries,
            [connectionId]: { ...prevEntry, status: "confirmed" },
          },
        }));
      }, 2500);
    }
  },
}));
