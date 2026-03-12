// Ensure bin is executable on Unix
if (process.platform !== 'win32') {
  const fs = require('fs');
  const path = require('path');
  const bin = path.join(__dirname, '..', 'bin', 'agent-rss');
  try { fs.chmodSync(bin, 0o755); } catch {}
}
