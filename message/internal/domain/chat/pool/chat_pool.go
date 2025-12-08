package pool

import "message/internal/domain/chat/job"

type ChatPool interface {
	AddTask(job *job.ChatLineJob)
}
