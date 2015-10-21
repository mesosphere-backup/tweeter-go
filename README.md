# Oinker-Go

Example Go (golang) web app with dependency injection and graceful shutdown. Acts like a mini Twitter clone.


## Dependencies

- [Facebook's Grace library](http://github.com/facebookgo/grace) - graceful shutdown
- [Inject](http://github.com/karlkfi/inject) - dependency injection
- [Humanize](http://github.com/dustin/go-humanize) - readable units

```
go get github.com/facebookgo/grace
go get github.com/karlkfi/inject
go get github.com/dustin/go-humanize
```


## Usage

1. Launch the server:

    ```
    go run main.go
    ```

    (ctrl-c to quit)

1. Home Page:

    ```
    $ curl http://localhost:8080/
    ```


## License

   Copyright 2015 Karl Isenberg

   Licensed under the [Apache License Version 2.0](LICENSE) (the "License");
   you may not use this project except in compliance with the License.

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
