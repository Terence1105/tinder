package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type HttpServer struct {
	cmd  *cobra.Command
	srv  *http.Server
	opts Options
}

func New(options ...Option) *HttpServer {
	h := &HttpServer{
		srv:  &http.Server{},
		opts: DefaultOptions(),
	}
	for _, opt := range options {
		opt(&h.opts)
	}
	h.newCommand()

	return h
}

func (h *HttpServer) Execute(ctx context.Context) error {
	return h.cmd.ExecuteContext(ctx)
}

func (h *HttpServer) newCommand() {
	cobra.OnInitialize(func() {
		viper.AutomaticEnv()
	})

	cmd := &cobra.Command{
		Run: h.run,
	}

	cmd.Flags().IntVarP(&h.opts.port, HttpServerPortFlag, "p", HttpServerPortVal, "http listen port")
	viper.BindPFlags(cmd.Flags())

	h.cmd = cmd
}

func (h *HttpServer) run(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	go func() {
		h.srv.Addr = fmt.Sprintf(":%d", h.opts.port)
		h.srv.Handler = h.opts.handler
		if h.opts.tls {
			// todo tls
			// h.server.ListenAndServeTLS()
			log.Fatalln("tls not implemented")
		} else {
			err := h.srv.ListenAndServe()
			if err != nil {
				// todo log
				log.Fatalln(err.Error())
			}
		}
	}()

	<-ctx.Done()
	shutdown, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()
	err := h.srv.Shutdown(shutdown)
	if err != nil {
		// todo log
		log.Println(err.Error())
	}

}
