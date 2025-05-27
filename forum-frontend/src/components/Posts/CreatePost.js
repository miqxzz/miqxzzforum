import React, { useState } from 'react';
import axios from 'axios';
import styled from 'styled-components';

const CreatePostContainer = styled.div`
  margin-bottom: 20px;
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 5px;
  background-color: #f9f9f9;
`;

const Input = styled.input`
  padding: 8px;
  margin-bottom: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 100%;
  box-sizing: border-box;
`;

const TextArea = styled.textarea`
  padding: 8px;
  margin-bottom: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 100%;
  box-sizing: border-box;
  height: 100px;
`;

const Button = styled.button`
  padding: 10px 15px;
  background-color: rgb(24, 255, 16);
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;

  &:hover {
    background-color: rgb(24, 234, 16);
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
            console.error('Error creating post:', error);
            alert('Error creating post. Please try again.');
        }
    };

    return (
        <CreatePostContainer>
            <h3>Create New Post</h3>
            <Input
                type="text"
                name="title"
                placeholder="Title"
                value={newPost.title} // Исправлено: Добавлено значение
                onChange={handleInputChange}
            />
            <TextArea
                name="content"
                placeholder="Content"
                value={newPost.content}
                onChange={handleInputChange}
            />
            <Button onClick={createPost}>Create Post</Button>
        </CreatePostContainer>
    );
};

export default CreatePost;