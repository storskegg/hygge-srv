package hygge_srv

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/storskegg/hygge-srv/internal/messages"

	"github.com/spf13/cobra"
	"github.com/storskegg/hygge-srv/internal/bridge"
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

	fmt.Println(brideLine.Message.Data.String())
}
