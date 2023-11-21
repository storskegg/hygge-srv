package hygge_srv

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/storskegg/hygge-srv/internal/bridge"
	"github.com/storskegg/hygge-srv/internal/messages"
)

func Execute() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			os.Exit(100)
		}
	}()

	if err := cmdRoot.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
}

var cmdRoot = &cobra.Command{
	Use:   "hygge-srv",
	Short: "A brief description of your application",
	RunE:  execRoot,
}

func execRoot(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	prometheus.MustRegister(gaugeHumi)
	prometheus.MustRegister(gaugeTemp)
	prometheus.MustRegister(gaugeBatt)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	log.Println("Connecting to lora bridge")

	bi, err := bridge.New(ctx)
	if err != nil {
		log.Println("Failed to connect to lora bridge")
		return err
	}
	defer bi.Close()

	log.Println("Connected to lora bridge")

	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-chSig
		log.Println("Got signal, shutting down")
		cancel()
	}()

	go bi.StartScanning()

	log.Println("Listening")
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Println(err)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		case line := <-bi.OutChan():
			go processLine(line)
		}
	}
}

func processLine(line string) {
	brideLine, err := messages.ParseBridgeLine(line)
	if err != nil {
		log.Println(err)
		return
	}

	gaugeHumi.Set(brideLine.Message.Data.Humidity)
	gaugeTemp.Set(brideLine.Message.Data.Temperature)
	gaugeBatt.Set(brideLine.Message.Data.Battery)

	fmt.Println(brideLine.Message.Data.String())
}
