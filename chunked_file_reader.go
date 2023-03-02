package main

import (
    "bufio"
    "encoding/hex"
    "fmt"
    "log"
    "os"
    "io"
)

func main() {

    f, err := os.Open("sid.jpg")

    if err != nil {
        log.Fatal(err)
    }

    defer f.Close()

    reader := bufio.NewReader(f)
    buf := make([]byte, 256)

    for {
        _, err := reader.Read(buf)

        if err != nil {

            if err != io.EOF {
                fmt.Println(err)
            }

            break
        }

        fmt.Printf("%s", hex.Dump(buf))
    }
}