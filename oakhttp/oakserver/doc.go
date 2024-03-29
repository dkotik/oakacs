/*

Package oakserver provides a standard library [http.Server] with conventional production defaults and a smooth configuration interface.

# Standard Usage

  err := oakserver.Run(context.Background())

# NGrok Usage

OakServer is easy to use with <https://ngrok.com> tunnel, which exposes your local server to the world. Use with caution. You should be fairly confident that your code is secure and will not leak data from your system or damage it.

  import (
    // ...
    "golang.ngrok.com/ngrok"
    "golang.ngrok.com/ngrok/config"
  )

  func main() {
    // ...
    tunnel, err := ngrok.Listen(ctx,
      config.HTTPEndpoint(),
      ngrok.WithAuthtokenFromEnv(),
    )
    if err != nil {
      panic(err)
    }

    fmt.Println("NGrok HTTP endpoint:", tunnel.URL())
    err := oakserver.Run(
      context.Background(),
      oakserver.WithListener(tunnel),
    )
    // ...
  }


*/
package oakserver
