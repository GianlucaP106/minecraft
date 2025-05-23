#version 330

uniform mat4 model;
uniform mat4 view;

in vec3 position;
in vec2 texCoords;

out vec2 TexCoords;

void main() {
    vec4 pos = model * vec4(position, 1);
    gl_Position = view * pos;
    TexCoords = texCoords;
}
