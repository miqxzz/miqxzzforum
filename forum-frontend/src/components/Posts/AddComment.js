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
    margin-top: 20px;
    padding: 15px;
    background: #f9f9f9;
    border-radius: 8px;
    border: 1px solid #e0e0e0;
`;

const TextArea = styled.textarea`
    width: 100%;
    min-height: 100px;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    margin-bottom: 10px;
    font-family: 'Montserrat', sans-serif;
    font-weight: 400;
    resize: vertical;

    &:focus {
        outline: none;
        border-color: #9b59b6;
        box-shadow: 0 0 0 2px rgba(155, 89, 182, 0.2);
    }
`;

const SubmitButton = styled.button`
    background-color: #9b59b6;
    color: white;
    padding: 10px 20px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
    font-weight: 500;
    transition: background-color 0.2s ease;

    &:hover {
        background-color: #8e44ad;
    }

    &:disabled {
        background-color: #ccc;
        cursor: not-allowed;
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
                <TextArea
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    placeholder="Add a comment..."
                    rows="3" // Allow for multiline input
                />
                <SubmitButton type="submit">Post Comment</SubmitButton>
            </CommentForm>
        </AddCommentContainer>
    );
};

export default AddComment;