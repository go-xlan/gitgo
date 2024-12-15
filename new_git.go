package gogitxexec

import (
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

type Gcm struct {
	Cmc *osexec.CMC
	Out []byte
	Erx error
	DBG bool
}

func New(path string) *Gcm {
	return newOK(osexec.NewCMC().WithPath(path), make([]byte, 0), false)
}

func newOK(cmc *osexec.CMC, data []byte, deb bool) *Gcm {
	if deb {
		if len(data) > 0 {
			zaplog.ZAPS.P3.SUG.Debugln("done", "message:", "\n"+eroticgo.GREEN.Sprint(string(data))+"\n", "-")
		} else {
			zaplog.ZAPS.P3.SUG.Debugln("done", "\n", "-")
		}
	}
	return &Gcm{
		Cmc: cmc,
		Out: data,
		Erx: nil,
		DBG: deb,
	}
}

func newWa(cmc *osexec.CMC, data []byte, erx error, deb bool) *Gcm {
	if deb {
		if len(data) > 0 {
			zaplog.ZAPS.P3.SUG.Errorln("wrong", eroticgo.RED.Sprint(erx), "message:", "\n"+eroticgo.RED.Sprint(string(data))+"\n", "-")
		} else {
			zaplog.ZAPS.P3.SUG.Errorln("wrong", eroticgo.RED.Sprint(erx))
		}
	}
	return &Gcm{
		Cmc: cmc,
		Out: data,
		Erx: erx,
		DBG: deb,
	}
}

// 这个函数不要使用导出的，因为日志里是跳过3层调用栈，假如使用导出的，跳出的栈的层数就不正确啦
func (G *Gcm) do(name string, args ...string) *Gcm {
	if G.Erx != nil {
		return G //当出错时就不要再往下执行，直接在这里拦住，这样整个链式的后续动作就都不执行
	}
	out, err := G.Cmc.Exec(name, args...)
	if err != nil {
		return newWa(G.Cmc, out, err, G.DBG)
	}
	return newOK(G.Cmc, out, G.DBG)
}

func (G *Gcm) UpdateCmc(update func(cmc *osexec.CMC)) *Gcm {
	update(G.Cmc)
	return G
}

func (G *Gcm) WithDebug() *Gcm {
	G.DBG = true
	return G
}

func (G *Gcm) ShowDebugMessage() *Gcm {
	switch {
	case G.Erx != nil && len(G.Out) > 0:
		zaplog.ZAPS.P1.SUG.Errorln("wrong", eroticgo.RED.Sprint(G.Erx), "message:", "\n"+eroticgo.RED.Sprint(string(G.Out))+"\n", "-")
	case G.Erx != nil:
		zaplog.ZAPS.P1.SUG.Errorln("wrong", eroticgo.RED.Sprint(G.Erx), "\n", "-")
	case len(G.Out) > 0:
		zaplog.ZAPS.P1.SUG.Debugln("done", "message:", "\n"+eroticgo.GREEN.Sprint(string(G.Out))+"\n", "-")
	default:
		zaplog.ZAPS.P1.SUG.Debugln("done", "\n", "-")
	}
	return G
}

func (G *Gcm) MustDone() *Gcm {
	if G.Erx != nil {
		if len(G.Out) > 0 {
			zaplog.ZAPS.P1.SUG.Panicln("wrong", eroticgo.RED.Sprint(G.Erx), "message:", "\n"+eroticgo.RED.Sprint(string(G.Out))+"\n", "-")
		} else {
			zaplog.ZAPS.P1.SUG.Panicln("wrong", eroticgo.RED.Sprint(G.Erx), "\n", "-")
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

func (G *Gcm) WhenExec(condition func(*Gcm) (bool, error), run func(*Gcm) *Gcm) *Gcm {
	if ok, err := condition(G); err != nil {
		return newWa(G.Cmc, nil, err, G.DBG)
	} else if ok {
		return run(G)
	}
	return G
}
