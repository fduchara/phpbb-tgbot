install GO

https://golang.org/doc/install


install trash (vendoring)

https://github.com/rancher/trash ( get -u github.com/rancher/trash )

install dependences: trash --directory -C src


assembly program: make build


edit config/config.yaml

run: bin/bot -c config/config.yaml

run debug : bin/bot -c config/config.yaml -v true
