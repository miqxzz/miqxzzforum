import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import styled from 'styled-components';
import CommentList from './CommentList';
import AddComment from './AddComment';

// === Styled Components ===

const Container = styled.div`
    max-width: 800px;
    margin: 20px auto;
    padding: 20px;
    background-color: #fff;
    border-radius: 8px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
`;

const PostItemContainer = styled.div`
    margin-bottom: 20px;
    padding: 15px;
    border: 1px solid #e0e0e0;
    border-radius: 6px;
    background-color: #f9f9f9;
`;

const PostTitle = styled.h3`
    color: #333;
    text-align: center;
    margin-bottom: 5px;
    font-size: 32px;
`;

const PostContent = styled.p`
    color: #555;
    font-size: 16px;
    margin-bottom: 10px;
    text-align: left;
`;

const PostAuthor = styled.small`
    color: #777;
    font-style: italic;
`;

const DeleteButton = styled.button`
    background-color: #f44336;
    color: white;
    padding: 8px 12px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s ease;
    margin-right: 10px;

    &:hover {
        background-color: #d32f2f;
    }
    margin-top: 10px;
`;

const EditButton = styled.button`
    background-color: #2196F3;
    color: white;
    padding: 8px 12px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s ease;

    &:hover {
        background-color: #0b7dda;
    }
    margin-top: 10px;
`;

const ButtonContainer = styled.div`
    display: flex;
`;

const LoadingMessage = styled.p`
    color: #555;
    font-style: italic;
    text-align: center;
`;

const ErrorMessage = styled.p`
    color: #d32f2f;
    text-align: center;
`;

const PaginationContainer = styled.div`
    display: flex;
    justify-content: center;
    margin: 20px 0;
    gap: 10px;
`;

const PaginationButton = styled.button`
    padding: 5px 10px;
    cursor: pointer;
    background-color: ${props => props.active ? '#2196F3' : '#f5f5f5'};
    border: 1px solid #ddd;
    border-radius: 4px;
    
    &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
`;

const PageSizeSelect = styled.select`
    padding: 5px;
    border-radius: 4px;
    border: 1px solid #ddd;
`;

// === PostList Component ===

