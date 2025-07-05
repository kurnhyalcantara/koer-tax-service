package log

import (
	"time"

	"go.uber.org/zap"
)

type LogPayload struct {
	ProcessID    string
	FunctionName string
	TaskID       uint64
	Message      string
	StartTime    *time.Time
	EndTime      *time.Time
	ResponseTime *time.Time
	Metadata     interface{}
}

// Debug implements LoggerCore.
func (l *Logger) Debug(payload LogPayload) {
	fields := []zap.Field{
		zap.String("product_name", l.productName),
		zap.String("service_name", l.serviceName),
		zap.String("host_name", l.hostName),
		zap.String("log_type", "application"),
	}

	// Merge data from logger state
	if l.functionName != "" {
		fields = append(fields, zap.String("function_name", l.functionName))
	}
	if l.processId != "" {
		fields = append(fields, zap.String("process_id", l.processId))
	}
	if l.taskId != 0 {
		fields = append(fields, zap.Uint64("task_id", l.taskId))
	}
	if l.startTime != nil {
		fields = append(fields, zap.Time("start_time", *l.startTime))
	}
	if l.endTime != nil {
		fields = append(fields, zap.Time("end_time", *l.endTime))
	}

	// Use payload override if available
	if payload.FunctionName != "" {
		fields = append(fields, zap.String("function_name", payload.FunctionName))
	}
	if payload.ProcessID != "" {
		fields = append(fields, zap.String("process_id", payload.ProcessID))
	}
	if payload.TaskID != 0 {
		fields = append(fields, zap.Uint64("task_id", payload.TaskID))
	}
	if payload.StartTime != nil {
		fields = append(fields, zap.Time("start_time", *payload.StartTime))
	}
	if payload.EndTime != nil {
		fields = append(fields, zap.Time("end_time", *payload.EndTime))
	}

	// Calculate response time
	start := payload.StartTime
	end := payload.EndTime
	if start == nil {
		start = l.startTime
	}
	if end == nil {
		end = l.endTime
	}
	if start != nil && end != nil {
		rt := end.Sub(*start)
		fields = append(fields, zap.Duration("response_time", rt))
	}

	// Add metadata
	if payload.Metadata != nil {
		fields = append(fields, zap.Any("metadata", payload.Metadata))
	}

	// Message is required
	msg := payload.Message
	if msg == "" {
		msg = "debug"
	}

	l.zapLog.Debug(msg, fields...)
}

// EndTime implements LoggerCore.
func (l *Logger) EndTime() LoggerCore {
	newLogger := *l
	now := time.Now()
	newLogger.endTime = &now
	return &newLogger
}

// Error implements LoggerCore.
func (l *Logger) Error(payload LogPayload) {
	fields := []zap.Field{
		zap.String("product_name", l.productName),
		zap.String("service_name", l.serviceName),
		zap.String("host_name", l.hostName),
		zap.String("log_type", "application"),
	}

	// Merge data from logger state
	if l.functionName != "" {
		fields = append(fields, zap.String("function_name", l.functionName))
	}
	if l.processId != "" {
		fields = append(fields, zap.String("process_id", l.processId))
	}
	if l.taskId != 0 {
		fields = append(fields, zap.Uint64("task_id", l.taskId))
	}
	if l.startTime != nil {
		fields = append(fields, zap.Time("start_time", *l.startTime))
	}
	if l.endTime != nil {
		fields = append(fields, zap.Time("end_time", *l.endTime))
	}

	// Use payload override if available
	if payload.FunctionName != "" {
		fields = append(fields, zap.String("function_name", payload.FunctionName))
	}
	if payload.ProcessID != "" {
		fields = append(fields, zap.String("process_id", payload.ProcessID))
	}
	if payload.TaskID != 0 {
		fields = append(fields, zap.Uint64("task_id", payload.TaskID))
	}
	if payload.StartTime != nil {
		fields = append(fields, zap.Time("start_time", *payload.StartTime))
	}
	if payload.EndTime != nil {
		fields = append(fields, zap.Time("end_time", *payload.EndTime))
	}

	// Calculate response time
	start := payload.StartTime
	end := payload.EndTime
	if start == nil {
		start = l.startTime
	}
	if end == nil {
		end = l.endTime
	}
	if start != nil && end != nil {
		rt := end.Sub(*start)
		fields = append(fields, zap.Duration("response_time", rt))
	}

	// Add metadata
	if payload.Metadata != nil {
		fields = append(fields, zap.Any("metadata", payload.Metadata))
	}

	// Message is required
	msg := payload.Message
	if msg == "" {
		msg = "error"
	}

	l.zapLog.Error(msg, fields...)
}

