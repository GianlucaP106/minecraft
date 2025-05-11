# Minecraft from scratch

A Minecraft clone from scratch with only modern OpenGL. No game engines or frameworks.

<img width="1612" alt="Screenshot 2025-05-11 at 12 48 01 AM" src="https://github.com/user-attachments/assets/66099f92-0b91-4eb3-bb00-945efd5f456f" />

<table>
  <tr>
    <td><img width="800" alt="Screenshot 2025-05-10 at 12 04 36 PM" src="https://github.com/user-attachments/assets/deaf40ca-8032-434d-bf84-c4806e01f8e0" /></td>
    <td><img width="800" alt="Screenshot 2025-05-06 at 10 01 25 PM" src="https://github.com/user-attachments/assets/382b472a-4c17-4bf1-b3bf-0a7c9013a9f7" /></td>
  </tr>
  <tr>
    <td><img width="800" alt="demo5" src="https://github.com/user-attachments/assets/848e97e3-5d9e-4cb3-9345-2478ca84424a" /></td>
    <td><img width="800" alt="demo7" src="https://github.com/user-attachments/assets/c617738b-34a6-4099-9012-e2541b2d108a" /></td>
  </tr>
</table>


https://github.com/user-attachments/assets/ba68da39-bf90-4e17-b6ad-3e84fae37e23

## Features

- Infinite and procedurally generated terrain using Perlin noise
- Physics engine with collision detection
- Day/night cycle
- Dynamic lighting with shadows
- Block placement and destruction
- Tree generation
- Basic cave systems
- Dynamic chunk loading/unloading based on player position
- Simple culling techniques for rendering optimization
- Simple inventory system
- Flying mode

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
