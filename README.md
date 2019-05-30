# SeqNo Generator

a sequence number generator based on Gorm.

# Why?

Need a human call number on some cases(Restaurant etc...).

# Examples

Base examples

```go

    db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?parseTime=true")
	if err != nil {
		fmt.Println(fmt.Sprintln(err))
	}

	defer db.Close()

	InitTable(db)

	num, err := NewSeqNoGenerator(db, "Order").SeqFormat("Order%05d").Step(5).Next()
	if err != nil {
		fmt.Println(fmt.Sprintln(err))
	}
	fmt.Printf("%s\n", num)
```

Locking is required for goroutine

```go

    db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?parseTime=true")
	if err != nil {
		fmt.Println(fmt.Sprintln(err))
	}

	defer db.Close()

	InitTable(db)

	arrs := [5]string{"test-1", "test-2", "test-3", "test-4", "test-5"}

	var mutex sync.Mutex

	for _, v := range arrs {
		fmt.Printf("start %s\n", v)
		go func(v string) {
			for n := 0; n < 100; n++ {
				mutex.Lock()
				num, err := NewSeqNoGenerator(db, v).Next()
				if err != nil {
					fmt.Println(fmt.Sprintln(err))
				}
				fmt.Printf("%s --- %s\n", v, num)
				mutex.Unlock()
			}
			time.Sleep(2 * time.Second)
			fmt.Printf("end %s\n", v)
		}(v)

	}

	time.Sleep(2 * time.Second)
	fmt.Println("end")
```

## Options

### Logic ID

`Logic ID` represents a concrete business logic(like `order`,`user`).

### Step

Incremental step. default is 1

### StartWith

Starting number. default is 0

### SeqFormat

go fmt format. default is %05d
