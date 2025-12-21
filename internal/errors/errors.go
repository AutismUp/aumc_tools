package errors

import "fmt"

// ConfigError represents a configuration-related error
type ConfigError struct {
	FilePath  string // Path to the config file
	Field     string // Field name that caused the error (if applicable)
	Operation string // Description of what was being done
	Err       error  // Underlying error
}

func (e *ConfigError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("config error in %s: field '%s' - %v", e.FilePath, e.Field, e.Err)
	}
	return fmt.Sprintf("config error in %s: %s - %v", e.FilePath, e.Operation, e.Err)
}

func (e *ConfigError) Cause() error {
	return e.Err
}

func NewConfigError(filePath, operation string, err error) *ConfigError {
	return &ConfigError{
		FilePath:  filePath,
		Operation: operation,
		Err:       err,
	}
}

func NewConfigFieldError(filePath, field string, err error) *ConfigError {
	return &ConfigError{
		FilePath: filePath,
		Field:    field,
		Err:      err,
	}
}

// PropertiesError represents a server.properties file operation error
type PropertiesError struct {
	FilePath  string // Path to the properties file
	Operation string // What operation failed (read, write, parse)
	Err       error  // Underlying error
}

func (e *PropertiesError) Error() string {
	return fmt.Sprintf("properties file error (%s): %s - %v", e.Operation, e.FilePath, e.Err)
}

func (e *PropertiesError) Cause() error {
	return e.Err
}

func NewPropertiesError(filePath, operation string, err error) *PropertiesError {
	return &PropertiesError{
		FilePath:  filePath,
		Operation: operation,
		Err:       err,
	}
}

// WorldError represents a world management error
type WorldError struct {
	WorldName string // Name of the world
	Operation string // What operation failed (create, delete, restore)
	Err       error  // Underlying error
}

func (e *WorldError) Error() string {
	return fmt.Sprintf("world error (%s '%s'): %v", e.Operation, e.WorldName, e.Err)
}

func (e *WorldError) Cause() error {
	return e.Err
}

func NewWorldError(worldName, operation string, err error) *WorldError {
	return &WorldError{
		WorldName: worldName,
		Operation: operation,
		Err:       err,
	}
}

// BuildError represents a jar building error
type BuildError struct {
	Version  string // Minecraft version being built
	BuildDir string // Build directory
	Phase    string // What phase of build failed
	Err      error  // Underlying error
}

func (e *BuildError) Error() string {
	if e.Phase != "" {
		return fmt.Sprintf("build error (version %s, phase '%s'): %v", e.Version, e.Phase, e.Err)
	}
	return fmt.Sprintf("build error (version %s in %s): %v", e.Version, e.BuildDir, e.Err)
}

func (e *BuildError) Cause() error {
	return e.Err
}

func NewBuildError(version, phase string, err error) *BuildError {
	return &BuildError{
		Version: version,
		Phase:   phase,
		Err:     err,
	}
}

// PublishError represents a jar publishing error
type PublishError struct {
	Filename string // Filename of the jar
	Repo     string // Git repository URL
	Phase    string // What phase of publish failed
	Err      error  // Underlying error
}

func (e *PublishError) Error() string {
	if e.Phase != "" {
		return fmt.Sprintf("publish error (phase '%s' for %s): %v", e.Phase, e.Filename, e.Err)
	}
	return fmt.Sprintf("publish error (file %s to %s): %v", e.Filename, e.Repo, e.Err)
}

func (e *PublishError) Cause() error {
	return e.Err
}

func NewPublishError(filename, phase string, err error) *PublishError {
	return &PublishError{
		Filename: filename,
		Phase:    phase,
		Err:      err,
	}
}

// SubprocessError represents an error from external command execution
type SubprocessError struct {
	Command  string // Command that was executed
	Reason   string // Why it failed (not found, timeout, etc.)
	ExitCode int    // Exit code if available
	Stderr   string // Stderr output if available
	Err      error  // Underlying error
}

func (e *SubprocessError) Error() string {
	if e.ExitCode != 0 {
		msg := fmt.Sprintf("subprocess error executing '%s': exit code %d", e.Command, e.ExitCode)
		if e.Stderr != "" {
			msg += fmt.Sprintf(" - %s", e.Stderr)
		}
		return msg
	}
	if e.Reason != "" {
		return fmt.Sprintf("subprocess error executing '%s': %s - %v", e.Command, e.Reason, e.Err)
	}
	return fmt.Sprintf("subprocess error executing '%s': %v", e.Command, e.Err)
}

func (e *SubprocessError) Cause() error {
	return e.Err
}

func NewSubprocessError(command, reason string, err error) *SubprocessError {
	return &SubprocessError{
		Command: command,
		Reason:  reason,
		Err:     err,
	}
}

func NewSubprocessExitError(command string, exitCode int, stderr string) *SubprocessError {
	return &SubprocessError{
		Command:  command,
		ExitCode: exitCode,
		Stderr:   stderr,
		Err:      fmt.Errorf("process exited with code %d", exitCode),
	}
}

// FileSystemError represents a filesystem operation error
type FileSystemError struct {
	Path      string // Path that was being accessed
	Operation string // What operation failed (read, write, create, delete)
	Err       error  // Underlying error
}

func (e *FileSystemError) Error() string {
	return fmt.Sprintf("filesystem error (%s): %s - %v", e.Operation, e.Path, e.Err)
}

func (e *FileSystemError) Cause() error {
	return e.Err
}

func NewFileSystemError(path, operation string, err error) *FileSystemError {
	return &FileSystemError{
		Path:      path,
		Operation: operation,
		Err:       err,
	}
}
