import { useState } from "react";
import { api, ApiError } from "../api/client";
import "./forms.css";

interface RouterPortFormProps {
  routerId: string;
  portName: string;
  currentNetworks: string[] | null;
  onClose: () => void;
  onSaved: () => void;
}

/**
 * Router PortへのIP割り当てフォーム。
 * 設計方針: 値入力系はonBlur即時反映ではなく明示的な保存ボタンにする。
 * CIDR等の複雑な入力を未完成のままAPIへ送ってしまうリスクを避けるため。
 */
export function RouterPortForm({
  routerId,
  portName,
  currentNetworks,
  onClose,
  onSaved,
}: RouterPortFormProps) {
  const [value, setValue] = useState(currentNetworks?.join(", ") ?? "");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSave = async () => {
    setSaving(true);
    setError(null);
    const networks = value
      .split(",")
      .map((s) => s.trim())
      .filter(Boolean);
    try {
      await api.assignRouterPortIp(routerId, portName, networks);
      onSaved();
      onClose();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "保存に失敗しました");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="form-panel-overlay" onClick={onClose}>
      <div className="form-panel" onClick={(e) => e.stopPropagation()}>
        <div className="form-panel__header">
          <span className="form-panel__kind">ROUTER PORT / IP割り当て</span>
          <button className="form-panel__close" onClick={onClose} aria-label="閉じる">
            ×
          </button>
        </div>
        <div className="form-panel__title">{portName}</div>

        <label className="form-field">
          <span className="form-field__label">Networks (CIDR, カンマ区切り)</span>
          <input
            className="form-field__input"
            value={value}
            onChange={(e) => setValue(e.target.value)}
            placeholder="例: 10.0.0.1/24"
            autoFocus
          />
        </label>

        {error && <div className="form-panel__error">{error}</div>}

        <div className="form-panel__actions">
          <button className="btn btn--ghost" onClick={onClose}>
            キャンセル
          </button>
          <button className="btn btn--primary" onClick={handleSave} disabled={saving}>
            {saving ? "保存中..." : "保存"}
          </button>
        </div>
      </div>
    </div>
  );
}
