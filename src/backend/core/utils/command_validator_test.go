package utils

import (
	"testing"
)

func TestCommandValidator_ValidateCommand(t *testing.T) {
	cv := NewCommandValidator()

	tests := []struct {
		name    string
		cmd     string
		args    []string
		wantErr bool
	}{
		{
			name:    "允許的命令 - ls",
			cmd:     "ls",
			args:    []string{"-la"},
			wantErr: false,
		},
		{
			name:    "允許的命令 - ps",
			cmd:     "ps",
			args:    []string{"-aux"},
			wantErr: false,
		},
		{
			name:    "允許的命令 - docker ps",
			cmd:     "docker",
			args:    []string{"ps"},
			wantErr: false,
		},
		{
			name:    "不允許的命令 - rm",
			cmd:     "rm",
			args:    []string{"-rf", "/"},
			wantErr: true,
		},
		{
			name:    "危險字符 - 分號",
			cmd:     "ls",
			args:    []string{"; rm -rf /"},
			wantErr: true,
		},
		{
			name:    "危險字符 - 管道",
			cmd:     "cat",
			args:    []string{"/etc/passwd | grep root"},
			wantErr: true,
		},
		{
			name:    "危險字符 - 命令替換",
			cmd:     "echo",
			args:    []string{"$(whoami)"},
			wantErr: true,
		},
		{
			name:    "文件路徑允許",
			cmd:     "cat",
			args:    []string{"/var/log/app.log"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cv.ValidateCommand(tt.cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandValidator_ExecuteCommand(t *testing.T) {
	cv := NewCommandValidator()

	tests := []struct {
		name    string
		cmd     string
		args    []string
		wantErr bool
	}{
		{
			name:    "執行安全命令 - whoami",
			cmd:     "whoami",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "執行不安全命令 - rm",
			cmd:     "rm",
			args:    []string{"-rf", "/tmp/test"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cv.ExecuteCommand(tt.cmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandValidator_AddAllowedCommand(t *testing.T) {
	cv := NewCommandValidator()

	// 添加新命令
	err := cv.AddAllowedCommand("custom", `^[-a-z]+$`)
	if err != nil {
		t.Errorf("AddAllowedCommand() error = %v", err)
	}

	// 驗證命令已添加
	if !cv.IsCommandAllowed("custom") {
		t.Error("Command 'custom' should be allowed after adding")
	}

	// 測試新命令
	err = cv.ValidateCommand("custom", []string{"-abc"})
	if err != nil {
		t.Errorf("ValidateCommand() for custom command error = %v", err)
	}
}

func TestCommandValidator_RemoveAllowedCommand(t *testing.T) {
	cv := NewCommandValidator()

	// 移除命令
	cv.RemoveAllowedCommand("ls")

	// 驗證命令已移除
	if cv.IsCommandAllowed("ls") {
		t.Error("Command 'ls' should not be allowed after removal")
	}

	// 測試移除後的命令
	err := cv.ValidateCommand("ls", []string{"-la"})
	if err == nil {
		t.Error("ValidateCommand() should return error for removed command")
	}
}

func TestCommandValidator_GetAllowedCommands(t *testing.T) {
	cv := NewCommandValidator()

	commands := cv.GetAllowedCommands()
	if len(commands) == 0 {
		t.Error("GetAllowedCommands() should return non-empty list")
	}

	// 驗證包含預期的命令
	found := false
	for _, cmd := range commands {
		if cmd == "ls" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GetAllowedCommands() should include 'ls'")
	}
}

func BenchmarkCommandValidator_ValidateCommand(b *testing.B) {
	cv := NewCommandValidator()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cv.ValidateCommand("ls", []string{"-la"})
	}
}

