package payload

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/haniz/byakugan/scanner/engine/logger"
)

// TemplateManager manages payload templates
type TemplateManager struct {
	generator   *Generator
	transformer *Transformer
	templates   map[string]*PayloadTemplate
	options     Options
	logger      logger.Logger
}

// NewTemplateManager creates a new template manager
func NewTemplateManager(opts Options) (*TemplateManager, error) {
	// Create generator with error handling
	generator, err := NewGenerator(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create generator: %w", err)
	}

	tm := &TemplateManager{
		generator:   generator,
		transformer: NewTransformer(),
		templates:   make(map[string]*PayloadTemplate),
		options:     opts,
		logger:      generator.logger, // Reuse generator's logger
	}

	return tm, nil
}

// LoadTemplatesFromDir loads all template files from a directory
func (tm *TemplateManager) LoadTemplatesFromDir(dir string) error {
	if tm.logger != nil {
		tm.logger.Info("Loading templates from directory", map[string]interface{}{
			"directory": dir,
		})
	}

	// Update options with template path
	tm.options.TemplatePath = dir

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			path := filepath.Join(dir, file.Name())
			if err := tm.LoadTemplateFile(path); err != nil {
				return err
			}
		}
	}

	return nil
}

// LoadTemplateFile loads a single template file
func (tm *TemplateManager) LoadTemplateFile(path string) error {
	if tm.logger != nil {
		tm.logger.Debug("Loading template file", map[string]interface{}{
			"path": path,
		})
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var tmpl PayloadTemplate
	if err := json.Unmarshal(data, &tmpl); err != nil {
		return err
	}

	return tm.generator.LoadTemplate(&tmpl)
}

// GeneratePayloads generates payloads from a template with transformations
func (tm *TemplateManager) GeneratePayloads(templateID string, data map[string]interface{}, transformations []string) ([]Payload, error) {
	if tm.logger != nil {
		tm.logger.Debug("Generating payloads", map[string]interface{}{
			"template_id":     templateID,
			"transformations": transformations,
		})
	}

	payloads, err := tm.generator.Generate(templateID, data)
	if err != nil {
		return nil, err
	}

	// Apply transformations if specified
	if len(transformations) > 0 {
		for i := range payloads {
			if err := tm.transformer.Transform(&payloads[i], transformations); err != nil {
				return nil, err
			}
		}
	}

	return payloads, nil
}
