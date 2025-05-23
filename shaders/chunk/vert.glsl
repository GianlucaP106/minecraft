#version 330

// world transformation
uniform mat4 model;

// perspective
uniform mat4 view;

// position of block being looked at
uniform vec3 lookedAtBlock;

// is looking at chunk
uniform bool isLooking;

// light matrix
uniform mat4 lightSpaceMatrix;

// position of vertex without tranform
in vec3 vert;

// normal vector of the vertex
in vec3 normal;

// texture coordinate in the atlas
in vec2 texCoord;

// outputs
out vec2 fragTexCoord;
out float selected;
out vec3 fragNorm;
out vec3 fragPos;
out vec4 fragPosLight;

void main() {
    // world pos of vertex
    vec4 pos = model * vec4(vert, 1);

    // bounding box of looked at block
    vec3 blockMin = vec3(lookedAtBlock);
    vec3 blockMax = blockMin + vec3(1.0);

    // is the block being looked at
    bool isSelected = pos.x >= blockMin.x && pos.x <= blockMax.x &&
            pos.y >= blockMin.y && pos.y <= blockMax.y &&
            pos.z >= blockMin.z && pos.z <= blockMax.z && isLooking;

    selected = isSelected ? 1.0 : 0.0;

    // apply special world transformation to normal vector
    // due to normal transformation issue
    fragNorm = mat3(transpose(inverse(model))) * normal;

    fragTexCoord = texCoord;
    fragPos = vec3(pos);
    fragPosLight = lightSpaceMatrix * pos;
    gl_Position = view * pos;
}
