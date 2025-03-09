package TaskManager

import (
	"GServer/Logger"
	"context"
	"sync"
	"time"
	"unsafe"
)

const (
	UNLIMITED_THREAD_COUNT = 0
	DISABLED_TASK_DELAY    = -1
)

type TaskCallback func(*Task)
type TaskManagerContextCancelCallback = context.CancelFunc

type TaskSafeLoopConditionFunction func(*TaskSafeLoop) bool
type TaskSafeLoopBodyFunction func(*TaskSafeLoop)

type taskSignalChannel chan bool

type Task struct {
	Id       uintptr
	Callback TaskCallback

	Started  bool
	Finished bool
	Joined   bool

	Delay time.Duration

	Manager *TaskManager

	SafeLoops map[uintptr]*TaskSafeLoop

	doneChannel taskSignalChannel

	waitingForDoneChannel bool
}

type TaskSafeLoop struct {
	Condition TaskSafeLoopConditionFunction
	Body      TaskSafeLoopBodyFunction

	Breaked bool

	AttachedTo *Task
}

type TaskManager struct {
	Name string

	TaskCount      int
	TasksStarted   int
	MaxmiumThreads int

	Tasks        []*Task
	StartedTasks []*Task

	Paused  bool
	Started bool
	Joined  bool

	Context       context.Context
	ContextCancel TaskManagerContextCancelCallback

	pauseChannel       taskSignalChannel
	waitForTaskChannel taskSignalChannel
	doneTaskChannel    taskSignalChannel
	joinChannel        taskSignalChannel

	waitingForPauseChannel    bool
	waitingForTaskChannel     bool
	waitingForDoneTaskChannel bool
	waitingForJoinChannel     bool

	mutex sync.Mutex
}

var Tasks map[string]*TaskManager = nil
var globalTasksMutex sync.Mutex

var MainContext context.Context = nil
var MainContextCancel context.CancelFunc = nil

func safeWaitForChannel(taskManager *TaskManager, channel taskSignalChannel, stillHasCondition func() bool) {
	for {
		if !taskManager.Started {
			return
		}

		if !stillHasCondition() {
			return
		}

		select {
		case <-channel:
			return
		case <-taskManager.Context.Done():
			taskManager.Stop()
			return
		case <-time.After(time.Second * 5):
		}
	}
}

func (this *Task) GetId() uintptr {
	return (uintptr)(unsafe.Pointer(this))
}

func (this *Task) Done() {
	if this.Finished || !this.Started {
		return
	}

	this.Manager.mutex.Lock()

	var safeLoops map[uintptr]*TaskSafeLoop = map[uintptr]*TaskSafeLoop{}

	for address, safeLoop := range this.SafeLoops {
		safeLoops[address] = safeLoop
	}

	for address, safeLoop := range safeLoops {
		safeLoop.Break()

		delete(this.SafeLoops, address)
	}

	this.Finished = true

	if this.Manager.TasksStarted > 0 {
		this.Manager.TasksStarted -= 1
	}

	if this.Manager.waitingForDoneTaskChannel {
		this.Manager.doneTaskChannel <- true
		this.Manager.waitingForDoneTaskChannel = false
	}

	for index, task := range this.Manager.StartedTasks {
		if task != this {
			continue
		}

		this.Manager.StartedTasks = append(this.Manager.StartedTasks[:index], this.Manager.StartedTasks[index+1:]...)
	}

	if this.Joined && this.waitingForDoneChannel {
		this.doneChannel <- true
		this.waitingForDoneChannel = false
	}

	this.Manager.mutex.Unlock()
}

func (this *Task) Join() {
	if this.Joined {
		return
	}

	this.Manager.mutex.Lock()

	this.Joined = true

	this.waitingForDoneChannel = true

	this.Manager.mutex.Unlock()

	safeWaitForChannel(this.Manager, this.doneChannel, func() bool { return this.waitingForDoneChannel })
}

func (this *Task) NewSafeLoop() *TaskSafeLoop {
	var safeLoop *TaskSafeLoop = new(TaskSafeLoop)

	safeLoop.Condition = func(safeLoop *TaskSafeLoop) bool { return false }
	safeLoop.Body = func(safeLoop *TaskSafeLoop) {}

	safeLoop.Breaked = false

	safeLoop.AttachedTo = this

	this.SafeLoops[(uintptr)(unsafe.Pointer(safeLoop))] = safeLoop

	return safeLoop
}

func (this *Task) SafeLoop(condition TaskSafeLoopConditionFunction, body TaskSafeLoopBodyFunction) *TaskSafeLoop {
	var safeLoop *TaskSafeLoop = this.NewSafeLoop()

	safeLoop.Condition = condition
	safeLoop.Body = body

	safeLoop.Start()

	return safeLoop
}

