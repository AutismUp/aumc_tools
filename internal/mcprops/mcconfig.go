package mcprops

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// MCConfig represents a Minecraft server.properties configuration file.
type MCConfig struct {
	Properties map[string]string
	order      []string
}

// LoadProperties parses the provided server.properties file and returns an MCConfig instance.
func LoadProperties(path string) (*MCConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open properties file: %w", err)
	}
	defer file.Close()

	properties := make(map[string]string)
	order := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if _, exists := properties[key]; !exists {
			order = append(order, key)
		}
		properties[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read properties file: %w", err)
	}

	return &MCConfig{Properties: properties, order: order}, nil
}

// UpdateProperty sets or updates a property value.
func (c *MCConfig) UpdateProperty(key, value string) {
	if c.Properties == nil {
		c.Properties = make(map[string]string)
	}

	if _, exists := c.Properties[key]; !exists {
		c.order = append(c.order, key)
	}

	c.Properties[key] = value
}

// WriteProperties writes the configuration back to a server.properties file.
func (c *MCConfig) WriteProperties(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create properties file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	if _, err := writer.WriteString("#Minecraft server properties\n"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	timestamp := fmt.Sprintf("%s EST %d", time.Now().Format("Mon Jan 02 15:04:05"), time.Now().Year())
	if _, err := writer.WriteString(fmt.Sprintf("#%s\n", timestamp)); err != nil {
		return fmt.Errorf("failed to write timestamp: %w", err)
	}

	written := make(map[string]bool)
	for _, key := range c.order {
		if value, ok := c.Properties[key]; ok {
			if err := writeProperty(writer, key, value); err != nil {
				return err
			}
			written[key] = true
		}
	}

	for key, value := range c.Properties {
		if written[key] {
			continue
		}
		if err := writeProperty(writer, key, value); err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush properties file: %w", err)
	}

	return nil
}

func writeProperty(writer *bufio.Writer, key, value string) error {
	if _, err := writer.WriteString(fmt.Sprintf("%s=%s\n", key, value)); err != nil {
		return fmt.Errorf("failed to write property %s: %w", key, err)
	}
	return nil
}
