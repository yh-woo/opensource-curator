package handler

import (
	"net/http"

	"github.com/hibiken/asynq"

	"github.com/woo/opensource-curator/internal/worker"
)

// TriggerCollect enqueues a collection task via the asynq client.
func TriggerCollect(client *asynq.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "POST required")
			return
		}

		task := asynq.NewTask(worker.TypeCollectAll, nil)
		info, err := client.Enqueue(task)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to enqueue collection: "+err.Error())
			return
		}

		writeJSON(w, http.StatusAccepted, Envelope{
			Data: map[string]string{
				"taskId": info.ID,
				"queue":  info.Queue,
				"status": "enqueued",
			},
			NextActions: []Action{
				{Rel: "health", Href: "/v1/health", Description: "Check API health"},
			},
		})
	}
}
