package hw06pipelineexecution

type (
	In  = <-chan any
	Out = In
	Bi  = chan any
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := processStage(in, done)
	for _, stage := range stages {
		if stage == nil {
			continue
		}
		out = processStage(stage(out), done)
	}
	return out
}

func processStage(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for range in { //nolint:revive
				// отбрасываем значения
			}
		}()
		for {
			select {
			case val, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- val:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()
	return out
}
