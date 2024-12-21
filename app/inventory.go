package app

type Inventory struct {
	// maps the blockType to count
	content map[string]int

	// selected block type
	selected string

	// TODO: hotbar
}

func newInventory() *Inventory {
	i := &Inventory{}
	i.content = make(map[string]int)
	return i
}

// Adds an item to the inventory.
func (i *Inventory) Add(blockType string, count int) {
	i.content[blockType] += count
}

// Grabs a number of the selected blockType and returns true if amount was deducted.
func (i *Inventory) Grab(count int) (bool, string) {
	stock, exists := i.content[i.selected]
	if exists && stock > count {
		i.content[i.selected] -= count
		return true, i.selected
	}
	return false, ""
}
