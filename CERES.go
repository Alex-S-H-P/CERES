package main

import (
    "os"
    "os/signal"

    "fmt"
    "syscall"

    "CERES/src/config"
    pprint "CERES/src/utils/printUtils"
    api "github.com/Alex-S-H-P/NOT_API"
)

func main() {

    conf := config.GetConfig("./var/config/general.yaml")
    pprint.PrintHLine('=')
    path_to_pipe := "var/" + conf.Path_to_pipe

    var availableMethods map[string]api.Method = map[string]api.Method{}

    pprint.PrintCentered("CERES")
    pprint.PrintHLine('=')

    fmt.Println("To close",
        pprint.Color("CERES", "cyan"), ", press",
        pprint.Color("[CRTL+C]", "bold red"))

    p := api.StartProcess("CERES",
                     availableMethods,
                     path_to_pipe + "CERES",
                     path_to_pipe + "__in__")

    c := make(chan os.Signal, 32)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    select {
    case <-c:
        fmt.Println("\tCauhgt a CRTL+C, closing")
        p.Stop()
        break
    case <-p.ListenToClosure():
        break
    }

    pprint.PrintHLine('_')
    fmt.Println(pprint.Color("Goodbye !", "cyan"))
}
