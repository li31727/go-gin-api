package generator_handler

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type gormExecuteRequest struct {
	Tables string `form:"tables"`
}

func (h *handler) GormExecute() core.HandlerFunc {
	dir, _ := os.Getwd()
	projectPath := strings.Replace(dir, "\\", "/", -1)
	gormgenSh := projectPath + "/scripts/gormgen.sh"
	gormgenBat := projectPath + "/scripts/gormgen.bat"

	return func(c core.Context) {
		req := new(gormExecuteRequest)
		if err := c.ShouldBindPostForm(req); err != nil {
			c.Payload("参数传递有误")
			return
		}

		mysqlConf := configs.Get().MySQL.Read
		shellPath := fmt.Sprintf("%s %s %s %s %s %s", gormgenSh, mysqlConf.Addr, mysqlConf.User, mysqlConf.Pass, mysqlConf.Name, req.Tables)
		batPath := fmt.Sprintf("%s %s %s %s %s %s", gormgenBat, mysqlConf.Addr, mysqlConf.User, mysqlConf.Pass, mysqlConf.Name, req.Tables)

		command := new(exec.Cmd)

		if runtime.GOOS == "windows" {
			command = exec.Command("cmd", "/C", batPath)
		} else {
			// runtime.GOOS = linux or darwin
			command = exec.Command("/bin/bash", "-c", shellPath)
		}

		var stderr bytes.Buffer
		command.Stderr = &stderr

		output, err := command.Output()
		if err != nil {
			c.Payload(stderr.String())
			return
		}

		c.Payload(string(output))
	}
}

func (h *handler) PgsqlGormExecute() core.HandlerFunc {
	dir, _ := os.Getwd()
	projectPath := strings.Replace(dir, "\\", "/", -1)
	gormgenSh := projectPath + "/scripts/pgsqlgormgen.sh"
	gormgenBat := projectPath + "/scripts/pgsqlgormgen.bat"

	return func(c core.Context) {
		req := new(gormExecuteRequest)
		if err := c.ShouldBindPostForm(req); err != nil {
			c.Payload("参数传递有误")
			return
		}

		pgsqlConf := configs.Get().PgSQL.Read
		shellPath := fmt.Sprintf("%s %s %s %s %s %s %s", gormgenSh, pgsqlConf.Addr, pgsqlConf.User, pgsqlConf.Pass, pgsqlConf.Name, pgsqlConf.Port, req.Tables)
		batPath := fmt.Sprintf("%s %s %s %s %s %s %s", gormgenBat, pgsqlConf.Addr, pgsqlConf.User, pgsqlConf.Pass, pgsqlConf.Name, pgsqlConf.Port, req.Tables)

		command := new(exec.Cmd)

		if runtime.GOOS == "windows" {
			command = exec.Command("cmd", "/C", batPath)
		} else {
			// runtime.GOOS = linux or darwin
			command = exec.Command("/bin/bash", "-c", shellPath)
		}

		var stderr bytes.Buffer
		command.Stderr = &stderr

		output, err := command.Output()
		if err != nil {
			c.Payload(stderr.String())
			return
		}

		c.Payload(string(output))
	}
}
