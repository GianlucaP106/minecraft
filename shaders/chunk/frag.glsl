#version 330

// block textures
uniform sampler2D tex;

// position of light source
uniform vec3 lightPos;

// level of light source
uniform float lightLevel;

// position of camera
uniform vec3 cameraPos;

// depth map
uniform sampler2D shadowMap;

// texture coordinate
in vec2 fragTexCoord;

// if this frag is selected
in float selected;

// normal vector
in vec3 fragNorm;

// world position
in vec3 fragPos;

// light world position
in vec4 fragPosLight;

// final color
out vec4 color;

float ShadowCalculation(vec4 fragPosLightSpace, vec3 lightDir) {
    // perspective divide
    vec3 projCoords = fragPosLightSpace.xyz / fragPosLightSpace.w;

    // transform to [0,1] range
    projCoords = projCoords * 0.5 + 0.5;

    // get closest depth value from light's perspective
    float closestDepth = texture(shadowMap, projCoords.xy).r;

    // get depth of current fragment from light's perspective
    float currentDepth = projCoords.z;

    // calc bias based on light angle
    float bias = max(0.03 * (1.0 - dot(fragNorm, lightDir)), 0.005);
    // float bias = 0.005;

    // check whether current frag pos is in shadow using 5x5 PCF
    float shadow = 0.0;
    vec2 texelSize = 1.0 / textureSize(shadowMap, 0);
    for(int x = -2; x <= 2; ++x) {
        for(int y = -2; y <= 2; ++y) {
            float pcfDepth = texture(shadowMap, projCoords.xy + vec2(x, y) * texelSize).r;
            shadow += currentDepth - bias > pcfDepth ? 1.0 : 0.0;
        }
    }
    shadow /= 25.0;

    return shadow;
}

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
    vec3 lightDir = normalize(lightPos - fragPos); // different than lightDirection since this takes view into account
    float diff = max(dot(norm, lightDir), 0.0);

    // specular lighting component
    vec3 viewDir = normalize(cameraPos - fragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), shininess);

    // shadow
    float shadow = ShadowCalculation(fragPosLight, lightDir);

    // combine
    vec3 diffuse = diff * lightColor;
    vec3 ambient = ambientStrength * lightColor;
    vec3 specular = specularStrength * spec * lightColor;
    vec4 total = vec4(ambient + (1.0 - shadow) * (diffuse + specular), 1.0);
    color = total * c;
}
