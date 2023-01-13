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

    conf := config.GetConfig("/var/config/general.yaml")
    pprint.PrintHLine('=')
    path_to_pipe := conf.Path_to_not + conf.Path_to_pipe

    var availableMethods map[string]api.Method = map[string]api.Method{}

    pprint.PrintCentered("CERES")
    pprint.PrintHLine('=')

    fmt.Println("To close CERES, press",
        pprint.Color("[CRTL+C]", "red"))

    p := api.StartProcess("CERES",
                     availableMethods,
                     path_to_pipe,
                     path_to_pipe + "__in__")

    c := make(chan os.Signal, 32)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    for {
        select {
        case <-c:
            p.Stop()
        case <-p.ListenToClosure():
        }
    }

    pprint.PrintHLine('_')
    fmt.Println(pprint.Color("Goodbye !", "cyan"))
}
