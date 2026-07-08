<script>
  import { marked } from 'marked';
  import mermaid from 'mermaid';
  import { onMount } from 'svelte';

  let { content, filepath = '' } = $props();

  let container;
  let ready = false;
  let renderSeq = 0; // abort stale renders on rapid updates
  let lightboxSvg = $state(null);

  onMount(() => {
    mermaid.initialize({
      startOnLoad: false,
      theme: 'dark',
      themeVariables: {
        background:        '#0d1117',
        primaryColor:      '#1f6feb',
        primaryTextColor:  '#c9d1d9',
        primaryBorderColor:'#30363d',
        lineColor:         '#58a6ff',
        secondaryColor:    '#161b22',
        tertiaryColor:     '#21262d',
        edgeLabelBackground:'#161b22',
        fontFamily:        'JetBrains Mono, Fira Code, monospace',
        fontSize:          '14px',
      },
      securityLevel: 'loose',
    });
    ready = true;
  });

  $effect(() => {
    if (!container || !ready || !content) return;
    render(content);
  });

  async function render(src) {
    const seq = ++renderSeq;

    // Render markdown HTML
    container.innerHTML = marked.parse(src, { breaks: true, gfm: true });

    // Rewrite relative image src → /api/image?path=<absolute>
    if (filepath) {
      const dir = filepath.substring(0, filepath.lastIndexOf('/') + 1);
      container.querySelectorAll('img').forEach((img) => {
        const src = img.getAttribute('src') ?? '';
        if (src && !src.startsWith('http') && !src.startsWith('//') && !src.startsWith('data:') && !src.startsWith('/')) {
          const parts = (dir + src).split('/');
          const resolved = [];
          for (const p of parts) {
            if (p === '..') resolved.pop();
            else if (p !== '.') resolved.push(p);
          }
          img.src = `/api/image?path=${encodeURIComponent(resolved.join('/'))}`;
        }
      });
    }

    // Find and replace mermaid code blocks with rendered SVG
    const blocks = container.querySelectorAll('code.language-mermaid');
    for (const block of blocks) {
      if (seq !== renderSeq) return; // content changed — abort

      const definition = block.textContent.trim();
      const id = `mermaid-${seq}-${Math.random().toString(36).slice(2, 8)}`;
      const wrapper = document.createElement('div');
      wrapper.className = 'mermaid-wrap';

      try {
        const { svg } = await mermaid.render(id, definition);
        wrapper.innerHTML = svg;
        wrapper.title = 'Click to expand';
        wrapper.addEventListener('click', () => {
          lightboxSvg = svg;
        });
      } catch (err) {
        wrapper.innerHTML = `<pre class="mermaid-err">Mermaid error:\n${err.message ?? err}</pre>`;
      }

      // block.parentElement is the <pre>
      block.parentElement?.replaceWith(wrapper);
    }
  }

  function closeLightbox() {
    lightboxSvg = null;
  }
</script>

<article class="markdown" bind:this={container}></article>

