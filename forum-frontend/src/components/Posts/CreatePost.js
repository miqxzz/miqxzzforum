import React, { useState } from 'react';
import axios from 'axios';
import styled from 'styled-components';

const CreatePostWrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 32px 0 24px 0;
`;

const CreatePostContainer = styled.div`
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 2px 8px #e1d5ee22;
  padding: 32px 40px 24px 40px;
  width: 100%;
  max-width: 900px;
  font-family: 'Montserrat', sans-serif;
`;

const CreateTitle = styled.h2`
  color: #a678f7;
  font-size: 1.7em;
  font-weight: 700;
  margin-bottom: 24px;
  text-align: center;
  font-family: 'Montserrat', sans-serif;
`;

const Input = styled.input`
  padding: 12px;
  margin-bottom: 16px;
  border: 2px solid #a678f7;
  border-radius: 8px;
  width: 100%;
  font-size: 1.1em;
  font-family: 'Montserrat', sans-serif;
  color: #3d2466;
  background: #fff;
  transition: border 0.2s;
  &:focus {
    border: 2px solid #8e44ad;
    outline: none;
  }
`;

const TextArea = styled.textarea`
  padding: 12px;
  margin-bottom: 16px;
  border: 2px solid #a678f7;
  border-radius: 8px;
  width: 100%;
  min-height: 100px;
  font-size: 1.1em;
  font-family: 'Montserrat', sans-serif;
  color: #3d2466;
  background: #fff;
  transition: border 0.2s;
  &:focus {
    border: 2px solid #8e44ad;
    outline: none;
  }
`;

const CreateButton = styled.button`
  background: #a678f7;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 14px 0;
  width: 100%;
  font-family: 'Montserrat', sans-serif;
  font-weight: 600;
  font-size: 1.1em;
  margin-top: 10px;
  transition: background 0.2s;
  box-shadow: 0 2px 8px #e1d5ee44;
  &:hover {
    background: #8e44ad;
  }
`;

const CreatePost = ({ onPostCreated }) => {
    const [newPost, setNewPost] = useState({ title: '', content: '' });
    const token = localStorage.getItem('token');

    const handleInputChange = (e) => {
        setNewPost({ ...newPost, [e.target.name]: e.target.value });
    };

    const createPost = async () => {
        try {
            const config = {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            };
            await axios.post('http://localhost:8081/posts', newPost, config);
            setNewPost({ title: '', content: '' });
            onPostCreated(); // Refresh the list
        } catch (error) {
            console.error('Ошибка создания поста:', error);
            alert('Ошибка создания поста. Пожалуйста, попробуйте снова.');
        }
    };

    return (
        <CreatePostWrapper>
            <CreatePostContainer>
                <CreateTitle>Создать новый пост</CreateTitle>
                <Input
                    type="text"
                    name="title"
                    placeholder="Заголовок"
                    value={newPost.title}
                    onChange={handleInputChange}
                />
                <TextArea
                    name="content"
                    placeholder="Содержание"
                    value={newPost.content}
                    onChange={handleInputChange}
                />
                <CreateButton onClick={createPost}>Создать пост</CreateButton>
            </CreatePostContainer>
        </CreatePostWrapper>
    );
};

export default CreatePost;