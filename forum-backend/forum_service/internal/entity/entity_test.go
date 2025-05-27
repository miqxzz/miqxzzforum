package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPost_Validate(t *testing.T) {
	tests := []struct {
		name    string
		post    Post
		wantErr bool
	}{
		{
			name: "valid post",
			post: Post{
				AuthorId: 1,
				Title:    "Test Title",
				Content:  "Test Content",
			},
			wantErr: false,
		},
		{
			name: "empty title",
			post: Post{
				AuthorId: 1,
				Title:    "",
				Content:  "Test Content",
			},
			wantErr: true,
		},
		{
			name: "empty content",
			post: Post{
				AuthorId: 1,
				Title:    "Test Title",
				Content:  "",
			},
			wantErr: true,
		},
		{
			name: "invalid author id",
			post: Post{
				AuthorId: 0,
				Title:    "Test Title",
				Content:  "Test Content",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.post.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestComment_Validate(t *testing.T) {
	tests := []struct {
		name    string
		comment Comment
		wantErr bool
	}{
		{
			name: "valid comment",
			comment: Comment{
				PostId:   1,
				AuthorId: 1,
				Content:  "Test Content",
			},
			wantErr: false,
		},
		{
			name: "empty content",
			comment: Comment{
				PostId:   1,
				AuthorId: 1,
				Content:  "",
			},
			wantErr: true,
		},
		{
			name: "invalid post id",
			comment: Comment{
				PostId:   0,
				AuthorId: 1,
				Content:  "Test Content",
			},
			wantErr: true,
		},
		{
			name: "invalid author id",
			comment: Comment{
				PostId:   1,
				AuthorId: 0,
				Content:  "Test Content",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestChatMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		message ChatMessage
		wantErr bool
	}{
		{
			name: "valid message",
			message: ChatMessage{
				UserID:    1,
				Username:  "testuser",
				Content:   "Test message",
				Timestamp: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty content",
			message: ChatMessage{
				UserID:    1,
				Username:  "testuser",
				Content:   "",
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty username",
			message: ChatMessage{
				UserID:    1,
				Username:  "",
				Content:   "Test message",
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid user id",
			message: ChatMessage{
				UserID:    0,
				Username:  "testuser",
				Content:   "Test message",
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
