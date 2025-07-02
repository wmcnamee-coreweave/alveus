package github

import "encoding/json"

// Workflow
// https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#about-yaml-syntax-for-workflows
type Workflow struct {
	Name        string         `json:"name,omitempty"`
	RunName     string         `json:"run-name,omitempty"`
	On          WorkflowOn     `json:"on"`
	Concurrency Concurrency    `json:"concurrency,omitempty"`
	Defaults    Defaults       `json:"defaults,omitempty"`
	Jobs        map[string]Job `json:"jobs"`
}

// WorkflowOn
// https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#on
type WorkflowOn struct {
	Call     OnCall     `json:"workflow_call,omitempty,omitzero"`
	Run      OnRun      `json:"workflow_run,omitempty,omitzero"`
	Dispatch OnDispatch `json:"workflow_dispatch,omitempty,omitzero"`
}

type OnCall struct {
	Inputs  map[string]CallInput  `json:"inputs,omitempty"`
	Outputs map[string]CallOutput `json:"outputs,omitempty"`
}

type CallInput struct {
	Description string        `json:"description,omitempty"`
	Default     string        `json:"default,omitempty"`
	Required    bool          `json:"required,omitempty"`
	Type        CallInputType `json:"type,omitempty"`
}

type CallOutput struct {
	Description string `json:"description,omitempty"`
	Value       string `json:"value,omitempty"`
}

type CallSecrets struct {
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

type CallInputType string

const (
	CallInputTypeString  CallInputType = "string"
	CallInputTypeBoolean CallInputType = "boolean"
	CallInputTypeNumber  CallInputType = "number"
)

type OnRun struct {
	Workflows []string `json:"workflows,omitempty"`
	Types     []string `json:"types,omitempty"`
	Branches  []string `json:"branches,omitempty"`
}

type OnDispatch struct {
	Inputs map[string]OnDispatchInput `json:"inputs,omitempty"`
}

type OnDispatchInput struct {
	Description string   `json:"description,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Default     string   `json:"default,omitempty"`
	Type        string   `json:"type,omitempty"`
	Options     []string `json:"options,omitempty"`
}

type OnDispatchInputType string

const (
	OnDispatchInputTypeString      OnDispatchInputType = "string"
	OnDispatchInputTypeBoolean     OnDispatchInputType = "boolean"
	OnDispatchInputTypeNumber      OnDispatchInputType = "number"
	OnDispatchInputTypeEnvironment OnDispatchInputType = "environment"
	OnDispatchInputTypeChoice      OnDispatchInputType = "choice"
)

// Concurrency
// https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#concurrency
type Concurrency struct{}

type Job struct {
	Name            string            `json:"name,omitempty"`
	Permissions     Permissions       `json:"permissions,omitempty,omitzero"`
	Needs           []string          `json:"needs,omitempty"`
	If              string            `json:"if,omitempty"`
	RunsOn          []string          `json:"runs-on,omitempty"`
	Environment     Environment       `json:"environment,omitempty,omitzero"`
	Concurrency     Concurrency       `json:"concurrency,omitempty,omitzero"`
	Outputs         map[string]string `json:"outputs,omitempty"`
	Env             map[string]string `json:"env,omitempty"`
	Defaults        Defaults          `json:"defaults,omitempty,omitzero"`
	Steps           []Step            `json:"steps,omitempty"`
	TimeoutMinutes  int               `json:"timeout-minutes,omitempty"`
	ContinueOnError bool              `json:"continue-on-error,omitempty"`
	Uses            string            `json:"uses,omitempty"`
	With            map[string]string `json:"with,omitempty"`
	Secrets         Secrets           `json:"secrets,omitempty,omitzero"`
}

type Secrets struct {
	Inherit bool
	Map     map[string]string
}

func (s *Secrets) MarshalJSON() ([]byte, error) {
	if s.Inherit {
		return json.Marshal(s.Inherit)
	} else {
		return json.Marshal(s.Map)
	}
}

func (s *Secrets) UnmarshalJSON(data []byte) error {
	return nil
}

type Permissions struct {
	Actions        AccessLevel `json:"actions,omitempty"`
	Attestations   AccessLevel `json:"attestations,omitempty"`
	Checks         AccessLevel `json:"checks,omitempty"`
	Contents       AccessLevel `json:"contents,omitempty"`
	Deployments    AccessLevel `json:"deployments,omitempty"`
	Discussions    AccessLevel `json:"discussions,omitempty"`
	IDToken        AccessLevel `json:"id-token,omitempty"`
	Issues         AccessLevel `json:"issues,omitempty"`
	Models         AccessLevel `json:"models,omitempty"`
	Packages       AccessLevel `json:"packages,omitempty"`
	Pages          AccessLevel `json:"pages,omitempty"`
	PullRequests   AccessLevel `json:"pull-requests,omitempty"`
	SecurityEvents AccessLevel `json:"security-events,omitempty"`
	Statuses       AccessLevel `json:"statuses,omitempty"`
}

type AccessLevel string

const (
	AccessLevelWrite AccessLevel = "write"
	AccessLevelRead  AccessLevel = "read"
	AccessLevelNone  AccessLevel = "none"
)

type Environment struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type Defaults struct {
	Run DefaultsRun `json:"run,omitempty"`
}

type DefaultsRun struct {
	Shell            Shell  `json:"shell,omitempty"`
	WorkingDirectory string `json:"working-directory,omitempty"`
}

type Step struct {
	ID               string            `json:"id,omitempty"`
	Name             string            `json:"name,omitempty"`
	If               string            `json:"if,omitempty"`
	Uses             string            `json:"uses,omitempty"`
	Run              string            `json:"run,omitempty"`
	WorkingDirectory string            `json:"working-directory,omitempty"`
	Shell            Shell             `json:"shell,omitempty"`
	With             map[string]string `json:"with,omitempty"`
	Env              map[string]string `json:"env,omitempty"`
	ContinueOnError  bool              `json:"continue-on-error,omitempty"`
	TimeoutMinutes   int               `json:"timeout-minutes,omitempty"`
}

type Shell string

const (
	ShellBash Shell = "bash"
)
