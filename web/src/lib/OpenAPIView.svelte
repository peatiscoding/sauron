<script>
  import { onMount } from 'svelte';
  import 'rapidoc';
  import jsyaml from 'js-yaml';

  let { content } = $props();
  let el;

  $effect(() => {
    if (!el || !content) return;
    loadSpec(content);
  });

  async function loadSpec(yaml) {
    await customElements.whenDefined('rapi-doc');
    try {
      const spec = jsyaml.load(yaml);
      el.loadSpec(spec);
    } catch (err) {
      console.error('OpenAPI parse error:', err);
      // Show raw error inside rapi-doc container as a fallback
      el.innerHTML = `
        <div style="padding:32px;color:#f85149;font-family:monospace;font-size:13px;">
          <strong>OpenAPI parse error:</strong><br><pre style="margin-top:8px;white-space:pre-wrap;">${err.message ?? err}</pre>
        </div>`;
    }
  }
</script>

<rapi-doc
  bind:this={el}
  theme="dark"
  bg-color="#0d1117"
  text-color="#c9d1d9"
  primary-color="#58a6ff"
  nav-bg-color="#161b22"
  nav-text-color="#8b949e"
  nav-hover-bg-color="#21262d"
  nav-hover-text-color="#c9d1d9"
  nav-accent-color="#58a6ff"
  render-style="focused"
  show-header="false"
  show-info="true"
  allow-search="true"
  allow-try="false"
  sort-tags="true"
  schema-style="tree"
  schema-expand-level="2"
  default-schema-tab="example"
  style="width:100%;height:100%;display:block;"
></rapi-doc>

<style>
  rapi-doc {
    width: 100%;
    height: 100%;
    display: block;
  }
</style>
