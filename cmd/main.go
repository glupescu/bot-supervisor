package main

import (
	"bot-supervisor/agent"
	"bot-supervisor/user"
)

func main() {
	var err error

	sysAgent := &agent.SysAgent{}
	sysAgentName := "Supervisor"
	if err = sysAgent.Init(sysAgentName,
		agent.NewBot("1234567890:AABBCCDDEEFFGGHHIIJJKKLLMMNNOOPPQQR")); err != nil {
		return
	}

	// TODO
	// leave empty to allow anyone who is not a bot
	allowedUsers := map[int64]user.Identity{
		/*
			1111111114: { // TODO user ID
				Type:      user.FullAccess,
				FirstName: "MyNameIs", // leave empty to only disallow via ID
			},
		*/
	}

	go agent.Run(sysAgent, allowedUsers)

	select {}
}
