# bot-supervisor

**bot-supervisor** allows you to run automation commands on a target machine and control it remotely through messaging apps like Telegram.

---

## Features

- Execute predefined system commands via Telegram messages.
- Remote control and monitoring via popular messaging platforms.
- Easy configuration and user access restrictions.

---

## Available System Commands

System commands are defined in [`agent/sys.go`](agent/sys.go) along with their message keys.

Example command configuration:

"temp": {
cmd: "sensors",
args: []string{},
sysFunc: sysExecutor,
},


Sending messages like `"temp"`, `"What temp"`, or `"Give me the temp"` to the bot triggers the execution of the `sensors` command on the target machine.

---

## Getting a Telegram Bot Token

Create a new Telegram bot and obtain its API token by following these steps:

1. **Open Telegram and find BotFather:**

   Search for `@BotFather` in the Telegram app (official Telegram bot management bot).

2. **Start a chat:**

   Click **Start** or send `/start` to interact with BotFather.

3. **Create a new bot:**

   Send `/newbot` and follow prompts to provide:
   - A display name for your bot.
   - A unique username ending with `bot` (e.g., `myexample_bot`).

4. **Receive your API token:**

   BotFather will generate a token (sample): `1234567890:AABBCCDDEEFFGGHHIIJJKKLLMMNNOOPPQQR`

Save this securely—it’s required for bot authentication.

---

## Running bot-supervisor

1. Insert your Telegram bot token in `main.go`.

2. Configure user access restrictions in `main.go` as needed.

3. Run: `make run`


This will build the bot, kill any running instance, and launch it in the background.

---

## Sample Conversation

User: yo

Bot: hey

User: what fpga

Bot: ??? write 'help' for available commands

User: help

Bot: Will execute help []...

Bot: Available commands:
temp
top
cpu
mem
speedtest
disk
gpu
help

User: what temp

Bot: Will execute sensors []...

Bot:
Adapter: Virtual device
temp1: +40.0°C
nvme-pci-0300
Adapter: PCI adapter
Composite: +55.9°C (low = -0.1°C, high = +99.8°C) (crit = +109.8°C)
...

---

## Project Structure

- Command definitions: [`agent/sys.go`](agent/sys.go)
- Main program: `main.go`
- Build and run: `Makefile`
