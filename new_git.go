package gitgo

import (
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

type Gcm struct {
	cmdConfig *osexec.CommandConfig
	output    []byte
	errorOnce error
	debugMode bool
}

func New(path string) *Gcm {
	return newOkGcm(osexec.NewCommandConfig().WithPath(path).WithDebugMode(debugModeOpen), make([]byte, 0), debugModeOpen)
}

func newOkGcm(cmdConfig *osexec.CommandConfig, output []byte, debugMode bool) *Gcm {
	if debugMode {
		if len(output) > 0 {
			zaplog.ZAPS.P3.SUG.Debugln("done", "message:", "\n"+eroticgo.GREEN.Sprint(string(output))+"\n", "-")
		} else {
			zaplog.ZAPS.P3.SUG.Debugln("done", "\n", "-")
		}
	}
	return &Gcm{
		cmdConfig: cmdConfig,
		output:    output,
		errorOnce: nil,
		debugMode: debugMode,
	}
}

func newWaGcm(cmdConfig *osexec.CommandConfig, output []byte, errorOnce error, debugMode bool) *Gcm {
	if debugMode {
		if len(output) > 0 {
			zaplog.ZAPS.P3.SUG.Errorln("wrong", eroticgo.RED.Sprint(errorOnce), "message:", "\n"+eroticgo.RED.Sprint(string(output))+"\n", "-")
		} else {
			zaplog.ZAPS.P3.SUG.Errorln("wrong", eroticgo.RED.Sprint(errorOnce))
		}
	}
	return &Gcm{
		cmdConfig: cmdConfig,
		output:    output,
		errorOnce: errorOnce,
		debugMode: debugMode,
	}
}

func (G *Gcm) Result() ([]byte, error) {
	return G.output, G.errorOnce
}

// 这个函数不要使用导出的，因为日志里是跳过3层调用栈，假如使用导出的，跳出的栈的层数就不正确啦
func (G *Gcm) do(name string, args ...string) *Gcm {
	if G.errorOnce != nil {
		return G //当出错时就不要再往下执行，直接在这里拦住，这样整个链式的后续动作就都不执行
	}
	output, err := G.cmdConfig.Exec(name, args...)
	if err != nil {
		return newWaGcm(G.cmdConfig, output, err, G.debugMode)
	}
	return newOkGcm(G.cmdConfig, output, G.debugMode)
}

func (G *Gcm) UpdateCommandConfig(updateConfig func(cfg *osexec.CommandConfig)) *Gcm {
	updateConfig(G.cmdConfig)
	return G
}

func (G *Gcm) WithDebug() *Gcm {
	return G.WithDebugMode(true)
}

func (G *Gcm) WithDebugMode(debugMode bool) *Gcm {
	G.debugMode = debugMode
	G.cmdConfig.WithDebugMode(debugMode)
	return G
}

func (G *Gcm) ShowDebugMessage() *Gcm {
	switch {
	case G.errorOnce != nil && len(G.output) > 0:
		zaplog.ZAPS.P1.SUG.Errorln("wrong", eroticgo.RED.Sprint(G.errorOnce), "message:", "\n"+eroticgo.RED.Sprint(string(G.output))+"\n", "-")
	case G.errorOnce != nil:
		zaplog.ZAPS.P1.SUG.Errorln("wrong", eroticgo.RED.Sprint(G.errorOnce), "\n", "-")
	case len(G.output) > 0:
		zaplog.ZAPS.P1.SUG.Debugln("done", "message:", "\n"+eroticgo.GREEN.Sprint(string(G.output))+"\n", "-")
	default:
		zaplog.ZAPS.P1.SUG.Debugln("done", "\n", "-")
	}
	return G
}

func (G *Gcm) MustDone() *Gcm {
	if G.errorOnce != nil {
		if len(G.output) > 0 {
			zaplog.ZAPS.P1.SUG.Panicln("wrong", eroticgo.RED.Sprint(G.errorOnce), "message:", "\n"+eroticgo.RED.Sprint(string(G.output))+"\n", "-")
		} else {
			zaplog.ZAPS.P1.SUG.Panicln("wrong", eroticgo.RED.Sprint(G.errorOnce), "\n", "-")
		}
	}
	return G
}

func (G *Gcm) When(condition func(*Gcm) bool, run func(*Gcm) *Gcm) *Gcm {
	if condition(G) {
		return run(G)
	}
	return G
}

func (G *Gcm) WhenThen(condition func(*Gcm) (bool, error), run func(*Gcm) *Gcm) *Gcm {
	if success, err := condition(G); err != nil {
		return newWaGcm(G.cmdConfig, []byte{}, err, G.debugMode)
	} else if success {
		return run(G)
	}
	return G
}
