package ui

import "charm.land/lipgloss/v2"

const (
	maxResponseLineLength = 80

	colorTextPrimary    = "#b7d0f8"
	colorTextSecondary  = "#b7f8f4"
	colorTextOKBadge    = "#052E16"
	colorTextErrorBadge = "#450A0A"
	colorTextHTML       = "#0369A1"
	colorTextBackground = "#5b435a"

	colorBorderPrimary = "#8921d8"
	colorBorderHeaders = "#3b82f6"
	colorBorderBody    = "#10b981"
	colorBorderOK      = "#15803D"
	colorBorderError   = "#B91C1C"
)

var (
	okBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorTextOKBadge)).
			Background(lipgloss.Color("#86EFAC")).
			Padding(0, 1)

	errBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorTextErrorBadge)).
			Background(lipgloss.Color("#FCA5A5")).
			Padding(0, 1)

	metaLabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorTextPrimary))

	metaValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary))

	panelStyle = lipgloss.NewStyle().
			Margin(0, 0, 0, 0).
			Padding(0, 1)

	tableCellStyle = lipgloss.NewStyle().
			Padding(0, 1)

	tableOddRowStyle = tableCellStyle.
				Foreground(lipgloss.Color(colorTextPrimary))

	tableEvenRowStyle = tableCellStyle.
				Foreground(lipgloss.Color(colorTextSecondary))

	headersBoxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary)).
			Padding(0, 1)

	headersBlockStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(colorBorderHeaders))

	bodyTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorTextPrimary))

	bodyBlockStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorBorderBody))
)