{#if lightboxSvg}
  <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
  <div class="lightbox" onclick={closeLightbox}>
    <div class="lightbox-inner" onclick={(e) => e.stopPropagation()}>
      <button class="lightbox-close" onclick={closeLightbox} aria-label="Close">✕</button>
      {@html lightboxSvg}
    </div>
  </div>
{/if}

<style>
  .markdown {
    padding: 40px 28px 80px;
    color: var(--text);
    font-family: var(--font-sans);
    font-size: 16px;
    line-height: 1.75;
  }

  /* ── Headings ───────────────────────────────────────────── */
  .markdown :global(h1),
  .markdown :global(h2),
  .markdown :global(h3),
  .markdown :global(h4),
  .markdown :global(h5),
  .markdown :global(h6) {
    color: var(--text);
    font-weight: 600;
    line-height: 1.3;
    margin: 1.6em 0 0.5em;
  }
  .markdown :global(h1) {
    font-size: 2em;
    border-bottom: 1px solid var(--border);
    padding-bottom: 0.3em;
    margin-top: 0;
  }
  .markdown :global(h2) {
    font-size: 1.5em;
    border-bottom: 1px solid var(--border);
    padding-bottom: 0.25em;
  }
  .markdown :global(h3)  { font-size: 1.25em; }
  .markdown :global(h4)  { font-size: 1.1em; }
  .markdown :global(h5),
  .markdown :global(h6)  { font-size: 1em; color: var(--text-muted); }

  /* ── Body text ──────────────────────────────────────────── */
  .markdown :global(p)       { margin: 0.8em 0; }
  .markdown :global(ul),
  .markdown :global(ol)      { margin: 0.8em 0; padding-left: 2em; }
  .markdown :global(li)      { margin: 0.3em 0; }
  .markdown :global(li > p)  { margin: 0.25em 0; }

  /* ── Links ──────────────────────────────────────────────── */
  .markdown :global(a)       { color: var(--accent); text-decoration: none; }
  .markdown :global(a:hover) { text-decoration: underline; }

  /* ── Inline code ────────────────────────────────────────── */
  .markdown :global(:not(pre) > code) {
    background:    var(--surface-2);
    border:        1px solid var(--border);
    border-radius: 3px;
    padding:       0.1em 0.4em;
    font-family:   var(--font-mono);
    font-size:     0.875em;
    color:         #ff7b72;
  }

  /* ── Code blocks ────────────────────────────────────────── */
  .markdown :global(pre) {
    background:    var(--surface);
    border:        1px solid var(--border);
    border-radius: var(--radius);
    padding:       16px 20px;
    overflow-x:    auto;
    margin:        1em 0;
    line-height:   1.55;
  }
  .markdown :global(pre code) {
    background: none;
    border:     none;
    padding:    0;
    font-size:  0.875em;
    color:      var(--text);
    font-family:var(--font-mono);
  }

  /* ── Blockquotes ────────────────────────────────────────── */
  .markdown :global(blockquote) {
    border-left:  3px solid var(--border);
    padding:      4px 16px;
    color:        var(--text-muted);
    margin:       1em 0;
    font-style:   italic;
  }

  /* ── Tables ─────────────────────────────────────────────── */
  .markdown :global(table) {
    width:          100%;
    border-collapse:collapse;
    margin:         1em 0;
    font-size:      0.9em;
    display:        block;
    overflow-x:     auto;
  }
  .markdown :global(th),
  .markdown :global(td) {
    border:  1px solid var(--border);
    padding: 8px 14px;
    text-align: left;
  }
  .markdown :global(th) {
    background:  var(--surface);
    font-weight: 600;
  }
  .markdown :global(tr:nth-child(even) td) {
    background: var(--surface);
  }

  /* ── Horizontal rule ────────────────────────────────────── */
  .markdown :global(hr) {
    border:     none;
    border-top: 1px solid var(--border);
    margin:     2em 0;
  }

  /* ── Images ─────────────────────────────────────────────── */
  .markdown :global(img) {
    max-width:    100%;
    border-radius:var(--radius);
    display:      block;
    margin:       1em 0;
  }

  /* ── Mermaid diagrams ───────────────────────────────────── */
  .markdown :global(.mermaid-wrap) {
    background:    var(--surface);
    border:        1px solid var(--border);
    border-radius: var(--radius);
    padding:       24px;
    margin:        1.2em 0;
    text-align:    center;
    overflow-x:    auto;
    cursor:        zoom-in;
    transition:    border-color 0.15s;
  }
  .markdown :global(.mermaid-wrap:hover) {
    border-color: var(--accent);
  }
  .markdown :global(.mermaid-wrap svg) {
    max-width: 100%;
    height:    auto;
  }
  .markdown :global(.mermaid-err) {
    color:       var(--red);
    font-family: var(--font-mono);
    font-size:   0.8em;
    text-align:  left;
    white-space: pre-wrap;
  }

  /* ── Lightbox ───────────────────────────────────────────── */
  .lightbox {
    position:        fixed;
    inset:           0;
    z-index:         1000;
    background:      rgba(0, 0, 0, 0.85);
    display:         flex;
    align-items:     center;
    justify-content: center;
    cursor:          zoom-out;
    backdrop-filter: blur(4px);
  }

  .lightbox-inner {
    position:      relative;
    background:    var(--surface);
    border:        1px solid var(--border);
    border-radius: var(--radius);
    padding:       48px 40px 40px;
    width:         92vw;
    height:        92vh;
    overflow:      auto;
    cursor:        default;
    display:       flex;
    align-items:   center;
    justify-content: center;
  }

  .lightbox-inner :global(svg) {
    width:      100%;
    height:     100%;
    max-width:  none;
    display:    block;
  }

  .lightbox-close {
    position:   absolute;
    top:        10px;
    right:      12px;
    background: none;
    border:     none;
    color:      var(--text-muted);
    font-size:  18px;
    cursor:     pointer;
    padding:    4px 8px;
    line-height: 1;
    border-radius: 4px;
    transition: color 0.15s, background 0.15s;
  }
  .lightbox-close:hover {
    color:      var(--text);
    background: var(--surface-2);
  }
</style>
