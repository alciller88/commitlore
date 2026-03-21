package renderer

import (
	"fmt"
	"strings"

	"github.com/alciller88/commitlore/internal/styles"
)

func buildCSS(theme styles.Theme) string {
	c := withDefaults(theme)
	var buf strings.Builder
	writeBaseCSS(&buf, c)
	writeTypeBadgeCSS(&buf, c)
	writeComponentCSS(&buf, c)
	writeCardStyleCSS(&buf, c)
	writeCustomCSS(&buf, theme)
	return buf.String()
}

func writeCustomCSS(buf *strings.Builder, theme styles.Theme) {
	if theme.CustomCSS != "" {
		buf.WriteString(theme.CustomCSS)
	}
	if !theme.Animations {
		buf.WriteString("* { animation: none !important; transition: none !important; }\n")
	}
}

func withDefaults(theme styles.Theme) styles.Theme {
	d := theme
	def(&d.Colors.Background, "#0d1117")
	def(&d.Colors.Text, "#c9d1d9")
	def(&d.Colors.Primary, "#58a6ff")
	def(&d.Colors.Secondary, "#79c0ff")
	def(&d.Colors.Surface, "#161b22")
	def(&d.Colors.Border, "#30363d")
	def(&d.Colors.Accent, "#1f6feb")
	def(&d.Typography.FontFamily, "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif")
	def(&d.Typography.FontSizeBase, "14px")
	def(&d.Typography.FontSizeH, "24px")
	def(&d.Typography.FontSizeCode, "12px")
	return d
}

func def(field *string, fallback string) {
	if *field == "" {
		*field = fallback
	}
}

func writeBaseCSS(buf *strings.Builder, t styles.Theme) {
	fmt.Fprintf(buf,
		"body { font-family: %s; font-size: %s; margin: 0; padding: 0; background: %s; color: %s; }\n",
		t.Typography.FontFamily, t.Typography.FontSizeBase, t.Colors.Background, t.Colors.Text)
	fmt.Fprintf(buf, ".content { max-width: 800px; margin: 0 auto; padding: 2rem; }\n")
	writeSiteHeaderCSS(buf, t)
	fmt.Fprintf(buf,
		"h1 { color: %s; border-bottom: 1px solid %s; padding-bottom: 0.5rem; font-size: %s; }\n",
		t.Colors.Primary, t.Colors.Border, t.Typography.FontSizeH)
	fmt.Fprintf(buf, "h2 { color: %s; margin-top: 1.5rem; }\n", t.Colors.Secondary)
	fmt.Fprintf(buf, "ul { list-style: none; padding-left: 0; }\n")
	fmt.Fprintf(buf, "li { padding: 0.3rem 0; border-bottom: 1px solid %s; }\n", t.Colors.Surface)
}

func writeSiteHeaderCSS(buf *strings.Builder, t styles.Theme) {
	fmt.Fprintf(buf,
		".site-header { width: 100%%; padding: 1rem 2rem; box-sizing: border-box; display: flex; align-items: center; margin-bottom: 1.5rem; border-bottom: 1px solid %s; }\n",
		t.Colors.Border)
	buf.WriteString(".site-header .logo-svg { display: inline-block; width: 100px; height: 100px; flex-shrink: 0; }\n")
	buf.WriteString(".site-header .logo-svg svg { width: 100px !important; height: 100px !important; }\n")
}

func writeTypeBadgeCSS(buf *strings.Builder, t styles.Theme) {
	buf.WriteString(".type-badge { display: inline-block; padding: 0.1rem 0.4rem; border-radius: 3px; font-size: 0.75em; font-weight: bold; margin-right: 0.3rem; }\n")
	fmt.Fprintf(buf, ".type-feat { background: %s33; color: %s; }\n", t.Colors.Accent, t.Colors.Primary)
	fmt.Fprintf(buf, ".type-fix { background: #3fb95033; color: #3fb950; }\n")
	fmt.Fprintf(buf, ".type-breaking { background: #f8514933; color: #f85149; }\n")
	fmt.Fprintf(buf, ".type-other { background: %s; color: %s; }\n", t.Colors.Surface, t.Colors.Text)
}

func writeComponentCSS(buf *strings.Builder, t styles.Theme) {
	fmt.Fprintf(buf, ".hash { color: %s; font-family: monospace; font-size: %s; }\n", t.Colors.Secondary, t.Typography.FontSizeCode)
	fmt.Fprintf(buf, ".author { color: #3fb950; }\n")
	fmt.Fprintf(buf, ".date { color: %s; font-size: 0.85em; }\n", t.Colors.Text)
	fmt.Fprintf(buf, ".footer { margin-top: 2rem; padding-top: 1rem; border-top: 1px solid %s; color: %s; font-size: 0.85em; }\n", t.Colors.Border, t.Colors.Text)
	fmt.Fprintf(buf, ".peak-bar { display: inline-block; background: %s; height: 0.8em; border-radius: 2px; margin-right: 0.5rem; vertical-align: middle; }\n", t.Colors.Accent)
	buf.WriteString(".logo { height: 100px; width: 100px; margin-right: 1rem; }\n")
	buf.WriteString(".header-image { width: 100%%; max-height: 200px; object-fit: cover; margin-bottom: 1rem; }\n")
	fmt.Fprintf(buf, ".narrative { margin-bottom: 2rem; line-height: 1.7; }\n")
	fmt.Fprintf(buf, ".narrative p { margin: 0.6rem 0; }\n")
	fmt.Fprintf(buf, ".data-section { border-top: 1px solid %s; padding-top: 1.5rem; margin-top: 1rem; }\n", t.Colors.Border)
}

func writeCardStyleCSS(buf *strings.Builder, t styles.Theme) {
	switch t.CardStyle {
	case "bordered":
		fmt.Fprintf(buf, "li { border: 1px solid %s; border-radius: 4px; padding: 0.5rem; margin-bottom: 0.3rem; }\n", t.Colors.Border)
	case "glassmorphism":
		fmt.Fprintf(buf, "li { background: %s80; backdrop-filter: blur(10px); border-radius: 8px; padding: 0.5rem; margin-bottom: 0.3rem; border: 1px solid %s40; }\n", t.Colors.Surface, t.Colors.Border)
	}
}
