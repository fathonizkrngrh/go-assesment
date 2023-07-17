package config

import "os"

var PublicKey string = os.Getenv("PUBLICKEY")
var PrivateKey string = "secret-private-key-for-paseto-ok" // must be 32 char exactly

//var PrivateKey string = os.Getenv("PRIVATEKEY")
