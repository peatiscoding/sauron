<script>
  import hljs from 'highlight.js';

  let { content } = $props();

  // Register Pkl language with highlight.js
  if (!hljs.getLanguage('pkl')) {
    hljs.registerLanguage('pkl', function () {
      const KEYWORDS = [
        'module', 'import', 'class', 'typealias', 'abstract', 'open',
        'local', 'hidden', 'fixed', 'const', 'new', 'amends', 'extends',
        'function', 'for', 'when', 'if', 'else', 'is', 'as', 'let',
        'throw', 'read', 'trace', 'in', 'out', 'outer', 'this', 'super',
      ];
      const LITERALS = ['true', 'false', 'null', 'nothing', 'unknown'];
      const BUILT_IN_TYPES = [
        'String', 'Int', 'Float', 'Boolean', 'Duration', 'DataSize',
        'List', 'Set', 'Map', 'Listing', 'Mapping', 'Dynamic', 'Pair',
        'Regex', 'Null', 'Nothing', 'Any', 'Module', 'Class', 'Number',
        'Int8', 'Int16', 'Int32', 'UInt', 'UInt8', 'UInt16', 'UInt32',
        'Collection', 'Comparable', 'Uri', 'Resource',
      ];

      const STRING_INTERPOLATION = {
        className: 'subst',
        begin: /\\\(/,
        end: /\)/,
        keywords: { keyword: KEYWORDS, literal: LITERALS },
        contains: [], // filled later
      };

      const STRING = {
        className: 'string',
        variants: [
          { begin: '"""', end: '"""', contains: [STRING_INTERPOLATION] },
          { begin: '"', end: '"', illegal: '\\n', contains: [
            { begin: /\\./ },
            STRING_INTERPOLATION,
          ]},
        ],
      };

      STRING_INTERPOLATION.contains = [STRING, { className: 'number', begin: /\b\d[\d_]*(\.\d[\d_]*)?\b/ }];

      return {
        name: 'Pkl',
        aliases: ['pkl'],
        keywords: {
          keyword: KEYWORDS,
          literal: LITERALS,
          type: BUILT_IN_TYPES,
        },
        contains: [
          hljs.C_LINE_COMMENT_MODE,
          { className: 'comment', begin: /\/\/\//, end: /$/, relevance: 10 },
          hljs.C_BLOCK_COMMENT_MODE,
          STRING,
          { className: 'number', begin: /\b0[xX][\da-fA-F_]+\b/ },
          { className: 'number', begin: /\b0[bB][01_]+\b/ },
          { className: 'number', begin: /\b0[oO][0-7_]+\b/ },
          { className: 'number', begin: /\b\d[\d_]*(\.\d[\d_]*)?(e[+-]?\d[\d_]*)?\b/, relevance: 0 },
          { className: 'meta', begin: /@\w+/ },
          { className: 'title.class', begin: /\b[A-Z]\w*/, relevance: 0 },
          { className: 'attr', begin: /\b[a-z_]\w*(?=\s*[=:{])/, relevance: 0 },
        ],
      };
    });
  }

  let highlighted = $derived(render(content));

  function render(src) {
    try {
      return hljs.highlight(src, { language: 'pkl' }).value;
    } catch (_) {
      return src.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    }
  }
</script>

<div class="pkl-wrap">
  <pre class="pkl-pre"><code class="hljs">{@html highlighted}</code></pre>
</div>

<style>
  .pkl-wrap {
    padding: 24px;
    overflow: auto;
    height: 100%;
    box-sizing: border-box;
  }

  .pkl-pre {
    margin: 0;
    font-family: var(--font-mono);
    font-size: 13px;
    line-height: 1.7;
    color: var(--text);
    white-space: pre;
    tab-size: 2;
  }

  .pkl-pre :global(.hljs-keyword)     { color: #ff7b72; }
  .pkl-pre :global(.hljs-literal)     { color: #79c0ff; }
  .pkl-pre :global(.hljs-type)        { color: #ffa657; }
  .pkl-pre :global(.hljs-title)       { color: #d2a8ff; }
  .pkl-pre :global(.hljs-string)      { color: #a5d6ff; }
  .pkl-pre :global(.hljs-subst)       { color: var(--text); }
  .pkl-pre :global(.hljs-number)      { color: #f8c555; }
  .pkl-pre :global(.hljs-comment)     { color: #8b949e; font-style: italic; }
  .pkl-pre :global(.hljs-meta)        { color: #d2a8ff; }
  .pkl-pre :global(.hljs-attr)        { color: #79c0ff; }

  :global([data-theme="light"]) .pkl-pre :global(.hljs-keyword)  { color: #cf222e; }
  :global([data-theme="light"]) .pkl-pre :global(.hljs-literal)  { color: #0550ae; }
  :global([data-theme="light"]) .pkl-pre :global(.hljs-type)     { color: #953800; }
  :global([data-theme="light"]) .pkl-pre :global(.hljs-title)    { color: #8250df; }
  :global([data-theme="light"]) .pkl-pre :global(.hljs-string)   { color: #0a3069; }
  :global([data-theme="light"]) .pkl-pre :global(.hljs-number)   { color: #0550ae; }
  :global([data-theme="light"]) .pkl-pre :global(.hljs-comment)  { color: #6e7781; }
  :global([data-theme="light"]) .pkl-pre :global(.hljs-meta)     { color: #8250df; }
  :global([data-theme="light"]) .pkl-pre :global(.hljs-attr)     { color: #0550ae; }
</style>
