package k8s

import (
	"net/http"

	apiutils "github.com/previousnext/m8s/api/utils"
	"github.com/previousnext/m8s/k8sclient"
	"github.com/previousnext/m8s/server/k8s/utils"
	"golang.org/x/net/websocket"
)

// Exec bash (shell) inside a container.
func (s Server) Exec(w http.ResponseWriter, r *http.Request) {
	pod, err := apiutils.Param(r, "pod")
	if err != nil {
		apiutils.Fatal(w, err)
		return
	}

	container, err := apiutils.Param(r, "container")
	if err != nil {
		apiutils.Fatal(w, err)
		return
	}

	client, config, err := k8sclient.New(s.Master, s.Config)
	if err != nil {
		apiutils.Fatal(w, err)
		return
	}

	wws := websocket.Handler(func(ws *websocket.Conn) {
		input := utils.PodExecInput{
			Client:    client,
			Config:    config,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			Reader:    ws,
			Writer:    ws,
			Namespace: s.Namespace,
			Pod:       pod,
			Container: container,
			Command: []string{
				"/bin/bash",
			},
		}

		err := utils.PodExec(input)
		if err != nil {
			apiutils.Fatal(w, err)
			return
		}
	})

	wws.ServeHTTP(w, r)
}