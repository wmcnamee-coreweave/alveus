package github

import (
	"encoding/json"
	"strings"

	"github.com/ghostsquad/alveus/internal/util"
)

// Workflow
// https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#about-yaml-syntax-for-workflows
type Workflow struct {
	Name        string         `json:"name"`
	RunName     string         `json:"run-name,omitempty,omitzero"`
	On          WorkflowOn     `json:"on"`
	Concurrency Concurrency    `json:"concurrency,omitempty,omitzero"`
	Defaults    Defaults       `json:"defaults,omitempty,omitzero"`
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
	Inputs  map[string]CallInput  `json:"inputs,omitempty,omitzero"`
	Outputs map[string]CallOutput `json:"outputs,omitempty,omitzero"`
}

type CallInput struct {
	Description string        `json:"description,omitempty,omitzero"`
	Default     string        `json:"default,omitempty,omitzero"`
	Required    bool          `json:"required"`
	Type        CallInputType `json:"type"`
}

type CallOutput struct {
	Description string `json:"description,omitempty,omitzero"`
	Value       string `json:"value"`
}

type CallSecrets struct {
	Description string `json:"description,omitempty,omitzero"`
	Required    bool   `json:"required"`
}

type CallInputType string

const (
	CallInputTypeString  CallInputType = "string"
	CallInputTypeBoolean CallInputType = "boolean"
	CallInputTypeNumber  CallInputType = "number"
)

type OnRun struct {
	Workflows []string `json:"workflows,omitempty,omitzero"`
	Types     []string `json:"types,omitempty,omitzero"`
	Branches  []string `json:"branches,omitempty,omitzero"`
}

type OnDispatch struct {
	Inputs map[string]OnDispatchInput `json:"inputs,omitempty,omitzero"`
}

