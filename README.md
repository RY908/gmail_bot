# gmail_bot

This is a line bot using golang and gmail api to notify users of unread messages.

# Usage

When a user send message "メールを確認", this bot search unread messages through gmail api and returns the messages' date, sender, and subject.
Also this bot regularly check if there are any unread messages everyday and send result if there are new messages.