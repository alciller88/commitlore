<script lang="ts">
  import { onMount } from 'svelte'
  import { ListStyles, ShowStyle, ImportStyle, ExportStyle, DeleteStyle, CreateStyle } from '../../bindings/github.com/alciller88/commitlore/app/styleapp.js'
  import { OpenFolderPicker } from '../../bindings/github.com/alciller88/commitlore/app/gitapp.js'

  export let activeRepo = ''

  type StyleInfo = {
    name: string
    builtIn: boolean
    description: string
    author: string
    version: string
  }

  let styles: StyleInfo[] = []
  let selectedStyle: any = null
  let selectedName = ''
  let loading = false
  let error = ''
  let showCreate = false
  let newName = ''
  let newDesc = ''
  let newAuthor = ''

  onMount(loadStyles)

  async function loadStyles() {
    loading = true
    error = ''
    try {
      const raw = await ListStyles()
      styles = JSON.parse(raw)
    } catch (e: any) {
      error = e?.message || 'Failed to load styles'
    } finally {
      loading = false
    }
  }

  async function selectStyle(name: string) {
    selectedName = name
    try {
      const raw = await ShowStyle(name)
      selectedStyle = JSON.parse(raw)
    } catch (e: any) {
      error = e?.message || 'Failed to load style details'
    }
  }

  async function importStyle() {
    error = ''
    try {
      const path = await OpenFolderPicker()
      if (path) {
        await ImportStyle(path)
        await loadStyles()
      }
    } catch (e: any) {
      error = e?.message || 'Import failed'
    }
  }

  async function exportStyle(name: string) {
    error = ''
    try {
      await ExportStyle(name, name + '.shipstyle')
    } catch (e: any) {
      error = e?.message || 'Export failed'
    }
  }

  async function deleteStyle(name: string) {
    error = ''
    try {
      await DeleteStyle(name)
      if (selectedName === name) {
        selectedStyle = null
        selectedName = ''
      }
      await loadStyles()
    } catch (e: any) {
      error = e?.message || 'Delete failed'
    }
  }

  async function createStyle() {
    if (!newName.trim()) return
    error = ''
    try {
      await CreateStyle(newName.trim(), newDesc.trim(), newAuthor.trim())
      showCreate = false
      newName = ''
      newDesc = ''
      newAuthor = ''
      await loadStyles()
    } catch (e: any) {
      error = e?.message || 'Create failed'
    }
  }
</script>

