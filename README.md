# Minecraft from scratch

A Minecraft clone written in Go using only OpenGL (no engines or third-party library/frameworks)

![Game Screenshot](path/to/screenshot.png) <!-- You may want to add a screenshot -->

## Features

- Procedurally generated terrain using Perlin noise
- Dynamic chunk loading/unloading based on player position
- Physics engine with collision detection
- Day/night cycle with dynamic lighting
- Simple inventory system
- Block placement and destruction
- Tree generation
- Flying mode
- Basic cave systems
- Face culling for rendering optimization

## Installation

```bash
git clone ...
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

- Custom physics engine with:
- Rigid body dynamics
- Custom collision detection and response
- Jump mechanics

## Dependencies

- github.com/go-gl/gl/v4.1-core/gl
- github.com/go-gl/glfw/v3.3/glfw
- github.com/go-gl/mathgl/mgl32

## Architecture

The game is built with a component-based architecture, with key systems including:

- World: (chunk loading, block updates)
- World Generator: (noise, terrain, tree, biome)
- Physics engine (collision, movement)
- Player (camera, rigid body)
- Entity management
