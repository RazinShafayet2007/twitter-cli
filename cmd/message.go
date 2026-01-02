package cmd

import (
	"fmt"
	"strings"

	"github.com/RazinShafayet2007/twitter-cli/internal/config"
	"github.com/RazinShafayet2007/twitter-cli/internal/display"
	"github.com/RazinShafayet2007/twitter-cli/internal/store"
	"github.com/spf13/cobra"
)

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Direct messaging",
	Long:  `Send and receive direct messages`,
}

var messageSendCmd = &cobra.Command{
	Use:   "send [username] [text]",
	Short: "Send a direct message",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		receiverUsername := args[0]
		text := strings.Join(args[1:], " ")

		// Check if logged in
		senderUsername, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in. Run: twt login <username>")
		}

		userStore := store.NewUserStore(DB)

		// Get sender
		sender, err := userStore.GetByUsername(senderUsername)
		if err != nil {
			return fmt.Errorf("failed to get sender: %w", err)
		}

		// Get receiver
		receiver, err := userStore.GetByUsername(receiverUsername)
		if err != nil {
			return fmt.Errorf("user @%s not found", receiverUsername)
		}

		// Can't message yourself
		if sender.ID == receiver.ID {
			return fmt.Errorf("you cannot message yourself")
		}

		// Send message
		messageStore := store.NewMessageStore(DB)
		_, err = messageStore.Send(sender.ID, receiver.ID, text)
		if err != nil {
			return err
		}

		fmt.Printf("Message sent to @%s\n", receiverUsername)
		return nil
	},
}

var messageInboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "View your inbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")

		// Check if logged in
		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in. Run: twt login <username>")
		}

		userStore := store.NewUserStore(DB)
		user, err := userStore.GetByUsername(username)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		// Get inbox
		messageStore := store.NewMessageStore(DB)
		messages, err := messageStore.GetInbox(user.ID, limit)
		if err != nil {
			return err
		}

		if len(messages) == 0 {
			fmt.Println("Your inbox is empty.")
			return nil
		}

		// Display messages
		fmt.Println("Inbox:")
		fmt.Println()
		for _, m := range messages {
			timeAgo := display.FormatTimeAgo(m.Message.CreatedAt)
			readStatus := ""
			if !m.Message.Read {
				readStatus = " 🔴" // Unread indicator
			}

			fmt.Printf("From @%s (%s)%s\n", m.SenderName, timeAgo, readStatus)
			fmt.Printf("%s\n", m.Message.Text)
			fmt.Println()
		}

		return nil
	},
}

var messageConversationCmd = &cobra.Command{
	Use:   "conversation [username]",
	Short: "View conversation with a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		otherUsername := args[0]
		limit, _ := cmd.Flags().GetInt("limit")

		// Check if logged in
		currentUsername, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in. Run: twt login <username>")
		}

		userStore := store.NewUserStore(DB)

		// Get current user
		currentUser, err := userStore.GetByUsername(currentUsername)
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}

		// Get other user
		otherUser, err := userStore.GetByUsername(otherUsername)
		if err != nil {
			return fmt.Errorf("user @%s not found", otherUsername)
		}

		// Get conversation
		messageStore := store.NewMessageStore(DB)
		messages, err := messageStore.GetConversation(currentUser.ID, otherUser.ID, limit)
		if err != nil {
			return err
		}

		if len(messages) == 0 {
			fmt.Printf("No messages with @%s yet.\n", otherUsername)
			return nil
		}

		// Mark as read
		if err := messageStore.MarkAsRead(currentUser.ID, otherUser.ID); err != nil {
			return err
		}

		// Display conversation
		fmt.Printf("Conversation with @%s:\n", otherUsername)
		fmt.Println()

		for _, m := range messages {
			timeAgo := display.FormatTimeAgo(m.Message.CreatedAt)

			if m.Message.SenderID == currentUser.ID {
				// Message from you
				fmt.Printf("[You → @%s] (%s)\n", m.ReceiverName, timeAgo)
			} else {
				// Message to you
				fmt.Printf("[@%s → You] (%s)\n", m.SenderName, timeAgo)
			}

			fmt.Printf("%s\n", m.Message.Text)
			fmt.Println()
		}

		return nil
	},
}

var messageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all conversations",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if logged in
		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in. Run: twt login <username>")
		}

		userStore := store.NewUserStore(DB)
		user, err := userStore.GetByUsername(username)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		// Get conversations
		messageStore := store.NewMessageStore(DB)
		conversations, err := messageStore.GetConversations(user.ID)
		if err != nil {
			return err
		}

		if len(conversations) == 0 {
			fmt.Println("No conversations yet.")
			return nil
		}

		// Display conversations
		fmt.Println("Conversations:")
		fmt.Println()

		for _, conv := range conversations {
			timeAgo := display.FormatTimeAgo(conv.LastMessageAt)
			unreadBadge := ""
			if conv.UnreadCount > 0 {
				unreadBadge = fmt.Sprintf(" (%d unread) 🔴", conv.UnreadCount)
			}

			fmt.Printf("@%s%s\n", conv.OtherUsername, unreadBadge)
			fmt.Printf("  Last: %s (%s)\n", truncate(conv.LastMessage, 50), timeAgo)
			fmt.Println()
		}

		return nil
	},
}

var messageUnreadCmd = &cobra.Command{
	Use:   "unread",
	Short: "Show unread message count",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if logged in
		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in. Run: twt login <username>")
		}

		userStore := store.NewUserStore(DB)
		user, err := userStore.GetByUsername(username)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		// Get unread count
		messageStore := store.NewMessageStore(DB)
		count, err := messageStore.GetUnreadCount(user.ID)
		if err != nil {
			return err
		}

		switch count {
		case 0:
			fmt.Println("No unread messages.")
		case 1:
			fmt.Println("You have 1 unread message.")
		default:
			fmt.Printf("You have %d unread messages.\n", count)
		}

		return nil
	},
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

var messageDeleteCmd = &cobra.Command{
	Use:   "delete [message_id]",
	Short: "Delete a message you sent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		messageID := args[0]

		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in")
		}

		userStore := store.NewUserStore(DB)
		user, err := userStore.GetByUsername(username)
		if err != nil {
			return err
		}

		messageStore := store.NewMessageStore(DB)
		if err := messageStore.DeleteMessage(messageID, user.ID); err != nil {
			return err
		}

		fmt.Println("Message deleted")
		return nil
	},
}

var messageSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search messages by text",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in")
		}

		userStore := store.NewUserStore(DB)
		user, err := userStore.GetByUsername(username)
		if err != nil {
			return err
		}

		messageStore := store.NewMessageStore(DB)
		messages, err := messageStore.SearchMessages(user.ID, query)
		if err != nil {
			return err
		}

		if len(messages) == 0 {
			fmt.Printf("No messages found matching '%s'\n", query)
			return nil
		}

		fmt.Printf("Found %d message(s) matching '%s':\n\n", len(messages), query)

		for _, m := range messages {
			timeAgo := display.FormatTimeAgo(m.Message.CreatedAt)
			fmt.Printf("[@%s → @%s] (%s)\n", m.SenderName, m.ReceiverName, timeAgo)
			fmt.Printf("%s\n", m.Message.Text)
			fmt.Println()
		}

		return nil
	},
}

func init() {
	// Add flags
	messageInboxCmd.Flags().Int("limit", 20, "Number of messages to show")
	messageConversationCmd.Flags().Int("limit", 50, "Number of messages to show")

	// Add subcommands
	messageCmd.AddCommand(messageSendCmd)
	messageCmd.AddCommand(messageInboxCmd)
	messageCmd.AddCommand(messageConversationCmd)
	messageCmd.AddCommand(messageListCmd)
	messageCmd.AddCommand(messageUnreadCmd)
	messageCmd.AddCommand(messageDeleteCmd)
	messageCmd.AddCommand(messageSearchCmd)

	rootCmd.AddCommand(messageCmd)
}