<div class="screen">
  <div class="header-row">
    <h1>Styles</h1>
    <div class="header-actions">
      <button class="tool-btn" on:click={importStyle}>Import</button>
      <button class="tool-btn" on:click={() => showCreate = !showCreate}>
        {showCreate ? 'Cancel' : 'Create'}
      </button>
    </div>
  </div>

  {#if error}
    <div class="banner error">{error}</div>
  {/if}

  {#if showCreate}
    <div class="create-form">
      <input type="text" bind:value={newName} placeholder="Style name" />
      <input type="text" bind:value={newDesc} placeholder="Description" />
      <input type="text" bind:value={newAuthor} placeholder="Author" />
      <button class="action-btn" on:click={createStyle} disabled={!newName.trim()}>Create style</button>
    </div>
  {/if}

  <div class="layout">
    <div class="style-list">
      {#each styles as s}
        <button
          class="style-item"
          class:active={selectedName === s.name}
          on:click={() => selectStyle(s.name)}
        >
          <div class="style-name">
            {s.name}
            <span class="badge" class:builtin={s.builtIn}>{s.builtIn ? 'built-in' : 'user'}</span>
          </div>
          <div class="style-desc">{s.description}</div>
          <div class="style-actions">
            <button class="small-btn" on:click|stopPropagation={() => exportStyle(s.name)}>Export</button>
            {#if !s.builtIn}
              <button class="small-btn danger" on:click|stopPropagation={() => deleteStyle(s.name)}>Delete</button>
            {/if}
          </div>
        </button>
      {/each}
      {#if styles.length === 0 && !loading}
        <p class="empty">No styles found.</p>
      {/if}
    </div>

    {#if selectedStyle}
      <div class="detail-panel">
        <h2>{selectedStyle.name}</h2>
        <div class="detail-row"><span class="label">Version:</span> {selectedStyle.version}</div>
        <div class="detail-row"><span class="label">Author:</span> {selectedStyle.author}</div>
        <div class="detail-row"><span class="label">Description:</span> {selectedStyle.description}</div>
        {#if selectedStyle.theme}
          <h3>Theme</h3>
          <div class="detail-row"><span class="label">Mode:</span> {selectedStyle.theme.mode || 'dark'}</div>
          {#if selectedStyle.theme.colors}
            <div class="color-row">
              {#each Object.entries(selectedStyle.theme.colors) as [key, val]}
                <div class="color-swatch">
                  <div class="swatch" style="background: {val}"></div>
                  <span class="color-name">{key}</span>
                </div>
              {/each}
            </div>
          {/if}
        {/if}
        {#if selectedStyle.llm_prompt}
          <h3>LLM Prompt</h3>
          <pre class="prompt-preview">{selectedStyle.llm_prompt}</pre>
        {/if}
      </div>
    {/if}
  </div>
</div>

<style>
  .screen { display: flex; flex-direction: column; gap: 16px; height: 100%; }
  .header-row { display: flex; justify-content: space-between; align-items: center; }
  h1 { color: #e6edf3; font-size: 22px; margin: 0; }
  h2 { color: #e6edf3; font-size: 18px; margin: 0 0 12px; }
  h3 { color: #8b949e; font-size: 13px; margin: 16px 0 8px; text-transform: uppercase; }
  .header-actions { display: flex; gap: 8px; }

  .tool-btn {
    padding: 6px 14px; background: #21262d; border: 1px solid #30363d;
    border-radius: 6px; color: #e6edf3; cursor: pointer; font-size: 13px; font-family: inherit;
  }
  .tool-btn:hover { border-color: #58a6ff; }

  .create-form {
    display: flex; gap: 8px; padding: 12px; background: #161b22;
    border: 1px solid #30363d; border-radius: 8px;
  }
  .create-form input {
    padding: 8px 12px; background: #0d1117; border: 1px solid #30363d;
    border-radius: 6px; color: #e6edf3; font-size: 13px; flex: 1;
    font-family: 'JetBrains Mono', monospace;
  }
  .create-form input:focus { outline: none; border-color: #58a6ff; }

  .action-btn {
    padding: 8px 16px; background: #238636; border: none; border-radius: 6px;
    color: #fff; font-size: 13px; cursor: pointer; font-family: inherit;
  }
  .action-btn:hover { background: #2ea043; }
  .action-btn:disabled { opacity: 0.5; cursor: not-allowed; }

  .layout { display: flex; gap: 16px; flex: 1; overflow: hidden; }
  .style-list {
    display: flex; flex-direction: column; gap: 4px; width: 320px;
    overflow-y: auto; flex-shrink: 0;
  }

  .style-item {
    display: flex; flex-direction: column; gap: 4px; padding: 10px 12px;
    background: transparent; border: 1px solid transparent; border-radius: 6px;
    color: #e6edf3; cursor: pointer; text-align: left; width: 100%; font-family: inherit;
    transition: background 0.15s;
  }
  .style-item:hover { background: #161b22; border-color: #30363d; }
  .style-item.active { background: #1f6feb22; border-color: #58a6ff; }

  .style-name { display: flex; align-items: center; gap: 8px; font-size: 14px; }
  .badge {
    font-size: 10px; padding: 2px 6px; border-radius: 10px;
    background: #30363d; color: #8b949e;
  }
  .badge.builtin { background: #1f6feb33; color: #58a6ff; }
  .style-desc { color: #8b949e; font-size: 12px; }

  .style-actions { display: flex; gap: 6px; margin-top: 4px; }
  .small-btn {
    font-size: 11px; padding: 2px 8px; background: none;
    border: 1px solid #30363d; border-radius: 4px; color: #8b949e;
    cursor: pointer; font-family: inherit;
  }
  .small-btn:hover { border-color: #58a6ff; color: #e6edf3; }
  .small-btn.danger:hover { border-color: #da3634; color: #f85149; }

  .detail-panel {
    flex: 1; padding: 16px; background: #161b22; border: 1px solid #30363d;
    border-radius: 8px; overflow-y: auto;
  }
  .detail-row { font-size: 13px; margin-bottom: 4px; }
  .detail-row .label { color: #8b949e; }
  .color-row { display: flex; gap: 8px; flex-wrap: wrap; }
  .color-swatch { display: flex; align-items: center; gap: 4px; }
  .swatch { width: 16px; height: 16px; border-radius: 3px; border: 1px solid #30363d; }
  .color-name { font-size: 11px; color: #8b949e; }

  .prompt-preview {
    font-family: 'JetBrains Mono', monospace; font-size: 12px;
    background: #0d1117; border: 1px solid #30363d; border-radius: 6px;
    padding: 12px; color: #e6edf3; white-space: pre-wrap; overflow: auto; max-height: 200px;
  }

  .empty { color: #8b949e; text-align: center; padding: 40px; }

  .banner.error {
    background: #da363433; border: 1px solid #da3634; color: #f85149;
    padding: 8px 12px; border-radius: 6px; font-size: 13px;
  }
</style>
