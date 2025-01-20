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
			if in == nil {
				return
			}
			if done != nil {
				select {
				case <-done:
					return
				case val, ok := <-in:
					if !ok {
						return
					}
					out <- val
				}
			} else {
				select {
				case val, ok := <-in:
					if !ok {
						return
					}
					out <- val
				default:
				}
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
