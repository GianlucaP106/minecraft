#version 330

uniform mat4 model;
uniform mat4 view;
uniform vec3 lookedAtBlock;
uniform bool isLooking;

in vec3 vert;
in vec3 color;

out vec4 fragColor;

void main() {
    // TODO: why?
    vec3 blockMin = vec3(lookedAtBlock-1.0);
    vec3 blockMax = blockMin + vec3(1.5);

    vec4 pos = model * vec4(vert, 1);

    if (pos.x >= blockMin.x && pos.x <= blockMax.x &&
        pos.y >= blockMin.y && pos.y <= blockMax.y &&
        pos.z >= blockMin.z && pos.z <= blockMax.z && isLooking) {

        fragColor = vec4(color, 1) * 0.3;
    } else {
        fragColor = vec4(color, 1);
    }
    gl_Position = view * pos;
}

