// Formats seconds into "m:ss" (e.g. 125 -> "2:05")
export function formatDuration(seconds) {
  if (!seconds || isNaN(seconds)) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

// Formats total seconds into a human summary like "2 hr 15 min" or "45 min"
export function formatTotalDuration(seconds) {
  if (!seconds || isNaN(seconds)) return '0 min'
  const hours = Math.floor(seconds / 3600)
  const mins = Math.floor((seconds % 3600) / 60)
  if (hours > 0) return `${hours} hr ${mins} min`
  return `${mins} min`
}

// Formats an ISO date string into "Mon DD, YYYY" (e.g. "Sep 10, 2026")
export function formatDate(isoString) {
  if (!isoString) return '-'
  return new Date(isoString).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}