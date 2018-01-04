package cli

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"

	"github.com/syd7/actionizer/database"
	"github.com/syd7/actionizer/models"
)

func ShowActions(ms *database.Database) {
	actions := ms.AllActions()
	fmt.Printf("A total of %d saved actions : \n", len(actions))
	for i := range actions {
		fmt.Println("- " + actions[i].Description)
	}
}
func ShowUsers(ms *database.Database) {
	users := ms.AllUsers()
	fmt.Printf("A total of %d saved users : \n", len(users))
	for i := range users {
		fmt.Println("- " + users[i].Fullname)
	}
}
func AddUser(ms *database.Database, Fullname string, ImageURL string) {
	if Fullname == "" {
		fmt.Println("User not added , please enter the username using -username==\"Name\"")
		return
	}
	if ImageURL == "" {
		fmt.Println("User not added , please enter the user image using -userimage==\"Image Url\"")
		return
	}
	var ID bson.ObjectId
	var user = models.User{ID, Fullname, ImageURL}
	err := ms.InsertUser(user)
	if err == nil {
		fmt.Printf("User %s has been added ! \n", Fullname)
	} else {
		fmt.Printf("An error occured when adding the user %s \n", Fullname)

	}
}

func DeleteUser(ms *database.Database, Fullname string) {
	if Fullname == "" {
		fmt.Println("User not added , please enter the username using -username==\"Name\"")
		return
	}

	err := ms.DeleteUser(Fullname)

	if err == nil {
		fmt.Printf("User %s has been deleted ! \n", Fullname)
	} else {
		fmt.Printf("An error occured when deleting the user %s \n", Fullname)
	}
	return
}

func DeleteAction(ms *database.Database, Description string) {
	if Description == "" {
		fmt.Println("User not added , please enter the username using -username==\"Name\"")
		return
	}

	err := ms.DeleteAction(Description)

	if err == nil {
		fmt.Printf("action \"%s\" has been deleted ! \n", Description)
	} else {
		fmt.Printf("An error occured when deleting the action \" %s \" \n", Description)
	}
	return
}

func AddAction(ms *database.Database, Description string) {
	if Description == "" {
		fmt.Println("Action not added , please enter the action description using -actiondesc==\"action\"")
		return
	}
	err := ms.InsertAction(models.Action{"", Description})
	if err == nil {
		fmt.Printf("Action \" %s \" has been added ! \n", Description)
	} else {
		fmt.Printf("An error occured when adding the user \" %s \" \n", Description)
	}
}
