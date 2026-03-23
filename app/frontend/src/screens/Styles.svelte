<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { ListStyles, GetStyleDetail, DeleteStyle, ImportStyle, ExportStyle, GetStyleTheme } from '../../bindings/github.com/alciller88/commitlore/app/styleapp.js'
  import { SetActiveStyle } from '../../bindings/github.com/alciller88/commitlore/app/configapp.js'
  import { OpenFolderPicker } from '../../bindings/github.com/alciller88/commitlore/app/gitapp.js'
  import { Browser } from '@wailsio/runtime'
  import { activeStyle } from '../lib/store'

  type StyleListItem = { name: string; builtIn: boolean; primary: string; accent: string; background: string }

  let styles: StyleListItem[] = []
  let selectedName = ''
  let detail: any = null
  let isBuiltIn = false
  let error = ''
  let success = ''
  let loading = false
  let confirmDelete = false
  let currentActiveStyle = 'formal'
  let llmOpen = false

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
    selectedName = name
    error = ''
    confirmDelete = false
    llmOpen = false
    try {
      detail = await GetStyleDetail(name)
      isBuiltIn = styles.find(s => s.name === name)?.builtIn ?? false
    } catch (e: any) {
      error = e?.message || 'Failed to load style'
      detail = null
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

  async function setActive() {
    if (!selectedName) return
    error = ''
    try {
      activeStyle.set(selectedName)
      await SetActiveStyle(selectedName)
      success = 'Active style changed.'
      setTimeout(() => success = '', 3000)
    } catch (e: any) {
      error = e?.message || 'Failed to set active style'
    }
  }

  function openMarketplace() {
    Browser.OpenURL('https://commitlore.dev/styles')
  }

  const colorFields = ['primary', 'secondary', 'background', 'surface', 'text', 'accent', 'border'] as const

  const uiLabelKeys: [string, string][] = [
    ['dashboard', 'Dashboard'],
    ['generate', 'Generate'],
    ['generateButton', 'Generate button'],
    ['story', 'Story'],
    ['storyButton', 'Story button'],
    ['history', 'History'],
    ['contributors', 'Contributors'],
    ['styles', 'Styles'],
    ['settings', 'Settings'],
  ]

  const iconKeys: [string, string][] = [
    ['feature', 'feature'],
    ['fix', 'fix'],
    ['breaking', 'breaking'],
    ['chore', 'chore'],
    ['docs', 'docs'],
    ['test', 'test'],
    ['bullet', 'bullet'],
  ]

  function hasCustomLabels(d: any): boolean {
    if (!d?.uiLabels) return false
    return uiLabelKeys.some(([k]) => d.uiLabels[k] && d.uiLabels[k] !== '')
  }

  function hasCustomIcons(d: any): boolean {
    if (!d?.icons) return false
    return iconKeys.some(([k]) => d.icons[k] && d.icons[k] !== '')
  }
</script>

<div class="two-col">
  <div class="left-col">
    <div class="left-header">
      <h1>Styles</h1>
    </div>
    <div class="style-cards">
      {#each styles as s}
        <button
          class="style-card"
          class:selected={selectedName === s.name}
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
    <div class="left-actions">
      <button class="left-btn" on:click={importStyle}>
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
        Import style
      </button>
      <button class="left-btn" on:click={openMarketplace}>
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14"><circle cx="12" cy="12" r="10"/><path d="M2 12h20"/><path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z"/></svg>
        Get more styles
      </button>
    </div>
  </div>

  <div class="right-col">
    {#if error}
      <div class="banner error">{error}</div>
    {/if}
    {#if success}
      <div class="banner success">{success}</div>
    {/if}

    {#if !detail}
      <div class="empty-state">
        <p>Select a style to preview</p>
      </div>
    {:else}
      <div class="detail-panel">
        <!-- Header -->
        <div class="detail-header">
          {#if detail.theme?.logo}
            <div class="detail-logo">{@html detail.theme.logo}</div>
          {/if}
          <div class="detail-meta">
            <div class="detail-name-row">
              <span class="detail-name">{detail.name}</span>
              {#if detail.version}
                <span class="detail-version">v{detail.version}</span>
              {/if}
            </div>
            {#if detail.author}
              <div class="detail-author">by {detail.author}</div>
            {/if}
          </div>
        </div>

        {#if detail.description}
          <div class="detail-desc">{detail.description}</div>
        {/if}

        <div class="section-sep"></div>

        <!-- Theme -->
        <div class="section">
          <div class="section-title">Theme</div>
          <div class="section-row">
            <span class="section-label">Colors</span>
            <div class="color-circles">
              {#each colorFields as cf}
                <span
                  class="color-circle"
                  style="background:{detail.theme?.colors?.[cf] || '#333'}"
                  title="{cf}: {detail.theme?.colors?.[cf] || 'not set'}"
                ></span>
              {/each}
            </div>
          </div>
          {#if detail.theme?.typography?.fontFamily}
            <div class="section-row">
              <span class="section-label">Font</span>
              <span class="font-preview" style="font-family:{detail.theme.typography.fontFamily}">The quick brown fox</span>
            </div>
          {/if}
          <div class="section-row">
            <span class="section-label">Mode</span>
            <span class="mode-badge">{detail.theme?.mode || 'dark'}</span>
          </div>
        </div>

        <!-- UI Labels -->
        {#if hasCustomLabels(detail)}
          <div class="section-sep"></div>
          <div class="section">
            <div class="section-title">UI Labels</div>
            {#each uiLabelKeys as [key, label]}
              {#if detail.uiLabels?.[key]}
                <div class="label-row">
                  <span class="label-default">{label}</span>
                  <span class="label-arrow">&rarr;</span>
                  <span class="label-custom">{detail.uiLabels[key]}</span>
                </div>
              {/if}
            {/each}
          </div>
        {/if}

        <!-- Icons -->
        {#if hasCustomIcons(detail)}
          <div class="section-sep"></div>
          <div class="section">
            <div class="section-title">Icons</div>
            <div class="icon-row">
              {#each iconKeys as [key, label]}
                {#if detail.icons?.[key]}
                  <span class="icon-item">{label}: <span class="icon-char">{detail.icons[key]}</span></span>
                {/if}
              {/each}
            </div>
          </div>
        {/if}

        <!-- LLM Prompt -->
        {#if detail.llmPrompt}
          <div class="section-sep"></div>
          <div class="section">
            <button class="section-toggle" on:click={() => llmOpen = !llmOpen}>
              <span class="section-title">LLM Prompt</span>
              <span class="toggle-arrow" class:open={llmOpen}>&rsaquo;</span>
            </button>
            {#if llmOpen}
              <textarea class="llm-preview" readonly rows="8">{detail.llmPrompt}</textarea>
            {/if}
          </div>
        {/if}

        <!-- Actions -->
        <div class="section-sep"></div>
        <div class="actions">
          <button class="action-btn primary" on:click={setActive} disabled={currentActiveStyle === selectedName}>
            {currentActiveStyle === selectedName ? 'Active' : 'Set as active'}
          </button>
          <button class="action-btn outline" on:click={exportStyle}>Export</button>
          {#if !isBuiltIn}
            {#if !confirmDelete}
              <button class="action-btn danger" on:click={() => confirmDelete = true}>Delete</button>
            {:else}
              <span class="confirm-text">Delete "{selectedName}"?</span>
              <button class="action-btn danger" on:click={doDelete}>Confirm</button>
              <button class="action-btn outline" on:click={() => confirmDelete = false}>Cancel</button>
            {/if}
          {/if}
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .two-col { display: flex; gap: 0; flex: 1; min-height: 0; }

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
  .mini-swatch { width: 16px; height: 16px; border-radius: 50%; border: 1px solid var(--cl-border, #30363d); }

  .left-actions {
    display: flex; flex-direction: column; gap: 4px;
    padding: 10px 12px 0 0; flex-shrink: 0;
  }
  .left-btn {
    display: flex; align-items: center; justify-content: center; gap: 6px;
    padding: 8px; background: transparent;
    border: 1px solid var(--cl-border, #30363d); border-radius: 6px;
    color: var(--cl-secondary, #8b949e); cursor: pointer; font-size: 12px; font-family: inherit;
  }
  .left-btn:hover { border-color: var(--cl-accent, #58a6ff); color: var(--cl-accent, #58a6ff); }

  .right-col {
    flex: 1; display: flex; flex-direction: column; gap: 0;
    padding: 0 0 0 16px; overflow-y: auto; min-width: 0;
  }

  .empty-state {
    flex: 1; display: flex; align-items: center; justify-content: center;
    color: var(--cl-secondary, #8b949e); font-size: 14px;
  }
  .empty-state p { margin: 0; }

  .detail-panel { display: flex; flex-direction: column; gap: 0; padding-right: 8px; }

  .detail-header { display: flex; align-items: center; gap: 12px; }
  .detail-logo { width: 48px; height: 48px; flex-shrink: 0; color: var(--cl-primary, #58a6ff); }
  .detail-logo :global(svg) { width: 48px; height: 48px; }
  .detail-meta { display: flex; flex-direction: column; gap: 2px; }
  .detail-name-row { display: flex; align-items: baseline; gap: 8px; }
  .detail-name { font-size: 18px; font-weight: 600; color: var(--cl-text, #e6edf3); }
  .detail-version {
    font-size: 11px; padding: 1px 6px; border-radius: 8px;
    background: var(--cl-border, #30363d); color: var(--cl-secondary, #8b949e);
    font-family: 'JetBrains Mono', monospace;
  }
  .detail-author { font-size: 12px; color: var(--cl-secondary, #8b949e); }
  .detail-desc { font-size: 13px; color: var(--cl-text, #e6edf3); margin-top: 8px; line-height: 1.5; }

  .section-sep {
    height: 1px; background: var(--cl-border, #30363d); margin: 14px 0;
  }

  .section { display: flex; flex-direction: column; gap: 6px; }
  .section-title {
    font-size: 11px; font-weight: 600; color: var(--cl-secondary, #8b949e);
    text-transform: uppercase; letter-spacing: 0.06em;
  }

  .section-row { display: flex; align-items: center; gap: 10px; }
  .section-label { font-size: 12px; color: var(--cl-secondary, #8b949e); width: 55px; flex-shrink: 0; }

  .color-circles { display: flex; gap: 6px; }
  .color-circle {
    width: 20px; height: 20px; border-radius: 50%;
    border: 1px solid var(--cl-border, #30363d); cursor: default;
  }

  .font-preview {
    font-size: 13px; color: var(--cl-text, #e6edf3);
    white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  }

  .mode-badge {
    font-size: 11px; padding: 1px 8px; border-radius: 8px;
    background: var(--cl-surface, #161b22); color: var(--cl-text, #e6edf3);
    border: 1px solid var(--cl-border, #30363d);
  }

  .label-row { display: flex; align-items: center; gap: 6px; font-size: 12px; }
  .label-default { color: var(--cl-secondary, #8b949e); width: 110px; flex-shrink: 0; }
  .label-arrow { color: var(--cl-border, #30363d); }
  .label-custom { color: var(--cl-text, #e6edf3); }

  .icon-row { display: flex; flex-wrap: wrap; gap: 10px; }
  .icon-item { font-size: 12px; color: var(--cl-secondary, #8b949e); }
  .icon-char { font-size: 15px; }

  .section-toggle {
    display: flex; align-items: center; gap: 6px; background: none; border: none;
    padding: 0; cursor: pointer; color: inherit; font-family: inherit;
  }
  .toggle-arrow {
    font-size: 14px; color: var(--cl-secondary, #8b949e);
    transition: transform 0.15s; display: inline-block;
  }
  .toggle-arrow.open { transform: rotate(90deg); }

  .llm-preview {
    width: 100%; padding: 8px; background: var(--cl-background, #0d1117);
    border: 1px solid var(--cl-border, #30363d); border-radius: 4px;
    color: var(--cl-secondary, #8b949e); font-size: 11px; line-height: 1.5;
    font-family: 'JetBrains Mono', monospace; resize: none;
  }

  .actions {
    display: flex; align-items: center; gap: 8px; flex-wrap: wrap;
  }

  .action-btn {
    padding: 7px 14px; border-radius: 6px; font-size: 12px;
    cursor: pointer; font-family: inherit; border: none;
  }
  .action-btn.primary {
    background: #238636; color: #fff;
  }
  .action-btn.primary:hover { background: #2ea043; }
  .action-btn.primary:disabled {
    background: var(--cl-surface, #161b22); color: var(--cl-secondary, #8b949e);
    cursor: default; border: 1px solid var(--cl-border, #30363d);
  }
  .action-btn.outline {
    background: transparent; border: 1px solid var(--cl-border, #30363d);
    color: var(--cl-secondary, #8b949e);
  }
  .action-btn.outline:hover { border-color: var(--cl-accent, #58a6ff); color: var(--cl-text, #e6edf3); }
  .action-btn.danger {
    background: transparent; border: 1px solid #da3634;
    color: #f85149;
  }
  .action-btn.danger:hover { background: #da363422; }

  .confirm-text { font-size: 12px; color: #f85149; }

  .banner { padding: 6px 10px; border-radius: 4px; font-size: 12px; flex-shrink: 0; }
  .banner.error { background: #da363433; border: 1px solid #da3634; color: #f85149; }
  .banner.success { background: #23863633; border: 1px solid #238636; color: #3fb950; }
</style>
