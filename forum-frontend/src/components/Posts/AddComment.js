import React, { useState } from 'react';
import axios from 'axios';
import styled from 'styled-components';

const AddCommentContainer = styled.div`
    margin-top: 20px;
    padding: 15px;
    border: 1px solid #e0e0e0;
    border-radius: 6px;
    background-color: #f9f9f9;
`;

const CommentForm = styled.form`
    display: flex;
    flex-direction: column;
`;

const CommentTextarea = styled.textarea`
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 4px;
    margin-bottom: 10px;
    font-size: 14px;
    resize: vertical;
`;

const CommentButton = styled.button`
    background-color: #4CAF50;
    color: white;
    padding: 10px 15px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
    transition: background-color 0.2s ease;

    &:hover {
        background-color: #3e8e41;
    }
`;

const AddComment = ({ postId, onCommentCreated }) => {
    const [content, setContent] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        const token = localStorage.getItem('token');

        if (!content.trim()) {
            alert('Please enter a comment.');
            return;
        }

        if (!token) {
            alert('You are not authenticated.');
            return;
        }

        try {
            await axios.post(`http://localhost:8081/posts/${postId}/comments`, { content }, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            setContent('');
            if (onCommentCreated) {
                onCommentCreated();
            }
        } catch (error) {
            console.error('Error creating comment:', error);
            alert('Failed to post comment.');
        }
    };

    return (
        <AddCommentContainer>
            <CommentForm onSubmit={handleSubmit}>
                <CommentTextarea
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    placeholder="Add a comment..."
                    rows="3" // Allow for multiline input
                />
                <CommentButton type="submit">Post Comment</CommentButton>
            </CommentForm>
        </AddCommentContainer>
    );
};

export default AddComment;