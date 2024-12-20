#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
flat in int selected;

out vec4 color;

void main() {
    if (selected == 1) {
        color = texture(tex, fragTexCoord);
        color = color * 0.6;
    } else {
        color = texture(tex, fragTexCoord);
    }
}
