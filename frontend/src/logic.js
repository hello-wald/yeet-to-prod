// Pure UI helpers — the frontend's tiny units under test.

// Full-screen background: green when safe, red when not, slate while unknown.
export function bgColor(safe) {
  if (safe === true) return "#16a34a"; // green
  if (safe === false) return "#dc2626"; // red
  return "#334155"; // slate (loading / error / unknown)
}

// Format a moment in a given IANA timezone as a readable date + 24h clock.
// e.g. "Wed, 12 Aug 2026, 10:00:00". Empty string if timeZone is bad/missing.
export function formatClock(date, timeZone) {
  if (!timeZone) return "";
  try {
    return new Intl.DateTimeFormat("en-GB", {
      timeZone,
      weekday: "short",
      day: "2-digit",
      month: "short",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: false,
    }).format(date);
  } catch {
    return "";
  }
}
