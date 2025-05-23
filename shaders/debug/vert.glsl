#version 330

in vec2 position;
in vec2 texCoords;

out vec2 TexCoords;

void main() {
    gl_Position = vec4(position, 0.0f, 1.0f);
    TexCoords = texCoords;
}
