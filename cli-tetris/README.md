# cli-2048

<p align="center">
A better 2048 to play in the terminal! <br/> Featuring high scores, game saving and resumption.
</p>

<p align="center">
  <img src="https://github.com/ekrekr/cli-2048/blob/main/gif.gif?raw=true" alt="demo gif" height="300">
</p>

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

If you get file write errors, then you will need to give the built binary file read and write permissions

```
$ chmod a+x <path to binary>
```