// Fatal implements LoggerCore.
func (l *Logger) Fatal(payload LogPayload) {
	fields := []zap.Field{
		zap.String("product_name", l.productName),
		zap.String("service_name", l.serviceName),
		zap.String("host_name", l.hostName),
		zap.String("log_type", "application"),
	}

	// Merge data from logger state
	if l.functionName != "" {
		fields = append(fields, zap.String("function_name", l.functionName))
	}
	if l.processId != "" {
		fields = append(fields, zap.String("process_id", l.processId))
	}
	if l.taskId != 0 {
		fields = append(fields, zap.Uint64("task_id", l.taskId))
	}
	if l.startTime != nil {
		fields = append(fields, zap.Time("start_time", *l.startTime))
	}
	if l.endTime != nil {
		fields = append(fields, zap.Time("end_time", *l.endTime))
	}

	// Use payload override if available
	if payload.FunctionName != "" {
		fields = append(fields, zap.String("function_name", payload.FunctionName))
	}
	if payload.ProcessID != "" {
		fields = append(fields, zap.String("process_id", payload.ProcessID))
	}
	if payload.TaskID != 0 {
		fields = append(fields, zap.Uint64("task_id", payload.TaskID))
	}
	if payload.StartTime != nil {
		fields = append(fields, zap.Time("start_time", *payload.StartTime))
	}
	if payload.EndTime != nil {
		fields = append(fields, zap.Time("end_time", *payload.EndTime))
	}

	// Calculate response time
	start := payload.StartTime
	end := payload.EndTime
	if start == nil {
		start = l.startTime
	}
	if end == nil {
		end = l.endTime
	}
	if start != nil && end != nil {
		rt := end.Sub(*start)
		fields = append(fields, zap.Duration("response_time", rt))
	}

	// Add metadata
	if payload.Metadata != nil {
		fields = append(fields, zap.Any("metadata", payload.Metadata))
	}

	// Message is required
	msg := payload.Message
	if msg == "" {
		msg = "fatal"
	}

	l.zapLog.Fatal(msg, fields...)
}

// Info implements LoggerCore.
func (l *Logger) Info(payload LogPayload) {
	fields := []zap.Field{
		zap.String("product_name", l.productName),
		zap.String("service_name", l.serviceName),
		zap.String("host_name", l.hostName),
		zap.String("log_type", "application"),
	}

	// Merge data from logger state
	if l.functionName != "" {
		fields = append(fields, zap.String("function_name", l.functionName))
	}
	if l.processId != "" {
		fields = append(fields, zap.String("process_id", l.processId))
	}
	if l.taskId != 0 {
		fields = append(fields, zap.Uint64("task_id", l.taskId))
	}
	if l.startTime != nil {
		fields = append(fields, zap.Time("start_time", *l.startTime))
	}
	if l.endTime != nil {
		fields = append(fields, zap.Time("end_time", *l.endTime))
	}

	// Use payload override if available
	if payload.FunctionName != "" {
		fields = append(fields, zap.String("function_name", payload.FunctionName))
	}
	if payload.ProcessID != "" {
		fields = append(fields, zap.String("process_id", payload.ProcessID))
	}
	if payload.TaskID != 0 {
		fields = append(fields, zap.Uint64("task_id", payload.TaskID))
	}
	if payload.StartTime != nil {
		fields = append(fields, zap.Time("start_time", *payload.StartTime))
	}
	if payload.EndTime != nil {
		fields = append(fields, zap.Time("end_time", *payload.EndTime))
	}

	// Calculate response time
	start := payload.StartTime
	end := payload.EndTime
	if start == nil {
		start = l.startTime
	}
	if end == nil {
		end = l.endTime
	}
	if start != nil && end != nil {
		rt := end.Sub(*start)
		fields = append(fields, zap.Duration("response_time", rt))
	}

	// Add metadata
	if payload.Metadata != nil {
		fields = append(fields, zap.Any("metadata", payload.Metadata))
	}

	// Message is required
	msg := payload.Message
	if msg == "" {
		msg = "info"
	}

	l.zapLog.Info(msg, fields...)
}

// QueueMessageError implements LoggerCore.
func (l *Logger) QueueMessageError(payload LogPayload) {
	fields := []zap.Field{
		zap.String("product_name", l.productName),
		zap.String("service_name", l.serviceName),
		zap.String("host_name", l.hostName),
		zap.String("log_type", "queue_event"),
	}

	// Merge data from logger state
	if l.functionName != "" {
		fields = append(fields, zap.String("function_name", l.functionName))
	}
	if l.processId != "" {
		fields = append(fields, zap.String("process_id", l.processId))
	}
	if l.taskId != 0 {
		fields = append(fields, zap.Uint64("task_id", l.taskId))
	}
	if l.startTime != nil {
		fields = append(fields, zap.Time("start_time", *l.startTime))
	}
	if l.endTime != nil {
		fields = append(fields, zap.Time("end_time", *l.endTime))
	}

	// Use payload override if available
	if payload.FunctionName != "" {
		fields = append(fields, zap.String("function_name", payload.FunctionName))
	}
	if payload.ProcessID != "" {
		fields = append(fields, zap.String("process_id", payload.ProcessID))
	}
	if payload.TaskID != 0 {
		fields = append(fields, zap.Uint64("task_id", payload.TaskID))
	}
	if payload.StartTime != nil {
		fields = append(fields, zap.Time("start_time", *payload.StartTime))
	}
	if payload.EndTime != nil {
		fields = append(fields, zap.Time("end_time", *payload.EndTime))
	}

	// Calculate response time
	start := payload.StartTime
	end := payload.EndTime
	if start == nil {
		start = l.startTime
	}
	if end == nil {
		end = l.endTime
	}
	if start != nil && end != nil {
		rt := end.Sub(*start)
		fields = append(fields, zap.Duration("response_time", rt))
	}

	// Add metadata
	if payload.Metadata != nil {
		fields = append(fields, zap.Any("metadata", payload.Metadata))
	}

	// Message is required
	msg := payload.Message
	if msg == "" {
		msg = "queue_error"
	}

	l.zapLog.Error(msg, fields...)
}

