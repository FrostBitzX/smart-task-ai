package service

import (
	"fmt"
	"strings"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/chats"
	taskEntity "github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/entity"
	"github.com/FrostBitzX/smart-task-ai/internal/utils"
)

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

	// Function calling instructions
	sb.WriteString("\n## การจัดการ Task\n")
	sb.WriteString("เมื่อผู้ใช้ต้องการจัดการ task ให้ตอบในรูปแบบ JSON ดังนี้:\n\n")

	sb.WriteString("### สร้าง Task ใหม่:\n")
	sb.WriteString("```json\n")
	sb.WriteString(`{"action": "create_task", "params": {"name": "ชื่อ task", "description": "รายละเอียด (optional)", "priority": "low|medium|high", "start_datetime": "2024-01-01T09:00:00Z (optional)", "end_datetime": "2024-01-01T17:00:00Z (optional)"}}`)
	sb.WriteString("\n```\n\n")

	sb.WriteString("### แก้ไข Task:\n")
	sb.WriteString("```json\n")
	sb.WriteString(`{"action": "update_task", "params": {"task_id": "task_xxx", "name": "ชื่อใหม่", "description": "รายละเอียดใหม่", "priority": "low|medium|high"}}`)
	sb.WriteString("\n```\n\n")

	sb.WriteString("### ลบ Task:\n")
	sb.WriteString("```json\n")
	sb.WriteString(`{"action": "delete_task", "params": {"task_id": "task_xxx"}}`)
	sb.WriteString("\n```\n\n")

	sb.WriteString("### ดู Task:\n")
	sb.WriteString("```json\n")
	sb.WriteString(`{"action": "get_task", "params": {"task_id": "task_xxx"}}`)
	sb.WriteString("\n```\n\n")

	sb.WriteString("### ดูรายการ Tasks:\n")
	sb.WriteString("```json\n")
	sb.WriteString(`{"action": "list_tasks", "params": {"status": "todo|in_progress|done (optional)"}}`)
	sb.WriteString("\n```\n\n")

	sb.WriteString("## กฎสำคัญ:\n")
	sb.WriteString("1. ถ้าผู้ใช้ต้องการจัดการ task ให้ตอบเป็น JSON เท่านั้น (ไม่ต้องมีข้อความอื่น)\n")
	sb.WriteString("2. ถ้าผู้ใช้ถามคำถามทั่วไปหรือต้องการคำแนะนำ ให้ตอบเป็นข้อความปกติ\n")
	sb.WriteString("3. priority ต้องเป็น: low, medium, หรือ high\n")
	sb.WriteString("4. datetime ต้องอยู่ในรูปแบบ RFC3339 (เช่น 2024-01-15T09:00:00Z)\n")
	sb.WriteString("5. ใช้ task_id จากรายการ tasks ด้านบนเมื่อต้องการแก้ไขหรือลบ\n")

	return sb.String()
}
