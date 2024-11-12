package gogitxexec

import (
	"github.com/yyle88/osexec"
	"github.com/yyle88/zaplog"
)

type GitCmx struct {
	Cmx *osexec.CMX
	Out []byte
	Erx error
	DBG bool
}

func New(path string) *GitCmx {
	return newOK(osexec.NewCMX().WithPath(path), make([]byte, 0), false)
}

func newOK(cmx *osexec.CMX, data []byte, deb bool) *GitCmx {
	if deb {
		if len(data) == 0 {
			zaplog.ZAPS.P3.SUG.Debugln("done", "\n", "-")
		} else {
			zaplog.ZAPS.P3.SUG.Debugln("done", "data:", "\n", string(data), "\n", "-")
		}
	}
	return &GitCmx{
		Cmx: cmx,
		Out: data,
		Erx: nil,
		DBG: deb,
	}
}

func newWa(cmx *osexec.CMX, data []byte, erx error, deb bool) *GitCmx {
	if deb {
		zaplog.ZAPS.P3.SUG.Errorln("wrong", erx, "data:", "\n", string(data), "\n", "-")
	}
	return &GitCmx{
		Cmx: cmx,
		Out: data,
		Erx: erx,
		DBG: deb,
	}
}

// 这个函数不要使用导出的，因为日志里是跳过3层调用栈，假如使用导出的，跳出的栈的层数就不正确啦
func (G *GitCmx) do(name string, args ...string) *GitCmx {
	if G.Erx != nil {
		return G //当出错时就不要再往下执行，直接在这里拦住，这样整个链式的后续动作就都不执行
	}
	out, err := G.Cmx.Exec(name, args...)
	if err != nil {
		return newWa(G.Cmx, out, err, G.DBG)
	}
	return newOK(G.Cmx, out, G.DBG)
}

func (G *GitCmx) UpdateCmx(update func(cmx *osexec.CMX)) *GitCmx {
	update(G.Cmx)
	return G
}

func (G *GitCmx) WithDebug() *GitCmx {
	G.DBG = true
	return G
}

func (G *GitCmx) ShowDebug() *GitCmx {
	switch {
	case G.Erx != nil && len(G.Out) > 0:
		zaplog.ZAPS.P1.SUG.Errorln("wrong", G.Erx, "data:", "\n", string(G.Out), "\n", "-")
	case G.Erx != nil:
		zaplog.ZAPS.P1.SUG.Errorln("wrong", G.Erx, "\n", "-")
	case len(G.Out) > 0:
		zaplog.ZAPS.P1.SUG.Debugln("done", "data:", "\n", string(G.Out), "\n", "-")
	default:
		zaplog.ZAPS.P1.SUG.Debugln("done", "\n", "-")
	}
	return G
}

func (G *GitCmx) MustDone() *GitCmx {
	if G.Erx != nil {
		zaplog.ZAPS.P1.SUG.Panicln("wrong", G.Erx, "data:", "\n", string(G.Out), "\n", "-")
	}
	return G
}

func (G *GitCmx) When(condition func(*GitCmx) bool, run func(*GitCmx) *GitCmx) *GitCmx {
	if condition(G) {
		return run(G)
	}
	return G
}
