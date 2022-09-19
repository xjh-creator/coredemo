package cobra

import (
	"coredemo/framework"
	"github.com/robfig/cron/v3"
	"log"
)

// SetContainer 设置服务容器
func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

// GetContainer 获取容器
func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}

// CommandSpec
type CommandSpec struct {
	Cmd  *Command
	Args []string
	Spec string
}

// AddCronCommand 创建一个 Cron 任务的
func (c *Command) AddCronCommand(spec string, cmd *Command, args ...string) {
	root := c.Root()
	if root.Cron == nil {
		// 初始化cron
		root.Cron = cron.New()
		root.CronSepcs = []CommandSpec{}
	}

	// 增加说明信息
	root.CronSepcs = append(root.CronSepcs, CommandSpec{
		Cmd:  cmd,
		Args: args,
		Spec: spec,
	})
	// 增加调用函数
	root.Cron.AddFunc(spec, func() {
		// 如果后续的command出现panic，这里要捕获
		var cronCmd Command
		ctx := root.Context()
		cronCmd = *cmd
		cronCmd.SetParentNull()
		cronCmd.args = []string{}
		err := cronCmd.ExecuteContext(ctx)
		if err != nil {
			log.Println(err)
		}
	})
}

// SetParentNull set parent
func (c *Command) SetParentNull() {
	c.parent = nil
}
