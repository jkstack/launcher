package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/jkframe/logging"
)

const port = 19999

var lastUpdate = time.Now()
var onStop = make(chan struct{})

func main() {
	go func() {
		tk := time.NewTicker(time.Second)
		for {
			<-tk.C
			if time.Since(lastUpdate) > time.Minute {
				logging.Error("no request")
				os.Exit(1)
			}
		}
	}()
	r := gin.Default()
	r.POST("/api/batch_run", batchRun)

	svr := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("listen: %s", err)
		}
	}()

	<-onStop
	svr.Shutdown(context.Background())
}

func batchRun(g *gin.Context) {
	lastUpdate = time.Now()
	var req []struct {
		Cmd  string   `json:"cmd" binding:"required"`
		Args []string `json:"args"`
	}
	err := g.ShouldBindJSON(&req)
	if err != nil {
		g.JSON(http.StatusOK, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	var cnt int
	for _, cmd := range req {
		ok := run(cmd.Cmd, cmd.Args)
		if ok {
			cnt++
		}
	}
	g.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": cnt,
	})
	onStop <- struct{}{}
}

func run(cmd string, args []string) bool {
	logging.Info("run: %s %v", cmd, args)
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Start() == nil
}
