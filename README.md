# GoDefenders

# Screenshots

![](/images/defenders_win.jpg)
![](/images/defenders_loss.jpg)
![](/images/DefendersDemo2.gif)

# How to Run

Clone the repo, then use `go run .` in the directory containing `main.go`.

```
git clone github.com/flycliff/godefenders
cd godefenders
go run .
```

Requires `github.com/jroimartin/gocui`, if you don't have it yet, run `go get github.com/jroimartin/gocui` before running.

# To-Do

- double check all thread safety. I got all big issues but still a bit rough around the edges.
- pretty up scoreboard, add debug info
- reduce this to 300 lines (personal challenge)
- Make a release with executables (personal challenge)