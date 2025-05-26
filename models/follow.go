package models

import "time"

type Follow struct {
	ID          int       `db:"id"`
	FollowerID  int       `db:"follower_id"`
	FollowingID int       `db:"following_id"`
	Accepted    bool      `db:"accepted"`
	CreatedAt   time.Time `db:"created_at"`
}
