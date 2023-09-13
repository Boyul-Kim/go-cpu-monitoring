package handlers

import (
	"bytes"
	"cpu-mon/database"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/ssh"
)

var (
	server1 string
)

type LinuxServer struct {
	Server1 string `json:"server1" bson:"server1"`
}

func RunSsh(c *fiber.Ctx) error {
	sshConfig := &ssh.ClientConfig{
		User:            "username",
		Auth:            []ssh.AuthMethod{ssh.Password("password")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	var wg sync.WaitGroup

	addressList := []string{"server:22"}
	for _, address := range addressList {
		wg.Add(1)
		go FetchCPU(sshConfig, address, &wg)
	}

	wg.Wait()
	fmt.Println(server1)

	return c.Status(200).JSON(fiber.Map{"server1": server1})
}

func FetchCPU(sshConfig *ssh.ClientConfig, address string, wg *sync.WaitGroup) {
	defer wg.Done()
	sshClient, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		fmt.Println("ssh dial err:", err)
	}

	defer sshClient.Close()

	sshSession, sessionErr := sshClient.NewSession()
	if sessionErr != nil {
		fmt.Println("session err", sessionErr)
	}

	defer sshSession.Close()

	sshSession.Stdin = bytes.NewBufferString("testing")
	sshSession.Stdout = &bytes.Buffer{}
	sshSession.Stderr = &bytes.Buffer{}

	cmd := `export TERM=xterm ; top -b -n1|grep -i "Cpu(s)"`
	if cmdErr := sshSession.Run(cmd); cmdErr != nil {
		fmt.Println("cmd error", cmd, cmdErr, sshSession.Stderr.(*bytes.Buffer).String())
	}

	result := sshSession.Stdout.(*bytes.Buffer).String()
	switch address {
	case "server1":
		server1 = address + " " + result
	}
}

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
}

func FetchAllUsers(c *fiber.Ctx) error {
	users := database.FetchCollection("users")

	filter := bson.M{}
	opts := options.Find().SetSkip(0).SetLimit(100)

	cursor, err := users.Find(c.Context(), filter, opts)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	result := make([]User, 0)
	err = cursor.All(c.Context(), &result)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"result": result})
}

func FetchUser(c *fiber.Ctx) error {
	requestBody := c.Request().Body()

	var requestedUser User
	json.Unmarshal([]byte(requestBody), &requestedUser)

	users := database.FetchCollection("users")

	filter := bson.M{"username": requestedUser.Username}

	var user User
	err := users.FindOne(c.Context(), filter).Decode(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"result": user})
}
