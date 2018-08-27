package main

func main() {
	//bc := NewBlockchain()

	//defer代码块会在函数调用链表中增加一个函数调用。这个函数调用不是普通的函数调用，而是会在函数正常返回，
	// 也就是return之后添加一个函数调用。defer通常用来释放函数内部变量,确保数据库关闭。
	//defer bc.db.Close()

	//test git
	cli := CLI{}
	cli.run()
}
