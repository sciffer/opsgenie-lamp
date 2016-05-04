package command

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"encoding/json"

	gcli "github.com/codegangsta/cli"
	"github.com/sciffer/opsgenie-go-sdk/escalation"
)
// Update escalation at OpsGenie.
func UpdateTeamEscalationAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}
	// Get escalation ID first
	req := escalation.GetEscalationRequest{}
	if val, success := getVal("name", c); success {
		req.Name = val + "_escalation"
	}
	resp, err := cli.Get(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
                os.Exit(1)
        }
	var ureq escalation.UpdateEscalationRequest
	if val,success := getVal("escalation", c); success {
		err := json.Unmarshal([]byte(val, &ureq))
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
	}
	ureq.Id = resp.Id
	uresp, err := cli.Update(ureq)
	if err != nil {
                fmt.Printf("%s\n", err.Error())
                os.Exit(1)
        } else {
		printVerboseMessage("Updated escalation:"+uresp.Status)
	}
}
/*
// CreateEscalationAction creates an escalation at OpsGenie.
func CreateEscalationAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}
	req := escalation.CreateEscalationRequest{}

	if val, success := getVal("name", c); success {
		req.Name = val
	}
	if c.IsSet("M") {
		req.Members = extractMembersFromCommand(c)
	}

	printVerboseMessage("Create escalation request prepared from flags, sending request to OpsGenie..")
c
	resp, err := cli.Create(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Escalation created successfully:")
	printVerboseMessage(resp.Id)
}

func extractMembersFromCommand(c *gcli.Context) []escalation.Member {
	var members []escalation.Member
	extraProps := c.StringSlice("D")
	for i := 0; i < len(extraProps); i++ {
		var member escalation.Member
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

// GetEscalationAction retrieves specified escalation details from OpsGenie.
func GetEscalationAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}
	req := escalations.GetEscalationRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}

	printVerboseMessage("Get escalation request prepared from flags, sending request to OpsGenie..")

	resp, err := cli.Get(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	outputFormat := strings.ToLower(c.String("output-format"))
	printVerboseMessage("Got Escalation successfully, and will print as " + outputFormat)
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

// AttachFileAction attaches a file to an escalation at OpsGenie.
func AttachFileAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := escalations.AttachFileEscalationRequest{}
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

	printVerboseMessage("File attached to escalation successfully.")
}

// AcknowledgeAction acknowledges an escalation at OpsGenie.
func UpdateEscalationAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := escalations.AcknowledgeEscalationRequest{}
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

	printVerboseMessage("Acknowledge escalation request prepared from flags, sending request to OpsGenie..")

	_, err = cli.Acknowledge(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	printVerboseMessage("Escalation acknowledged successfully.")
}

// RenotifyAction re-notifies recipients at OpsGenie.
func RenotifyAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := escalations.RenotifyEscalationRequest{}
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

// TakeOwnershipAction takes the ownership of an escalation at OpsGenie.
func TakeOwnershipAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := escalations.TakeOwnershipEscalationRequest{}
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

// AssignOwnerAction assigns the specified user as the owner of the escalation at OpsGenie.
func AssignOwnerAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := escalations.AssignOwnerEscalationRequest{}
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

// AddEscalationAction adds a escalation to an escalation at OpsGenie.
func AddEscalationAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := escalations.AddEscalationEscalationRequest{}
	if val, success := getVal("id", c); success {
		req.ID = val
	}
	if val, success := getVal("alias", c); success {
		req.Alias = val
	}
	if val, success := getVal("escalation", c); success {
		req.Escalation = val
	}
	req.User = grabUsername(c)
	if val, success := getVal("source", c); success {
		req.Source = val
	}
	if val, success := getVal("note", c); success {
		req.Note = val
	}

	printVerboseMessage("Add escalation request prepared from flags, sending request to OpsGenie.."

	_, err = cli.AddEscalation(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Escalation added successfully.")
}
*/
// Update escalation at OpsGenie.
func UpdateEscalationAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := escalation.UpdateEscalationRequest{}
	if val, success := getVal("id", c); success {
		req.Id = val
	}
	if val, success := getVal("name", c); success {
		req.Name = val
	}
	if c.IsSet("R") {
                req.Rules = extractRulesFromCommand(c)
        }

	printVerboseMessage("Update request prepared from flags, sending request to OpsGenie..")

	_, err = cli.Update(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Updated successfully.")
}

func extractRulesFromCommand(c *gcli.Context) []escalation.Rule {
	var rules []escalation.Rule
	extraProps := c.StringSlice("D")
	for i := 0; i < len(extraProps); i++ {
		var rule escalation.Rule
		prop := extraProps[i]
		if !isEmpty("D", prop, c) && strings.Contains(prop, "=") {
			for _, keypair := range strings.Split(prop[1:], ",") {
				value := strings.Split(keypair, "=")
				switch {
				case strings.ToLower(value[0]) == "delay":
					if delay,err := strconv.Atoi(value[1]); err != nil {
						fmt.Printf("String %s is not a number, cannot fill delay field! exiting abnormally!",value[1])
						os.Exit(1)
					} else {
						rule.Delay = delay
					}
				case strings.ToLower(value[0]) == "notify":
					rule.Notify = value[1]
				case strings.ToLower(value[0]) == "type":
					rule.NotifyType = value[1]
				case strings.ToLower(value[0]) == "cond":
					rule.NotifyCondition = value[1]
				default:
					fmt.Printf("Rule parameters should be one of delay,notify,type or cond, I go: %s\n", value[0])
					gcli.ShowCommandHelp(c, c.Command.Name)
					os.Exit(1)
				}
			}
			rules[i] = rule
		} else {
			fmt.Printf("Rule parameters should have the value of the form a=b, but got: %s\n", prop)
			gcli.ShowCommandHelp(c, c.Command.Name)
			os.Exit(1)
		}
	}
	return rules
}
/*
// AddTagsAction adds tags to an escalation at OpsGenie.
func AddTagsAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := escalations.AddTagsEscalationRequest{}
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

// AddNoteAction adds a note to an escalation at OpsGenie.
func AddNoteAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := escalations.AddNoteEscalationRequest{}

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

// ExecuteActionAction executes a custom action on an escalation at OpsGenie.
func ExecuteActionAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := escalations.ExecuteActionEscalationRequest{}

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

// CloseEscalationAction closes an escalation at OpsGenie.
func CloseEscalationAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := escalations.CloseEscalationRequest{}
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

	printVerboseMessage("Close escalation request prepared from flags, sending request to OpsGenie..")

	_, err = cli.Close(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	printVerboseMessage("Escalation closed successfully.")
}

// DeleteEscalationAction deletes an escalation at OpsGenie.
func DeleteEscalationAction(c *gcli.Context) {
	cli, err := NewEscalationClient(c)
	if err != nil {
		os.Exit(1)
	}

	req := escalations.DeleteEscalationRequest{}
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

	printVerboseMessage("Delete escalation request prepared from flags, sending request to OpsGenie..")

	_, err = cli.Delete(req)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	printVerboseMessage("Escalation deleted successfully.")
}
*/
