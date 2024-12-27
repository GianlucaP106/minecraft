#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 color;

void main() {
    color = texture(tex, fragTexCoord);
    if (color.a < 0.1) {
        discard;
    }
}
