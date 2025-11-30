// bar_top.go – FINAL, Fyne 2.7.1 compatible, fixed GlowCypherButton
package main

import (
    "image/color"
    "math"
    "time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"
)

var themeSelector *widget.Select

// === ALL CUSTOM THEMES ===
type draculaTheme struct{ customTheme }
type nordTheme struct{ customTheme }
type oledTheme struct{ customTheme }
type gruvboxTheme struct{ customTheme }
type cypherpunkTheme struct{ customTheme }

// Dracula
func (draculaTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
    switch name {
    case theme.ColorNameBackground: return color.RGBA{40, 42, 54, 255}
    case theme.ColorNameForeground: return color.RGBA{248, 248, 242, 255}
    case theme.ColorNamePrimary:    return color.RGBA{189, 147, 249, 255}
    default:
        return theme.DarkTheme().Color(name, theme.VariantDark)
    }
}

// Nord
func (nordTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
    switch name {
    case theme.ColorNameBackground: return color.RGBA{46, 52, 64, 255}
    case theme.ColorNameForeground: return color.RGBA{216, 222, 233, 255}
    case theme.ColorNamePrimary:    return color.RGBA{136, 192, 208, 255}
    default:
        return theme.DarkTheme().Color(name, theme.VariantDark)
    }
}

// OLED Black
func (oledTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
    if name == theme.ColorNameBackground {
        return color.RGBA{0, 0, 0, 255}
    }
    return theme.DarkTheme().Color(name, theme.VariantDark)
}

// Gruvbox Dark
func (gruvboxTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
    switch name {
    case theme.ColorNameBackground: return color.RGBA{40, 40, 40, 255}
    case theme.ColorNameForeground: return color.RGBA{235, 219, 178, 255}
    case theme.ColorNamePrimary:    return color.RGBA{215, 153, 85, 255}
    default:
        return theme.DarkTheme().Color(name, theme.VariantDark)
    }
}

// CYPHERPUNK – toxic green on pure black
func (cypherpunkTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
    switch name {
    case theme.ColorNameBackground:
        return color.RGBA{0, 0, 0, 255} // black background
    case theme.ColorNameOverlayBackground:
        return color.RGBA{18, 18, 18, 255} // black background
    case theme.ColorNameInputBackground:
        return color.RGBA{18, 18, 18, 255} // background of dropdown field + menu (need to be tested)   
    case theme.ColorNameMenuBackground:
        return color.RGBA{36, 36, 36, 255} // background of the dropdown popup list
    case theme.ColorNameForeground:
        return color.RGBA{0, 255, 70, 255} // neon green labels/text
    case theme.ColorNamePrimary:
        return color.RGBA{18, 18, 18, 255} // highlights
    case theme.ColorNameFocus:
        return color.RGBA{18, 18, 18, 255} // color on hover
    case theme.ColorNameButton:
        return color.RGBA{18, 18, 18, 255} // green button background
    case theme.ColorNameHover:
        return color.RGBA{36, 36, 36, 255} // background hover of buttons
    case theme.ColorNameDisabled:
        return color.RGBA{0, 255, 70, 255}
    default:
        return theme.DarkTheme().Color(name, theme.VariantDark)
    }
}

// === Custom GlowCypherButton with pulsing neon border ===
type GlowCypherButton struct {
    widget.Button
    glowRect *canvas.Rectangle
    ticker   *time.Ticker
    alpha    float64
}

func NewGlowCypherButton(label string, tapped func()) *GlowCypherButton {
    b := &GlowCypherButton{}
    b.ExtendBaseWidget(b)
    b.Text = label
    b.OnTapped = tapped

    b.glowRect = canvas.NewRectangle(color.RGBA{0, 255, 80, 0})
    b.glowRect.FillColor = color.Transparent
    b.glowRect.StrokeWidth = 6
    b.glowRect.StrokeColor = color.RGBA{0, 255, 80, 255}
    b.glowRect.Hide()

    return b
}

