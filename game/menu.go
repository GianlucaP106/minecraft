package game

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Menu struct {
	db *Database
}

func newMenu(db *Database) *Menu {
	return &Menu{
		db,
	}
}

// Returns the world id that is selected.
// Returns -1 for a new world.
func (m *Menu) Run() *WorldEntity {
	worlds := m.db.Worlds()
	for {
		fmt.Println("Available Worlds:")
		if len(worlds) == 0 {
			fmt.Println("No worlds yet")
		}

		fmt.Printf("(%d) %s\n", 0, "Create a new World")

		for i, we := range worlds {
			fmt.Printf("(%d) %s", i+1, we.name)
		}

		fmt.Print("Enter: ")
		var idx int
		if _, err := fmt.Scan(&idx); err != nil {
			fmt.Printf("Invalid input")
			continue
		}

		if idx == 0 {
			fmt.Print("New World Name: ")
			reader := bufio.NewReader(os.Stdin)
			name, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}

			return m.db.World(m.db.CreateWorld(name))
		}

		if idx > 0 && idx <= len(worlds) {
			return worlds[idx-1]
		}

		fmt.Println("Invalid input")
	}
}
