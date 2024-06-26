package config

import (
    "os"
    "path/filepath"
)

var (
    CAFile = configfile("ca.pem")
    ServerCertFile = configfile("server.pem")
    ServerKeyFile = configfile("server-key.pem")
    RootClientCertFile = configfile("root-client.pem")
    RootClientKeyFile = configfile("root-client-key.pem")
    NobodyClientCertFile = configfile("nobody-client.pem")
    NobodyClientKeyFile = configfile("nobody-client-key.pem")
    ACLModelFile = configfile("model.conf")
    ACLPolicyFile = configfile("policy.csv")
)

func configfile(filename string) string {
    if dir := os.Getenv("CONFIG_DIR");dir != ""{
        return filepath.Join(dir,filename)
    }
    homeDir,err := os.UserHomeDir()
    if err != nil {
        panic(err)
    }
    return filepath.Join(homeDir,".proglog",filename)
}