func (b *GlowCypherButton) CreateRenderer() fyne.WidgetRenderer {
    base := widget.NewButton(b.Text, b.OnTapped)
    baseRenderer := base.CreateRenderer()

    // Set button text black
    if txt, ok := baseRenderer.Objects()[1].(*canvas.Text); ok {
        txt.Color = color.Black
    }

    stack := container.NewMax(base, b.glowRect)
    b.startGlowAnimation()

    return &glowButtonRenderer{
        base:     baseRenderer,
        glowRect: b.glowRect,
        stack:    stack,
    }
}

// Fixed renderer struct
type glowButtonRenderer struct {
    base     fyne.WidgetRenderer
    glowRect *canvas.Rectangle
    stack    *fyne.Container
}

func (r *glowButtonRenderer) Layout(size fyne.Size) {
    r.base.Layout(size)
    r.glowRect.Resize(size)
}

func (r *glowButtonRenderer) MinSize() fyne.Size {
    return r.base.MinSize()
}

func (r *glowButtonRenderer) Objects() []fyne.CanvasObject {
    return r.stack.Objects
}

func (r *glowButtonRenderer) Refresh() {
    r.base.Refresh()
    canvas.Refresh(r.stack)
}

//  Add this to implement fyne.WidgetRenderer
func (r *glowButtonRenderer) Destroy() {}

// Glow animation
func (b *GlowCypherButton) startGlowAnimation() {
    if b.ticker != nil {
        return
    }

    b.ticker = time.NewTicker(time.Millisecond * 50)
    go func() {
        for range b.ticker.C {
            b.alpha += 0.05
            if b.alpha > math.Pi*2 {
                b.alpha = 0
            }
            a := uint8((math.Sin(b.alpha)+1)/2*100 + 50) // 50–150 alpha
            b.glowRect.StrokeColor = color.RGBA{0, 255, 80, a}
            canvas.Refresh(b.glowRect)
        }
    }()
}

func (b *GlowCypherButton) StopGlow() {
    if b.ticker != nil {
        b.ticker.Stop()
        b.ticker = nil
    }
}

// === Top bar with theme selector ===
func topbar() *fyne.Container {
    if themeSelector == nil {
        options := []string{"Light", "Dark", "System", "Dracula", "Nord", "OLED Black", "Gruvbox", "Cypherpunk"}

        themeSelector = widget.NewSelect(options, func(selected string) {
            program.preferences.SetString("AppTheme", selected)

            switch selected {
            case "Light":
                program.application.Settings().SetTheme(theme.LightTheme())
            case "Dark":
                program.application.Settings().SetTheme(theme.DarkTheme())
            case "System":
                program.application.Settings().SetTheme(theme.DefaultTheme())
            case "Dracula":
                program.application.Settings().SetTheme(&draculaTheme{})
            case "Nord":
                program.application.Settings().SetTheme(&nordTheme{})
            case "OLED Black":
                program.application.Settings().SetTheme(&oledTheme{})
            case "Gruvbox":
                program.application.Settings().SetTheme(&gruvboxTheme{})
            case "Cypherpunk":
                program.application.Settings().SetTheme(&cypherpunkTheme{})
            }
        })

        // Apply saved theme on startup
        saved := program.preferences.StringWithFallback("AppTheme", "System")
        themeSelector.SetSelected(saved)

        switch saved {
        case "Light":
            program.application.Settings().SetTheme(theme.LightTheme())
        case "Dark":
            program.application.Settings().SetTheme(theme.DarkTheme())
        case "System":
            program.application.Settings().SetTheme(theme.DefaultTheme())
        case "Dracula":
            program.application.Settings().SetTheme(&draculaTheme{})
        case "Nord":
            program.application.Settings().SetTheme(&nordTheme{})
        case "OLED Black":
            program.application.Settings().SetTheme(&oledTheme{})
        case "Gruvbox":
            program.application.Settings().SetTheme(&gruvboxTheme{})
        case "Cypherpunk":
            program.application.Settings().SetTheme(&cypherpunkTheme{})
        }
    }

    rightSide := container.NewHBox(
        widget.NewSeparator(),
        widget.NewLabel("Theme:"),
        themeSelector,
    )

    return container.NewCenter(
        container.NewHBox(
            program.labels.height,
            program.labels.connection,
            program.labels.loggedin,
            program.labels.ws_server,
            program.labels.rpc_server,
            container.NewMax(),
            rightSide,
        ),
    )
}
