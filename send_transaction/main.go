package main

import (
	"fmt"
	"github.com/truechain/truechain-engineering-code/rpc"
	"os"
	"strconv"
	"sync"
)

var Count int64

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("invalid args : %s count [\"ip:port\"]\n", os.Args[0])
		return
	}

	count, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("count err")
		return
	}

	ip := "127.0.0.1:8899"
	if len(os.Args) == 3 {
		ip = os.Args[2]
	}

	send(count, ip)

}

func send(count int, ip string) {
	client, err := rpc.Dial("http://" + ip)
	if err != nil {
		fmt.Println("Dail:", ip, err.Error())
		return
	}
	var account []string
	err = client.Call(&account, "etrue_accounts")
	if err != nil {
		fmt.Println("etrue_accounts Error", err.Error())
		return
	}
	if len(account) == 0 {
		fmt.Println("no account")
		return
	}
	fmt.Println("account:", account)

	//解锁账户
	var result string = ""

	err = client.Call(&result, "etrue_getBalance", account[0], "latest")

	if err != nil {
		fmt.Println("etrue_getBalance Error:", err)
		return
	} else {
		fmt.Println("etrue_getBalance Ok:", result)
	}

	var reBool bool

	err = client.Call(&reBool, "personal_unlockAccount", account[0], "admin", 90)
	if err != nil {
		fmt.Println("personal_unlockAccount Error:", err.Error())
		return
	} else {
		fmt.Println("personal_unlockAccount Ok", reBool)
	}

	waitGroup := &sync.WaitGroup{}
	//发送交易
	for a := 0; a < count; a++ {
		waitGroup.Add(1)
		go sendTransaction(client, account, waitGroup)
	}

	fmt.Println("Complete", count)
	waitGroup.Wait()
}

func sendTransaction(client *rpc.Client, account []string, wait *sync.WaitGroup) {
	defer wait.Done()
	map_data := make(map[string]string)
	map_data["from"] = account[0]
	map_data["to"] = account[1]
	map_data["value"] = "0x2100"
	var result string
	client.Call(&result, "etrue_sendTransaction", map_data)
	if result != "" {
		Count += 1
	}
}
