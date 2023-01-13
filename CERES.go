package main

import (
    api "github.com/Alex-S-H-P/NOT_API"
    "CERES/src/config"
)

func main() {

    conf := config.GetConfig("/var/config/general.yaml")
    path_to_pipe := conf.Path_to_not + conf.Path_to_pipe

    var availableMethods map[string]api.Method = make(map[string]api.Method)

    api.StartProcess("CERES",
                     availableMethods,
                     path_to_pipe,
                     path_to_pipe + "__in__")
}
