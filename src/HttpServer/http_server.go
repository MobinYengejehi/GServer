package HttpServer

import (
	"GServer/Config"
	"GServer/Logger"
	"GServer/TaskManager"
	HTTP "net/http"
	"time"
)

var Tasks *TaskManager.TaskManager = nil

func serverListen(serverHostAddress string) {
	err := HTTP.ListenAndServe(serverHostAddress, nil)

	if err == nil {
		return
	}

	Logger.ERROR("Couldn't start HTTP server. [Mesasge: ", err.Error(), "]")
}

func Initialize() {
	var serverHostAddress string = Config.Main.HttpHostAddress

	Logger.INFO("Starting HTTP server on '" + serverHostAddress + "' ...")

	Tasks = TaskManager.CreateTaskManager("HTTP_SERVER", TaskManager.UNLIMITED_THREAD_COUNT)

	HTTP.HandleFunc("/", h_NotFound)
	HTTP.HandleFunc("/add", h_Add)

	Tasks.AddTask(func(task *TaskManager.Task) {
		serverListen(serverHostAddress)
	})

	Tasks.Start()

	time.Sleep(time.Second * 1)

	Logger.INFO("HTTP server started and listening to requests on '" + serverHostAddress + "'.")
}

func Uninitialize() {
	Logger.INFO("Uninitializing HTTP server ...")

	TaskManager.DeleteTaskManager(Tasks.Name)

	Logger.INFO("HTTP server uninitialized.")
}