type OnDispatchInput struct {
	Description string              `json:"description,omitempty,omitzero"`
	Required    bool                `json:"required"`
	Default     string              `json:"default,omitempty,omitzero"`
	Type        OnDispatchInputType `json:"type,omitempty,omitzero"`
	Options     []string            `json:"options,omitempty,omitzero"`
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
type Concurrency struct {
	Group            string `json:"group,omitempty,omitzero"`
	CancelInProgress bool   `json:"cancel-in-progress,omitempty,omitzero"`
}

type Job struct {
	Name            string            `json:"name,omitempty,omitzero"`
	Permissions     Permissions       `json:"permissions,omitempty,omitzero"`
	Needs           []string          `json:"needs,omitempty,omitzero"`
	If              string            `json:"if,omitempty,omitzero"`
	RunsOn          []string          `json:"runs-on,omitempty,omitzero"`
	Environment     Environment       `json:"environment,omitempty,omitzero"`
	Concurrency     Concurrency       `json:"concurrency,omitempty,omitzero"`
	Outputs         map[string]string `json:"outputs,omitempty,omitzero"`
	Env             map[string]string `json:"env,omitempty,omitzero"`
	Defaults        Defaults          `json:"defaults,omitempty,omitzero"`
	Steps           []Step            `json:"steps,omitempty,omitzero"`
	TimeoutMinutes  int               `json:"timeout-minutes,omitempty,omitzero"`
	ContinueOnError bool              `json:"continue-on-error,omitempty,omitzero"`
	Uses            string            `json:"uses,omitempty,omitzero"`
	With            map[string]any    `json:"with,omitempty,omitzero"`
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
	Actions        AccessLevel `json:"actions,omitempty,omitzero"`
	Attestations   AccessLevel `json:"attestations,omitempty,omitzero"`
	Checks         AccessLevel `json:"checks,omitempty,omitzero"`
	Contents       AccessLevel `json:"contents,omitempty,omitzero"`
	Deployments    AccessLevel `json:"deployments,omitempty,omitzero"`
	Discussions    AccessLevel `json:"discussions,omitempty,omitzero"`
	IDToken        AccessLevel `json:"id-token,omitempty,omitzero"`
	Issues         AccessLevel `json:"issues,omitempty,omitzero"`
	Models         AccessLevel `json:"models,omitempty,omitzero"`
	Packages       AccessLevel `json:"packages,omitempty,omitzero"`
	Pages          AccessLevel `json:"pages,omitempty,omitzero"`
	PullRequests   AccessLevel `json:"pull-requests,omitempty,omitzero"`
	SecurityEvents AccessLevel `json:"security-events,omitempty,omitzero"`
	Statuses       AccessLevel `json:"statuses,omitempty,omitzero"`
}

type AccessLevel string

const (
	AccessLevelWrite AccessLevel = "write"
	AccessLevelRead  AccessLevel = "read"
	AccessLevelNone  AccessLevel = "none"
)

type Environment struct {
	Name string `json:"name,omitempty,omitzero"`
	URL  string `json:"url,omitempty,omitzero"`
}

type Defaults struct {
	Run DefaultsRun `json:"run,omitempty,omitzero"`
}

type DefaultsRun struct {
	Shell            Shell  `json:"shell,omitempty,omitzero"`
	WorkingDirectory string `json:"working-directory,omitempty"`
}

type Step struct {
	ID               string            `json:"id,omitempty,omitzero"`
	Name             string            `json:"name,omitempty,omitzero"`
	If               string            `json:"if,omitempty,omitzero"`
	Uses             string            `json:"uses,omitempty,omitzero"`
	Run              string            `json:"run,omitempty,omitzero"`
	WorkingDirectory string            `json:"working-directory,omitempty,omitzero"`
	Shell            Shell             `json:"shell,omitempty,omitzero"`
	With             map[string]any    `json:"with,omitempty,omitzero"`
	Env              map[string]string `json:"env,omitempty,omitzero"`
	ContinueOnError  bool              `json:"continue-on-error,omitempty,omitzero"`
	TimeoutMinutes   int               `json:"timeout-minutes,omitempty,omitzero"`
	Container        Container         `json:"container,omitempty,omitzero"`
	Strategy         Strategy          `json:"strategy,omitempty,omitzero"`
}

type Shell string

const (
	ShellBash Shell = "bash"
)

type Container struct {
	Image       string               `json:"image,omitempty"`
	Env         map[string]string    `json:"env,omitempty"`
	Ports       map[string]int       `json:"ports,omitempty"`
	Volumes     []string             `json:"volumes,omitempty"`
	Credentials ContainerCredentials `json:"credentials,omitempty,omitzero"`
	Options     string               `json:"options,omitempty"`
}

type ContainerCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Strategy struct {
	Matrix      Matrix `json:"matrix,omitempty"`
	FailFast    bool   `json:"fail-fast,omitempty"`
	MaxParallel int    `json:"max-parallel,omitempty"`
}

type Matrix struct {
	Map     map[string]ListOrExpression
	Include []map[string]ListOrExpression
	Exclude []map[string]ListOrExpression
}

func (m *Matrix) UnmarshalJSON(data []byte) error {
	return nil
}

func (m *Matrix) MarshalJSON() ([]byte, error) {
	return nil, nil
}

type ListOrExpression struct {
	ListValue  []StringOrInt
	Expression string
}

func (l *ListOrExpression) UnmarshalJSON(data []byte) error {
	return nil
}

func (l *ListOrExpression) MarshalJSON() ([]byte, error) {
	return nil, nil
}

type StringOrInt struct {
	StringValue string `json:"string_value,omitempty"`
	IntValue    int    `json:"int_value,omitempty"`
}

func (l *StringOrInt) UnmarshalJSON(data []byte) error {
	return nil
}

func (l *StringOrInt) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (w *Workflow) GetFullFilename() string {
	newName := strings.Map(func(r rune) rune {
		switch {
		case '0' <= r && r <= '9':
			fallthrough
		case 'A' <= r && r <= 'Z':
			fallthrough
		case 'a' <= r && r <= 'z':
			return r
		default:
			return '-'
		}
	}, w.Name)

	newName = strings.Trim(newName, "-")
	newName = util.RemoveDupOf(newName, '-')
	newName = "./.github/workflows/" + strings.ToLower(newName+".yml")

	return newName
}
