package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func terminate(in In, done In, stage Stage) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			if in == nil || done == nil {
				return
			}
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				out <- val
			}
		}
	}()
	return stage(out)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = terminate(out, done, stage)
	}
	return out
}
