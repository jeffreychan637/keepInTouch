package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
)

func checkError(err error) {
    if (err != nil) {
        panic(err)
    }
}

func skipLines(scanner *bufio.Scanner, linesToSkip int) {
    for i := 0; i < linesToSkip; i++ {
        scanner.Scan()
    }
}

func getNextLineText(scanner *bufio.Scanner) string {
    scanner.Scan()
    return scanner.Text()
}

func readConfigFile(configFileName string) {
// return statement
//(string, string, int, map[int]string, []string) {
    configFile, err := os.Open(configFileName)
    checkError(err)
    defer configFile.Close()
    
    scanner := bufio.NewScanner(configFile)
    skipLines(scanner, 1) // API comment
    twilioAPIKey := getNextLineText(scanner)
    
    skipLines(scanner, 2) // New Line + Phone Number comment
    phoneNumber := getNextLineText(scanner)
    
    skipLines(scanner, 2) // New Line + Friend Count comment
    friendsCount := getNextLineText(scanner)
    friendsCountInt, err := strconv.ParseInt(friendsCount, 0, 64)
    checkError(err)
    
    skipLines(scanner, 2) // New Line + Friend List comment
    friendsList := getFriendsList(scanner, int(friendsCountInt)) // Assuming reasonable number of friends
    
    fmt.Println(twilioAPIKey)
    fmt.Println(phoneNumber)
    fmt.Println(friendsCount)
    fmt.Println(friendsList)

    checkError(scanner.Err())
    
    
    
    
	fmt.Println(configFileName)
// 	twilioApiKey := "API Key"
//	phoneNumber := "123-456-7890"
//	friendsList := [
/* 		1 : "John Smith", 
		2 : "Harry Potter",
	]
	alreadySelected := [1]
	return twilioApiKey, phoneNumber, friendsList, alreadySelected */
}

func getFriendsList(scanner *bufio.Scanner, friendsCountInt int) []string {
    var friendsList []string
    var friend string
    for i := 0; i < friendsCountInt; i++ {
        scanner.Scan()
        friend = scanner.Text()
        friendsList = append(friendsList, friend)
    }
    return friendsList
}

func selectRandomFriend() {
	
}

func sendTwilioMessage() {
	
}

// Can have goroutine that runs this at same time as twilio message being sent
func cleanUpFile() {
	
}

func main() {
    if len(os.Args) >= 2 {
        configFileName := os.Args[1] // First arg (zero index) is path to file
        fmt.Println("Starting the script!")
        readConfigFile(configFileName)
    } else {
        fmt.Println("Please include the config file as the first command line argument to the script")
    }
	
}

