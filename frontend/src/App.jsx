import { useEffect, useState, useCallback } from "react";
import { bgColor, formatClock } from "./logic.js";

const API = import.meta.env.VITE_API_URL || "http://localhost:8080";
const DEFAULT_COUNTRY = import.meta.env.VITE_DEFAULT_COUNTRY || "ID";
const COUNTRIES = ["ID", "IN", "CN", "US", "AE", "JP"]; // AE = Dubai/UAE

// Client-side fallback so the clock shows immediately, even while loading or on
// error. Backend's `timezone` (authoritative) is preferred when present.
const COUNTRY_TZ = {
  ID: "Asia/Jakarta",
  IN: "Asia/Kolkata",
  CN: "Asia/Shanghai",
  US: "America/New_York",
  AE: "Asia/Dubai",
  JP: "Asia/Tokyo",
};

export default function App() {
  const [country, setCountry] = useState(DEFAULT_COUNTRY);
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);
  const [now, setNow] = useState(() => new Date());

  const check = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(`${API}/should-i-deploy?country=${country}`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      setData(await res.json());
    } catch (e) {
      setError(e.message);
      setData(null);
    } finally {
      setLoading(false);
    }
  }, [country]);

  useEffect(() => {
    check();
  }, [check]);

  // tick the clock every second
  useEffect(() => {
    const id = setInterval(() => setNow(new Date()), 1000);
    return () => clearInterval(id);
  }, []);

  const safe = error ? null : data ? data.safe : null;
  const tz = data?.timezone || COUNTRY_TZ[country];
  const clock = formatClock(now, tz);

  return (
    <div className="screen" style={{ backgroundColor: bgColor(safe) }}>
      {clock && (
        <div className="clock">
          {data?.country || country} · {clock}
        </div>
      )}
      <div className="message">
        {loading ? "…" : error ? `⚠️ ${error}` : data?.message}
      </div>
      <div className="controls">
        <select value={country} onChange={(e) => setCountry(e.target.value)}>
          {COUNTRIES.map((c) => (
            <option key={c} value={c}>
              {c}
            </option>
          ))}
        </select>
        <button onClick={check}>check again</button>
      </div>
    </div>
  );
}
