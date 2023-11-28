package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)

	go func() {
		<-done
		close(out)
	}()

	go func() {
		stageOut := in
		for _, stage := range stages {
			stageOut = stage(stageOut)
		}
		for v := range stageOut {
			select {
			case <-done:
				return
			default:
				out <- v
			}
		}
		close(out)
	}()

	return out
}
