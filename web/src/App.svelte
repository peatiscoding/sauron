<script>
  import { onMount, tick } from 'svelte';
  import MarkdownView from './lib/MarkdownView.svelte';
  import OpenAPIView from './lib/OpenAPIView.svelte';
  import JSONView from './lib/JSONView.svelte';

  let connected   = $state(false);
  let content     = $state('');
  let filetype    = $state('');
  let filename    = $state('');
  let filepath    = $state('');
  let gitStatus   = $state('');
  let anchorRatio = $state(null);
  let viewEl;

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
    <span class="dot" class:live={connected} title={connected ? 'Connected' : 'Reconnecting…'}></span>
    <span class="fname">{filename || 'No file focused'}</span>
    {#if filetype}
      <span class="badge">{filetype}</span>
    {/if}
    {#if gitStatus}
      <span class="badge git-badge git-{gitStatus}">{gitStatus}</span>
    {/if}
    <span class="spacer"></span>
    <span class="brand">◉ SAURON</span>
  </header>

  <main class="view" bind:this={viewEl}>
    {#if !content}
      <div class="idle">
        <span class="eye">◉</span>
        <p>Watching for file focus…</p>
        <p class="hint">
          Open a <code>.md</code>, <code>.yaml</code>, or <code>.json</code> in nvim,
          or let Claude edit a file.
        </p>
      </div>
    {:else if filetype === 'markdown'}
      <MarkdownView {content} {filepath} />
    {:else if filetype === 'yaml'}
      <OpenAPIView {content} />
    {:else if filetype === 'json'}
      <JSONView {content} />
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

  .dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--red);
    flex-shrink: 0;
    transition: background 0.4s;
  }
  .dot.live { background: var(--green); }

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
    font-size: 11px;
    font-weight: 700;
    letter-spacing: 3px;
    color: var(--text-muted);
    opacity: 0.4;
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
