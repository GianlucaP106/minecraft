package game

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	dsn string
	db  *sql.DB
}

func newDatabase(dsn string) *Database {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database %s", dsn)
		return nil
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Failed to enable foreign keys:", err)
	}

	return &Database{
		dsn,
		db,
	}
}

func (d *Database) Drop() {
	dropTables := `
		DROP TABLE IF EXISTS blocks;
		DROP TABLE IF EXISTS chunks;
		DROP TABLE IF EXISTS worlds
	`

	_, err := d.db.Exec(dropTables)
	if err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
		return
	}
}

func (d *Database) Migrate() {
	createTables := `
	CREATE TABLE IF NOT EXISTS worlds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		inventory TEXT NOT NULL,
		player_x REAL NOT NULL,
		player_y REAL NOT NULL,
		player_z REAL NOT NULL
	);

	CREATE TABLE IF NOT EXISTS chunks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		world_id INTEGER NOT NULL,
		x INTEGER NOT NULL,
		y INTEGER NOT NULL,
		z INTEGER NOT NULL,
		UNIQUE (world_id, x, y, z),
		FOREIGN KEY (world_id) REFERENCES worlds (id)
	);

	CREATE TABLE IF NOT EXISTS blocks (
		chunk_id INTEGER NOT NULL,
		i INTEGER NOT NULL,
		j INTEGER NOT NULL,
		k INTEGER NOT NULL,
		block_type TEXT NOT NULL,
		active INTEGER NOT NULL,
		FOREIGN KEY (chunk_id) REFERENCES chunks (id),
		PRIMARY KEY (chunk_id, i, j, k)
	)
	`

	_, err := d.db.Exec(createTables)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
		return
	}
}

type (
	WorldEntity struct {
		id                        int
		name                      string
		inventory                 string
		playerX, playerY, playerZ float32
	}
	ChunkEntity struct {
		id       int
		world_id int
		x, y, z  int
	}
	BlockEntity struct {
		chunkId   int
		i, j, k   int
		blockType string
		active    bool
	}
)

func (d *Database) World(id int) *WorldEntity {
	res := d.db.QueryRow("SELECT id, name, inventory, player_x, player_y, player_z FROM worlds WHERE id = ?", id)
	if res == nil {
		return nil
	}

	var world WorldEntity
	if err := res.Scan(&world.id, &world.name, &world.inventory, &world.playerX, &world.playerY, &world.playerZ); err != nil {
		return nil
	}

	return &world
}

func (d *Database) Worlds() []*WorldEntity {
	res, err := d.db.Query("SELECT id, name, inventory, player_x, player_y, player_z FROM worlds")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer res.Close()

	out := []*WorldEntity{}
	for res.Next() {
		var w WorldEntity
		if err := res.Scan(&w.id, &w.name, &w.inventory, &w.playerX, &w.playerY, &w.playerZ); err != nil {
			log.Fatal(err)
		}

		out = append(out, &w)
	}

	if err := res.Err(); err != nil {
		log.Fatal(err)
		return nil
	}

	return out
}

func (d *Database) CreateWorld(name string) int {
	r, err := d.db.Exec(
		"INSERT INTO worlds (name, inventory, player_x, player_y, player_z) VALUES (?, ?, ?, ?, ?)",
		name,
		"{}",
		startPosition.X(),
		startPosition.Y(),
		startPosition.Z(),
	)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	id, err := r.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return int(id)
}

func (d *Database) UpdatePosition(worldId int, x, y, z float32) {
	_, err := d.db.Exec(`
		UPDATE worlds
		SET player_x = ?, player_y = ?, player_z = ?
		WHERE id = ?
	`, x, y, z, worldId)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func (d *Database) UpdateInventory(worldId int, content map[string]int) {
	jsonString, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = d.db.Exec(`
		UPDATE worlds
		SET inventory = ?
		WHERE id = ?
	`, string(jsonString), worldId)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func (w *WorldEntity) Inventory() map[string]int {
	var invMap map[string]int
	if err := json.Unmarshal([]byte(w.inventory), &invMap); err != nil {
		panic(err)
	}
	return invMap
}

func (d *Database) FindChunk(worldId, x, y, z int) *ChunkEntity {
	res := d.db.QueryRow("SELECT id, world_id, x, y, z FROM chunks WHERE world_id = ? AND x = ? AND y = ? AND z = ?", worldId, x, y, z)
	if res == nil {
		return nil
	}

	chunk := &ChunkEntity{}
	if err := res.Scan(&chunk.id, &chunk.world_id, &chunk.x, &chunk.y, &chunk.z); err != nil {
		return nil
	}

	return chunk
}

func (d *Database) CreateChunk(worldId, x, y, z int) int {
	r, err := d.db.Exec("INSERT INTO chunks (world_id, x, y, z) VALUES (?, ?, ?, ?)", worldId, x, y, z)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	id, err := r.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return 0
	}
	return int(id)
}

func (d *Database) Block(chunkId, i, j, k int) *BlockEntity {
	res := d.db.QueryRow("SELECT chunk_id, i, j, k, block_type, active FROM blocks WHERE chunk_id = ? AND i = ? AND j = ? AND k = ?", chunkId, i, j, k)
	if res == nil {
		return nil
	}

	block := &BlockEntity{}
	if err := res.Scan(&block.chunkId, &block.i, &block.j, &block.k, &block.blockType, &block.active); err != nil {
		return nil
	}

	return block
}

func (d *Database) Blocks(chunkId int) []*BlockEntity {
	res, err := d.db.Query("SELECT chunk_id, i, j, k, block_type, active FROM blocks WHERE chunk_id = ?", chunkId)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer res.Close()

	var out []*BlockEntity
	for res.Next() {
		var block BlockEntity
		if err := res.Scan(&block.chunkId, &block.i, &block.j, &block.k, &block.blockType, &block.active); err != nil {
			log.Fatal(err)
			return nil
		}

		out = append(out, &block)
	}

	if err := res.Err(); err != nil {
		log.Fatal(err)
		return nil
	}

	return out
}

func (d *Database) CreateBlock(chunkId, i, j, k int, blockType string, active bool) {
	activeVal := 0
	if active {
		activeVal = 1
	}

	_, err := d.db.Exec("INSERT INTO blocks (chunk_id, i, j, k, block_type, active) VALUES (?, ?, ?, ?, ?, ?)", chunkId, i, j, k, blockType, activeVal)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func (d *Database) UpdateBlock(b *BlockEntity) {
	activeVal := 0
	if b.active {
		activeVal = 1
	}

	_, err := d.db.Exec(`
		UPDATE blocks
		SET block_type = ?, active = ?
		WHERE chunk_id = ? AND i = ? AND j = ? AND k = ?
	`, b.blockType, activeVal, b.chunkId, b.i, b.j, b.k)
	if err != nil {
		log.Fatal(err)
		return
	}
}
