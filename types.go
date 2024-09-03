package main

type TriggerType string

const (
	TriggerTypeScheduled TriggerType = "scheduled"
)

// The type used for the sourceDetails of a trigger whose sourceType is "scheduled"
type ScheduledSourceDetails struct {
	Cron     string      `json:"cron,omitempty"`     // must be a cron expression
	Interval int         `json:"interval,omitempty"` // in minutes (not yet implemented)
	Once     string      `json:"once,omitempty"`     // an ISO-format date (not yet implemented)
	Body     interface{} `json:"body,omitempty"`     // optional body to use when posting the function
}

type TriggerSpec struct {
	Name             string                 `yaml:"name"`                       // The name of the trigger. Must be unique within the namespace.
	Type             TriggerType            `yaml:"type"`                       // Currently, the one supported value "scheduled" is required.
	ScheduledDetails ScheduledSourceDetails `yaml:"scheduledDetails,omitempty"` // Optional details for scheduled trigger
	Enabled          bool                   `yaml:"enabled,omitempty"`          // Assumed true if omitted
}

type ActionSpec struct {
	Name    string `yaml:"name"`              // The name of the action
	Package string `yaml:"package,omitempty"` // The name of the package where action appears ('default' if no package)

	Sequence    []string          `yaml:"sequence,omitempty"`    // Indicates that this action is a sequence and provides its components.  Mutually exclusive with the 'exec' options
	Web         interface{}       `yaml:"web,omitempty"`         // like --web on the CLI; expands to multiple annotations.  Project reader assigns true unless overridden.
	WebSecure   interface{}       `yaml:"webSecure,omitempty"`   // like --web-secure on the CLI.  False unless overridden
	Annotations map[string]string `yaml:"annotations,omitempty"` // 'web' and 'webSecure' are merged with what's here iff present
	Parameters  map[string]string `yaml:"parameters,omitempty"`  // Bound parameters for the action passed in the usual way
	Environment map[string]string `yaml:"environment,omitempty"` // Bound parameters for the action destined to go in the environment
	Limits      Limits            `yaml:"limits,omitempty"`      // Action limits (time, memory, logs)
	RemoteBuild bool              `yaml:"remoteBuild,omitempty"` // States that the build (if any) must be done remotely
	LocalBuild  bool              `yaml:"localBuild,omitempty"`  // States that the build (if any) must be done locally
	Triggers    []TriggerSpec     `yaml:"triggers,omitempty"`    // Triggers for the function if any
}

type BindingSpec struct {
	Name      string // the name of the package to which the present package is bound
	Namespace string // the namespace of the package to which the present package is bound
}

type PackageSpec struct {
	Name                string                 `yaml:"name"`
	Actions             []ActionSpec           `yaml:"actions"`
	Annotations         map[string]string      `yaml:"annotations"`
	Parameters          map[string]interface{} `yaml:"parameters"`
	Environment         map[string]string      `yaml:"environment"`
	Shared              bool                   `yaml:"shared"`
	Web                 interface{}            `yaml:"web"`
	DeployedDuringBuild bool                   `yaml:"deployed_during_build"`
	Binding             BindingSpec            `yaml:"binding"`
}

type ProjectSpec struct {
	Packages        []PackageSpec     `yaml:"packages,omitempty"`        // The packages found in the package directory
	TargetNamespace string            `yaml:"targetNamespace,omitempty"` // The namespace to which we are deploying
	Parameters      map[string]string `yaml:"parameters,omitempty"`      // Parameters to apply to all packages in the project
	Environment     map[string]string `yaml:"environment,omitempty"`     // Environment to apply to all packages in the project
}
