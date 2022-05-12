package main

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	/* 这是我的第一个简单的程序 */
	fmt.Println("Hello, World!")

	var i int
	var f float64
	var b bool
	var s string
	var ss *int
	fmt.Printf("%v %v %v %q %v\n", i, f, b, s, ss)

	str := "this is str"
	fmt.Println(str)

	_变量()
	_指针()
	_结构体()
	_range()
	_map()
	_类型转换()
	_接口()
	_error()
	_线程()
	_缓冲区()
	_反射()
	_并发()
	_并发_多处理器()
	_竞争状态()
}

var (
	number int64
	wg2    sync.WaitGroup
)

func _竞争状态() {
	println("=========_竞争状态================")
	wg2.Add(4)
	go incCounter()
	go incCounter()
	go incCounter()
	go incCounter()
	fmt.Println("等待")
	wg2.Wait()
	fmt.Println("结束：", number)
}
func incCounter() {
	defer wg2.Done()
	for n := 0; n < 2; n++ {
		value := number
		// 使用 sync/atomic 这个包的 AddInt64() 来解决竞争问题
		atomic.AddInt64(&number, 1)

		// 当前 goroutine 从线程退出，并放回队列
		runtime.Gosched()

		value++

		number = value
	}
}

var wg sync.WaitGroup

