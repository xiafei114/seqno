# SeqNo Generator

一个基于 Gorm 的发号器

# Why?

顺序号发号器，根据步长顺序递增，使用 format 格式化输出

# Examples

简单使用

```go

	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?parseTime=true")
	if err != nil {
		fmt.Println(fmt.Sprintln(err))
	}

	defer db.Close()

	InitTable(db)

	num, err := NewSeqNoGenerator(db, "test-0").Next()
	if err != nil {
		fmt.Println(fmt.Sprintln(err))
	}
	fmt.Printf("%s\n", num)
```

使用协程或者 web 发号时需要加锁

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
				num, err := NewSeqNoGenerator(db, v).Step(2).Next()
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

Step 步长 . 默认 1

### StartWith

开始值. 默认为 0，第一个号发号后为 1

### SeqFormat

format 格式化输出. 默认 %05d
