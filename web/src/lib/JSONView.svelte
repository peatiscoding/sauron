<script>
  let { content } = $props();

  let rendered = $derived(highlight(content));
  let error     = $derived(parseError(content));

  function parseError(src) {
    try { JSON.parse(src); return null; }
    catch (e) { return e.message; }
  }

  function escape(s) {
    return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
  }

  function highlight(src) {
    let parsed;
    try { parsed = JSON.parse(src); }
    catch (_) {
      // Show raw with basic escaping if invalid JSON
      return `<span class="err">${escape(src)}</span>`;
    }
    return renderValue(parsed, 0);
  }

  function indent(depth) {
    return '  '.repeat(depth);
  }

  function renderValue(val, depth) {
    if (val === null)            return `<span class="null">null</span>`;
    if (val === true)            return `<span class="bool">true</span>`;
    if (val === false)           return `<span class="bool">false</span>`;
    if (typeof val === 'number') return `<span class="num">${val}</span>`;
    if (typeof val === 'string') return `<span class="str">"${escape(val)}"</span>`;
    if (Array.isArray(val))      return renderArray(val, depth);
    if (typeof val === 'object') return renderObject(val, depth);
    return escape(String(val));
  }

  function renderArray(arr, depth) {
    if (arr.length === 0) return `<span class="brace">[]</span>`;
    const inner = arr.map((v, i) => {
      const comma = i < arr.length - 1 ? '<span class="punct">,</span>' : '';
      return `${indent(depth + 1)}${renderValue(v, depth + 1)}${comma}`;
    }).join('\n');
    return `<span class="brace">[</span>\n${inner}\n${indent(depth)}<span class="brace">]</span>`;
  }

  function renderObject(obj, depth) {
    const keys = Object.keys(obj);
    if (keys.length === 0) return `<span class="brace">{}</span>`;
    const inner = keys.map((k, i) => {
      const comma = i < keys.length - 1 ? '<span class="punct">,</span>' : '';
      return `${indent(depth + 1)}<span class="key">"${escape(k)}"</span><span class="colon">: </span>${renderValue(obj[k], depth + 1)}${comma}`;
    }).join('\n');
    return `<span class="brace">{</span>\n${inner}\n${indent(depth)}<span class="brace">}</span>`;
  }
</script>

<div class="json-wrap">
  {#if error}
    <div class="parse-error">JSON parse error: {error}</div>
  {/if}
  <pre class="json-pre">{@html rendered}</pre>
</div>

<style>
  .json-wrap {
    padding: 24px;
    overflow: auto;
    height: 100%;
    box-sizing: border-box;
  }

  .parse-error {
    background: color-mix(in srgb, var(--red) 15%, transparent);
    border: 1px solid var(--red);
    border-radius: 6px;
    padding: 8px 14px;
    margin-bottom: 16px;
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--red);
  }

  .json-pre {
    margin: 0;
    font-family: var(--font-mono);
    font-size: 13px;
    line-height: 1.7;
    color: var(--text);
    white-space: pre;
    tab-size: 2;
  }

  .json-pre :global(.key)   { color: #79b8ff; }
  .json-pre :global(.str)   { color: #9ecbff; }
  .json-pre :global(.num)   { color: #f8c555; }
  .json-pre :global(.bool)  { color: #f97583; }
  .json-pre :global(.null)  { color: #f97583; opacity: 0.7; }
  .json-pre :global(.brace) { color: var(--text-muted); }
  .json-pre :global(.colon) { color: var(--text-muted); }
  .json-pre :global(.punct) { color: var(--text-muted); }
  .json-pre :global(.err)   { color: var(--red); }
</style>
