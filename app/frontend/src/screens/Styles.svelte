<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { ListStyles, GetStyleDetail, SaveStyleDetail, DeleteStyle, ImportStyle, ExportStyle, GetStyleTheme } from '../../bindings/github.com/alciller88/commitlore/app/styleapp.js'
  import { GetActiveStyle } from '../../bindings/github.com/alciller88/commitlore/app/configapp.js'
  import { OpenFolderPicker } from '../../bindings/github.com/alciller88/commitlore/app/gitapp.js'
  import { activeStyle } from '../lib/store'

  type StyleListItem = { name: string; builtIn: boolean; primary: string; accent: string; background: string }

  let styles: StyleListItem[] = []
  let selectedName = ''
  let detail: any = null
  let isBuiltIn = false
  let activeTab = 'colors'
  let error = ''
  let success = ''
  let loading = false
  let creating = false
  let newName = ''
  let confirmDelete = false
  let currentActiveStyle = 'formal'

  const unsubStyle = activeStyle.subscribe(v => { currentActiveStyle = v })
  onDestroy(() => unsubStyle())

  onMount(loadStyles)

  async function loadStyles() {
    loading = true
    error = ''
    try {
      const raw = await ListStyles()
      const list = JSON.parse(raw)
      const enriched: StyleListItem[] = []
      for (const s of list) {
        try {
          const t = await GetStyleTheme(s.name)
          enriched.push({ name: s.name, builtIn: s.builtIn, primary: t.primary, accent: t.accent, background: t.background })
        } catch {
          enriched.push({ name: s.name, builtIn: s.builtIn, primary: '#58a6ff', accent: '#58a6ff', background: '#0d1117' })
        }
      }
      styles = enriched
    } catch (e: any) {
      error = e?.message || 'Failed to load styles'
    } finally {
      loading = false
    }
  }

  async function selectStyle(name: string) {
    creating = false
    selectedName = name
    error = ''
    try {
      detail = await GetStyleDetail(name)
      isBuiltIn = styles.find(s => s.name === name)?.builtIn ?? false
      activeTab = 'colors'
    } catch (e: any) {
      error = e?.message || 'Failed to load style'
      detail = null
    }
  }

  function startCreate() {
    creating = true
    selectedName = ''
    newName = ''
    detail = {
      name: '', version: '1.0.0', description: '', author: '',
      llmPrompt: '', tags: [], previewUrl: '', homepage: '',
      templates: { header: '# Changelog {{.Version}}\n', feature: '- {{.Message}}\n', fix: '- {{.Message}}\n', breaking: '- {{.Message}}\n', footer: '', storyIntro: '', storyMilestone: '', storyPeak: '', storyContributor: '', storyFooter: '' },
      vocabulary: {},
      theme: { mode: 'dark', colors: { primary: '#58a6ff', secondary: '#8b949e', background: '#0d1117', surface: '#161b22', text: '#e6edf3', accent: '#58a6ff', border: '#30363d' }, typography: { fontFamily: 'system-ui, sans-serif', fontSizeBase: '14px', fontSizeHeader: '24px', fontSizeCode: '13px' }, headerImage: '', logo: '', cardStyle: 'minimal', animations: false, customCss: '', windowControls: { default: '#666666', close: '#FF5F57', minimize: '#FEBC2E', maximize: '#28C840' } },
      terminal: { colors: { header: '', feature: '', fix: '', breaking: '', footer: '' }, decorators: { separator: '', bullet: '', indent: '' }, density: 'normal' }
    }
    isBuiltIn = false
    activeTab = 'colors'
  }

  async function saveStyle() {
    error = ''
    success = ''
    if (creating) {
      if (!newName.trim() || !/^[a-zA-Z0-9_-]+$/.test(newName.trim())) {
        error = 'Name must contain only letters, numbers, hyphens, underscores'
        return
      }
      if (styles.some(s => s.name === newName.trim())) {
        error = 'A style with this name already exists'
        return
      }
      detail.name = newName.trim()
    }
    try {
      await SaveStyleDetail(detail)
      success = 'Style saved.'
      setTimeout(() => success = '', 3000)
      if (creating) {
        creating = false
        selectedName = detail.name
      }
      await loadStyles()
    } catch (e: any) {
      error = e?.message || 'Save failed'
    }
  }

  async function doDelete() {
    error = ''
    try {
      await DeleteStyle(selectedName)
      detail = null
      selectedName = ''
      confirmDelete = false
      await loadStyles()
    } catch (e: any) {
      error = e?.message || 'Delete failed'
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

  async function exportStyle() {
    if (!selectedName) return
    error = ''
    try {
      await ExportStyle(selectedName, selectedName + '.shipstyle')
      success = 'Exported.'
      setTimeout(() => success = '', 3000)
    } catch (e: any) {
      error = e?.message || 'Export failed'
    }
  }

  let vocabKey = ''
  let vocabVal = ''
  function addVocab() {
    if (!vocabKey.trim()) return
    if (!detail.vocabulary) detail.vocabulary = {}
    detail.vocabulary[vocabKey.trim()] = vocabVal.trim()
    detail = detail
    vocabKey = ''
    vocabVal = ''
  }
  function removeVocab(key: string) {
    delete detail.vocabulary[key]
    detail = detail
  }

  const tabs = [
    { id: 'colors', label: 'Colors' },
    { id: 'typography', label: 'Typography' },
    { id: 'images', label: 'Images' },
    { id: 'templates', label: 'Templates' },
    { id: 'advanced', label: 'Advanced' },
  ]

  const colorFields = ['primary', 'secondary', 'background', 'surface', 'text', 'accent', 'border'] as const
  const templateFields = [
    { key: 'header', label: 'Header', hint: '{{.Version}}, {{.Date}}' },
    { key: 'feature', label: 'Feature', hint: '{{.Message}}, {{.Hash}}, {{.Author}}' },
    { key: 'fix', label: 'Fix', hint: '{{.Message}}, {{.Hash}}, {{.Author}}' },
    { key: 'breaking', label: 'Breaking', hint: '{{.Message}}, {{.Hash}}' },
    { key: 'footer', label: 'Footer', hint: '{{.Date}}' },
    { key: 'storyIntro', label: 'Story Intro', hint: '{{.FirstAuthor}}, {{.TotalCommits}}' },
    { key: 'storyMilestone', label: 'Story Milestone', hint: '{{.Tag}}, {{.Date}}' },
    { key: 'storyPeak', label: 'Story Peak', hint: '{{.Month}}, {{.Count}}' },
    { key: 'storyContributor', label: 'Story Contributor', hint: '{{.Name}}, {{.Date}}' },
    { key: 'storyFooter', label: 'Story Footer', hint: '' },
  ]
</script>

<div class="two-col">
  <div class="left-col">
    <div class="left-header">
      <h1>Styles</h1>
      <button class="small-action" on:click={importStyle} title="Import .shipstyle file">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
      </button>
    </div>
    <div class="style-cards">
      {#each styles as s}
        <button
          class="style-card"
          class:selected={selectedName === s.name && !creating}
          class:is-active={currentActiveStyle === s.name}
          on:click={() => selectStyle(s.name)}
        >
          <div class="card-top">
            <span class="card-name">{s.name}</span>
            <span class="card-badge" class:builtin={s.builtIn}>{s.builtIn ? 'built-in' : 'user'}</span>
          </div>
          <div class="card-swatches">
            <span class="mini-swatch" style="background:{s.primary}"></span>
            <span class="mini-swatch" style="background:{s.accent}"></span>
            <span class="mini-swatch" style="background:{s.background}"></span>
          </div>
        </button>
      {/each}
    </div>
    <button class="new-style-btn" on:click={startCreate}>
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14"><path d="M12 5v14M5 12h14"/></svg>
      New style
    </button>
  </div>

  <div class="right-col">
    {#if error}
      <div class="banner error">{error}</div>
    {/if}
    {#if success}
      <div class="banner success">{success}</div>
    {/if}

    {#if !detail && !creating}
      <div class="empty-editor">
        <p>Select a style to view or edit</p>
      </div>
    {:else if detail}
      {#if isBuiltIn}
        <div class="banner info">Built-in styles cannot be modified</div>
      {/if}

      {#if creating}
        <div class="name-field">
          <label>Name</label>
          <input type="text" bind:value={newName} placeholder="my-style-name" />
        </div>
      {:else}
        <div class="editor-header">{detail.name}</div>
      {/if}

      <div class="tab-bar">
        {#each tabs as tab}
          <button class="tab" class:active={activeTab === tab.id} on:click={() => activeTab = tab.id}>{tab.label}</button>
        {/each}
      </div>

      <div class="tab-content">
        {#if activeTab === 'colors'}
          <div class="field-group">
            <label>Mode</label>
            <div class="radio-group">
              <label><input type="radio" bind:group={detail.theme.mode} value="dark" disabled={isBuiltIn} /> Dark</label>
              <label><input type="radio" bind:group={detail.theme.mode} value="light" disabled={isBuiltIn} /> Light</label>
            </div>
          </div>
          {#each colorFields as cf}
            <div class="color-field">
              <label>{cf}</label>
              <div class="color-input-row">
                <input type="color" value={detail.theme.colors[cf] || '#000000'} on:input={(e) => { detail.theme.colors[cf] = e.currentTarget.value; detail = detail }} disabled={isBuiltIn} />
                <input type="text" bind:value={detail.theme.colors[cf]} disabled={isBuiltIn} class="hex-input" />
              </div>
            </div>
          {/each}

        {:else if activeTab === 'typography'}
          <div class="field-group">
            <label>Font Family</label>
            <input type="text" bind:value={detail.theme.typography.fontFamily} disabled={isBuiltIn} />
            <div class="font-preview" style="font-family:{detail.theme.typography.fontFamily || 'sans-serif'}">The quick brown fox jumps over the lazy dog</div>
          </div>
          <div class="field-group">
            <label>Base Font Size</label>
            <input type="text" bind:value={detail.theme.typography.fontSizeBase} disabled={isBuiltIn} placeholder="14px" />
          </div>
          <div class="field-group">
            <label>Header Font Size</label>
            <input type="text" bind:value={detail.theme.typography.fontSizeHeader} disabled={isBuiltIn} placeholder="24px" />
          </div>
          <div class="field-group">
            <label>Code Font Size</label>
            <input type="text" bind:value={detail.theme.typography.fontSizeCode} disabled={isBuiltIn} placeholder="13px" />
          </div>

        {:else if activeTab === 'images'}
          <div class="field-group">
            <label>Logo</label>
            {#if detail.theme.logo}
              <img src={detail.theme.logo} alt="logo" class="img-preview-sm" />
            {/if}
            <input type="text" bind:value={detail.theme.logo} disabled={isBuiltIn} placeholder="URL or base64" />
            {#if !isBuiltIn}
              <button class="small-action" on:click={() => { detail.theme.logo = ''; detail = detail }}>Clear</button>
            {/if}
          </div>
          <div class="field-group">
            <label>Header Image</label>
            {#if detail.theme.headerImage}
              <img src={detail.theme.headerImage} alt="header" class="img-preview-wide" />
            {/if}
            <input type="text" bind:value={detail.theme.headerImage} disabled={isBuiltIn} placeholder="URL or base64" />
            {#if !isBuiltIn}
              <button class="small-action" on:click={() => { detail.theme.headerImage = ''; detail = detail }}>Clear</button>
            {/if}
          </div>
          <div class="field-group">
            <label>Card Style</label>
            <div class="radio-group">
              <label><input type="radio" bind:group={detail.theme.cardStyle} value="minimal" disabled={isBuiltIn} /> Minimal</label>
              <label><input type="radio" bind:group={detail.theme.cardStyle} value="bordered" disabled={isBuiltIn} /> Bordered</label>
              <label><input type="radio" bind:group={detail.theme.cardStyle} value="glassmorphism" disabled={isBuiltIn} /> Glassmorphism</label>
            </div>
          </div>
          <div class="field-group">
            <label>Animations</label>
            <label class="toggle-label">
              <input type="checkbox" bind:checked={detail.theme.animations} disabled={isBuiltIn} />
              {detail.theme.animations ? 'On' : 'Off'}
            </label>
          </div>
          <div class="field-group">
            <label>Window Controls</label>
            <div class="inline-fields">
              {#each [['default', 'Default'], ['close', 'Close'], ['minimize', 'Minimize'], ['maximize', 'Maximize']] as [key, label]}
                <div class="mini-field">
                  <span class="mini-label">{label}</span>
                  <div class="color-input-row">
                    <input type="color" value={detail.theme.windowControls?.[key] || '#666666'} on:input={(e) => { if (!detail.theme.windowControls) detail.theme.windowControls = {}; detail.theme.windowControls[key] = e.currentTarget.value; detail = detail }} disabled={isBuiltIn} />
                    <input type="text" value={detail.theme.windowControls?.[key] || ''} on:input={(e) => { if (!detail.theme.windowControls) detail.theme.windowControls = {}; detail.theme.windowControls[key] = e.currentTarget.value; detail = detail }} disabled={isBuiltIn} class="hex-input" />
                  </div>
                </div>
              {/each}
            </div>
          </div>

        {:else if activeTab === 'templates'}
          {#each templateFields as tf}
            <div class="field-group">
              <label>{tf.label} {#if tf.hint}<span class="hint">{tf.hint}</span>{/if}</label>
              <textarea bind:value={detail.templates[tf.key]} disabled={isBuiltIn} rows="3" class="mono"></textarea>
            </div>
          {/each}
          <div class="field-group">
            <label>Vocabulary</label>
            <div class="vocab-table">
              {#if detail.vocabulary}
                {#each Object.entries(detail.vocabulary) as [k, v]}
                  <div class="vocab-row">
                    <span class="vocab-key">{k}</span>
                    <span class="vocab-val">{v}</span>
                    {#if !isBuiltIn}
                      <button class="vocab-del" on:click={() => removeVocab(k)}>x</button>
                    {/if}
                  </div>
                {/each}
              {/if}
              {#if !isBuiltIn}
                <div class="vocab-add">
                  <input type="text" bind:value={vocabKey} placeholder="key" />
                  <input type="text" bind:value={vocabVal} placeholder="value" />
                  <button class="small-action" on:click={addVocab}>Add</button>
                </div>
              {/if}
            </div>
          </div>

        {:else if activeTab === 'advanced'}
          <div class="field-group">
            <label>Terminal Colors (ANSI names)</label>
            <div class="inline-fields">
              {#each ['header', 'feature', 'fix', 'breaking', 'footer'] as tc}
                <div class="mini-field">
                  <span class="mini-label">{tc}</span>
                  <input type="text" bind:value={detail.terminal.colors[tc]} disabled={isBuiltIn} placeholder="e.g. cyan" />
                </div>
              {/each}
            </div>
          </div>
          <div class="field-group">
            <label>Terminal Decorators</label>
            <div class="inline-fields">
              <div class="mini-field">
                <span class="mini-label">separator</span>
                <input type="text" bind:value={detail.terminal.decorators.separator} disabled={isBuiltIn} />
              </div>
              <div class="mini-field">
                <span class="mini-label">bullet</span>
                <input type="text" bind:value={detail.terminal.decorators.bullet} disabled={isBuiltIn} />
              </div>
              <div class="mini-field">
                <span class="mini-label">indent</span>
                <input type="text" bind:value={detail.terminal.decorators.indent} disabled={isBuiltIn} />
              </div>
            </div>
          </div>
          <div class="field-group">
            <label>Terminal Density</label>
            <div class="radio-group">
              <label><input type="radio" bind:group={detail.terminal.density} value="compact" disabled={isBuiltIn} /> Compact</label>
              <label><input type="radio" bind:group={detail.terminal.density} value="normal" disabled={isBuiltIn} /> Normal</label>
              <label><input type="radio" bind:group={detail.terminal.density} value="verbose" disabled={isBuiltIn} /> Verbose</label>
            </div>
          </div>
          <div class="field-group">
            <label>Custom CSS</label>
            <textarea bind:value={detail.theme.customCss} disabled={isBuiltIn} rows="6" class="mono" placeholder="Additional CSS injected at end of HTML reports"></textarea>
          </div>
          <div class="field-group">
            <label>LLM Prompt
              <span class="hint-warn">Injected into LLM calls. Keep it focused on tone and format.</span>
            </label>
            <textarea bind:value={detail.llmPrompt} disabled={isBuiltIn} rows="6" class="mono"></textarea>
          </div>
        {/if}
      </div>

      {#if !isBuiltIn}
        <div class="editor-actions">
          <button class="action-btn" on:click={saveStyle}>Save</button>
          {#if !creating && selectedName}
            <button class="small-action" on:click={exportStyle}>Export</button>
            {#if !confirmDelete}
              <button class="delete-btn" on:click={() => confirmDelete = true}>Delete</button>
            {:else}
              <span class="confirm-text">Delete "{selectedName}"?</span>
              <button class="delete-btn" on:click={doDelete}>Confirm</button>
              <button class="small-action" on:click={() => confirmDelete = false}>Cancel</button>
            {/if}
          {/if}
        </div>
      {/if}
    {/if}
  </div>
</div>

<style>
  .two-col { display: flex; gap: 0; height: 100%; }

  .left-col {
    width: 260px; flex-shrink: 0; display: flex; flex-direction: column;
    border-right: 1px solid var(--cl-border, #30363d);
  }
  .left-header {
    display: flex; justify-content: space-between; align-items: center;
    padding: 0 12px 10px 0;
  }
  h1 { color: var(--cl-text, #e6edf3); font-size: 20px; margin: 0; }

  .style-cards { flex: 1; overflow-y: auto; display: flex; flex-direction: column; gap: 2px; padding-right: 12px; }

  .style-card {
    display: flex; flex-direction: column; gap: 4px; padding: 8px 10px;
    background: transparent; border: none; border-left: 3px solid transparent;
    border-radius: 0; color: var(--cl-text, #e6edf3); cursor: pointer;
    text-align: left; width: 100%; font-family: inherit; transition: background 0.12s;
  }
  .style-card:hover { background: var(--cl-surface, #161b22); }
  .style-card.selected { background: color-mix(in srgb, var(--cl-accent, #58a6ff) 10%, transparent); }
  .style-card.is-active { border-left-color: var(--cl-accent, #58a6ff); }

  .card-top { display: flex; align-items: center; gap: 6px; }
  .card-name { font-size: 13px; font-weight: 500; }
  .card-badge {
    font-size: 9px; padding: 1px 5px; border-radius: 8px;
    background: var(--cl-border, #30363d); color: var(--cl-secondary, #8b949e);
  }
  .card-badge.builtin { background: color-mix(in srgb, var(--cl-accent, #58a6ff) 20%, transparent); color: var(--cl-accent, #58a6ff); }

  .card-swatches { display: flex; gap: 4px; }
  .mini-swatch { width: 12px; height: 12px; border-radius: 50%; border: 1px solid var(--cl-border, #30363d); }

  .new-style-btn {
    display: flex; align-items: center; justify-content: center; gap: 6px;
    padding: 10px; margin: 8px 12px 0 0; background: transparent;
    border: 1px dashed var(--cl-border, #30363d); border-radius: 6px;
    color: var(--cl-secondary, #8b949e); cursor: pointer; font-size: 13px; font-family: inherit;
    flex-shrink: 0;
  }
  .new-style-btn:hover { border-color: var(--cl-accent, #58a6ff); color: var(--cl-accent, #58a6ff); }

  .right-col {
    flex: 1; display: flex; flex-direction: column; gap: 10px;
    padding: 0 0 0 16px; overflow-y: auto; min-width: 0;
  }

  .empty-editor {
    flex: 1; display: flex; align-items: center; justify-content: center;
    color: var(--cl-secondary, #8b949e); font-size: 14px;
  }
  .empty-editor p { margin: 0; }

  .editor-header { font-size: 18px; font-weight: 600; color: var(--cl-text, #e6edf3); }

  .name-field { display: flex; flex-direction: column; gap: 4px; }
  .name-field label { color: var(--cl-secondary, #8b949e); font-size: 11px; text-transform: uppercase; }
  .name-field input {
    padding: 7px 10px; background: var(--cl-background, #0d1117); border: 1px solid var(--cl-border, #30363d);
    border-radius: 6px; color: var(--cl-text, #e6edf3); font-size: 14px; font-family: 'JetBrains Mono', monospace;
    max-width: 300px;
  }
  .name-field input:focus { outline: none; border-color: var(--cl-accent, #58a6ff); }

  .tab-bar { display: flex; gap: 0; border-bottom: 1px solid var(--cl-border, #30363d); flex-shrink: 0; }
  .tab {
    padding: 8px 14px; background: transparent; border: none; border-bottom: 2px solid transparent;
    color: var(--cl-secondary, #8b949e); cursor: pointer; font-size: 12px; font-family: inherit;
    transition: color 0.12s;
  }
  .tab:hover { color: var(--cl-text, #e6edf3); }
  .tab.active { color: var(--cl-accent, #58a6ff); border-bottom-color: var(--cl-accent, #58a6ff); }

  .tab-content { flex: 1; overflow-y: auto; display: flex; flex-direction: column; gap: 10px; padding: 4px 0; }

  .field-group { display: flex; flex-direction: column; gap: 4px; }
  .field-group > label { color: var(--cl-secondary, #8b949e); font-size: 11px; text-transform: uppercase; display: flex; gap: 6px; align-items: center; }

  .hint { font-size: 10px; color: var(--cl-secondary, #8b949e); opacity: 0.7; text-transform: none; }
  .hint-warn { font-size: 10px; color: #D97706; text-transform: none; }

  input[type="text"], textarea, select {
    padding: 6px 8px; background: var(--cl-background, #0d1117); border: 1px solid var(--cl-border, #30363d);
    border-radius: 4px; color: var(--cl-text, #e6edf3); font-size: 12px; font-family: inherit;
  }
  input[type="text"]:focus, textarea:focus { outline: none; border-color: var(--cl-accent, #58a6ff); }
  input:disabled, textarea:disabled, select:disabled { opacity: 0.6; cursor: not-allowed; }
  .mono { font-family: 'JetBrains Mono', monospace; }

  .color-field { display: flex; align-items: center; gap: 8px; }
  .color-field label { width: 80px; color: var(--cl-secondary, #8b949e); font-size: 11px; text-transform: uppercase; flex-shrink: 0; }
  .color-input-row { display: flex; align-items: center; gap: 6px; }
  input[type="color"] { width: 28px; height: 28px; padding: 0; border: 1px solid var(--cl-border, #30363d); border-radius: 4px; cursor: pointer; background: none; }
  input[type="color"]:disabled { cursor: not-allowed; }
  .hex-input { width: 80px; font-family: 'JetBrains Mono', monospace; }

  .radio-group { display: flex; gap: 12px; font-size: 12px; color: var(--cl-text, #e6edf3); }
  .radio-group label { display: flex; align-items: center; gap: 4px; cursor: pointer; }
  .toggle-label { display: flex; align-items: center; gap: 6px; font-size: 12px; color: var(--cl-text, #e6edf3); cursor: pointer; }

  .font-preview {
    padding: 8px; background: var(--cl-surface, #161b22); border: 1px solid var(--cl-border, #30363d);
    border-radius: 4px; font-size: 13px; color: var(--cl-text, #e6edf3);
  }

  .img-preview-sm { width: 64px; height: 64px; object-fit: contain; border: 1px solid var(--cl-border, #30363d); border-radius: 4px; }
  .img-preview-wide { width: 200px; height: 64px; object-fit: cover; border: 1px solid var(--cl-border, #30363d); border-radius: 4px; }

  .inline-fields { display: flex; gap: 8px; flex-wrap: wrap; }
  .mini-field { display: flex; flex-direction: column; gap: 2px; }
  .mini-label { font-size: 10px; color: var(--cl-secondary, #8b949e); }
  .mini-field input { width: 90px; }

  .vocab-table { display: flex; flex-direction: column; gap: 4px; }
  .vocab-row { display: flex; align-items: center; gap: 8px; font-size: 12px; font-family: 'JetBrains Mono', monospace; }
  .vocab-key { color: var(--cl-accent, #58a6ff); min-width: 80px; }
  .vocab-val { color: var(--cl-text, #e6edf3); flex: 1; }
  .vocab-del { background: none; border: none; color: #f85149; cursor: pointer; font-size: 14px; padding: 0 4px; }
  .vocab-add { display: flex; gap: 6px; align-items: center; margin-top: 4px; }
  .vocab-add input { width: 100px; }

  .editor-actions {
    display: flex; align-items: center; gap: 8px; padding: 8px 0; border-top: 1px solid var(--cl-border, #30363d);
    flex-shrink: 0;
  }
  .action-btn {
    padding: 8px 16px; background: #238636; border: none; border-radius: 6px;
    color: #fff; font-size: 13px; cursor: pointer; font-family: inherit;
  }
  .action-btn:hover { background: #2ea043; }

  .small-action {
    padding: 4px 10px; background: transparent; border: 1px solid var(--cl-border, #30363d);
    border-radius: 4px; color: var(--cl-secondary, #8b949e); cursor: pointer; font-size: 12px; font-family: inherit;
    display: flex; align-items: center; gap: 4px;
  }
  .small-action:hover { border-color: var(--cl-accent, #58a6ff); color: var(--cl-text, #e6edf3); }

  .delete-btn {
    padding: 4px 10px; background: transparent; border: 1px solid #da3634;
    border-radius: 4px; color: #f85149; cursor: pointer; font-size: 12px; font-family: inherit;
  }
  .delete-btn:hover { background: #da363422; }

  .confirm-text { font-size: 12px; color: #f85149; }

  .banner { padding: 6px 10px; border-radius: 4px; font-size: 12px; flex-shrink: 0; }
  .banner.error { background: #da363433; border: 1px solid #da3634; color: #f85149; }
  .banner.success { background: #23863633; border: 1px solid #238636; color: #3fb950; }
  .banner.info { background: color-mix(in srgb, var(--cl-accent, #58a6ff) 10%, transparent); border: 1px solid var(--cl-accent, #58a6ff); color: var(--cl-accent, #58a6ff); }
</style>