const PostList = () => {
    const [posts, setPosts] = useState([]);
    const [isAdmin, setIsAdmin] = useState(false);
    const [userId, setUserId] = useState(null);
    const [forceUpdate, setForceUpdate] = useState(false);
    const [editingPostId, setEditingPostId] = useState(null);
    const [editTitle, setEditTitle] = useState('');
    const [editContent, setEditContent] = useState('');
    const [pagination, setPagination] = useState({
        page: 1,
        limit: 5,
        total: 0
    });

    const fetchPosts = useCallback(async () => {
        try {
            const response = await axios.get('http://localhost:8081/posts', {
                params: {
                    page: pagination.page,
                    limit: pagination.limit
                }
            });
            setPosts(response.data.posts || []);
            setPagination(prev => ({
                ...prev,
                total: response.data.total || 0
            }));
        } catch (error) {
            console.error('Error fetching posts:', error);
        }
    }, [pagination.page, pagination.limit]);

    useEffect(() => {
        const storedRole = localStorage.getItem('role');
        const storedUserId = localStorage.getItem('userID');
        setIsAdmin(storedRole === 'admin');
        setUserId(storedUserId ? parseInt(storedUserId) : null);
        fetchPosts();
        const intervalId = setInterval(fetchPosts, 5000);
        return () => clearInterval(intervalId);
    }, [fetchPosts]);

    const handleDeletePost = async (postId) => {
        const token = localStorage.getItem('token');

        if (!token) {
            alert('You are not authenticated.');
            return;
        }

        try {
            await axios.delete(`http://localhost:8081/posts/${postId}`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            fetchPosts(); // Обновляем список с учетом пагинации
        } catch (error) {
            console.error('Error deleting post:', error);
            if (error.response && error.response.status === 403) {
                alert('You are not authorized to delete this post.');
            } else {
                alert('Failed to delete post.');
            }
        }
    };

    const handleEditPost = (post) => {
        setEditingPostId(post.id);
        setEditTitle(post.title);
        setEditContent(post.content);
    };

    const handleCancelEdit = () => {
        setEditingPostId(null);
        setEditTitle('');
        setEditContent('');
    };

    const handleSaveEdit = async (postId) => {
        const token = localStorage.getItem('token');

        if (!token) {
            alert('You are not authenticated.');
            return;
        }

        try {
            await axios.put(`http://localhost:8081/posts/${postId}`, {
                title: editTitle,
                content: editContent
            }, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            setEditingPostId(null);
            fetchPosts();
        } catch (error) {
            console.error('Error updating post:', error);
            if (error.response && error.response.status === 403) {
                alert('You are not authorized to edit this post.');
            } else {
                alert('Failed to update post.');
            }
        }
    };

    const handleCommentCreated = () => {
        setForceUpdate(prev => !prev);
    };

    const handlePageChange = (newPage) => {
        setPagination(prev => ({
            ...prev,
            page: newPage
        }));
    };

    const handleLimitChange = (e) => {
        setPagination({
            page: 1, // Сбрасываем на первую страницу при изменении количества
            limit: parseInt(e.target.value),
            total: pagination.total
        });
    };

    if (posts === null) {
        return <LoadingMessage>Loading posts...</LoadingMessage>;
    }

    if (posts.length === 0 && !isAdmin) {
        return <LoadingMessage>No posts available.</LoadingMessage>;
    }

    const totalPages = Math.ceil(pagination.total / pagination.limit);

    return (
        <Container>
            {posts.length === 0 && isAdmin ? (
                <LoadingMessage>No posts yet.</LoadingMessage>
            ) : (
                <>
                    {posts.map(post => (
                        <PostItemContainer key={post.id}>
                            {editingPostId === post.id ? (
                                <>
                                    <input 
                                        type="text" 
                                        value={editTitle} 
                                        onChange={(e) => setEditTitle(e.target.value)} 
                                        style={{width: '100%', marginBottom: '10px', fontSize: '32px', textAlign: 'center'}}
                                    />
                                    <textarea 
                                        value={editContent} 
                                        onChange={(e) => setEditContent(e.target.value)} 
                                        style={{width: '100%', minHeight: '100px', marginBottom: '10px'}}
                                    />
                                    <ButtonContainer>
                                        <button onClick={() => handleSaveEdit(post.id)}>Save</button>
                                        <button onClick={handleCancelEdit}>Cancel</button>
                                    </ButtonContainer>
                                </>
                            ) : (
                                <>
                                    <PostTitle>{post.title}</PostTitle>
                                    <PostContent>{post.content}</PostContent>
                                    <PostAuthor>Posted by: {post.username || `User ID: ${post.author_id}`}</PostAuthor>
                                    <CommentList postId={post.id} />
                                    <AddComment postId={post.id} onCommentCreated={handleCommentCreated} />
                                    {(isAdmin || userId === post.author_id) && (
                                        <ButtonContainer>
                                            <EditButton onClick={() => handleEditPost(post)}>
                                                Edit
                                            </EditButton>
                                            <DeleteButton onClick={() => handleDeletePost(post.id)}>
                                                Delete
                                            </DeleteButton>
                                        </ButtonContainer>
                                    )}
                                </>
                            )}
                        </PostItemContainer>
                    ))}

                    <PaginationContainer>
                        <PageSizeSelect 
                            value={pagination.limit} 
                            onChange={handleLimitChange}
                        >
                            <option value={5}>5 per page</option>
                            <option value={10}>10 per page</option>
                            <option value={20}>20 per page</option>
                        </PageSizeSelect>

                        <PaginationButton
                            onClick={() => handlePageChange(1)}
                            disabled={pagination.page <= 1}
                        >
                            First
                        </PaginationButton>

                        <PaginationButton
                            onClick={() => handlePageChange(pagination.page - 1)}
                            disabled={pagination.page <= 1}
                        >
                            Previous
                        </PaginationButton>

                        <span>Page {pagination.page} of {totalPages}</span>

                        <PaginationButton
                            onClick={() => handlePageChange(pagination.page + 1)}
                            disabled={pagination.page >= totalPages}
                        >
                            Next
                        </PaginationButton>

                        <PaginationButton
                            onClick={() => handlePageChange(totalPages)}
                            disabled={pagination.page >= totalPages}
                        >
                            Last
                        </PaginationButton>
                    </PaginationContainer>
                </>
            )}
        </Container>
    );
};

export default PostList;