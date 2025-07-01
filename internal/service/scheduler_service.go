// internal/service/scheduler_service.go (最终版)
package service

import (
	"gotasksys/internal/model"
	"gotasksys/internal/repository"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// 将调度器实例和它的任务ID映射设为全局，以便其他服务可以访问
var (
	cronScheduler *cron.Cron
	jobIDs        map[string]cron.EntryID
)

// InitScheduler 初始化并启动所有定时任务
func InitScheduler() {
	log.Println("Initializing scheduler...")
	jobIDs = make(map[string]cron.EntryID)
	// 使用带秒级的解析器，以支持更灵活的测试
	cronScheduler = cron.New(cron.WithSeconds())

	periodicTasks, err := repository.ListAllActivePeriodicTasks()
	if err != nil {
		log.Printf("Error fetching periodic tasks on init: %v", err)
		return
	}

	for _, pt := range periodicTasks {
		AddJob(pt)
	}

	go cronScheduler.Start()
	log.Println("Scheduler started.")
}

// AddJob 向运行中的调度器添加一个新作业
func AddJob(pt model.PeriodicTask) {
	// 使用 'pt' 的副本，以避免闭包问题
	taskToSchedule := pt

	id, err := cronScheduler.AddFunc(taskToSchedule.CronExpression, func() {
		log.Printf("Running periodic task: %s", taskToSchedule.Title)
		createTaskFromRule(taskToSchedule)
	})

	if err != nil {
		log.Printf("Error scheduling task '%s': %v", taskToSchedule.Title, err)
		return
	}
	jobIDs[pt.ID.String()] = id
	log.Printf("Scheduled task '%s' (ID: %s) with schedule: %s", pt.Title, pt.ID.String(), pt.CronExpression)
}

// RemoveJob 从运行中的调度器移除一个作业
func RemoveJob(periodicTaskID string) {
	if entryID, ok := jobIDs[periodicTaskID]; ok {
		cronScheduler.Remove(entryID)
		delete(jobIDs, periodicTaskID)
		log.Printf("Unscheduled task for rule ID: %s", periodicTaskID)
	}
}

// createTaskFromRule 是实际执行创建任务的函数
func createTaskFromRule(pt model.PeriodicTask) {
	// 为任务标题添加日期戳，方便识别
	taskTitle := pt.Title + " - " + time.Now().Format("2006-01-02")

	newTask := model.Task{
		Title:          taskTitle,
		Description:    pt.Description,
		Status:         "in_pool",
		Priority:       pt.DefaultPriority,
		Effort:         pt.DefaultEffort,
		OriginalEffort: pt.DefaultEffort,
		TaskTypeID:     pt.DefaultTaskTypeID,
		CreatorID:      pt.CreatedByID,
		AssigneeID:     pt.DefaultAssigneeID,
	}
	if err := repository.CreateTask(&newTask); err != nil {
		log.Printf("Error creating task from periodic rule '%s': %v", pt.Title, err)
	}
}
