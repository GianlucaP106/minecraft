#version 330

uniform sampler2D tex;

in vec2 TexCoords;

out vec4 color;

void main() {
    vec4 c = texture(tex, TexCoords);
    if (c.a < 0.1) {
        discard;
    }

    color = c;
}