// QueueMessageInfo implements LoggerCore.
func (l *Logger) QueueMessageInfo(payload LogPayload) {
	fields := []zap.Field{
		zap.String("product_name", l.productName),
		zap.String("service_name", l.serviceName),
		zap.String("host_name", l.hostName),
		zap.String("log_type", "queue_event"),
	}

	// Merge data from logger state
	if l.functionName != "" {
		fields = append(fields, zap.String("function_name", l.functionName))
	}
	if l.processId != "" {
		fields = append(fields, zap.String("process_id", l.processId))
	}
	if l.taskId != 0 {
		fields = append(fields, zap.Uint64("task_id", l.taskId))
	}
	if l.startTime != nil {
		fields = append(fields, zap.Time("start_time", *l.startTime))
	}
	if l.endTime != nil {
		fields = append(fields, zap.Time("end_time", *l.endTime))
	}

	// Use payload override if available
	if payload.FunctionName != "" {
		fields = append(fields, zap.String("function_name", payload.FunctionName))
	}
	if payload.ProcessID != "" {
		fields = append(fields, zap.String("process_id", payload.ProcessID))
	}
	if payload.TaskID != 0 {
		fields = append(fields, zap.Uint64("task_id", payload.TaskID))
	}
	if payload.StartTime != nil {
		fields = append(fields, zap.Time("start_time", *payload.StartTime))
	}
	if payload.EndTime != nil {
		fields = append(fields, zap.Time("end_time", *payload.EndTime))
	}

	// Calculate response time
	start := payload.StartTime
	end := payload.EndTime
	if start == nil {
		start = l.startTime
	}
	if end == nil {
		end = l.endTime
	}
	if start != nil && end != nil {
		rt := end.Sub(*start)
		fields = append(fields, zap.Duration("response_time", rt))
	}

	// Add metadata
	if payload.Metadata != nil {
		fields = append(fields, zap.Any("metadata", payload.Metadata))
	}

	// Message is required
	msg := payload.Message
	if msg == "" {
		msg = "queue_info"
	}

	l.zapLog.Info(msg, fields...)
}

// Warn implements LoggerCore.
func (l *Logger) Warn(payload LogPayload) {
	fields := []zap.Field{
		zap.String("product_name", l.productName),
		zap.String("service_name", l.serviceName),
		zap.String("host_name", l.hostName),
		zap.String("log_type", "application"),
	}

	// Merge data from logger state
	if l.functionName != "" {
		fields = append(fields, zap.String("function_name", l.functionName))
	}
	if l.processId != "" {
		fields = append(fields, zap.String("process_id", l.processId))
	}
	if l.taskId != 0 {
		fields = append(fields, zap.Uint64("task_id", l.taskId))
	}
	if l.startTime != nil {
		fields = append(fields, zap.Time("start_time", *l.startTime))
	}
	if l.endTime != nil {
		fields = append(fields, zap.Time("end_time", *l.endTime))
	}

	// Use payload override if available
	if payload.FunctionName != "" {
		fields = append(fields, zap.String("function_name", payload.FunctionName))
	}
	if payload.ProcessID != "" {
		fields = append(fields, zap.String("process_id", payload.ProcessID))
	}
	if payload.TaskID != 0 {
		fields = append(fields, zap.Uint64("task_id", payload.TaskID))
	}
	if payload.StartTime != nil {
		fields = append(fields, zap.Time("start_time", *payload.StartTime))
	}
	if payload.EndTime != nil {
		fields = append(fields, zap.Time("end_time", *payload.EndTime))
	}

	// Calculate response time
	start := payload.StartTime
	end := payload.EndTime
	if start == nil {
		start = l.startTime
	}
	if end == nil {
		end = l.endTime
	}
	if start != nil && end != nil {
		rt := end.Sub(*start)
		fields = append(fields, zap.Duration("response_time", rt))
	}

	// Add metadata
	if payload.Metadata != nil {
		fields = append(fields, zap.Any("metadata", payload.Metadata))
	}

	// Message is required
	msg := payload.Message
	if msg == "" {
		msg = "warn"
	}

	l.zapLog.Warn(msg, fields...)
}
