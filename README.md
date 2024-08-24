# Progress

A simple, lightweight progress bar for the terminal written in Go, cause i just wanted a simple plug and play progress bar that doesn't come with a big overhead.

### Installation

```sh
go get -u github.com/coalaura/progress
```

### Usage

The progress bar comes with a few themes built-in, but you can also define your own.

```go
bar := progress.NewProgressBar("Downloading", 100) // Uses progress.Theme.Default

// Themes can be passed like this:
// bar := progress.NewProgressBarWithTheme("Downloading", 100, progress.ThemeGradientUnicode)

bar.Start() // Starts the draw goroutine

for i := 0; i < 100; i++ {
    bar.Increment() // Increments the progress bar by 1

    // Or increment by a specific amount
    // bar.IncrementBy(5)
}

bar.Stop() // Stops the draw goroutine and waits for it to finish, also prints the final progress bar

// Or abort the draw goroutine (still waits for it to finish, but doesn't print the final progress bar)
// bar.Abort()
```

### Themes

```go
// progress.ThemeDefault
// Label [===>    ]  69.0%

// progress.ThemeDots
// Label [:::o....]  69.0%

// progress.ThemeHash
// Label [####>    ]  69.0%

// progress.ThemeBlocksUnicode
// Label [████▌    ]  69.0%

// progress.ThemeGradientUnicode
// Label [▓▓▓▓█░░░░]  69.0%

// Or create your own theme:
theme := progress.NewProgressTheme(" ", "=", ">")
```