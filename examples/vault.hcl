listener "tcp" {
  tls_disable = 1
  address = "127.0.0.1:8200"
}

storage "inmem" {
}

disable_mlock = true
