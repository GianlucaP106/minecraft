#version 330

uniform mat4 model;

in vec3 vert;
in vec2 texCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = texCoord;
    gl_Position = model * vec4(vert, 1);
}
