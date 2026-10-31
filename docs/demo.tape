# Generate docs/demo.gif with VHS: https://github.com/charmbracelet/vhs
#   vhs docs/demo.tape
# Or regenerate the committed GIF with the small Go helper used in CI docs.

Output docs/demo.gif
Set FontSize 14
Set Width 900
Set Height 400
Set Theme "Catppuccin Mocha"

Type "go run ./cmd/sieve serve -config config.demo.yaml"
Enter
Sleep 1s

Type `curl -s -w "\nHTTP %{http_code}\n" localhost:8080/v1/examples/scraped_page/scan -X POST`
Enter
Sleep 2s
