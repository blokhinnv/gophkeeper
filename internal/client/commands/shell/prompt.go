package shell

import (
	"log"

	"github.com/manifoldco/promptui"

	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

// promptText prompts the user for a text input and returns the text entered.
func promptText(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}
	text, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return text
}

// selectItem prompts the user to select an item from a list of items and returns the selected item.
func selectItem(label string, items any) string {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

// getMetadata prompts the user for metadata keys and values until an empty key
// is entered and returns a map of metadata.
func getMetadata() models.Metadata {
	md := make(models.Metadata)
	for {
		k := promptText("Metadata key: ")
		if k == "" {
			break
		}
		v := promptText("Metadata value: ")
		md[k] = v
	}
	return md
}
