package game

// Holds count of held blocks.
type Inventory struct {
	// maps the blockType to count
	content map[string]int
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

// Returns count of the item.
func (i *Inventory) Count(blockType string) int {
	return i.content[blockType]
}

// Grabs a number of the selected blockType and returns true if amount was deducted.
func (i *Inventory) Grab(blockType string, count int) bool {
	stock, exists := i.content[blockType]
	if exists && stock >= count {
		i.content[blockType] -= count
		return true
	}
	return false
}
