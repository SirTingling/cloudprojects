package tui

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	// "cloudprojects/current-converter/api"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConversionParams represents the conversion details entered by the user.
type ConversionParams struct {
	Amount       float64
	CurrencyFrom string
	CurrencyTo   string
}

// Item represents a currency option.
type Item struct {
	Code string
	Name string
}

func (i Item) Title() string       { return fmt.Sprintf("%s %s", currencySymbol(i.Code), i.Code) }
func (i Item) Description() string { return i.Name }
func (i Item) FilterValue() string { return i.Code }

func currencySymbol(code string) string {
	switch code {
	case "USD":
		return "$"
	case "GBP":
		return "£"
	case "EUR":
		return "€"
	case "JPY":
		return "¥"
	default:
		return ""
	}
}

type model struct {
	stage         int // Tracks the current question (0: base currency, 1: target currency, 2: amount)
	list          list.Model
	textInput     textinput.Model
	currencyFrom  string
	currencyTo    string
	amount        float64
	isCustomInput bool // Tracks whether the user is entering a custom currency
	finished      bool
	err           error
}

// Styling variables
var (
	questionStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#A020F0")).Bold(true)
	highlightStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	cursorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF69B4"))
	unselectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
)

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "enter":
			if m.isCustomInput {
				// Handle custom currency input
				input := strings.ToUpper(strings.TrimSpace(m.textInput.Value()))
				if !isAlphabetic(input) || len(input) != 3 {
					m.err = fmt.Errorf("invalid currency code: must be 3 letters")
					return m, nil
				}
				if m.stage == 0 {
					m.currencyFrom = input
					m.isCustomInput = false
					m.textInput.Reset()
					m.stage++
					m.list.ResetSelected()
					return m, nil
				} else if m.stage == 1 {
					m.currencyTo = input
					m.isCustomInput = false
					m.textInput.Reset()
					m.stage++
					m.textInput.Placeholder = "Enter amount (e.g., 100)"
					m.textInput.Focus()
					return m, nil
				}
			} else if m.stage == 2 {
				// Handle amount input
				amountStr := strings.TrimSpace(m.textInput.Value())
				amount, err := strconv.ParseFloat(amountStr, 64)
				if err != nil {
					m.err = fmt.Errorf("invalid amount: %v", err)
					return m, nil
				}
				m.amount = amount
				m.finished = true
				return m, tea.Quit
			} else {
				// Handle list selection
				selectedItem := m.list.SelectedItem().(Item)
				if selectedItem.Code == "OTHER" {
					m.isCustomInput = true
					m.textInput.Placeholder = "Enter currency code (e.g., USD)"
					m.textInput.Focus()
					return m, textinput.Blink
				}
				if m.stage == 0 {
					m.currencyFrom = selectedItem.Code
					m.stage++
					m.list.ResetSelected()
					return m, nil
				} else if m.stage == 1 {
					m.currencyTo = selectedItem.Code
					m.stage++
					m.textInput.Placeholder = "Enter amount (e.g., 100)"
					m.textInput.Focus()
					return m, nil
				}
			}
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	// Handle text input for custom currency or amount input
	if m.isCustomInput || m.stage == 2 {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	// Handle list updates for currency selection
	if m.stage == 0 || m.stage == 1 {
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress any key to continue.\n", m.err)
	}

	if m.finished {
		// Return an empty string when finished to avoid redundant output.
		return ""
	}

	switch m.stage {
	case 0:
		if m.isCustomInput {
			return questionStyle.Render("Enter your custom base currency code (e.g., USD):\n\n") + m.textInput.View()
		}
		return questionStyle.Render("What is your base currency?\n\n") + m.list.View()
	case 1:
		if m.isCustomInput {
			return questionStyle.Render("Enter your custom target currency code (e.g., EUR):\n\n") + m.textInput.View()
		}
		return questionStyle.Render("What do you want to convert to?\n\n") + m.list.View()
	case 2:
		return questionStyle.Render("How much to convert?\n\n") + m.textInput.View()
	default:
		return ""
	}
}

func RunTUI() (ConversionParams, error) {
	currencyList := []list.Item{
		Item{Code: "USD", Name: "United States Dollar"},
		Item{Code: "GBP", Name: "British Pound"},
		Item{Code: "EUR", Name: "Euro"},
		Item{Code: "JPY", Name: "Japanese Yen"},
		Item{Code: "OTHER", Name: "Type a custom currency"},
	}

	// Create the list model
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = highlightStyle
	delegate.Styles.NormalTitle = unselectedStyle
	delegate.Styles.SelectedDesc = highlightStyle
	delegate.Styles.NormalDesc = unselectedStyle

	listModel := list.New(currencyList, delegate, 30, 10)
	listModel.SetShowStatusBar(false)
	listModel.SetFilteringEnabled(false)
	listModel.DisableQuitKeybindings()
	listModel.SetShowHelp(false)

	// Text input for custom currency and amount
	textInput := textinput.New()
	textInput.CursorStyle = cursorStyle

	// Initialize the TUI model
	initialModel := &model{
		stage:     0,
		list:      listModel,
		textInput: textInput,
	}

	p := tea.NewProgram(initialModel)
	finalModel, err := p.Run()
	if err != nil {
		return ConversionParams{}, fmt.Errorf("error running TUI: %v", err)
	}

	fm := finalModel.(*model)

	if fm.err != nil {
		return ConversionParams{}, fm.err
	}

	return ConversionParams{
		Amount:       fm.amount,
		CurrencyFrom: fm.currencyFrom,
		CurrencyTo:   fm.currencyTo,
	}, nil
}

func isAlphabetic(input string) bool {
	for _, r := range input {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
