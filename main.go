package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

var (
	term = termenv.EnvColorProfile()
	help = makeFgStyle("241")
	activeText = makeFgStyle("251")
	inactiveText = makeFgStyle("8")
	dotActive = colorFg(" ● ", "10")
	dotInactive = colorFg(" ○ ", "7")
	returnText =  makeFgBgStyle("251", "236")
)

type model struct {
    choices  []string           // items on the to-do list
    cursor   int                // which to-do list item our cursor is pointing at
    selected map[int]struct{}   // which to-do items are selected
    screen string
}

// Top row
var TOP_ROW_CHARS = []byte{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p'}

// Home row
var HOME_ROW_CHARS = []byte{'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l'}

// Bottom row
var BOTTOM_ROW_CHARS = []byte{'z', 'x', 'c', 'v', 'b', 'n', 'm'}

// All keys
var ALL_KEYBOARD_CHARS = append(append(append([]byte{}, TOP_ROW_CHARS...), HOME_ROW_CHARS...), BOTTOM_ROW_CHARS...)

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		choices:  []string{"Top Row", "Home Row", "Bottom Row"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
        screen: "select",
	}
}

func (activeModel model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return nil
}

func (activeModel model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    // Is it a key press?
    case tea.KeyMsg:

        // Cool, what was the actual key pressed?
        switch msg.String() {

        // These keys should exit the program.
        case "ctrl+c", "q":
            return activeModel, tea.Quit

        // The "up" and "k" keys move the cursor up
        case "up", "k":
            if activeModel.cursor > 0 {
                activeModel.cursor--
            }

        // The "down" and "j" keys move the cursor down
        case "down", "j":
            if activeModel.cursor < len(activeModel.choices)-1 {
                activeModel.cursor++
            }

        // The "enter" key and the spacebar (a literal space) toggle
        // the selected state for the item that the cursor is pointing at.
        case " ":
            if activeModel.screen == "select" {
                _, ok := activeModel.selected[activeModel.cursor]
                if ok {
                    delete(activeModel.selected, activeModel.cursor)
                } else {
                    activeModel.selected[activeModel.cursor] = struct{}{}
                }
            }
        
        case "enter":
            if activeModel.screen == "select" {
                activeModel.screen = "output"
            } else if activeModel.screen == "output" {
                activeModel.screen = "select"
            }
			return activeModel, tea.Quit
        }
        
    }

    // Return the updated model to the Bubble Tea runtime for processing.
    // Note that we're not returning a command.
    return activeModel, nil
}

func (activeModel model) View() string {
    // The header
    uiText := ""

    switch activeModel.screen {
        case "select":
            uiText += activeText("Select the rows to include\n\n")

            // Iterate over our choices
            for i, choice := range activeModel.choices {

                // // Is the cursor pointing at this choice?
                // cursor := output.String("a") // no cursor
				// cursor.Foreground(output.Color("1"))
                // if activeModel.cursor == i {
                //     cursor.Foreground(output.Color("241"))
                // }

                // Is this choice selected?
                checked := dotInactive // not selected
                if _, ok := activeModel.selected[i]; ok {
                    checked = dotActive // selected!
                }
				
				if activeModel.cursor == i {
                    choice = activeText(choice)
                } else {
					choice = inactiveText(choice)
				}

                // Render the row
                // Add the text to the UI
				currentText := checked + choice + "\n"

				

                uiText += currentText
            }

            // The help footer
			uiText += help("\nPress q or ctrl+c to quit.\n")
			uiText += help("Press enter to view chars\n")
			uiText += help("Press space to toggle value\n")

        case "output":
            _, topRow := activeModel.selected[0]
            _, homeRow := activeModel.selected[1]
            _, bottomRow := activeModel.selected[2]

            uiText += "Here are your chars for MonkeyType\n\n"

            var allowList []byte
            var denyList []byte
            
            // Allow List
            if topRow {
                allowList = append(allowList, TOP_ROW_CHARS...)
            }
            if homeRow {
                allowList = append(allowList, HOME_ROW_CHARS...)
            }
            if bottomRow {
                allowList = append(allowList, BOTTOM_ROW_CHARS...)
            }
            
            // Deny List
            if !topRow {
                denyList = append(denyList, TOP_ROW_CHARS...)
            }
            if !homeRow {
                denyList = append(denyList, HOME_ROW_CHARS...)
            }
            if !bottomRow {
                denyList = append(denyList, BOTTOM_ROW_CHARS...)
            }


            uiText += fmt.Sprintf("Allow: %s\n", returnText(addSpaces(allowList)))
            uiText += fmt.Sprintf("Deny: %s\n", returnText(addSpaces(denyList)))

    }

    

    
    // Send the UI for rendering
    return uiText
}

func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}

func addSpaces(list []byte) string {
	var strBuilder strings.Builder
	strBuilder.WriteString(" ")
	for _, v := range list {
		strBuilder.WriteByte(v)
		strBuilder.WriteString(" ")
	}

	return strBuilder.String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

// Return a function that will colorize the background of a given string.
func makeFgBgStyle(FgColor string, BgColor string) func(string) string {
	return termenv.Style{}.Background(term.Color(BgColor)).Foreground(term.Color(FgColor)).Styled
}

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}