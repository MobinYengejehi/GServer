package main

import (
	"GServer/Config"
	"GServer/Crawler"
	"GServer/HttpServer"
	"GServer/Logger"
	"GServer/TaskManager"
	"GServer/YTS"
	"context"
	"fmt"
	"time"
)

func main() {
	fmt.Println(
		" ██████╗ ███████╗███████╗██████╗ ██╗   ██╗███████╗██████╗ \n" +
			"██╔════╝ ██╔════╝██╔════╝██╔══██╗██║   ██║██╔════╝██╔══██╗\n" +
			"██║  ███╗███████╗█████╗  ██████╔╝██║   ██║█████╗  ██████╔╝\n" +
			"██║   ██║╚════██║██╔══╝  ██╔══██╗╚██╗ ██╔╝██╔══╝  ██╔══██╗\n" +
			"╚██████╔╝███████║███████╗██║  ██║ ╚████╔╝ ███████╗██║  ██║\n" +
			" ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝\n",
	)

	TaskManager.Initialize()
	Config.Initialize()
	HttpServer.Initialize()
	Crawler.Initialize()

	{
		ytsClient := YTS.NewClient(context.WithoutCancel(TaskManager.MainContext), time.Second*5)

		if ytsClient != nil {
			params := YTS.NewMoviesListParameters()

			params.Limit = 50

			movies, err, movieCount := ytsClient.GetMovieList(params)

			if err != nil {
				Logger.ERROR("yts err is : ", err.Error())
			} else {
				Logger.INFO("yts movies are : ", movies[0].Title, " | ", movies[0].Torrents[0].MainFile, " | ", movieCount)
			}
		}
	}

	Logger.INFO("getting film info done!")

	TaskManager.Wait()

	Crawler.Uninitialize()
	HttpServer.Uninitialize()
	Config.Uninitialize()
	TaskManager.Uninitialize()
}
