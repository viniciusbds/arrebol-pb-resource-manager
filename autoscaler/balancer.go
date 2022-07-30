package autoscaler

import (
	"fmt"
)

type Balancer interface {
	Check(qs QueueState) (int, error)
}

type DefaultBalancer struct{}

func NewDefaultBalancer() Balancer {
	return &DefaultBalancer{}
}

// O balanceador padrão implementa uma estrategia simples, que é de
// equalizar o numero de tasks read to run ao numero de workers.
// portanto temos três cenários possíveis:
// tasks > workers:
// 		nesse caso é retornado um número positivo que indicará o numero
//		de workers que precisam ser criados para equalizar a workload ao
//		numero de recursos.
// tasks < workers:
//		nesse caso é retornado um número negativo que indicará o número de
//		de workers que precisam ser removidos para equalizar a workload ao
//		numero de recursos.
// tasks == workers:
//		nesse caso é retornado 0, indicando que a workload está equalziada
//		com os recursos disponiveis.
//
// OBS: outros mecanismos mais sofisticados poderão ser anexados ao ResourceManager
// por meio da interface Balancer
func (b *DefaultBalancer) Check(qs QueueState) (int, error) {
	fmt.Printf("Checking %v using Default balancer\n", qs)
	return qs.NumReadToRunTasks - qs.NumWorkers, nil
}
