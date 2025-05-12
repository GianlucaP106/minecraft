# Minecraft from scratch ![minecraft](https://github.com/user-attachments/assets/b3f022b3-17d2-49c8-a4ad-be2d3433ef17) 

A Minecraft clone from scratch with only OpenGL. No game engines or frameworks.

<https://github.com/user-attachments/assets/ba68da39-bf90-4e17-b6ad-3e84fae37e23>

<img width="1612" alt="Screenshot 2025-05-11 at 12 48 01 AM" src="https://github.com/user-attachments/assets/66099f92-0b91-4eb3-bb00-945efd5f456f" />

<table>
  <tr>
    <td><img width="800" alt="Screenshot 2025-05-11 at 9 50 57 PM" src="https://github.com/user-attachments/assets/96b171da-a1ea-4b5d-a908-870494763269" /></td>
    <td><img width="800" alt="Screenshot 2025-05-11 at 9 50 12 PM" src="https://github.com/user-attachments/assets/aaeb708a-b788-4d41-8d73-05d69bf5feda" /></td>
  </tr>
  <tr>
    <td><img width="800" alt="Screenshot 2025-05-11 at 9 52 20 PM" src="https://github.com/user-attachments/assets/62e14b50-ac5d-4c25-9ce2-d06d52eb1ebb" /></td>
    <td><img width="800" alt="demo7" src="https://github.com/user-attachments/assets/c617738b-34a6-4099-9012-e2541b2d108a" /></td>
  </tr>
</table>

## ✨ Features

- 🌍 Infinite & procedurally generated terrain using Perlin noise
- ⚙️ Physics engine with collision detection and response
- 💡 Dynamic lighting with shadows and day/night cycle 🌞🌙
- 🧱 Block placement and destruction
- 🌳 Tree generation & basic cave systems 🕳️
- 📦 Dynamic chunk loading/unloading based on player position
- 🎯 Frustum culling for rendering optimization
- 🎒 Simple inventory system with hotbar (1–9)
- 🕹️ Flying mode for creative exploration
- 🗺️ Biome-based terrain variation

---

## 🛠️ Installation

```bash
# Clone the repository
git clone https://github.com/GianlucaP106/minecraft minecraft && cd minecraft

# Run the game (requires Go)
go run .
````

📦 *Make sure you have Go installed: [https://go.dev/dl/](https://go.dev/dl/)*

---

## 🎮 Controls

| Action      | Key/Mouse          |
| ----------- | ------------------ |
| Move        | `W`, `A`, `S`, `D` |
| Jump        | `Space`            |
| Toggle Fly  | `F`                |
| Look Around | `Mouse`            |
| Break Block | `Left Click`       |
| Place Block | `Right Click`      |
| Select Item | `1-9`              |

---

## 🧪 Technical Highlights

### 🖼️ Graphics

- Uses **OpenGL 4.1**
- Custom **shader programs** for blocks, UI, and lighting
- **Frustum culling** for performance optimization

### 🌄 World Generation

- Multi-octave **Perlin noise** for terrain shaping
- Biome system to vary landscape types
- **Procedural caves** and tree generation
- Real-time chunk loading/unloading

### ⚙️ Physics

- Custom physics engine with **rigid body dynamics**
- Block-based **collision detection**
- Jumping & flying mechanics

---

## 🧩 Architecture Overview

![architecture](https://github.com/user-attachments/assets/9945151b-daf4-4918-b670-24881ceb35a4)

The engine follows a **component-based** design. Key systems include:

- 🎮 Game: Core game loop and simulation
- 🌍 World system: Chunk loading, block updates
- 🧬 Generator: Terrain, trees, caves, biomes
- ⚙️ Physics engine: Collision, movement, response
- 🧑 Player: Camera, controls, raycasting

---

## 📦 Dependencies

This project uses ONLY the following Go packages:

- [`go-gl/gl`](https://pkg.go.dev/github.com/go-gl/gl/v4.1-core/gl)
- [`go-gl/glfw`](https://pkg.go.dev/github.com/go-gl/glfw/v3.3/glfw)
- [`go-gl/mathgl`](https://pkg.go.dev/github.com/go-gl/mathgl/mgl32)
- [`mattn/go-sqlite3`](https://pkg.go.dev/github.com/mattn/go-sqlite3)
  
---

## ⚠️ Disclaimer

This project is developed for **educational and non-commercial purposes only**.  
It is **not affiliated with or endorsed by Mojang, Microsoft, or any related entities**.  
