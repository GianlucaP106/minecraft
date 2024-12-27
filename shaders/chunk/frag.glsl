#version 330

// block textures
uniform sampler2D tex;

// position of light source
uniform vec3 lightPos;

// level of light source
uniform float lightLevel;

// position of camera
uniform vec3 cameraPos;

// texture coordinate
in vec2 fragTexCoord;

// if this frag is selected
in float selected;

// normal vector
in vec3 fragNorm;

// world position
in vec3 fragPos;

// final color
out vec4 color;

void main() {
    vec4 c = texture(tex, fragTexCoord);
    // make transparent
    if (c.a < 0.1) {
        discard;
    }

    // make darker when selected
    if (selected == 1.0) {
        c = c * 0.6;
    }

    // lighting parameters
    float ambientStrength = 0.5;
    float specularStrength = 0.25;
    float shininess = 8;
    vec3 lightColor = vec3(1.0, 1.0, 1.0);
    lightColor = lightColor * lightLevel;

    // diffuse lighting component
    vec3 norm = normalize(fragNorm);
    vec3 lightDir = normalize(lightPos - fragPos);
    float diff = max(dot(norm, lightDir), 0.0);

    // specular lighting component
    vec3 viewDir = normalize(cameraPos - fragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), shininess);

    // combine
    vec3 diffuse = diff * lightColor;
    vec3 ambient = ambientStrength * lightColor;
    vec3 specular = specularStrength * spec * lightColor;
    vec4 total = vec4(ambient + diffuse + specular, 1.0);
    color = total * c;
}
