#version 330

uniform sampler2D tex;

in vec4 fragColor;

out vec4 color;

void main() {
    color = fragColor;
}
