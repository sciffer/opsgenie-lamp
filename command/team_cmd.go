package command

import (
	"fmt"
	"os"
	"strings"

	gcli "github.com/codegangsta/cli"
	"github.com/sciffer/opsgenie-go-sdk/team"
)

// CreateTeamAction creates an team at OpsGenie.
func CreateTeamAction(c *gcli.Context) {
	// Create users first
	printVerboseMessage("Create team users first")
	CreateTeamUsersAction(c)
	// Create team for the users
	printVerboseMessage("Create team")
	cli, err := NewTeamClient(c)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	req := team.CreateTeamRequest{}
	if val, success := getVal("name", c); success {
		req.Name = val
	}
	if c.IsSet("M") {
		req.Members = extractMembersFromCommand(c)
	}
	resp, err := cli.Create(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Team created successfully:"+resp.Id)
	// Update team schedule
	printVerboseMessage("Enabling and updating the default team schedule")
	UpdateTeamScheduleAction(c)
	// Update team escalation
	if _, success := getVal("escalation", c); success {
		printVerboseMessage("Updating the default team escalation")
		UpdateTeamEscalationAction(c)
	}
}

func extractMembersFromCommand(c *gcli.Context) []team.Member {
	extraProps := c.StringSlice("M")
	members := make([]team.Member,len(extraProps))
	for i := 0; i < len(extraProps); i++ {
		var member team.Member
		prop := extraProps[i]
		if !isEmpty("M", prop, c) && strings.Contains(prop, "=") {
			for _, keypair := range strings.Split(prop, ",") {
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
// ListTeamAction retrieves specified team details from OpsGenie.
func ListTeamAction(c *gcli.Context) {
	cli, err := NewTeamClient(c)
	if err != nil {
		os.Exit(1)
	}
	req := team.ListTeamsRequest{}

	printVerboseMessage("List team request prepared from flags, sending request to OpsGenie..")

	resp, err := cli.List(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	outputFormat := strings.ToLower(c.String("output-format"))
	printVerboseMessage("Got Teams successfully, and will print as " + outputFormat)
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
/*
// AttachFileAction attaches a file to an team at OpsGenie.
func AttachFileAction(c *gcli.Context) {
	cli, err := NewTeamClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := teams.AttachFileTeamRequest{}
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

	printVerboseMessage("File attached to team successfully.")
}

// AcknowledgeAction acknowledges an team at OpsGenie.
func UpdateTeamAction(c *gcli.Context) {
	cli, err := NewTeamClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := teams.AcknowledgeTeamRequest{}
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

	printVerboseMessage("Acknowledge team request prepared from flags, sending request to OpsGenie..")

	_, err = cli.Acknowledge(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	printVerboseMessage("Team acknowledged successfully.")
}

// RenotifyAction re-notifies recipients at OpsGenie.
func RenotifyAction(c *gcli.Context) {
	cli, err := NewTeamClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := teams.RenotifyTeamRequest{}
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

// TakeOwnershipAction takes the ownership of an team at OpsGenie.
func TakeOwnershipAction(c *gcli.Context) {
	cli, err := NewTeamClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := teams.TakeOwnershipTeamRequest{}
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

// AssignOwnerAction assigns the specified user as the owner of the team at OpsGenie.
func AssignOwnerAction(c *gcli.Context) {
	cli, err := NewTeamClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := teams.AssignOwnerTeamRequest{}
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

// AddTeamAction adds a team to an team at OpsGenie.
func AddTeamAction(c *gcli.Context) {
	cli, err := NewTeamClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := teams.AddTeamTeamRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	if val, success := getVal("team", c); success {
		req.Team = val
	}
	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}
	if val, success := getVal("note", c); success {
		req.Note = val
	}

	printVerboseMessage("Add team request prepared from flags, sending request to OpsGenie..")

	_, err = cli.AddTeam(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Team added successfully.")
}

// AddRecipientAction adds recipient to an team at OpsGenie.
func AddRecipientAction(c *gcli.Context) {
	cli, err := NewTeamClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := teams.AddRecipientTeamRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	if val, success := getVal("recipient", c); success {
		req.Recipient = val
	}
	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}
	if val, success := getVal("note", c); success {
		req.Note = val
	}

	printVerboseMessage("Add recipient request prepared from flags, sending request to OpsGenie..")

	_, err = cli.AddRecipient(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Recipient added successfully.")
}

// AddTagsAction adds tags to an team at OpsGenie.
func AddTagsAction(c *gcli.Context) {
	cli, err := NewTeamClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := teams.AddTagsTeamRequest{}
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

// AddNoteAction adds a note to an team at OpsGenie.
func AddNoteAction(c *gcli.Context) {
	cli, err := NewTeamClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := teams.AddNoteTeamRequest{}

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

// ExecuteActionAction executes a custom action on an team at OpsGenie.
func ExecuteActionAction(c *gcli.Context) {
	cli, err := NewTeamClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := teams.ExecuteActionTeamRequest{}

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

// CloseTeamAction closes an team at OpsGenie.
func CloseTeamAction(c *gcli.Context) {
	cli, err := NewTeamClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := teams.CloseTeamRequest{}
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

	printVerboseMessage("Close team request prepared from flags, sending request to OpsGenie..")

	_, err = cli.Close(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Team closed successfully.")
}

// DeleteTeamAction deletes an team at OpsGenie.
func DeleteTeamAction(c *gcli.Context) {
	cli, err := NewTeamClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := teams.DeleteTeamRequest{}
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

	printVerboseMessage("Delete team request prepared from flags, sending request to OpsGenie..")

	_, err = cli.Delete(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	printVerboseMessage("Team deleted successfully.")
}
*/
