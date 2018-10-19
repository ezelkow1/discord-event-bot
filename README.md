# discord-event-bot
A discord bot that can parse a steam groups event page, display the schedule, and send out reminders of upcoming events

Run it with event-bot -c conf.json with a properly edited configuration file:
{
      "Token": "Enter your bot token here",
	     "BroadcastChannel": "Channel ID number of your broadcast channel (enable dev options in discord, right click on channel and copy ID",
        "EventURL": "http://steamcommunity.com/groups/GROUP_NAME_GOES_HERE#events"
}