func (this *TaskSafeLoop) Break() {
	this.Breaked = true
}

func (this *TaskSafeLoop) Start() bool {
	this.Breaked = false

	var succeed = false

	var task *Task = this.AttachedTo

	if task == nil {
		goto Cleanup
	}

Loop:
	{
		if !task.Started || task.Finished {
			goto Cleanup
		}

		if this.Breaked {
			succeed = true
			goto Cleanup
		}

		if !this.Condition(this) {
			succeed = true
			goto Cleanup
		}

		this.Body(this)

		goto Loop
	}

Cleanup:
	{
		for address, safeLoop := range task.SafeLoops {
			if safeLoop != this {
				continue
			}

			delete(task.SafeLoops, address)

			break
		}

		return succeed
	}
}

func (this *TaskManager) Add(xA ...int) {
	var x int = 1

	if len(xA) > 0 {
		x = xA[0]
	}

	this.TaskCount += x
}

func (this *TaskManager) AddTask(callback TaskCallback) *Task {
	this.mutex.Lock()

	var lastTaskCount int = len(this.Tasks)

	var task *Task = new(Task)

	task.Id = task.GetId()
	task.Callback = callback

	task.Started = false
	task.Finished = false
	task.Joined = false

	task.Delay = DISABLED_TASK_DELAY

	task.Manager = this

	task.SafeLoops = map[uintptr]*TaskSafeLoop{}

	task.doneChannel = make(taskSignalChannel, 1)

	task.waitingForDoneChannel = false

	this.Tasks = append(this.Tasks, task)

	this.Add()

	if lastTaskCount < 1 && this.waitingForTaskChannel {
		this.waitForTaskChannel <- true
		this.waitingForTaskChannel = false
	}

	this.mutex.Unlock()

	return task
}

func (this *TaskManager) AddTaskWithDelay(task TaskCallback, delay time.Duration) *Task {
	var t *Task = this.AddTask(task)

	t.Delay = delay

	return t
}

func (this *TaskManager) resumeProcessorThread() {
	this.mutex.Lock()

	if this.waitingForPauseChannel {
		this.pauseChannel <- false
		this.waitingForPauseChannel = false
	}

	if this.waitingForTaskChannel {
		this.waitForTaskChannel <- false
		this.waitingForTaskChannel = false
	}

	if this.waitingForDoneTaskChannel {
		this.doneTaskChannel <- false
		this.waitingForDoneTaskChannel = false
	}

	this.mutex.Unlock()
}

func (this *TaskManager) Pause() {
	this.mutex.Lock()

	this.Paused = true

	this.mutex.Unlock()

	if this.Paused {
		return
	}

	this.resumeProcessorThread()
}

func (this *TaskManager) processTasks() {
	for this.Started {
		this.mutex.Lock()

		paused := this.Paused
		taskCount := len(this.Tasks)
		reachedMaximum := this.MaxmiumThreads != UNLIMITED_THREAD_COUNT && this.TasksStarted >= this.MaxmiumThreads && len(this.StartedTasks) >= this.MaxmiumThreads

		this.mutex.Unlock()

		if paused {
			this.mutex.Lock()

			this.waitingForPauseChannel = true

			this.mutex.Unlock()

			safeWaitForChannel(this, this.pauseChannel, func() bool { return this.waitingForPauseChannel })

			continue
		}

		if taskCount < 1 {
			this.mutex.Lock()

			this.waitingForTaskChannel = true

			this.mutex.Unlock()

			safeWaitForChannel(this, this.waitForTaskChannel, func() bool { return this.waitingForTaskChannel })

			continue
		}

		if reachedMaximum {
			this.mutex.Lock()

			this.waitingForDoneTaskChannel = true

			this.mutex.Unlock()

			safeWaitForChannel(this, this.doneTaskChannel, func() bool { return this.waitingForDoneTaskChannel })

			continue
		}

		var task *Task = this.Tasks[0]

		if task.Delay != DISABLED_TASK_DELAY {
			time.Sleep(task.Delay)
		}

		go (func() {
			if !this.Started || this.Paused {
				this.mutex.Lock()

				for index, t := range this.StartedTasks {
					if t != task {
						continue
					}

					this.StartedTasks = append(this.StartedTasks[:index], this.StartedTasks[index+1:]...)
				}

				this.Tasks = append(this.Tasks, task)

				if this.TasksStarted > 0 {
					this.TasksStarted -= 1
				}

				this.mutex.Unlock()

				task.Done()

				return
			}

			this.mutex.Lock()

			task.Started = true

			this.mutex.Unlock()

			task.Callback(task)

			task.Done()
		})()

		this.mutex.Lock()

		this.StartedTasks = append(this.StartedTasks, task)
		this.Tasks = this.Tasks[1:]

		this.TasksStarted += 1

		this.mutex.Unlock()
	}
}

