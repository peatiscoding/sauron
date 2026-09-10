<script>
  import { onMount, tick } from 'svelte';
  import MarkdownView from './lib/MarkdownView.svelte';
  import OpenAPIView from './lib/OpenAPIView.svelte';
  import JSONView from './lib/JSONView.svelte';
  import PklView from './lib/PklView.svelte';
  import SauronEye from './lib/SauronEye.svelte';

  let connected   = $state(false);
  let content     = $state('');
  let filetype    = $state('');
  let filename    = $state('');
  let filepath    = $state('');
  let gitStatus   = $state('');
  let anchorRatio = $state(null);
  let viewEl;
  let markdownView;

  // ── Theme ──────────────────────────────────────────────
  function getInitialTheme() {
    const stored = localStorage.getItem('depict-theme');
    if (stored === 'light' || stored === 'dark') return stored;
    return window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
  }
  let theme = $state(getInitialTheme());
  $effect(() => {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('depict-theme', theme);
  });
  $effect(() => {
    document.title = content ? `S - ${filename}` : 'Sauron';
  });

  // ── Dynamic favicon ────────────────────────────────────
  const eyeShards = [
    '18.5,3.1 21.5,3.1 20,15',    '22.9,3.2 25.8,4 21.3,15.2',
    '27.2,4.6 29.7,6.1 22.5,15.7','30.9,7 33,9.1 23.5,16.5',
    '33.9,10.3 35.4,12.8 24.3,17.5','36,14.2 36.7,17 24.8,18.7',
    '36.9,18.5 36.9,21.5 25,20',  '36.7,23 36,25.8 24.8,21.3',
    '35.4,27.2 33.9,29.7 24.3,22.5','33,30.9 30.9,33 23.5,23.5',
    '29.7,33.9 27.2,35.4 22.5,24.3','25.8,36 22.9,36.7 21.3,24.8',
    '21.5,36.9 18.5,36.9 20,25',  '17.1,36.7 14.2,36 18.7,24.8',
    '12.8,35.4 10.3,33.9 17.5,24.3','9.1,33 7,30.9 16.5,23.5',
    '6.1,29.7 4.6,27.2 15.7,22.5','4,25.8 3.2,22.9 15.2,21.3',
    '3.1,21.5 3.1,18.5 15,20',    '3.2,17.1 4,14.2 15.2,18.7',
    '4.6,12.8 6.1,10.3 15.7,17.5','7,9.1 9.1,7 16.5,16.5',
    '10.3,6.1 12.8,4.6 17.5,15.7','14.2,4 17.1,3.2 18.7,15.2',
  ];
  function buildFaviconSvg(open) {
    const polys = eyeShards.map((pts, i) =>
      `<polygon points="${pts}" fill="${i % 2 ? '#ff4400' : '#bb0000'}"/>`
    ).join('');
    const pupil = `<ellipse cx="20" cy="20" rx="2.3" ry="9.8" fill="#050000"/>` +
      `<ellipse cx="20" cy="20" rx="3.1" ry="10.5" fill="none" stroke="#ffe066" stroke-width="1.1"/>`;
    const scaleY = open ? 1 : 0.05;
    return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 40">` +
      `<g transform="translate(20,20) scale(1,${scaleY}) translate(-20,-20)">` +
      `${polys}${pupil}</g></svg>`;
  }
  $effect(() => {
    const el = document.getElementById('favicon');
    if (el) el.href = 'data:image/svg+xml,' + encodeURIComponent(buildFaviconSvg(connected));
  });
  function toggleTheme() { theme = theme === 'dark' ? 'light' : 'dark'; }

  onMount(() => { connect(); });

  function scrollToRatio(ratio) {
    if (!viewEl) return;
    tick().then(() => {
      viewEl.scrollTo({
        top: Math.min(1, Math.max(0, ratio)) * (viewEl.scrollHeight - viewEl.clientHeight),
        behavior: 'smooth',
      });
    });
  }

  async function exportPDF() {
    if (markdownView?.awaitRender) await markdownView.awaitRender();
    window.print();
  }

  function downloadSource() {
    const mimeMap = {
      yaml: 'application/yaml',
      json: 'application/json',
      markdown: 'text/markdown',
      mdx: 'text/markdown',
      pkl: 'text/plain',
    };
    const fallbackName = {
      yaml: 'spec.yaml',
      json: 'data.json',
      markdown: 'document.md',
      mdx: 'document.mdx',
      pkl: 'config.pkl',
    };
    const blob = new Blob([content], { type: mimeMap[filetype] || 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename || fallbackName[filetype] || 'file.txt';
    a.click();
    URL.revokeObjectURL(url);
  }

  async function exportHTML() {
    if (markdownView?.awaitRender) await markdownView.awaitRender();

    // Gather all stylesheets into inline styles
    const styles = [...document.styleSheets].map(sheet => {
      try { return [...sheet.cssRules].map(r => r.cssText).join('\n'); }
      catch (_) { return ''; }
    }).join('\n');

    const rendered = viewEl?.innerHTML || '';
    const baseName = (filename || 'export').replace(/\.[^.]+$/, '');

    const html = `<!DOCTYPE html>
<html data-theme="${theme}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>${baseName}</title>
<style>${styles}</style>
</head>
<body>
<div id="app"><main class="view" style="overflow:visible;height:auto;">${rendered}</main></div>
</body>
</html>`;

    const blob = new Blob([html], { type: 'text/html' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${baseName}.html`;
    a.click();
    URL.revokeObjectURL(url);
  }

  function connect() {
    const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(`${proto}//${location.host}/ws`);

    ws.onopen  = () => { connected = true; };
    ws.onclose = () => { connected = false; setTimeout(connect, 2000); };
    ws.onerror = () => ws.close();
    ws.onmessage = (e) => {
      try {
        const msg = JSON.parse(e.data);
        if (msg.type === 'focus') {
          if (msg.filename !== filename) anchorRatio = null;
          content   = msg.content;
          filetype  = msg.filetype;
          filename  = msg.filename;
          filepath  = msg.filepath || '';
          gitStatus = msg.gitStatus || '';
          if (anchorRatio !== null) scrollToRatio(anchorRatio);
        }
        if (msg.type === 'anchor') {
          anchorRatio = msg.line / msg.total;
          scrollToRatio(anchorRatio);
        }
      } catch (_) {}
    };
  }
</script>

<div class="layout">
  <header class="statusbar">
    <span class="fname">{filename || 'No file focused'}</span>
    {#if filetype}
      <span class="badge">{filetype}</span>
    {/if}
    {#if gitStatus}
      <span class="badge git-badge git-{gitStatus}">{gitStatus}</span>
    {/if}
    <span class="spacer"></span>
    {#if content}
      <button class="pdf-btn" onclick={downloadSource} title="Download {filetype.toUpperCase()}">⬇ {filetype.toUpperCase()}</button>
    {/if}
    {#if content}
      <button class="pdf-btn" onclick={exportHTML} title="Export to HTML">⬇ HTML</button>
      <button class="pdf-btn" onclick={exportPDF} title="Export to PDF">⬇ PDF</button>
    {/if}
    <button class="theme-btn" onclick={toggleTheme}
            title={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
            aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}>
      {#if theme === 'dark'}
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24"
             fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="4"/>
          <line x1="12" y1="2" x2="12" y2="4"/><line x1="12" y1="20" x2="12" y2="22"/>
          <line x1="2" y1="12" x2="4" y2="12"/><line x1="20" y1="12" x2="22" y2="12"/>
          <line x1="4.93" y1="4.93" x2="6.34" y2="6.34"/><line x1="17.66" y1="17.66" x2="19.07" y2="19.07"/>
          <line x1="4.93" y1="19.07" x2="6.34" y2="17.66"/><line x1="17.66" y1="6.34" x2="19.07" y2="4.93"/>
        </svg>
      {:else}
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24"
             fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
        </svg>
      {/if}
    </button>
    <span class="brand"><SauronEye open={connected} /> SAURON</span>
  </header>

  <main class="view" bind:this={viewEl}>
    {#if !content}
      <div class="idle">
        <span class="eye">◉</span>
        <p>Watching for file focus…</p>
        <p class="hint">
          Open a <code>.md</code>, <code>.mdx</code>, <code>.yaml</code>, <code>.json</code>, or <code>.pkl</code> in nvim,
          or let Claude edit a file.
        </p>
      </div>
    {:else if filetype === 'markdown' || filetype === 'mdx'}
      <MarkdownView {content} {filepath} {theme} {filetype} bind:this={markdownView} />
    {:else if filetype === 'yaml'}
      <OpenAPIView {content} />
    {:else if filetype === 'json'}
      <JSONView {content} />
    {:else if filetype === 'pkl'}
      <PklView {content} />
    {:else}
      <pre class="raw">{content}</pre>
    {/if}
  </main>
</div>

<style>
  .layout {
    height: 100vh;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  /* ── Status bar ─────────────────────────────────────────── */
  .statusbar {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 5px 16px;
    background: var(--surface);
    border-bottom: 1px solid var(--border);
    font-size: 12px;
    flex-shrink: 0;
    user-select: none;
  }

  /* ── Theme toggle button ────────────────────────────────── */
  .theme-btn {
    background:      none;
    border:          1px solid var(--border);
    border-radius:   4px;
    color:           var(--text-muted);
    width:           24px;
    height:          22px;
    display:         flex;
    align-items:     center;
    justify-content: center;
    padding:         0;
    cursor:          pointer;
    flex-shrink:     0;
    transition:      color 0.15s, border-color 0.15s;
  }
  .theme-btn:hover {
    color:        var(--accent);
    border-color: var(--accent);
  }

  .fname {
    font-family: var(--font-mono);
    color: var(--text);
    font-size: 12px;
  }

  .badge {
    background: var(--surface-2);
    border: 1px solid var(--border);
    border-radius: 4px;
    padding: 1px 7px;
    font-size: 11px;
    color: var(--accent);
    font-family: var(--font-mono);
  }

  .git-badge {
    color: var(--text-muted);
    border-color: var(--border);
  }
  .git-staged    { color: var(--green);  border-color: var(--green);  }
  .git-modified  { color: var(--yellow); border-color: var(--yellow); }
  .git-untracked { color: var(--text-muted); border-color: var(--border); }

  .spacer { flex: 1; }

  .brand {
    display:     flex;
    align-items: center;
    gap:         6px;
    font-size:   11px;
    font-weight: 700;
    letter-spacing: 3px;
    color:       var(--text-muted);
    opacity:     0.4;
  }

  /* ── Main view ──────────────────────────────────────────── */
  .view {
    flex: 1;
    overflow: auto;
  }

  /* ── Idle state ─────────────────────────────────────────── */
  .idle {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    gap: 10px;
    color: var(--text-muted);
    text-align: center;
  }

  .eye {
    font-size: 48px;
    opacity: 0.15;
    display: block;
    margin-bottom: 8px;
  }

  .hint {
    font-size: 13px;
    opacity: 0.6;
  }

  .hint code {
    background: var(--surface-2);
    border: 1px solid var(--border);
    padding: 1px 5px;
    border-radius: 3px;
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--accent);
  }

  /* ── PDF export button ──────────────────────────────────── */
  .pdf-btn {
    background:    none;
    border:        1px solid var(--border);
    border-radius: 4px;
    color:         var(--text-muted);
    font-family:   var(--font-mono);
    font-size:     11px;
    height:        22px;
    padding:       0 8px;
    cursor:        pointer;
    transition:    color 0.15s, border-color 0.15s;
    flex-shrink:   0;
  }
  .pdf-btn:hover {
    color:         var(--accent);
    border-color:  var(--accent);
  }

  /* ── Print layout ───────────────────────────────────────── */
  @media print {
    .statusbar { display: none; }
    .layout    { height: auto; overflow: visible; }
    .view      { overflow: visible; height: auto; flex: none; }
  }

  /* ── Raw fallback ───────────────────────────────────────── */
  .raw {
    padding: 24px;
    font-family: var(--font-mono);
    font-size: 13px;
    color: var(--text);
    white-space: pre-wrap;
    word-break: break-word;
    line-height: 1.6;
  }
</style>
