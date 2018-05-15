package main

import (
    "os"
    "fmt"
    "path/filepath"
    "strings"
    "crypto/sha256"
    "io"
)

// calculate sha256 checksum for each file and return byte slice
func sha256sum(filepath string) []byte {
    f, err := os.Open(filepath)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    h := sha256.New()
    io.Copy(h, f)
    return h.Sum(nil)
}

// find number of occurences for each hash in slice and return map of hash:count
func duplicateCounter(sl []string) map[string]int {
    n := make(map[string]int)

    for _, sl := range sl {
        if _, exists := n[sl]; !exists {
            n[sl] = 1
        } else {
            n[sl] += 1
        }
    }
    return n
}

func main() {
    // abort if no root dir provided as argument
    if len(os.Args) < 2 {
        panic("No root dir specified.\nUsage: go run dupsearch root_folder")
    }

    rootdir := os.Args[1]
    //fmt.Println("rootdir=", rootdir)
    if _, err := os.Stat(rootdir); err != nil {
        if os.IsNotExist(err) {
            panic("No such directory!")
        }
    }

    cksumsMap := make(map[string]string) //map of filepath:checksum for reference
    cksumsList := []string{} //slice of checksums to work on

    filepath.Walk(rootdir, func(path string, info os.FileInfo, err error) error {
        // ignore directories
        if info.IsDir() {
            //fmt.Println("---FOUND DIR!!", path)
            return nil
        }
        // ignore hidden files
        if strings.HasPrefix(info.Name(), ".") {
            //fmt.Println("---FOUND HIDDEN FILE!!", path)
            return nil
        }
        //fmt.Println(info.Name())
        //fmt.Println(path)

        // ignore symlinks
        if info.Mode() & os.ModeSymlink != 0 {
            //fmt.Println("---FOUND SYMLINK!!", path)
            return nil
        }

        //ignore empty files
        if info.Size() == 0 {
            //fmt.Println("---FOUND EMPTY FILE!!", path)
            return nil
        }

        bs := sha256sum(path) // call sha256 func to make hash for each file
        s := string(bs[:]) //convert byte slice to string
        //fmt.Println(bs)
        //fmt.Printf("%s --> %x\n", path, bs)

        cksumsMap[path] = s // store in map: path/filename + checksum
        cksumsList = append(cksumsList, s) // create a slice of checksums

        return nil
    }) //flepath.Walk

    //fmt.Println("-----------filepath & checksums map")
    //for k, v := range cksumsMap {
    //    fmt.Printf("%s --> %x\n", k, v)
    //}

    //fmt.Println("-----------slice of checksums")
    //for _, h := range cksumsList {
    //    fmt.Printf("%x\n", h)
    //}

    cksumsCount := duplicateCounter(cksumsList)

    var i int

    // iterate over checksum:counts map and print all files that have
    // matching checksum values in the original filepath:checksum map
    for k, v := range cksumsCount {
        if v > 1 { // ignore files that have only 1 occurence
            //fmt.Printf("%x --> %d\n", k, v)
            fmt.Printf("These %d files are identical:\n", v)
            for key, val := range cksumsMap {
                if val == k {
                    fmt.Printf("%s\n", key)
                    i++
                }
            }
            fmt.Println("----  ****  ----\n")
        }
    }
    fmt.Println("Total identical files found:", i)
} //main()

