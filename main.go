package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "strconv"
    "time"
)

const contactedFileName = "contacted"
const done = true

var twilioAPIKey, phoneNumber, selectedFriend string
var numFriends, selectedIndex int
var friendsList []string
var contactedMap map[int]bool = make(map[int]bool)

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

func checkFileExistsOrCreate(fileName string) {
    if _, err := os.Stat(fileName); os.IsNotExist(err) {
        file, err := os.Create(fileName)
        checkError(err)
        file.Close()
    }
}

func readFile(fileName string, readerFunc func(scanner *bufio.Scanner), channel chan bool) {
    file, err := os.Open(fileName)
    checkError(err)
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    readerFunc(scanner)
    checkError(scanner.Err())
    channel <- done
}

func readConfigFile(scanner *bufio.Scanner) {    
    skipLines(scanner, 1) // Phone Number comment
    phoneNumber = getNextLineText(scanner)
    
    skipLines(scanner, 2) // New Line + Friend Count comment
    numFriendsString := getNextLineText(scanner)
    numFriendsInt64, err := strconv.ParseInt(numFriendsString, 0, 64)
    checkError(err)
    numFriends = int(numFriendsInt64) // Assuming reasonable number of friends
    
    skipLines(scanner, 2) // New Line + Friend List comment
    friendsList = getFriendsList(scanner, numFriends)
    
    fmt.Println(phoneNumber)
    fmt.Println(numFriends)
    fmt.Println(friendsList)
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

func readContactedFile(scanner *bufio.Scanner) {
    for scanner.Scan() {
        num, err := strconv.ParseInt(scanner.Text(), 0, 64)
        checkError(err)
        contactedMap[int(num)] = true
       // contactedList = append(contactedList, int(num))  // Assuming reasonable number of friends
    }
}

func processContactedFile(channel chan bool) {
    checkFileExistsOrCreate(contactedFileName)
    readFile(contactedFileName, readContactedFile, channel)
}

func readSecretFile(scanner *bufio.Scanner) {
    skipLines(scanner, 1) // API comment
    twilioAPIKey = getNextLineText(scanner)
    
    fmt.Println(twilioAPIKey)
}

//  This function selects a random friend that hasn't been contacted yet. 
//  It chooses a random index and makes sure that that index is not in the
//  contacted map. If it is, then it chooses another random index.
//  This isn't the most efficient algorithm, but for the amount of
//  friends we expect to be working with (< 200), this is more than fine.
func selectRandomFriend() {
    random := rand.New(rand.NewSource(time.Now().UnixNano())) // use time as random seed
    inMap := true
	for inMap {
	    selectedIndex = random.Intn(numFriends)
	    inMap = contactedMap[selectedIndex]
	}
 	selectedFriend = friendsList[selectedIndex]
}

func sendTwilioMessage() {
	
}

// Can have goroutine that runs this at same time as twilio message being sent
func cleanUpFile() {
	
}

func main() {
    if len(os.Args) >= 3 {
        // First arg (zero index) is path to file
        configFileName := os.Args[1]
        secretFileName := os.Args[2]
        fmt.Println("Starting the script!")
        
        channel := make(chan bool)
        go processContactedFile(channel)
        go readFile(configFileName, readConfigFile, channel)
        go readFile(secretFileName, readSecretFile, channel)
        for i := 0; i < 3; i++ { // wait till all 3 files have been read
            <- channel
        }
        fmt.Println("Done processing all files")
        
        selectRandomFriend()
        fmt.Printf("Today's friend is %s\n", selectedFriend)
        
    } else {
        fmt.Println("Usage: keepInTouch config_file secret_file")
    }
	
}

