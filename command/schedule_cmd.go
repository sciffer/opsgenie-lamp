package command


import (
	"fmt"
	"os"
//	"strings"

	gcli "github.com/codegangsta/cli"
	"github.com/sciffer/opsgenie-go-sdk/schedule"
)

// UpdateScheduleAction creates an schedule at OpsGenie.
func UpdateTeamScheduleAction(c *gcli.Context) {
        cli, err := NewScheduleClient(c)
        if err != nil {
                fmt.Printf("%s\n", err.Error())
                os.Exit(1)
        }
        req := schedule.GetScheduleRequest{}
	if val, success := getVal("name", c); success {
                req.Name = val + "_schedule"
        }
        resp, err := cli.Get(req)
        if err != nil {
                printVerboseMessage(err.Error())
                os.Exit(1)
        }
        ureq := schedule.UpdateScheduleRequest{}
        ureq.Id = resp.Id
        ureq.Enabled = true
        if val, success := getVal("timezone", c); success {
                ureq.Timezone = val
        }
        uresp, err := cli.Update(ureq)
        if err != nil {
                fmt.Printf("%s\n", err.Error())
                os.Exit(1)
        } else {
                printVerboseMessage("Schedule updated successfully:")
                printVerboseMessage(uresp.Status)
        }
}
/*
// CreateScheduleAction creates an schedule at OpsGenie.
func CreateScheduleAction(c *gcli.Context) {
	cli, err := NewScheduleClient(c)
	if err != nil {
		os.Exit(1)
	}
	req := schedule.CreateScheduleRequest{}

	if val, success := getVal("name", c); success {
		req.Name = val
	}
	if c.IsSet("M") {
		req.Members = extractMembersFromCommand(c)
	}

	printVerboseMessage("Create schedule request prepared from flags, sending request to OpsGenie..")
c
	resp, err := cli.Create(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Schedule created successfully:")
	printVerboseMessage(resp.Id)
}

func extractMembersFromCommand(c *gcli.Context) []schedule.Member {
	var members []schedule.Member
	extraProps := c.StringSlice("D")
	for i := 0; i < len(extraProps); i++ {
		var member schedule.Member
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

// GetScheduleAction retrieves specified schedule details from OpsGenie.
func GetScheduleAction(c *gcli.Context) {
	cli, err := NewScheduleClient(c)
	if err != nil {
		os.Exit(1)
	}
	req := schedules.GetScheduleRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}

	printVerboseMessage("Get schedule request prepared from flags, sending request to OpsGenie..")

	resp, err := cli.Get(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	outputFormat := strings.ToLower(c.String("output-format"))
	printVerboseMessage("Got Schedule successfully, and will print as " + outputFormat)
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

// AttachFileAction attaches a file to an schedule at OpsGenie.
func AttachFileAction(c *gcli.Context) {
	cli, err := NewScheduleClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := schedules.AttachFileScheduleRequest{}
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

	printVerboseMessage("File attached to schedule successfully.")
}

// AcknowledgeAction acknowledges an schedule at OpsGenie.
func UpdateScheduleAction(c *gcli.Context) {
	cli, err := NewScheduleClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := schedules.AcknowledgeScheduleRequest{}
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

	printVerboseMessage("Acknowledge schedule request prepared from flags, sending request to OpsGenie..")

	_, err = cli.Acknowledge(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	printVerboseMessage("Schedule acknowledged successfully.")
}

// RenotifyAction re-notifies recipients at OpsGenie.
func RenotifyAction(c *gcli.Context) {
	cli, err := NewScheduleClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := schedules.RenotifyScheduleRequest{}
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

// TakeOwnershipAction takes the ownership of an schedule at OpsGenie.
func TakeOwnershipAction(c *gcli.Context) {
	cli, err := NewScheduleClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := schedules.TakeOwnershipScheduleRequest{}
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

// AssignOwnerAction assigns the specified user as the owner of the schedule at OpsGenie.
func AssignOwnerAction(c *gcli.Context) {
	cli, err := NewScheduleClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := schedules.AssignOwnerScheduleRequest{}
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

// AddScheduleAction adds a schedule to an schedule at OpsGenie.
func AddScheduleAction(c *gcli.Context) {
	cli, err := NewScheduleClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := schedules.AddScheduleScheduleRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	if val, success := getVal("schedule", c); success {
		req.Schedule = val
	}
	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}
	if val, success := getVal("note", c); success {
		req.Note = val
	}

	printVerboseMessage("Add schedule request prepared from flags, sending request to OpsGenie.."

	_, err = cli.AddSchedule(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Schedule added successfully.")
}

// Update schedule at OpsGenie.
func UpdateAction(c *gcli.Context) {
	cli, err := NewScheduleClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := schedule.UpdateScheduleRequest{}
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

func extractRotationsFromCommand(c *gcli.Context) []schedule.Rotation {
	var members []schedule.Rotation
	extraProps := c.StringSlice("D")
	for i := 0; i < len(extraProps); i++ {
		var rotation schedule.Rotation
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
func getRestrictionssFromCommand(list string) []schedule.Restriction {
	var restrictions []schedule.Restriction
	for num, rest := range strings.split(list, "&") {
		var restriction schedule.Restriction
		for _, rest := range
	}
	return restrictions
}

// AddTagsAction adds tags to an schedule at OpsGenie.
func AddTagsAction(c *gcli.Context) {
	cli, err := NewScheduleClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := schedules.AddTagsScheduleRequest{}
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

// AddNoteAction adds a note to an schedule at OpsGenie.
func AddNoteAction(c *gcli.Context) {
	cli, err := NewScheduleClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := schedules.AddNoteScheduleRequest{}

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

// ExecuteActionAction executes a custom action on an schedule at OpsGenie.
func ExecuteActionAction(c *gcli.Context) {
	cli, err := NewScheduleClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := schedules.ExecuteActionScheduleRequest{}

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

// CloseScheduleAction closes an schedule at OpsGenie.
func CloseScheduleAction(c *gcli.Context) {
	cli, err := NewScheduleClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := schedules.CloseScheduleRequest{}
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

	printVerboseMessage("Close schedule request prepared from flags, sending request to OpsGenie..")

	_, err = cli.Close(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Schedule closed successfully.")
}

// DeleteScheduleAction deletes an schedule at OpsGenie.
func DeleteScheduleAction(c *gcli.Context) {
	cli, err := NewScheduleClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := schedules.DeleteScheduleRequest{}
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

	printVerboseMessage("Delete schedule request prepared from flags, sending request to OpsGenie..")

	_, err = cli.Delete(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	printVerboseMessage("Schedule deleted successfully.")
}
*/
