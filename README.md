# Minecraft from scratch

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

## Features

- Infinite and procedurally generated terrain using Perlin noise
- Physics engine with collision detection
- Dynamic lighting with shadows
- Block placement and destruction
- Tree generation
- Basic cave systems
- Dynamic chunk loading/unloading based on player position
- Simple culling techniques for rendering optimization
- Simple inventory system
- Flying mode
- Day/night cycle

## Installation

```bash
# clone
git clone https://github.com/GianlucaP106/minecraft minecraft && cd minecraft

# run game (requires go)
go run .
```

## Controls

- WASD - Movement
- Space - Jump
- F - Toggle flying mode
- Mouse - Look around
- Left Click - Break block
- Right Click - Place block
- 1-9 - Select inventory slot

## Technical Highlights

### Graphics

- Written in OpenGL 4.1
- Custom shader programs for blocks, UI elements, and effects
- View frustum culling for performance optimization

### World Generation

- Multi-octave Perlin noise for terrain generation
- Biome system affecting terrain height and features
- Procedural cave system generation
- Dynamic tree placement based on biome

### Physics

- Custom physics engine
- Rigid body dynamics
- Custom collision detection and response
- Jump mechanics

## Architecture

![architecture](https://github.com/user-attachments/assets/9945151b-daf4-4918-b670-24881ceb35a4)

The game is built with a component-based architecture, with key systems including:

- World: (chunk loading, block updates)
- World Generator: (noise, terrain, tree, biome)
- Physics engine and Rigid Body (collision, movement)
- Player, Camera and Ray
- Chunk and Block
  
## Dependencies

- github.com/go-gl/gl/v4.1-core/gl
- github.com/go-gl/glfw/v3.3/glfw
- github.com/go-gl/mathgl/mgl32
