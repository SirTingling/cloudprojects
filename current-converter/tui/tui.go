package tui

import (
	"fmt"
	"strconv"

	"cloudprojects/current-converter/api"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

type model struct {
	inputs      []textinput.Model
	focused     int
	currencyMap map[string]bool
	done        bool
	err         error
}

type ConversionParams struct {
	Amount       float64
	CurrencyFrom string
	CurrencyTo   string
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.focused < len(m.inputs)-1 {
				m.focused++
				cmd := m.inputs[m.focused].Focus()
				m.inputs[m.focused-1].Blur()
				return m, cmd
			}

			// Validate inputs
			_, err := strconv.ParseFloat(m.inputs[0].Value(), 64)
			if err != nil {
				m.err = fmt.Errorf("invalid amount: %v", err)
				return m, nil
			}

			if !m.currencyMap[m.inputs[1].Value()] || !m.currencyMap[m.inputs[2].Value()] {
				m.err = fmt.Errorf("unsupported currency")
				return m, nil
			}

			m.done = true
			return m, tea.Quit
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress any key to exit.\n", m.err)
	}

	if m.done {
		return "Processing...\n"
	}

	var output string
	for i, input := range m.inputs {
		if i == m.focused {
			output += fmt.Sprintf("• %s\n", input.View())
		} else {
			output += fmt.Sprintf("  %s\n", input.View())
		}
	}
	output += "\nPress Enter to submit, or Ctrl+C to exit.\n"
	return output
}

func RunTUI(rates api.CurrencyData) (ConversionParams, error) {
	currencyMap := make(map[string]bool)
	for currency := range rates.Rates {
		currencyMap[currency] = true
	}

	inputs := make([]textinput.Model, 3)
	for i := range inputs {
		inputs[i] = textinput.New()
	}

	inputs[0].Placeholder = "Amount (e.g., 100)"
	inputs[1].Placeholder = "From Currency (e.g., USD)"
	inputs[2].Placeholder = "To Currency (e.g., EUR)"
	inputs[0].Focus()

	initialModel := model{
		inputs:      inputs,
		focused:     0,
		currencyMap: currencyMap,
	}

	p := tea.NewProgram(initialModel)
	finalModel, err := p.Run()
	if err != nil {
		return ConversionParams{}, fmt.Errorf("error running TUI: %v", err)
	}

	fm := finalModel.(model)
	amount, err := strconv.ParseFloat(fm.inputs[0].Value(), 64)
	if err != nil {
		return ConversionParams{}, fmt.Errorf("invalid amount input")
	}

	return ConversionParams{
		Amount:       amount,
		CurrencyFrom: fm.inputs[1].Value(),
		CurrencyTo:   fm.inputs[2].Value(),
	}, nil
}
