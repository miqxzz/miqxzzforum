import React, { useState } from 'react';
import axios from 'axios';
import styled from 'styled-components';

const AddCommentContainer = styled.div`
    margin-top: 24px;
    padding: 18px 14px 14px 14px;
    border: 1.5px solid rgba(162, 89, 255, 0.3);
    border-radius: 12px;
    background-color: #f6f2fa;
`;

const CommentForm = styled.form`
    display: flex;
    flex-direction: column;
    gap: 10px;
`;

const CommentTextarea = styled.textarea`
    padding: 12px;
    border: 1.2px solid #a259ff;
    border-radius: 6px;
    margin-bottom: 0;
    font-size: 15px;
    resize: vertical;
    background: #fff;
    font-family: 'Montserrat', Arial, sans-serif;
    min-height: 60px;
    &::placeholder {
        font-family: 'Montserrat', Arial, sans-serif;
        color: #bfa6e6;
        opacity: 1;
    }
`;

const CommentButton = styled.button`
    background-color: #a259ff;
    color: white;
    padding: 12px 0;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-size: 16px;
    transition: background-color 0.2s ease;
    font-family: 'Montserrat', Arial, sans-serif;
    width: 100%;
    margin-top: 4px;
    &:hover {
        background-color: #6c2eb7;
    }
`;

const AddComment = ({ postId, onCommentCreated }) => {
    const [content, setContent] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        const token = localStorage.getItem('token');

        if (!content.trim()) {
            alert('Пожалуйста, введите комментарий.');
            return;
        }

        if (!token) {
            alert('Вы не авторизованы.');
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
            console.error('Ошибка при создании комментария:', error);
            alert('Не удалось отправить комментарий.');
        }
    };

    return (
        <AddCommentContainer>
            <CommentForm onSubmit={handleSubmit}>
                <CommentTextarea
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    placeholder="Добавить комментарий..."
                    rows="3" // Allow for multiline input
                />
                <CommentButton type="submit">Отправить комментарий</CommentButton>
            </CommentForm>
        </AddCommentContainer>
    );
};

export default AddComment;