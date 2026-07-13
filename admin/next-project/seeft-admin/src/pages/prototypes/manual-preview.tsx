import React, { useState } from "react";

// 入力された Google ドキュメントの URL を iframe 用の /preview URL に変換する
// （/export?format=pdf は Content-Disposition: attachment でダウンロードになるため、表示には /preview を使う）
const toEmbedPreviewUrl = (originalUrl: string): string | null => {
  const trimmed = originalUrl.trim();
  if (!trimmed) return null;

  try {
    const url = new URL(trimmed);

    if (url.hostname.includes("docs.google.com")) {
      // パスからドキュメントID部分だけ残し、/preview に統一（iframe 内でその場表示される）
      const match = url.pathname.match(/\/document\/d\/([^/]+)/);
      if (match) {
        return `https://docs.google.com/document/d/${match[1]}/preview`;
      }
    }

    return trimmed;
  } catch {
    return null;
  }
};

const ManualPreviewPrototype: React.FC = () => {
  const [inputUrl, setInputUrl] = useState("");
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handlePreview = () => {
    const converted = toEmbedPreviewUrl(inputUrl);
    if (!converted) {
      setError("有効な URL を入力してください。");
      setPreviewUrl(null);
      setIsOpen(false);
      return;
    }

    setError(null);
    setPreviewUrl(converted);
    setIsOpen(true);
  };

  return (
    <main style={{ maxWidth: 960, margin: "40px auto", padding: "0 16px" }}>
      <h1 style={{ fontSize: 24, fontWeight: 600, marginBottom: 16 }}>
        Google ドキュメント プレビュー（プロトタイプ）
      </h1>

      <p style={{ marginBottom: 12, fontSize: 14, color: "#444" }}>
        共有リンク（edit や view の URL）を入力すると、
        /preview に書き換えて iframe 内にその場で表示します（ダウンロードされません）。
        既存機能とは独立したサンプルページです。
      </p>

      <div style={{ marginBottom: 16 }}>
        <label style={{ display: "block", fontSize: 14, marginBottom: 4 }}>
          Google ドキュメントの URL
        </label>
        <input
          type="text"
          value={inputUrl}
          onChange={(e) => setInputUrl(e.target.value)}
          placeholder="https://docs.google.com/document/d/..."
          style={{
            width: "100%",
            padding: "8px 12px",
            borderRadius: 4,
            border: "1px solid #ccc",
            fontSize: 14,
          }}
        />
      </div>

      <button
        type="button"
        onClick={handlePreview}
        style={{
          padding: "6px 16px",
          borderRadius: 4,
          border: "none",
          backgroundColor: "#2563eb",
          color: "#fff",
          fontSize: 14,
          cursor: "pointer",
        }}
      >
        プレビューを開く
      </button>

      {error && (
        <p style={{ marginTop: 12, color: "#b91c1c", fontSize: 14 }}>{error}</p>
      )}

      {previewUrl && (
        <section style={{ marginTop: 24 }}>
          <button
            type="button"
            onClick={() => setIsOpen((prev) => !prev)}
            style={{
              padding: "4px 12px",
              borderRadius: 4,
              border: "1px solid #ccc",
              backgroundColor: "#f3f4f6",
              fontSize: 13,
              cursor: "pointer",
            }}
          >
            {isOpen ? "プレビューを閉じる" : "プレビューを表示"}
          </button>

          {isOpen && (
            <div style={{ marginTop: 16 }}>
              <iframe
                src={previewUrl}
                style={{
                  width: "100%",
                  height: "600px",
                  border: "1px solid #e5e7eb",
                }}
              />
            </div>
          )}
        </section>
      )}
    </main>
  );
};

export default ManualPreviewPrototype;

