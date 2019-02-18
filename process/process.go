package process

import (
	"context"
	"io"
	"log"
	"os/exec"
	"sync"
)

type Process struct {
	cmd *exec.Cmd
	out <-chan []byte
}

func Start(ctx context.Context, name string, arg ...string) (*Process, error) {
	cmd := exec.CommandContext(ctx, name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	p := &Process{cmd: cmd, out: merge(read(stdout), read(stderr))}
	return p, cmd.Start()
}

func (p *Process) Output() <-chan []byte { return p.out }
func (p *Process) Wait() error           { return p.cmd.Wait() }

// copied from io.Copy
func read(r io.Reader) <-chan []byte {
	out := make(chan []byte)
	go func() {
		buf := make([]byte, 32*1024)
		for {
			nr, err := r.Read(buf)
			if nr > 0 {
				out <- buf[:nr]
			}
			if err == io.EOF {
				close(out)
				return
			} else if err != nil {
				// TODO: store err in process?
				log.Print(err)
			}
		}
	}()
	return out
}

// copied from https://blog.golang.org/pipelines
func merge(cs ...<-chan []byte) <-chan []byte {
	var wg sync.WaitGroup
	out := make(chan []byte)

	output := func(c <-chan []byte) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
