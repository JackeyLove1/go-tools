package test

func FailOnErr(err error) {
    if err != nil {
        panic(err)
    }
}
