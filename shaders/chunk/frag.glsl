#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
flat in int selected;

out vec4 color;

void main() {
    color = texture(tex, fragTexCoord);
    color = texture(tex, fragTexCoord);
    if (selected == 1) {
        color = color * 0.6;
    }
}
