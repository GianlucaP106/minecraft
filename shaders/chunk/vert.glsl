#version 330

uniform mat4 model;
uniform mat4 view;
uniform vec3 lookedAtBlock;
uniform bool isLooking;

in vec3 vert;
in vec2 texCoord;

out vec2 fragTexCoord;
out vec2 selected;

void main() {
    vec3 blockMin = vec3(lookedAtBlock);
    vec3 blockMax = blockMin + vec3(1.0);

    vec4 pos = model * vec4(vert, 1);
    bool isSelected = pos.x >= blockMin.x && pos.x <= blockMax.x &&
        pos.y >= blockMin.y && pos.y <= blockMax.y &&
        pos.z >= blockMin.z && pos.z <= blockMax.z && isLooking;

    // TODO: find better way to do this
    if (isSelected) {
        selected = vec2(1.0);
    } else {
        selected = vec2(0.0);
    }


    fragTexCoord = texCoord;
    gl_Position = view * pos;
}

