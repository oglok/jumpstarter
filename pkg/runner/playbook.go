package runner

import (
	"fmt"

	"github.com/creasty/defaults"
)

// yaml parser

type JumpstarterPlaybook struct {
	Name          string            `yaml:"name"`
	Tags          []string          `yaml:"tags"`
	Drivers       []string          `yaml:"drivers"`
	ExpectTimeout uint              `yaml:"expect-timeout"`
	Tasks         []JumpstarterTask `yaml:"tasks"`
	Cleanup       []JumpstarterTask `yaml:"cleanup"`
}

type JumpstarterTask struct {
	// name of the task
	Name                  string                     `yaml:"name"`
	SetDiskImage          *SetDiskImageTask          `yaml:"set-disk-image,omitempty"`
	Expect                *ExpectTask                `yaml:"expect,omitempty"`
	Send                  *SendTask                  `yaml:"send,omitempty"`
	Storage               *StorageTask               `yaml:"storage,omitempty"`
	Power                 *PowerTask                 `yaml:"power,omitempty"`
	Reset                 *ResetTask                 `yaml:"reset,omitempty"`
	Pause                 *PauseTask                 `yaml:"pause,omitempty"`
	WriteAnsibleInventory *WriteAnsibleInventoryTask `yaml:"write-ansible-inventory,omitempty"`
	LocalShell            *LocalShell                `yaml:"local-shell,omitempty"`
	parent                *JumpstarterPlaybook
}

type SetDiskImageTask struct {
	Image         string `yaml:"image"`
	AttachStorage bool   `yaml:"attach_storage"`
	OffsetGB      uint   `yaml:"offset-gb"`
}

type ExpectTask struct {
	This         string `yaml:"this"`
	Fatal        string `yaml:"fatal"`
	Echo         bool   `default:"true" yaml:"echo"`
	DebugEscapes bool   `default:"true" yaml:"debug_escapes"`
	Timeout      uint   `yaml:"timeout"`
}

type ResetTask struct {
	TimeMs uint `yaml:"time_ms"`
}

type PauseTask struct {
	Seconds uint `yaml:"seconds"`
}

type WriteAnsibleInventoryTask struct {
	Filename string `default:"inventory" yaml:"filename"`
	User     string `default:"root" yaml:"user"`
	SshKey   string `yaml:"ssh_key"`
}

type LocalShell struct {
	Script string `yaml:"script"`
}

func (e *ExpectTask) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(e)
	type plain ExpectTask
	if err := unmarshal((*plain)(e)); err != nil {
		return err
	}

	return nil
}

type SendTask struct {
	This         []string `yaml:"this"`
	DelayMs      uint     `default:"100" yaml:"delay_ms"`
	Echo         bool     `default:"true" yaml:"echo"`
	DebugEscapes bool     `default:"true" yaml:"debug_escapes"`
}

func (s *SendTask) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(s)
	type plain SendTask
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}

type StorageTask struct {
	Attached bool `yaml:"attached"`
}

type PowerTask struct {
	Action string `yaml:"action"`
}

// a type enum with changed, ok, error
type TaskStatus int

const (
	Changed TaskStatus = iota
	Ok
	Fatal
)

type TaskResult struct {
	status TaskStatus
	err    error
}

func (p *JumpstarterTask) getName() string {
	if p.Name != "" {
		return p.Name
	}

	switch {
	case p.SetDiskImage != nil:
		return "set-disk-image"
	case p.Expect != nil:
		return fmt.Sprintf("expect: %q", p.Expect.This) // we should add a getName method instead
	case p.Send != nil:
		return "send"
	case p.Storage != nil:
		return "storage"
	case p.Power != nil:
		return "power"
	case p.Reset != nil:
		return "reset"
	case p.Pause != nil:
		return "pause"
	case p.WriteAnsibleInventory != nil:
		return "write-ansible-inventory"
	case p.LocalShell != nil:
		return "local-shell"
	default:
		return "unknown"
	}
}
