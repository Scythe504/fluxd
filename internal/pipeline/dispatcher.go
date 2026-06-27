package pipeline

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/scythe504/kronos/internal/nodes"
)

func (p *Pipeline) Start(ctx context.Context) {
	nodeCfg := nodes.GetNodeConfig(ctx)
	machineID := nodeCfg.MachineID

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		tasks, err := p.db.GetTasks(ctx, machineID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				time.Sleep(1 * time.Second)
				continue
			}

			log.Println("ERR(Dispatcher): ", err)
			time.Sleep(2 * time.Second)
			continue
		}
		for _, task := range tasks {
			go func() {
				adapted, err := AdaptTask(task)
				if err != nil {
					log.Println(err)
					return
				}

				if err := p.Enqueue(ctx, task.PayloadSlug, adapted); err != nil {
					log.Println("ERR(Enqueued): ", err)
				}
			}()
		}
	}

}
