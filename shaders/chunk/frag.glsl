#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
in vec2 selected;

out vec4 color;

void main() {
    color = texture(tex, fragTexCoord);
    if (selected.x == 1.0) {
        color = color * 0.6;
    }
}
