import React, { useState } from 'react';
import axios from 'axios';
import styled from 'styled-components';

const CreatePostContainer = styled.div`
  margin-bottom: 20px;
  padding: 15px;
  border: 1px solid #a259ff;
  border-radius: 5px;
  background-color: #f9f9ff;
`;

const Input = styled.input`
  padding: 8px;
  margin-bottom: 10px;
  border: 1px solid #a259ff;
  border-radius: 4px;
  width: 100%;
  box-sizing: border-box;
  font-family: 'Montserrat', Arial, sans-serif;
`;

const TextArea = styled.textarea`
  padding: 8px;
  margin-bottom: 10px;
  border: 1px solid #a259ff;
  border-radius: 4px;
  width: 100%;
  box-sizing: border-box;
  height: 100px;
  font-family: 'Montserrat', Arial, sans-serif;
`;

const Button = styled.button`
  padding: 10px 15px;
  background-color: #a259ff;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
  font-family: 'Montserrat', Arial, sans-serif;

  &:hover {
    background-color: #6c2eb7;
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
            console.error('Ошибка при создании поста:', error);
            alert('Ошибка при создании поста. Пожалуйста, попробуйте снова.');
        }
    };

    return (
        <CreatePostContainer>
            <h3 style={{ color: '#a259ff' }}>Создать новый пост</h3>
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
            <Button onClick={createPost}>Создать пост</Button>
        </CreatePostContainer>
    );
};

export default CreatePost;