package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/application/task"
	taskEntity "github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/entity"
	taskSvc "github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/service"
	"github.com/FrostBitzX/smart-task-ai/internal/utils"
	"github.com/google/uuid"
)

// FunctionExecutor handles execution of AI function calls for task management
type FunctionExecutor struct {
	taskService *taskSvc.TaskService
}

// NewFunctionExecutor creates a new function executor
func NewFunctionExecutor(taskService *taskSvc.TaskService) *FunctionExecutor {
	return &FunctionExecutor{
		taskService: taskService,
	}
}

// FunctionCall represents a function call from AI
type FunctionCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// FunctionResult represents the result of a function execution
type FunctionResult struct {
	Success bool        `json:"success"`
	Data    any         `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Action  *TaskAction `json:"action,omitempty"`
}

// CreateTaskArgs represents arguments for creating a task
type CreateTaskArgs struct {
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	Priority      string  `json:"priority"`
	StartDateTime *string `json:"start_datetime,omitempty"`
	EndDateTime   *string `json:"end_datetime,omitempty"`
	Location      *string `json:"location,omitempty"`
}

// UpdateTaskArgs represents arguments for updating a task
type UpdateTaskArgs struct {
	TaskID        string  `json:"task_id"`
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	Priority      string  `json:"priority"`
	StartDateTime *string `json:"start_datetime,omitempty"`
	EndDateTime   *string `json:"end_datetime,omitempty"`
	Location      *string `json:"location,omitempty"`
}

// DeleteTaskArgs represents arguments for deleting a task
type DeleteTaskArgs struct {
	TaskID string `json:"task_id"`
}

// GetTaskArgs represents arguments for getting a task
type GetTaskArgs struct {
	TaskID string `json:"task_id"`
}

// ListTasksArgs represents arguments for listing tasks
type ListTasksArgs struct {
	Status *string `json:"status,omitempty"`
}

// ExecuteFunction executes a function call and returns the result
func (fe *FunctionExecutor) ExecuteFunction(ctx context.Context, projectID uuid.UUID, call *FunctionCall) *FunctionResult {
	switch call.Name {
	case "create_task":
		return fe.executeCreateTask(ctx, projectID, call.Arguments)
	case "update_task":
		return fe.executeUpdateTask(ctx, call.Arguments)
	case "delete_task":
		return fe.executeDeleteTask(ctx, call.Arguments)
	case "get_task":
		return fe.executeGetTask(ctx, call.Arguments)
	case "list_tasks":
		return fe.executeListTasks(ctx, projectID, call.Arguments)
	default:
		return &FunctionResult{
			Success: false,
			Error:   fmt.Sprintf("unknown function: %s", call.Name),
		}
	}
}

func (fe *FunctionExecutor) executeCreateTask(ctx context.Context, projectID uuid.UUID, args json.RawMessage) *FunctionResult {
	var createArgs CreateTaskArgs
	if err := json.Unmarshal(args, &createArgs); err != nil {
		return &FunctionResult{
			Success: false,
			Error:   fmt.Sprintf("invalid arguments: %v", err),
		}
	}

	req := &task.CreateTaskRequest{
		Name:          createArgs.Name,
		Description:   createArgs.Description,
		Priority:      createArgs.Priority,
		StartDateTime: createArgs.StartDateTime,
		EndDateTime:   createArgs.EndDateTime,
		Location:      createArgs.Location,
	}

	createdTask, err := fe.taskService.CreateTask(ctx, projectID, req)
	if err != nil {
		return &FunctionResult{
			Success: false,
			Error:   fmt.Sprintf("failed to create task: %v", err),
		}
	}

	taskID := utils.ShortUUIDWithPrefix(createdTask.ID, taskEntity.TaskIDPrefix)

	return &FunctionResult{
		Success: true,
		Data: map[string]interface{}{
			"task_id": taskID,
			"name":    createdTask.Name,
			"status":  createdTask.Status,
		},
		Action: &TaskAction{
			Type:   "created",
			TaskID: createdTask.ID,
			Name:   createdTask.Name,
		},
	}
}

func (fe *FunctionExecutor) executeUpdateTask(ctx context.Context, args json.RawMessage) *FunctionResult {
	var updateArgs UpdateTaskArgs
	if err := json.Unmarshal(args, &updateArgs); err != nil {
		return &FunctionResult{
			Success: false,
			Error:   fmt.Sprintf("invalid arguments: %v", err),
		}
	}

	taskID, err := utils.ParseID(updateArgs.TaskID, taskEntity.TaskIDPrefix)
	if err != nil {
		return &FunctionResult{
			Success: false,
			Error:   fmt.Sprintf("invalid task ID: %v", err),
		}
	}

	req := &task.UpdateTaskRequest{
		Name:          updateArgs.Name,
		Description:   updateArgs.Description,
		Priority:      updateArgs.Priority,
		StartDateTime: updateArgs.StartDateTime,
		EndDateTime:   updateArgs.EndDateTime,
		Location:      updateArgs.Location,
	}

	updatedTask, err := fe.taskService.UpdateTask(ctx, taskID, req)
	if err != nil {
		return &FunctionResult{
			Success: false,
			Error:   fmt.Sprintf("failed to update task: %v", err),
		}
	}

	return &FunctionResult{
		Success: true,
		Data: map[string]interface{}{
			"task_id": updateArgs.TaskID,
			"name":    updatedTask.Name,
			"status":  updatedTask.Status,
		},
		Action: &TaskAction{
			Type:   "updated",
			TaskID: taskID,
			Name:   updatedTask.Name,
		},
	}
}

func (fe *FunctionExecutor) executeDeleteTask(ctx context.Context, args json.RawMessage) *FunctionResult {
	var deleteArgs DeleteTaskArgs
	if err := json.Unmarshal(args, &deleteArgs); err != nil {
		return &FunctionResult{
			Success: false,
			Error:   fmt.Sprintf("invalid arguments: %v", err),
		}
	}

	taskID, err := utils.ParseID(deleteArgs.TaskID, taskEntity.TaskIDPrefix)
	if err != nil {
		return &FunctionResult{
			Success: false,
			Error:   fmt.Sprintf("invalid task ID: %v", err),
		}
	}

	// Get task name before deletion for the action
	existingTask, err := fe.taskService.GetTaskByID(ctx, taskID)
	if err != nil {
		return &FunctionResult{
			Success: false,
			Error:   fmt.Sprintf("task not found: %v", err),
		}
	}

	taskName := existingTask.Name

	if err := fe.taskService.DeleteTask(ctx, taskID); err != nil {
		return &FunctionResult{
			Success: false,
			Error:   fmt.Sprintf("failed to delete task: %v", err),
		}
	}

	return &FunctionResult{
		Success: true,
		Data: map[string]interface{}{
			"task_id": deleteArgs.TaskID,
			"deleted": true,
		},
		Action: &TaskAction{
			Type:   "deleted",
			TaskID: taskID,
			Name:   taskName,
		},
	}
}

func (fe *FunctionExecutor) executeGetTask(ctx context.Context, args json.RawMessage) *FunctionResult {
	var getArgs GetTaskArgs
	if err := json.Unmarshal(args, &getArgs); err != nil {
		return &FunctionResult{
			Success: false,
			Error:   fmt.Sprintf("invalid arguments: %v", err),
		}
	}

	taskID, err := utils.ParseID(getArgs.TaskID, taskEntity.TaskIDPrefix)
	if err != nil {
		return &FunctionResult{
			Success: false,
			Error:   fmt.Sprintf("invalid task ID: %v", err),
		}
	}

	tsk, err := fe.taskService.GetTaskByID(ctx, taskID)
	if err != nil {
		return &FunctionResult{
			Success: false,
			Error:   fmt.Sprintf("task not found: %v", err),
		}
	}

	return &FunctionResult{
		Success: true,
		Data:    formatTaskForAI(tsk),
	}
}

func (fe *FunctionExecutor) executeListTasks(ctx context.Context, projectID uuid.UUID, args json.RawMessage) *FunctionResult {
	var listArgs ListTasksArgs
	if args != nil && len(args) > 0 {
		json.Unmarshal(args, &listArgs)
	}

	tasks, err := fe.taskService.ListTasksByProject(ctx, projectID)
	if err != nil {
		return &FunctionResult{
			Success: false,
			Error:   fmt.Sprintf("failed to list tasks: %v", err),
		}
	}

	// Filter by status if specified
	if listArgs.Status != nil && *listArgs.Status != "" {
		filtered := make([]*taskEntity.Task, 0)
		for _, t := range tasks {
			if strings.EqualFold(t.Status, *listArgs.Status) {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}

	taskList := make([]map[string]interface{}, len(tasks))
	for i, t := range tasks {
		taskList[i] = formatTaskForAI(t)
	}

	return &FunctionResult{
		Success: true,
		Data: map[string]interface{}{
			"tasks": taskList,
			"count": len(taskList),
		},
	}
}

func formatTaskForAI(t *taskEntity.Task) map[string]interface{} {
	taskID := utils.ShortUUIDWithPrefix(t.ID, taskEntity.TaskIDPrefix)

	result := map[string]interface{}{
		"task_id":    taskID,
		"name":       t.Name,
		"status":     t.Status,
		"priority":   t.Priority,
		"created_at": t.CreatedAt.Format(time.RFC3339),
	}

	if t.Description != nil {
		result["description"] = *t.Description
	}
	if t.StartDateTime != nil {
		result["start_datetime"] = *t.StartDateTime
	}
	if t.EndDateTime != nil {
		result["end_datetime"] = *t.EndDateTime
	}
	if t.Location != nil {
		result["location"] = *t.Location
	}

	return result
}
