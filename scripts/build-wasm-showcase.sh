#!/usr/bin/env bash
#
# Build the ImmyGo showcase as a Wasm bundle suitable for embedding on
# a static website. Output goes to the directory passed as the first
# argument (default: ./website-wasm/).
#
# Usage:
#   scripts/build-wasm-showcase.sh             # → ./website-wasm/
#   scripts/build-wasm-showcase.sh public/demo # → ./public/demo/
#
# The output directory contains:
#   index.html      — minimal page that boots the wasm bundle
#   showcase.wasm   — compiled showcase (~20 MB raw, ~5 MB gzipped)
#   wasm_exec.js    — Go runtime support shim
#
# Serve the directory with any static-file server, e.g.:
#   cd website-wasm && python3 -m http.server 8000
#
# The yzma local-LLM provider is unavailable on wasm (it needs CGO).
# All other ImmyGo features — including the node canvas — work in the
# browser.

set -euo pipefail

DIST="${1:-website-wasm}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"

mkdir -p "$DIST"

GOROOT="$(go env GOROOT)"
WASM_EXEC=""
for candidate in "$GOROOT/lib/wasm/wasm_exec.js" "$GOROOT/misc/wasm/wasm_exec.js"; do
    if [ -f "$candidate" ]; then
        WASM_EXEC="$candidate"
        break
    fi
done
if [ -z "$WASM_EXEC" ]; then
    echo "error: could not find wasm_exec.js under $GOROOT" >&2
    exit 1
fi

echo "Building showcase.wasm…"
GOOS=js GOARCH=wasm go build -o "$DIST/showcase.wasm" "$ROOT/examples/ui-showcase"
cp "$WASM_EXEC" "$DIST/wasm_exec.js"

cat > "$DIST/index.html" <<'HTML'
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>ImmyGo — Showcase</title>
    <style>
        html, body {
            margin: 0;
            padding: 0;
            height: 100%;
            background: #f3f3f3;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif;
            overflow: hidden;
        }
        #status {
            position: fixed;
            inset: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #555;
            background: #f3f3f3;
            transition: opacity 200ms ease-out;
            pointer-events: none;
        }
        #status.hidden { opacity: 0; }
        #status .label { font-size: 14px; }
        #status .err { color: #c4271c; padding: 16px 24px; max-width: 600px; }
    </style>
</head>
<body>
    <div id="status"><div class="label">Loading ImmyGo showcase…</div></div>
    <script src="wasm_exec.js"></script>
    <script>
        const status = document.getElementById('status');
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("showcase.wasm"), go.importObject)
            .then(result => {
                status.classList.add('hidden');
                setTimeout(() => status.remove(), 250);
                go.run(result.instance);
            })
            .catch(err => {
                status.innerHTML = '<div class="err">Failed to load: ' + err.message + '</div>';
            });
    </script>
</body>
</html>
HTML

echo
echo "Built to: $DIST/"
ls -lh "$DIST/"
echo
echo "Serve with: cd $DIST && python3 -m http.server 8000"
