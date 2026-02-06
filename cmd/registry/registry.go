package registry

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

type ID string

// Spec describes a command and where it belongs in the tree.
type Spec struct {
	ID       ID
	ParentID ID
	Source   string
	Build    func() *cli.Command
}

// Registry stores command specs and can build a command tree.
type Registry struct {
	specs []Spec
	byID  map[ID]Spec
}

func New() *Registry {
	return &Registry{byID: make(map[ID]Spec)}
}

func (r *Registry) Add(spec Spec) error {
	id := normalizeID(spec.ID)
	if id == "" {
		return errors.New("command spec id is required")
	}
	if spec.Build == nil {
		return fmt.Errorf("command spec %q has nil builder", id)
	}
	parentID := normalizeID(spec.ParentID)
	if parentID == id {
		return fmt.Errorf("command spec %q cannot reference itself as parent", id)
	}
	if _, exists := r.byID[id]; exists {
		return fmt.Errorf("duplicate command spec id %q", id)
	}

	spec.ID = id
	spec.ParentID = parentID
	spec.Source = strings.TrimSpace(spec.Source)
	r.specs = append(r.specs, spec)
	r.byID[id] = spec
	return nil
}

func (r *Registry) AddAll(specs ...Spec) error {
	for _, spec := range specs {
		if err := r.Add(spec); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) Build() ([]*cli.Command, error) {
	if err := r.validateParentGraph(); err != nil {
		return nil, err
	}

	built := make(map[ID]*cli.Command, len(r.specs))
	childCounts := make(map[ID]int)
	for _, spec := range r.specs {
		if spec.ParentID != "" {
			childCounts[spec.ParentID]++
		}
	}

	for _, spec := range r.specs {
		cmd := spec.Build()
		if cmd == nil {
			return nil, fmt.Errorf("command spec %q built nil command", spec.ID)
		}
		cmd.Name = strings.TrimSpace(cmd.Name)
		if cmd.Name == "" {
			return nil, fmt.Errorf("command spec %q built command with empty name", spec.ID)
		}
		if err := validateOwnTokens(cmd); err != nil {
			return nil, fmt.Errorf("command spec %q invalid tokens: %w", spec.ID, err)
		}
		if childCounts[spec.ID] > 0 && len(cmd.Commands) > 0 {
			return nil, fmt.Errorf("command spec %q has both registry children and preconfigured child commands", spec.ID)
		}
		built[spec.ID] = cmd
	}

	siblingTokens := make(map[ID]map[string]ID)
	checkSiblingTokens := func(parentID ID, cmd *cli.Command, spec Spec) error {
		tokens := []string{cmd.Name}
		tokens = append(tokens, cmd.Aliases...)
		if _, ok := siblingTokens[parentID]; !ok {
			siblingTokens[parentID] = make(map[string]ID)
		}
		for _, token := range tokens {
			token = normalizeToken(token)
			if token == "" {
				continue
			}
			if prev, exists := siblingTokens[parentID][token]; exists {
				if prev == spec.ID {
					continue
				}
				parent := string(parentID)
				if parent == "" {
					parent = "<root>"
				}
				return fmt.Errorf("command token %q conflicts between specs %q (%s) and %q (%s) under parent %q", token, prev, sourceForID(r.byID[prev]), spec.ID, sourceForID(spec), parent)
			}
			siblingTokens[parentID][token] = spec.ID
		}
		return nil
	}

	roots := make([]*cli.Command, 0, len(r.specs))
	for _, spec := range r.specs {
		cmd := built[spec.ID]
		parentID := spec.ParentID
		if err := checkSiblingTokens(parentID, cmd, spec); err != nil {
			return nil, err
		}

		if parentID == "" {
			roots = append(roots, cmd)
			continue
		}

		parent, ok := built[parentID]
		if !ok {
			return nil, fmt.Errorf("command spec %q references unknown parent %q", spec.ID, parentID)
		}
		parent.Commands = append(parent.Commands, cmd)
	}

	return roots, nil
}

func (r *Registry) validateParentGraph() error {
	for _, spec := range r.specs {
		current := spec
		seen := map[ID]bool{spec.ID: true}
		for current.ParentID != "" {
			parent, ok := r.byID[current.ParentID]
			if !ok {
				return fmt.Errorf("command spec %q references unknown parent %q", current.ID, current.ParentID)
			}
			if seen[parent.ID] {
				return fmt.Errorf("cycle detected in command specs at %q", parent.ID)
			}
			seen[parent.ID] = true
			current = parent
		}
	}
	return nil
}

func normalizeID(id ID) ID {
	return ID(strings.TrimSpace(string(id)))
}

func normalizeToken(token string) string {
	return strings.ToLower(strings.TrimSpace(token))
}

func validateOwnTokens(cmd *cli.Command) error {
	seen := map[string]bool{}
	tokens := append([]string{cmd.Name}, cmd.Aliases...)
	for _, token := range tokens {
		normalized := normalizeToken(token)
		if normalized == "" {
			return fmt.Errorf("empty command token")
		}
		if seen[normalized] {
			return fmt.Errorf("duplicate token %q", normalized)
		}
		seen[normalized] = true
	}
	return nil
}

func sourceForID(spec Spec) string {
	if spec.Source == "" {
		return "unspecified"
	}
	return spec.Source
}
