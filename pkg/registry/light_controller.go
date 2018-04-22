package registry

import "fmt"

func (c *ControllerInstance) SetLightIntensity(secret Secret, intensity float32) error {

	if intensity > 1.0 || intensity < 0 {
		return fmt.Errorf("error: intensity[%.2f] is allowed to be [0, 1]", intensity)
	}

	lightId := c.SecretToId[secret]
	if lightId == "" {
		return fmt.Errorf("error: secret[%s] not in use", secret)
	}

	light := c.IdToLight[lightId]
	if light == nil {
		return fmt.Errorf("error: secret[%s] in use but mapped to nil light", secret)
	}

	light.Intensity = intensity // TODO: this should be a method on light so it has a chance to act

	c.LightBoard.SetIntensity(light.Index, light.Intensity)

	return nil
}
