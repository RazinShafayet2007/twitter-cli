package errors

import "fmt"

// NotLoggedInError returns a helpful error for when user isn't logged in
func NotLoggedInError() error {
	return fmt.Errorf("not logged in. Run: twt login <username>")
}

// UserNotFoundError returns a formatted user not found error
func UserNotFoundError(username string) error {
	return fmt.Errorf("user @%s not found. Create with: twt user create %s", username, username)
}

// PostNotFoundError returns a post not found error
func PostNotFoundError(postID string) error {
	return fmt.Errorf("post %s not found", postID)
}

// AlreadyFollowingError returns already following error with suggestion
func AlreadyFollowingError(username string) error {
	return fmt.Errorf("already following @%s", username)
}

// NotFollowingError returns not following error
func NotFollowingError(username string) error {
	return fmt.Errorf("not following @%s. Nothing to unfollow", username)
}

// AlreadyLikedError returns already liked error
func AlreadyLikedError(postID string) error {
	return fmt.Errorf("already liked post %s", postID)
}

// NotLikedError returns not liked error
func NotLikedError(postID string) error {
	return fmt.Errorf("not liked post %s. Nothing to unlike", postID)
}

// CannotFollowSelfError returns cannot follow yourself error
