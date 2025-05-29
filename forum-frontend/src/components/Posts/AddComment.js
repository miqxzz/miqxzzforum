import React, { useState } from 'react';
import axios from 'axios';
import styled from 'styled-components';

const AddCommentContainer = styled.div`
    margin-top: 20px;
    padding: 15px;
    border: 1px solid #e0e0e0;
    border-radius: 6px;
    background-color: #f8f4fc;
    font-family: 'Montserrat', sans-serif;
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

const AddCommentButton = styled.button`
    background: linear-gradient(90deg, #9b59b6 0%, #8e44ad 100%);
    color: #fff;
    border: none;
    border-radius: 6px;
    padding: 8px 16px;
    font-family: 'Montserrat', sans-serif;
    font-weight: 600;
    transition: background 0.3s, box-shadow 0.3s;
    box-shadow: 0 2px 8px #e1d5ee44;
    &:hover {
        background: linear-gradient(90deg, #8e44ad 0%, #9b59b6 100%);
        box-shadow: 0 4px 16px #c8a2e8aa;
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
            console.error('Ошибка создания комментария:', error);
            alert('Не удалось отправить комментарий.');
        }
    };

    return (
        <AddCommentContainer>
            <CommentForm onSubmit={handleSubmit}>
                <TextArea
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    placeholder="Добавить комментарий..."
                    rows="3" // Allow for multiline input
                />
                <SubmitButton type="submit">Отправить комментарий</SubmitButton>
            </CommentForm>
        </AddCommentContainer>
    );
};

export default AddComment;