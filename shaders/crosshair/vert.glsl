#version 330

uniform mat4 model;

in vec3 vert;
in vec3 color;

out vec4 fragColor;

void main() {
    fragColor = vec4(color, 1.0);
    gl_Position = model * vec4(vert, 1);
}

