<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte'
  import { OpenFolderPicker } from '../../bindings/github.com/alciller88/commitlore/app/gitapp.js'
  import { GenerateStory } from '../../bindings/github.com/alciller88/commitlore/app/storyapp.js'
  import { ListStyles } from '../../bindings/github.com/alciller88/commitlore/app/styleapp.js'
  import { GetLLMConfig } from '../../bindings/github.com/alciller88/commitlore/app/configapp.js'

  export let activeRepo = ''

  const dispatch = createEventDispatcher()

  let repo = ''
  let from = ''
  let styleName = 'formal'
  let styles: Array<{name: string}> = []
  let loading = false
  let error = ''
  let htmlResult = ''
  let llmStatus = ''

  $: repo = activeRepo || repo

  onMount(() => {
    loadStyles()
    loadLLMStatus()
  })

  async function loadStyles() {
    try {
      const raw = await ListStyles()
      styles = JSON.parse(raw)
    } catch {}
  }

  async function loadLLMStatus() {
    try {
      const cfg = await GetLLMConfig()
      if (cfg.provider && cfg.keyConfigured) {
        llmStatus = `${cfg.provider} (${cfg.model || 'default'})`
      } else if (cfg.provider) {
        llmStatus = `${cfg.provider} — key not configured`
      } else {
        llmStatus = 'not configured'
      }
    } catch {
      llmStatus = 'not configured'
    }
  }

  async function pickRepo() {
    const path = await OpenFolderPicker()
    if (path) {
      repo = path
      dispatch('repoSelected', path)
    }
  }

  async function tellStory() {
    if (!repo) return
    loading = true
    error = ''
    htmlResult = ''
    try {
      htmlResult = await GenerateStory(repo, from, styleName)
    } catch (e: any) {
      error = e?.message || 'Story generation failed'
    } finally {
      loading = false
    }
  }

  function copyText() {
    const tmp = document.createElement('div')
    tmp.innerHTML = htmlResult
    navigator.clipboard.writeText(tmp.innerText)
  }

  function saveAsFile() {
    const blob = new Blob([htmlResult], { type: 'text/html' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'story.html'
    a.click()
    URL.revokeObjectURL(url)
  }
</script>

<div class="screen">
  <h1>Tell the Story</h1>

  {#if error}
    <div class="banner error">{error}</div>
  {/if}

  <div class="form">
    <div class="field">
      <label>Repository</label>
      <div class="repo-input">
        <input type="text" bind:value={repo} placeholder="Local path or owner/repo" />
        <button class="pick-btn" on:click={pickRepo}>Browse</button>
      </div>
    </div>

    <div class="row">
      <div class="field">
        <label>From</label>
        <input type="text" bind:value={from} placeholder="Starting commit or tag" />
      </div>
      <div class="field">
        <label>Style</label>
        <select bind:value={styleName}>
          {#each styles as s}
            <option value={s.name}>{s.name}</option>
          {/each}
        </select>
      </div>
    </div>

    <div class="action-row">
      <button class="action-btn" on:click={tellStory} disabled={loading || !repo}>
        {#if loading}
          <span class="spinner"></span> Generating...
        {:else}
          Tell the story
        {/if}
      </button>
      <span class="llm-status">LLM: {llmStatus}</span>
    </div>
  </div>

  {#if htmlResult}
    <div class="result">
      <div class="result-toolbar">
        <button class="tool-btn" on:click={copyText}>Copy text</button>
        <button class="tool-btn" on:click={saveAsFile}>Save as .html</button>
      </div>
      <iframe class="preview" srcdoc={htmlResult} sandbox="" title="Story preview"></iframe>
    </div>
  {/if}
</div>

<style>
  .screen { display: flex; flex-direction: column; gap: 20px; height: 100%; }
  h1 { color: #e6edf3; font-size: 22px; margin: 0; }
  .form { display: flex; flex-direction: column; gap: 14px; }
  .field { display: flex; flex-direction: column; gap: 4px; flex: 1; }
  .field label { color: #8b949e; font-size: 12px; text-transform: uppercase; }
  .row { display: flex; gap: 12px; }

  input, select {
    padding: 8px 12px; background: #0d1117; border: 1px solid #30363d;
    border-radius: 6px; color: #e6edf3; font-size: 14px;
    font-family: 'JetBrains Mono', monospace;
  }
  input:focus, select:focus { outline: none; border-color: #58a6ff; }
  select { cursor: pointer; }

  .repo-input { display: flex; gap: 8px; }
  .repo-input input { flex: 1; }
  .pick-btn {
    padding: 8px 14px; background: #21262d; border: 1px solid #30363d;
    border-radius: 6px; color: #e6edf3; cursor: pointer; font-family: inherit; font-size: 13px;
  }
  .pick-btn:hover { border-color: #58a6ff; }

  .action-row { display: flex; align-items: center; gap: 16px; }
  .action-btn {
    padding: 10px 20px; background: #238636; border: none; border-radius: 6px;
    color: #fff; font-size: 14px; cursor: pointer;
    display: flex; align-items: center; gap: 8px; font-family: inherit;
  }
  .action-btn:hover { background: #2ea043; }
  .action-btn:disabled { opacity: 0.5; cursor: not-allowed; }

  .llm-status { color: #8b949e; font-size: 12px; }

  .spinner {
    display: inline-block; width: 14px; height: 14px;
    border: 2px solid #ffffff44; border-top-color: #fff;
    border-radius: 50%; animation: spin 0.6s linear infinite;
  }
  @keyframes spin { to { transform: rotate(360deg); } }

  .result { display: flex; flex-direction: column; gap: 8px; flex: 1; min-height: 300px; }
  .result-toolbar { display: flex; gap: 8px; }
  .tool-btn {
    padding: 6px 12px; background: #21262d; border: 1px solid #30363d;
    border-radius: 6px; color: #e6edf3; cursor: pointer; font-size: 12px; font-family: inherit;
  }
  .tool-btn:hover { border-color: #58a6ff; }

  .preview {
    flex: 1; border: 1px solid #30363d; border-radius: 6px;
    background: #fff; min-height: 300px;
  }

  .banner.error {
    background: #da363433; border: 1px solid #da3634; color: #f85149;
    padding: 8px 12px; border-radius: 6px; font-size: 13px;
  }
</style>
