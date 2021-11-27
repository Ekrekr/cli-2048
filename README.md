# cli-2048

Play 2048 in your terminal! Featuring high scores, game saving and resumption.

<!-- TODO: Add gif -->

## Play ASAP

```
$ git clone https://github.com/Ekrekr/cli-2048.git ~/Documents
...
$ cd ~/cli-2048
$ go run cli-2048.go
```

## Install to Path (recommended)

```
$ go build cli-2048.go
```

Then either install the binary
```
$ sudo mv cli-2048 /usr/bin/2048
```

or set a launch alias (use .bashrc if you don't use zsh)

```
$ echo "alias 2048=~/Documents/cli-2048" >> ~/.zshrc
$ source ~/.zshrc
```
