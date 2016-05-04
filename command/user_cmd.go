package command

import (
	"fmt"
	"os"
//	"strings"

	gcli "github.com/codegangsta/cli"
	"github.com/sciffer/opsgenie-go-sdk/user"
)

// CreateTeamUsersAction creates an user at OpsGenie.
func CreateTeamUsersAction(c *gcli.Context) {
        // Create users first
        cli, err := NewUserClient(c)
        if err != nil {
                fmt.Printf("ERROR: %s\n", err.Error())
                os.Exit(1)
        }
        if c.IsSet("M") {
                members := extractMembersFromCommand(c)
                userRequest := user.CreateUserRequest{}
                for _, member := range members {
                        if member.User  != "" {
				printVerboseMessage("user creation call:"+member.User)
                                userRequest.Username = member.User
                                userRequest.Fullname = member.User
                                if member.Role != "" {
                                        userRequest.Role = member.Role
                                } else {
                                        userRequest.Role = "user"
                                }
                        } else {
                                continue
                        }
                        resp, err := cli.Create(userRequest)
                        if err == nil {
                                printVerboseMessage("User creation status:"+resp.Status)
                        } else {
                                printVerboseMessage("User creation failed:"+err.Error())
                        }
                }
        }
}
/*
// CreateUserAction creates an user at OpsGenie.
func CreateUserAction(c *gcli.Context) {
	cli, err := NewUserClient(c)
	if err != nil {
		os.Exit(1)
	}
	req := user.CreateUserRequest{}

	if val, success := getVal("name", c); success {
		req.Name = val
	}
	if c.IsSet("M") {
		req.Members = extractMembersFromCommand(c)
	}

	printVerboseMessage("Create user request prepared from flags, sending request to OpsGenie..")
c
	resp, err := cli.Create(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("User created successfully:")
	printVerboseMessage(resp.Id)
}

func extractMembersFromCommand(c *gcli.Context) []user.Member {
	var members []user.Member
	extraProps := c.StringSlice("D")
	for i := 0; i < len(extraProps); i++ {
		var member user.Member
		prop := extraProps[i]
		if !isEmpty("D", prop, c) && strings.Contains(prop, "=") {
			for _, keypair := range strings.Split(prop[1:], ",") {
				value := strings.Split(keypair, "=")
				switch {
				case strings.ToLower(value[0]) == "user":
					member.User = value[1]
				case strings.ToLower(value[0]) == "role":
					member.Role = value[1]
				default:
					fmt.Printf("Member parameters should be one of user or role, I go: %s\n", value[0])
					gcli.ShowCommandHelp(c, c.Command.Name)
					os.Exit(1)
				}
			}
			members[i] = member
		} else {
			fmt.Printf("Member parameters should have the value of the form a=b, but got: %s\n", prop)
			gcli.ShowCommandHelp(c, c.Command.Name)
			os.Exit(1)
		}
	}
	return members
}

// GetUserAction retrieves specified user details from OpsGenie.
func GetUserAction(c *gcli.Context) {
	cli, err := NewUserClient(c)
	if err != nil {
		os.Exit(1)
	}
	req := users.GetUserRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}

	printVerboseMessage("Get user request prepared from flags, sending request to OpsGenie..")

	resp, err := cli.Get(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	outputFormat := strings.ToLower(c.String("output-format"))
	printVerboseMessage("Got User successfully, and will print as " + outputFormat)
	switch outputFormat {
	case "yaml":
		output, err := resultToYAML(resp)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("%s\n", output)
	default:
		isPretty := c.IsSet("pretty")
		output, err := resultToJSON(resp, isPretty)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("%s\n", output)
	}
}

// AttachFileAction attaches a file to an user at OpsGenie.
func AttachFileAction(c *gcli.Context) {
	cli, err := NewUserClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := users.AttachFileUserRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	if val, success := getVal("attachment", c); success {
		f, err := os.Open(val)
		defer f.Close()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
		req.Attachment = f
	}

	if val, success := getVal("indexFile", c); success {
		req.IndexFile = val
	}

	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}
	if val, success := getVal("note", c); success {
		req.Note = val
	}

	printVerboseMessage("Attach request prepared from flags, sending request to OpsGenie..")

	_, err = cli.AttachFile(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	printVerboseMessage("File attached to user successfully.")
}

// AcknowledgeAction acknowledges an user at OpsGenie.
func UpdateUserAction(c *gcli.Context) {
	cli, err := NewUserClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := users.AcknowledgeUserRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}
	if val, success := getVal("note", c); success {
		req.Note = val
	}

	printVerboseMessage("Acknowledge user request prepared from flags, sending request to OpsGenie..")

	_, err = cli.Acknowledge(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	printVerboseMessage("User acknowledged successfully.")
}

// RenotifyAction re-notifies recipients at OpsGenie.
func RenotifyAction(c *gcli.Context) {
	cli, err := NewUserClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := users.RenotifyUserRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	if val, success := getVal("recipients", c); success {
		req.Recipients = strings.Split(val, ",")
	}
	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}
	if val, success := getVal("note", c); success {
		req.Note = val
	}

	printVerboseMessage("Renotify request prepared from flags, sending request to OpsGenie..")

	_, err = cli.Renotify(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Renotified successfully.")
}

// TakeOwnershipAction takes the ownership of an user at OpsGenie.
func TakeOwnershipAction(c *gcli.Context) {
	cli, err := NewUserClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := users.TakeOwnershipUserRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}
	if val, success := getVal("note", c); success {
		req.Note = val
	}

	printVerboseMessage("Take ownership request prepared from flags, sending request to OpsGenie..")

	_, err = cli.TakeOwnership(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Ownership taken successfully.")
}

// AssignOwnerAction assigns the specified user as the owner of the user at OpsGenie.
func AssignOwnerAction(c *gcli.Context) {
	cli, err := NewUserClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := users.AssignOwnerUserRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	if val, success := getVal("owner", c); success {
		req.Owner = val
	}
	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}
	if val, success := getVal("note", c); success {
		req.Note = val
	}

	printVerboseMessage("Assign ownership request prepared from flags, sending request to OpsGenie..")

	_, err = cli.AssignOwner(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	printVerboseMessage("Ownership assigned successfully.")
}

// AddUserAction adds a user to an user at OpsGenie.
func AddUserAction(c *gcli.Context) {
	cli, err := NewUserClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := users.AddUserUserRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	if val, success := getVal("user", c); success {
		req.User = val
	}
	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}
	if val, success := getVal("note", c); success {
		req.Note = val
	}

	printVerboseMessage("Add user request prepared from flags, sending request to OpsGenie.."

	_, err = cli.AddUser(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("User added successfully.")
}

// Update user at OpsGenie.
func UpdateAction(c *gcli.Context) {
	cli, err := NewUserClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := user.UpdateUserRequest{}
	if val, success := getVal("id", c); success {
		req.Id = val
	}
	if val, success := getVal("name", c); success {
		req.Name = val
	}
	if val, success := getVal("timezone", c); success {
		req.Timezone = val
	}
	if val, success := getVal("enabled", c); success {
		req.Timezone = val
	}
	if c.IsSet("R") {
                req.Rotations = extractRotationsFromCommand(c)
        }

	printVerboseMessage("Update request prepared from flags, sending request to OpsGenie..")

	_, err = cli.AddRecipient(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Updated successfully.")
}

func extractRotationsFromCommand(c *gcli.Context) []user.Rotation {
	var members []user.Rotation
	extraProps := c.StringSlice("D")
	for i := 0; i < len(extraProps); i++ {
		var rotation user.Rotation
		prop := extraProps[i]
		if !isEmpty("D", prop, c) && strings.Contains(prop, "=") {
			for _, keypair := range strings.Split(prop[1:], ",") {
				value := strings.Split(keypair, "=")
				switch {
				case strings.ToLower(value[0]) == "start":
					rotation.StartDate = value[1]
				case strings.ToLower(value[0]) == "end":
					rotation.EndDate = value[1]
				case strings.ToLower(value[0]) == "len":
					rotation.RotationLength = value[1]
				case strings.ToLower(value[0]) == "type":
					rotation.RotationType = value[1]
				case strings.ToLower(value[0]) == "name":
					rotation.Name = value[1]
				case strings.ToLower(value[0]) == "part":
					rotation.Participants = strings.split(value[1], "&")
				case strings.ToLower(value[0]) == "res":
					rotation.NotifyCondition = getRestrictionsFromCommand(value[1])
				default:
					fmt.Printf("Rotation parameters are invalid, I got: %s\n", value[0])
					gcli.ShowCommandHelp(c, c.Command.Name)
					os.Exit(1)
				}
			}
			rotations[i] = rotation
		} else {
			fmt.Printf("Rotation parameters should have the value of the form a=b, but got: %s\n", prop)
			gcli.ShowCommandHelp(c, c.Command.Name)
			os.Exit(1)
		}
	}
	return rotations
}
func getRestrictionssFromCommand(list string) []user.Restriction {
	var restrictions []user.Restriction
	for num, rest := range strings.split(list, "&") {
		var restriction user.Restriction
		for _, rest := range
	}
	return restrictions
}

// AddTagsAction adds tags to an user at OpsGenie.
func AddTagsAction(c *gcli.Context) {
	cli, err := NewUserClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := users.AddTagsUserRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	if val, success := getVal("tags", c); success {
		req.Tags = strings.Split(val, ",")
	}
	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}
	if val, success := getVal("note", c); success {
		req.Note = val
	}

	printVerboseMessage("Add tag request prepared from flags, sending request to OpsGenie..")

	_, err = cli.AddTags(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Tags added successfully.")
}

// AddNoteAction adds a note to an user at OpsGenie.
func AddNoteAction(c *gcli.Context) {
	cli, err := NewUserClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := users.AddNoteUserRequest{}

	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}
	if val, success := getVal("note", c); success {
		req.Note = val
	}

	printVerboseMessage("Add note request prepared from flags, sending request to OpsGenie..")

	_, err = cli.AddNote(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Note added successfully.")
}

// ExecuteActionAction executes a custom action on an user at OpsGenie.
func ExecuteActionAction(c *gcli.Context) {
	cli, err := NewUserClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := users.ExecuteActionUserRequest{}

	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	action, success := getVal("action", c)
	if success {
		req.Action = action
	}
	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}
	if val, success := getVal("note", c); success {
		req.Note = val
	}

	printVerboseMessage("Execute action request prepared from flags, sending request to OpsGenie..")

	resp, err := cli.ExecuteAction(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Action" + action + "executed successfully.")
	fmt.Printf("result=%s\n", resp.Result)
}

// CloseUserAction closes an user at OpsGenie.
func CloseUserAction(c *gcli.Context) {
	cli, err := NewUserClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := users.CloseUserRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}
	if val, success := getVal("note", c); success {
		req.Note = val
	}
	if val, success := getVal("notify", c); success {
		req.Notify = strings.Split(val, ",")
	}

	printVerboseMessage("Close user request prepared from flags, sending request to OpsGenie..")

	_, err = cli.Close(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("User closed successfully.")
}

// DeleteUserAction deletes an user at OpsGenie.
func DeleteUserAction(c *gcli.Context) {
	cli, err := NewUserClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := users.DeleteUserRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}

	printVerboseMessage("Delete user request prepared from flags, sending request to OpsGenie..")

	_, err = cli.Delete(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	printVerboseMessage("User deleted successfully.")
}
*/
