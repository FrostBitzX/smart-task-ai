package service

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/chats"
	taskEntity "github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/entity"
	"github.com/FrostBitzX/smart-task-ai/internal/utils"
)

//go:embed instruction.txt
var instructionPrompt string

// PromptBuilder defines the interface for building system prompts
type PromptBuilder interface {
	BuildSystemPrompt(config *chats.AIConfig, tasks []*taskEntity.Task) string
}

// promptBuilder implements the PromptBuilder interface
type promptBuilder struct{}

// NewPromptBuilder creates a new prompt builder
func NewPromptBuilder() PromptBuilder {
	return &promptBuilder{}
}

// BuildSystemPrompt builds a system prompt for the AI assistant
// It includes the AI config settings and current task list
func (p *promptBuilder) BuildSystemPrompt(config *chats.AIConfig, tasks []*taskEntity.Task) string {
	var sb strings.Builder

	// AI assistant introduction
	sb.WriteString("คุณเป็นผู้ช่วย AI สำหรับจัดการ task\n")

	// Include AI config settings
	sb.WriteString(fmt.Sprintf("Style: %s\n", config.ChatStyle))
	sb.WriteString(fmt.Sprintf("ความเชี่ยวชาญ: %s\n", strings.Join(config.DomainKnowledge, ", ")))
	sb.WriteString(fmt.Sprintf("ภาษา: %s\n\n", config.Language))

	// Include task list in prompt with task IDs
	sb.WriteString("Tasks ปัจจุบันใน project:\n")
	if len(tasks) == 0 {
		sb.WriteString("- ไม่มี task\n")
	} else {
		for _, task := range tasks {
			taskID := utils.ShortUUIDWithPrefix(task.ID, taskEntity.TaskIDPrefix)
			sb.WriteString(fmt.Sprintf("- [%s] %s (status: %s, priority: %s)", taskID, task.Name, task.Status, task.Priority))
			if task.Description != nil && *task.Description != "" {
				sb.WriteString(fmt.Sprintf(" - %s", *task.Description))
			}
			if task.StartDateTime != nil {
				sb.WriteString(fmt.Sprintf(" เริ่ม: %s", *task.StartDateTime))
			}
			if task.EndDateTime != nil {
				sb.WriteString(fmt.Sprintf(" สิ้นสุด: %s", *task.EndDateTime))
			}
			sb.WriteString("\n")
		}
	}

	// Response format instructions from external file
	sb.WriteString("\n")
	sb.WriteString(instructionPrompt)

	return sb.String()
}
