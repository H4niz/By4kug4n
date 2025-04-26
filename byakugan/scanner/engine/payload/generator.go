package payload

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/haniz/byakugan/scanner/engine/logger"
	"gopkg.in/yaml.v2"
)

// Generator handles payload generation
type Generator struct {
	templates map[string]*PayloadTemplate
	options   Options
	cache     *Cache // Change to custom Cache type
	logger    logger.Logger
	mu        sync.RWMutex
	config    *Config
}

type GeneratorMetrics struct {
	TotalGenerated   int64         `json:"total_generated"`
	FailedGeneration int64         `json:"failed_generation"`
	GenerationTime   time.Duration `json:"generation_time"`
	TemplateCache    int           `json:"template_cache"`
}

// NewGenerator creates a new payload generator
func NewGenerator(opts Options) (*Generator, error) {
	cfg := &Config{
		LogPath: opts.LogPath,
		Generator: GeneratorConfig{
			TemplatePath:   opts.TemplatePath,
			CacheSize:      opts.CacheSize,
			MaxPayloadSize: int(opts.MaxSize),
		},
	}

	logConfig := logger.LogConfig{
		Component:  "payload_generator",
		Level:      opts.LogLevel,
		LogPath:    opts.LogPath,
		FilePath:   filepath.Join(opts.LogPath, "generator.log"),
		MaxSize:    100 * 1024 * 1024,
		MaxAge:     24 * time.Hour,
		Compress:   true,
		BufferSize: 4096,
		Console:    true,
		AsyncWrite: true,
		FormatJSON: true,
	}

	l, err := logger.NewBaseLogger(logConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return &Generator{
		config:    cfg,
		logger:    l,
		templates: make(map[string]*PayloadTemplate),
		cache:     NewCache(cfg.Generator.CacheSize),
		options:   opts,
	}, nil
}

// LoadTemplates loads all templates from the template directory
func (g *Generator) LoadTemplates() error {
	files, err := os.ReadDir(g.options.TemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read template directory: %w", err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml") {
			if err := g.loadTemplateFile(filepath.Join(g.options.TemplatePath, f.Name())); err != nil {
				return fmt.Errorf("failed to load template %s: %w", f.Name(), err)
			}
		}
	}

	return nil
}

// LoadTemplate loads a payload template
func (g *Generator) LoadTemplate(tmpl *PayloadTemplate) error {
	if tmpl.ID == "" {
		return fmt.Errorf("template ID is empty")
	}

	// Parse template with custom functions
	t := template.New(tmpl.ID).Funcs(g.getTemplateFuncs())

	if tmpl.Template != "" {
		if _, err := t.Parse(tmpl.Template); err != nil {
			return fmt.Errorf("invalid template syntax: %w", err)
		}
	}

	// Handle multiple templates
	if len(tmpl.Templates) > 0 {
		for i, templateStr := range tmpl.Templates {
			if _, err := t.New(fmt.Sprintf("%s-%d", tmpl.ID, i)).Parse(templateStr); err != nil {
				return fmt.Errorf("invalid template syntax in template %d: %w", i, err)
			}
		}
	}

	g.templates[tmpl.ID] = tmpl
	return nil
}

// Generate creates payloads from a template
func (g *Generator) Generate(templateID string, data map[string]interface{}) ([]Payload, error) {
	// Step 1: Validate inputs
	if templateID == "" {
		return nil, fmt.Errorf("template ID cannot be empty")
	}
	if data == nil {
		data = make(map[string]interface{})
	}

	// Step 2: Check cache
	cacheKey := g.generateCacheKey(templateID, data)
	if payloads, exists := g.getFromCache(cacheKey); exists {
		return payloads, nil
	}

	// Step 3: Get and validate template
	g.mu.RLock()
	tmpl, exists := g.templates[templateID]
	g.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("template not found: %s", templateID)
	}

	if err := g.validateTemplate(tmpl); err != nil {
		return nil, fmt.Errorf("invalid template %s: %w", templateID, err)
	}

	// Step 4: Validate template data
	if err := g.validateTemplateData(tmpl, data); err != nil {
		return nil, fmt.Errorf("invalid template data for %s: %w", templateID, err)
	}

	// Step 5: Generate payloads
	var payloads []Payload
	var errs []error

	// Handle single template
	if tmpl.Template != "" {
		payload, err := g.generateFromTemplate(tmpl, tmpl.Template, data)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to generate from main template: %w", err))
		} else {
			payloads = append(payloads, payload)
		}
	}

	// Handle multiple templates
	if len(tmpl.Templates) > 0 {
		for i, templateStr := range tmpl.Templates {
			if templateStr == "" {
				continue
			}
			payload, err := g.generateFromTemplate(tmpl, templateStr, data)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to generate from template %d: %w", i, err))
				continue
			}
			payloads = append(payloads, payload)
		}
	}

	// Handle pattern-based template
	if tmpl.Pattern != "" {
		payload, err := g.generateFromPattern(tmpl, data)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to generate from pattern: %w", err))
		} else {
			payloads = append(payloads, payload)
		}
	}

	// Step 6: Validate results
	if len(payloads) == 0 {
		if len(errs) > 0 {
			return nil, fmt.Errorf("failed to generate payloads: %v", errs)
		}
		return nil, fmt.Errorf("no payloads generated from template %s", templateID)
	}

	// Step 7: Post-process payloads
	processedPayloads, err := g.postProcessPayloads(payloads)
	if err != nil {
		return nil, fmt.Errorf("failed to post-process payloads: %w", err)
	}

	// Step 8: Cache results
	g.cachePayloads(cacheKey, processedPayloads)

	return processedPayloads, nil
}