func _并发() {
	println("=========_并发================")
	// 计数器加2，表示需要等待2个goroutine执行完毕
	wg.Add(2)
	// 设置当前使用逻辑处理器的数量
	runtime.GOMAXPROCS(1)
	// 查看当前使用的逻辑处理器数量
	//fmt.Println(runtime.GOMAXPROCS(0)) // 1
	go printPrime("A")
	go printPrime("B")

	fmt.Println("等待中")
	wg.Wait()
	fmt.Println("结束")

	// 没有起作用
	fmt.Printf("CUP:%s\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg.Add(2)
	go printLetter('a')
	go printLetter('A')
	fmt.Println("等待")
	wg.Wait()
	fmt.Println("\n结束")
}
func _并发_多处理器() {
	println("=========_并发_多处理器================")
	fmt.Printf("CUP:%s\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg.Add(2)
	go printLetter('a')
	go printLetter('A')
	fmt.Println("等待")
	wg.Wait()
	fmt.Println("\n结束")
}

// 打印字母
func printLetter(firstLetter rune) {
	defer wg.Done()
	for count := 0; count < 3; count++ {
		for char := firstLetter; char < firstLetter+26; char++ {
			fmt.Printf("%c ", char)
		}
	}
}

// 查找素数
func printPrime(prefix string) {
	defer wg.Done()
next:
	for n := 2; n < 5000; n++ {
		for m := 2; m < n; m++ {
			if n%m == 0 {
				continue next
			}
		}
		fmt.Printf("%s: %d\n", prefix, n)
	}
}

func _反射() {
	println("=========_反射================")
	var f_val float32 = 32
	var b_val bool = false
	var s_val string = "this is string"
	var map_val = make(map[string]string)
	map_val["aaa"] = "this is map value"
	fmt.Println(reflect.TypeOf(f_val))
	fmt.Println(reflect.TypeOf(b_val))
	fmt.Println(reflect.TypeOf(s_val))
	fmt.Println(reflect.TypeOf(map_val))
	println("========= valueof================")
	fmt.Println(reflect.ValueOf(f_val))
	fmt.Println(reflect.ValueOf(b_val))
	fmt.Println(reflect.ValueOf(s_val))
	fmt.Println(reflect.ValueOf(map_val))
	println("========= setValue================")
	setValue(&f_val)
	setValue(&b_val)
	setValue(&s_val)
	setValue(&map_val)
	println("========= after setValue================")
	fmt.Println(reflect.ValueOf(f_val))
	fmt.Println(reflect.ValueOf(b_val))
	fmt.Println(reflect.ValueOf(s_val))
	fmt.Println(reflect.ValueOf(map_val))

}

func setValue(i interface{}) {
	value := reflect.ValueOf(i)
	fmt.Printf("value[%s]:%s\n", reflect.TypeOf(i), value)
	switch value.Elem().Kind() {
	case reflect.Bool:
		fmt.Println("is bool")
		value.Elem().SetBool(true)
	case reflect.String:
		fmt.Println("is string")
		value.Elem().SetString("new String...")
	case reflect.Float32:
		fmt.Println("is float")
		value.Elem().SetFloat(22222)
	case reflect.Map:
		fmt.Println("is map")
	default:
		fmt.Println("no match...")
	}
}

func _缓冲区() {
	println("=========_缓冲区================")
	c := make(chan int, 2)
	c <- 1
	c <- 2
	fmt.Println(cap(c))
	fmt.Println(<-c)
	fmt.Println(<-c)
	println("=========_fibonacci 1================")
	c = make(chan int, 10)
	_fibonacci(cap(c), c)
	println("=========_fibonacci 2================")
	for i := range c {
		fmt.Println(i)
	}
	println("=========_fibonacci 3================")
}

func _fibonacci(len int, ch chan int) {
	x, y := 0, 1
	for i := 0; i < len; i++ {
		ch <- x
		x, y = y, x+y
	}
	close(ch)
}

func say(msg string) {
	fmt.Printf("say:%s\n", msg)
}
func _线程() {
	println("=========_线程================")
	say("a")
	go say("bbbbbbbbb")
	go say("cdddddd")
	say("d")
	time.Sleep(100 * 10000)
	say("e")
	println("=========channel1================")
	var s = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	time.Sleep(100 * 100)
	rs, re := <-c, <-c
	fmt.Printf("1:%s,2:%s\n", rs, re)
	println("=========channel2================")
}
func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // 把 sum 发送到通道 c
}

func _error() {
	println("=========_error================")
	// 正常情况
	if result, errorMsg := Divide(100, 10); errorMsg == "" {
		fmt.Println("100/10 = ", result)
	}
	// 当除数为零的时候会返回错误信息
	if _, errorMsg := Divide(100, 0); errorMsg != "" {
		fmt.Println("errorMsg is: ", errorMsg)
	}
}

// 定义一个 DivideError 结构
type DivideError struct {
	dividee int
	divider int
}

// 实现 `error` 接口
func (de *DivideError) Error() string {
	strFormat := `
    Cannot proceed, the divider is zero.
    dividee: %d
    divider: 0
`
	return fmt.Sprintf(strFormat, de.dividee)
}

// 定义 `int` 类型除法运算的函数
func Divide(varDividee int, varDivider int) (result int, errorMsg string) {
	if varDivider == 0 {
		dData := DivideError{
			dividee: varDividee,
			divider: varDivider,
		}
		errorMsg = dData.Error()
		return
	} else {
		return varDividee / varDivider, ""
	}

}

type Phone interface {
	call()
}
type Nokia struct {
}
type Iphone struct {
}

func (nokia Nokia) call() {
	fmt.Println("this is nokia")
}
func (iphone Iphone) call() {
	fmt.Println("this is iphone")
}

func _接口() {
	println("=========_接口================")
	var phone Phone
	phone = new(Nokia)
	phone.call()
	phone = new(Iphone)
	phone.call()
}

func _类型转换() {
	println("=========_类型转换================")
	var sum int = 17
	var count int = 5
	var mean float32

	mean = float32(sum) / float32(count)
	fmt.Printf("mean 的值为: %f\n", mean)
}

func _map() {
	println("=========_map================")
	var cMap map[int]string
	cMap = make(map[int]string)
	fmt.Println(cMap)
	cMap[1] = "aaaa"
	cMap[2] = "2222"
	fmt.Println(cMap)

	/* 创建map */
	countryCapitalMap := map[string]string{"France": "Paris", "Italy": "Rome", "Japan": "Tokyo", "India": "New delhi"}

	fmt.Println("原始地图")
	/* 打印地图 */
	for country := range countryCapitalMap {
		fmt.Println(country, "首都是", countryCapitalMap[country])
	}
	/*删除元素*/
	delete(countryCapitalMap, "France")
	fmt.Println("法国条目被删除")
	fmt.Println("删除元素后地图")
	/*打印地图*/
	for country := range countryCapitalMap {
		fmt.Println(country, "首都是", countryCapitalMap[country])
	}
}
func _range() {
	println("=========_Range================")
	var intAarry = []int{1, 2, 3, 4, 5}
	for key, value := range intAarry {
		fmt.Printf("key:%s value:%s \n", key, value)
	}

	println("=========_Range，读取Key================")
	// 读取 key
	for key := range intAarry {
		fmt.Printf("key is: %d\n", key)
	}
	println("=========_Range 读取 value================")
	// 读取 value
	for _, value := range intAarry {
		fmt.Printf("value is: %f\n", value)
	}

}

func _切片() {
	println("=========_切片================")
	var sli = make([]int, 3, 5)
	fmt.Printf("sli:%s, len:%d, cap:%d \n", sli, len(sli), cap(sli))
	/* 创建切片 */
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	printSlice(numbers)
	new1 := numbers[0:3]
	printSlice(new1)

	var numbers2 []int
	printSlice(numbers2)

	/* 允许追加空切片 */
	numbers2 = append(numbers2, 0)
	printSlice(numbers2)
	numbers2 = append(numbers2, 0, 1, 2, 3)
	printSlice(numbers2)

	/* 创建切片 numbers1 是之前切片的两倍容量*/
	numbers22 := make([]int, len(numbers2), (cap(numbers2))*2)

	/* 拷贝 numbers 的内容到 numbers1 */
	copy(numbers2, numbers22)
	printSlice(numbers22)

}

func printSlice(x []int) {
	fmt.Printf("len=%d cap=%d slice=%v\n", len(x), cap(x), x)
}

func _变量() {
	var a, b, c string = "1", "2", "3"
	print(a, b, c)

	d, e, f := "aaaa", "bbb", ""
	print(d, e, f, "\n")
}

func _指针() {
	println("=========_指针================")
	var a int = 10
	fmt.Println("指针地址：%s, %x", &a, &a)

	a = 20
	var ap *int
	ap = &a
	fmt.Printf("a的指针地址:%s\n", &a)
	fmt.Printf("ap的指针地址:%s\n", &ap)
	fmt.Printf("ap的值:%d\n", *ap)
}

type Books struct {
	title   string
	author  string
	subject string
	book_id int
}

func _结构体() {
	println("=========_结构体================")
	var Book1 Books /* 声明 Book1 为 Books 类型 */
	var Book2 Books /* 声明 Book2 为 Books 类型 */

	/* book 1 描述 */
	Book1.title = "Go 语言"
	Book1.author = "www.runoob.com"
	Book1.subject = "Go 语言教程"
	Book1.book_id = 6495407

	/* book 2 描述 */
	Book2.title = "Python 教程"
	Book2.author = "www.runoob.com"
	Book2.subject = "Python 语言教程"
	Book2.book_id = 6495700

	/* 打印 Book1 信息 */
	printBook(Book1)

	/* 打印 Book2 信息 */
	printBook(Book2)

	var struct_pointer *Books
	struct_pointer = &Book1
	println(struct_pointer.title)
}

func printBook(book Books) {
	fmt.Printf("Book title : %s\n", book.title)
	fmt.Printf("Book author : %s\n", book.author)
	fmt.Printf("Book subject : %s\n", book.subject)
	fmt.Printf("Book book_id : %d\n", book.book_id)
}
