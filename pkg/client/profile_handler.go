package client

import (
	"fmt"
	"nutritionapp/pkg/server"
)

func (c *Client) handleProfile(args []string) {
	if len(args) == 0 {
		resp, err := makeRequestTyped[server.ProfileResponseData](c, server.ReqGetProfile, nil)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		c.displayProfile(*resp)
		return
	} else if args[0] == "create" {
		c.createProfile()
	} else {
		fmt.Println("Unknown profile command. Use 'help' for usage.")
	}
}

func (c *Client) createProfile() {
	fmt.Println("\nCreating new profile:")
	data := server.CreateProfileData{}

	fmt.Print("First Name: ")
	data.FirstName = c.readString()

	fmt.Print("Last Name: ")
	data.LastName = c.readString()

	fmt.Print("Age: ")
	data.Age = c.readInt()

	fmt.Print("Weight (kg): ")
	data.Weight = c.readFloat()

	fmt.Print("Height (cm): ")
	data.Height = c.readFloat()

	fmt.Print("Gender (male/female): ")
	data.Gender = c.readString()

	fmt.Print("Goal (weight loss/muscle gain/maintenance): ")
	data.Goal = c.readString()

	_, err := makeRequest(c, server.ReqCreateProfile, data)
	if err != nil {
		fmt.Printf("Error creating profile: %s\n", err)
		return
	}

	fmt.Println("\nProfile created successfully!")
}

func (c *Client) displayProfile(profile server.ProfileResponseData) {
	fmt.Println("\n=== Profile ===")
	fmt.Printf("Name: %s %s\n", profile.FirstName, profile.LastName)
	fmt.Printf("Age: %d\n", profile.Age)
	fmt.Printf("Weight: %.1f kg\n", profile.Weight)
	fmt.Printf("Height: %.1f cm\n", profile.Height)
	fmt.Printf("Gender: %s\n", profile.Gender)
	fmt.Printf("Goal: %s\n", profile.Goal)
	fmt.Printf("BMI: %.1f\n", profile.BMI)
	fmt.Printf("Estimated Body Fat: %.1f%%\n", profile.BodyFatPerc)
}