// Add helper methods
func (g *Generator) validateTemplate(tmpl *PayloadTemplate) error {
	if tmpl == nil {
		return fmt.Errorf("template is nil")
	}

	if tmpl.ID == "" {
		return fmt.Errorf("template ID is empty")
	}

	if tmpl.Type == "" {
		return fmt.Errorf("template type is empty")
	}

	// Check for content
	if tmpl.Template == "" && len(tmpl.Templates) == 0 && tmpl.Pattern == "" {
		return fmt.Errorf("template has no content")
	}

	return nil
}

func (g *Generator) generateFromPattern(tmpl *PayloadTemplate, data map[string]interface{}) (Payload, error) {
	// Generate payload from pattern with variable substitution
	pattern := tmpl.Pattern
	for k, v := range data {
		pattern = strings.ReplaceAll(pattern, fmt.Sprintf("{{%s}}", k), fmt.Sprint(v))
	}

	return Payload{
		ID:          fmt.Sprintf("%s-pattern-%d", tmpl.ID, time.Now().UnixNano()),
		Type:        tmpl.Type,
		Value:       pattern,
		Headers:     make(map[string]string),
		Metadata:    data,
		Description: tmpl.Description,
		Encoded:     false,
	}, nil
}

func (g *Generator) postProcessPayloads(payloads []Payload) ([]Payload, error) {
	processed := make([]Payload, len(payloads))
	for i, p := range payloads {
		// Add timestamp
		p.Metadata["generated_at"] = time.Now()

		// Add unique identifier
		p.ID = fmt.Sprintf("%s-%d", p.ID, i)

		// Validate payload value
		if p.Value == "" {
			continue
		}

		processed[i] = p
	}
	return processed, nil
}

func (g *Generator) generateFromTemplate(tmpl *PayloadTemplate, templateStr string, data map[string]interface{}) (Payload, error) {
	// Create new template with functions
	t := template.New(tmpl.ID).Funcs(g.getTemplateFuncs())

	// Parse template
	if _, err := t.Parse(templateStr); err != nil {
		return Payload{}, fmt.Errorf("template parsing error: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return Payload{}, fmt.Errorf("template execution error: %w", err)
	}

	payload := Payload{
		ID:          fmt.Sprintf("%s-%d", tmpl.ID, time.Now().UnixNano()),
		Type:        tmpl.Type,
		Value:       buf.String(),
		Metadata:    data,
		Description: tmpl.Description,
	}

	return payload, nil
}

func (g *Generator) mergeTemplateData(tmpl *PayloadTemplate, data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Add default variables
	for k, v := range tmpl.Variables {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}

	// Override with provided data
	for k, v := range data {
		result[k] = v
	}

	return result
}

func (g *Generator) validateTemplateData(tmpl *PayloadTemplate, data map[string]interface{}) error {
	if data == nil {
		return fmt.Errorf("template data is nil")
	}

	// Validate required fields based on template type
	switch tmpl.Type {
	case TypeSQLi:
		if _, ok := data["table"]; !ok {
			return fmt.Errorf("missing required field 'table' for SQL injection template")
		}
		if _, ok := data["columns"]; !ok {
			return fmt.Errorf("missing required field 'columns' for SQL injection template")
		}
	}

	return nil
}

func (g *Generator) getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"join": strings.Join,
		"index": func(arr []string, i int) string {
			if i >= 0 && i < len(arr) {
				return arr[i]
			}
			return ""
		},
		"concat":  strings.Join,
		"replace": strings.ReplaceAll,
		"lower":   strings.ToLower,
		"upper":   strings.ToUpper,
	}
}

// GetTemplate retrieves a template by ID
func (g *Generator) GetTemplate(id string) (*PayloadTemplate, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	tmpl, exists := g.templates[id]
	return tmpl, exists
}

// Cache operations
func (g *Generator) getFromCache(key string) ([]Payload, bool) {
	return g.cache.Get(key)
}

func (g *Generator) cachePayloads(key string, payloads []Payload) {
	g.cache.Set(key, payloads)
}

// GetTemplatePath returns the template directory path
func (g *Generator) GetTemplatePath() string {
	return g.options.TemplatePath
}

// GetTemplates returns a list of available template names
func (g *Generator) GetTemplates() []string {
	templates := make([]string, 0, len(g.templates))
	for name := range g.templates {
		templates = append(templates, name)
	}
	return templates
}

// GetTemplateObjects returns a list of all template objects
func (g *Generator) GetTemplateObjects() []*PayloadTemplate {
	g.mu.RLock()
	defer g.mu.RUnlock()

	templates := make([]*PayloadTemplate, 0, len(g.templates))
	for _, tmpl := range g.templates {
		templates = append(templates, tmpl)
	}
	return templates
}

func (g *Generator) loadTemplateFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	var templateFile struct {
		Templates []*PayloadTemplate `yaml:"templates"`
	}

	if err := yaml.Unmarshal(data, &templateFile); err != nil {
		return fmt.Errorf("failed to parse template file: %w", err)
	}

	for _, tmpl := range templateFile.Templates {
		if err := g.LoadTemplate(tmpl); err != nil {
			return fmt.Errorf("failed to load template %s: %w", tmpl.ID, err)
		}
	}

	return nil
}

// Helper method to generate cache key
func (g *Generator) generateCacheKey(templateID string, data map[string]interface{}) string {
	// Create a stable cache key from template ID and data
	key := templateID
	if len(data) > 0 {
		keys := make([]string, 0, len(data))
		for k := range data {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			key += fmt.Sprintf("|%s=%v", k, data[k])
		}
	}
	return key
}