func (this *TaskManager) Start() {
	this.mutex.Lock()

	var lastPauseState bool = this.Paused

	this.Paused = false

	this.mutex.Unlock()

	if !lastPauseState {
		return
	}

	this.resumeProcessorThread()

	if this.Started {
		return
	}

	go this.processTasks()

	this.mutex.Lock()

	this.Started = true

	this.mutex.Unlock()
}

func (this *TaskManager) Stop() {
	if !this.Started {
		return
	}

	this.mutex.Lock()

	this.Started = false

	startedTasks := []*Task{}

	startedTasks = append(startedTasks, this.StartedTasks...)

	this.mutex.Unlock()

	for _, task := range startedTasks {
		task.Done()
	}

	this.mutex.Lock()

	if this.waitingForJoinChannel {
		this.joinChannel <- true
		this.waitingForJoinChannel = false
	}

	this.mutex.Unlock()
}

func (this *TaskManager) WaitForTasks() {
	for {
		this.mutex.Lock()

		pending := len(this.Tasks) + len(this.StartedTasks)

		this.mutex.Unlock()

		if pending == 0 {
			break
		}

		time.Sleep(time.Millisecond * 100)
	}
}

func (this *TaskManager) Join() {
	if this.Joined || !this.Started {
		return
	}

	this.mutex.Lock()

	this.Joined = true

	this.waitingForJoinChannel = true

	this.mutex.Unlock()

	safeWaitForChannel(this, this.joinChannel, func() bool { return this.waitingForJoinChannel })
}

func ExistsTaskManager(name string) bool {
	globalTasksMutex.Lock()

	_, exists := Tasks[name]

	defer globalTasksMutex.Unlock()

	return exists
}

func DeleteTaskManager(name string) {
	if !ExistsTaskManager(name) {
		return
	}

	globalTasksMutex.Lock()

	var taskManager *TaskManager = Tasks[name]

	taskManager.Stop()

	delete(Tasks, name)

	defer globalTasksMutex.Unlock()
}

func GetTaskManager(name string) *TaskManager {
	globalTasksMutex.Lock()
	defer globalTasksMutex.Unlock()

	return Tasks[name]
}

func CreateTaskManager(name string, maximumThreads int) *TaskManager {
	DeleteTaskManager(name)

	globalTasksMutex.Lock()

	var taskManager *TaskManager = new(TaskManager)

	taskManager.Name = name

	taskManager.TaskCount = 0
	taskManager.TasksStarted = 0
	taskManager.MaxmiumThreads = maximumThreads

	taskManager.Tasks = []*Task{}
	taskManager.StartedTasks = []*Task{}

	taskManager.Paused = true
	taskManager.Started = false
	taskManager.Joined = false

	ctx, ctxCancel := context.WithCancel(MainContext)

	taskManager.Context = ctx
	taskManager.ContextCancel = func() {
		taskManager.Stop()

		ctxCancel()
	}

	taskManager.pauseChannel = make(taskSignalChannel, 1)
	taskManager.waitForTaskChannel = make(taskSignalChannel, 1)
	taskManager.doneTaskChannel = make(taskSignalChannel, 1)
	taskManager.joinChannel = make(taskSignalChannel, 1)

	taskManager.waitingForPauseChannel = false
	taskManager.waitingForTaskChannel = false
	taskManager.waitingForDoneTaskChannel = false
	taskManager.waitingForJoinChannel = false

	taskManager.mutex = sync.Mutex{}

	Tasks[name] = taskManager

	defer globalTasksMutex.Unlock()

	return taskManager
}

func CreateTaskManagerWithContext(ctx context.Context, name string, maximumThreads int) *TaskManager {
	var taskManager *TaskManager = CreateTaskManager(name, maximumThreads)

	tCtx, tCtxCancel := context.WithCancel(ctx)

	taskManager.Context = tCtx
	taskManager.ContextCancel = func() {
		taskManager.Stop()

		tCtxCancel()
	}

	return taskManager
}

func Wait() {
	for _, taskManager := range Tasks {
		taskManager.Join()
	}
}

func Initialize() {
	Logger.INFO("Initializing task manager...")

	Tasks = make(map[string]*TaskManager)

	MainContext, MainContextCancel = context.WithCancel(context.Background())

	globalTasksMutex = sync.Mutex{}

	Logger.INFO("Task manager initialized.")
}

func Uninitialize() {
	Logger.INFO("Uninitializing task manager...")

	globalTasksMutex.Lock()

	tasks := make(map[string]*TaskManager)

	for taskName, task := range tasks {
		tasks[taskName] = task
	}

	for taskName := range tasks {
		DeleteTaskManager(taskName)
	}

	defer globalTasksMutex.Unlock()

	Logger.INFO("Task manager uninitialized.")
}
