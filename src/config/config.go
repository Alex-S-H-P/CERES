package config

import (
    "gopkg.in/yaml.v3"
    "os"
)

type CERES_Config struct {
    Version string

    Path_to_not  string
    Path_to_pipe string
}

func GetConfig(path string) *CERES_Config {
    file_content, err := os.ReadFile(path)
    if err != nil {panic(err)}



    config := new(CERES_Config)
    err = yaml.Unmarshal(file_content, config)
    if err != nil {
        panic(err)
    }
    return config
}
