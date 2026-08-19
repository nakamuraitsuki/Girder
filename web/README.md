# girder-web (Phase 1: トポロジー編集キャンバス)

React Flow (@xyflow/react) + msw によるフロントエンド基盤。
実バックエンド (Go) 不要でキャンバス操作を確認できます。

## セットアップ

```bash
npm install
npm run dev
```

## 構成

- `src/types/topology.ts` — Girder API実データ形状に準拠した型定義
- `src/mocks/` — msw ハンドラ + fixture (3層Webアプリ構成: web/app/db + external/internal switch + router + gateway)
- `src/api/client.ts` — 薄いAPIクライアント。Phase2で実バックエンドに繋ぐ際はここ(または`VITE_API_BASE_URL`)だけ変更すればよい
- `src/store/topologyStore.ts` — zustand store。楽観的更新 (pending → confirmed/error) を管理
- `src/nodes/` — VM / Switch / Router / Gateway のカスタムノード
- `src/forms/` — Router PortへのIP割り当てフォーム (明示的な保存ボタン方式)
- `src/Canvas.tsx` — React Flow本体の配線。ノード間の線引き→関係操作APIの呼び出し、edgeのstatus別色分け(pending=黄, error=赤+2.5秒後に消える, confirmed=グレー)

## 操作方法

- VMのNICハンドル → Switch へドラッグ: `ConnectNICtoSwitch` 相当
- Switch → Router へドラッグ: `ConnectSwitchToRouter` 相当
- Gateway → Router へドラッグ: `ConnectGatewayToRouter` 相当
- 許可されない組み合わせ (例: VM→Router) はトースト表示で拒否
- 確立済みの線をダブルクリック: 切断確認 → Disconnect系API呼び出し
- Router Portの行をクリック: IP割り当てフォームを開く
- ノード位置はブラウザの localStorage に保持 (Girder Core側の状態ではない、UIのみのレイアウト情報)

## 未決事項 (Issueに記載の通り、着手時に決定)

- フォーム submit タイミング → 本実装では「値入力系(IP割り当て)は明示的な保存ボタン」「接続操作(線引き)は即時反映」で分離
- ポーリング vs 手動リフレッシュ → 現状は接続操作後に `GET /api/topology` を再取得する方式。他クライアントの変更検知が必要になった場合はポーリングを追加
- エラー表示/ロールバックUX → pending中は黄色破線、失敗時は赤+エラーメッセージをedge labelに表示し2.5秒でフェードアウト

## Phase 2 (Ubuntu実機に戻ってから)

- `VITE_API_BASE_URL` を実Goサーバーに向ける、または `msw` の起動を無効化してビルド成果物をGoサーバーへ埋め込み配信
- コンソール(VNC/シリアル)機能は別Issueで着手
