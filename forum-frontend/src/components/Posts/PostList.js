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

const PostCard = styled.div`
    background: white;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 20px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    border: 1px solid #e0e0e0;
    transition: transform 0.2s ease;

    &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
`;

const PostTitle = styled.h2`
    color: #8e44ad;
    margin: 0 0 10px 0;
    font-size: 1.5rem;
    font-weight: 600;
`;

const PostMeta = styled.div`
    color: #666;
    font-size: 0.9rem;
    margin-bottom: 15px;
    font-weight: 400;
`;

const PostContent = styled.p`
    color: #333;
    line-height: 1.6;
    margin-bottom: 15px;
    font-weight: 400;
`;

const PostActions = styled.div`
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 15px;
    padding-top: 15px;
    border-top: 1px solid #eee;
`;

const ActionButton = styled.button`
    background: none;
    border: none;
    color: #9b59b6;
    cursor: pointer;
    padding: 5px 10px;
    font-size: 0.9rem;
    transition: color 0.2s ease;
    font-weight: 500;

    &:hover {
        color: #8e44ad;
    }
`;

const DeleteButton = styled(ActionButton)`
    color: #e74c3c;
    &:hover {
        color: #c0392b;
    }
`;

const EditButton = styled.button`
    background-color: #2196F3;
    color: white;
    padding: 8px 12px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s ease;
    font-weight: 500;

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
    font-weight: 400;
`;

const ErrorMessage = styled.p`
    color: #d32f2f;
    text-align: center;
    font-weight: 500;
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
                        <PostCard key={post.id}>
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
                                    <PostActions>
                                        <ButtonContainer>
                                            <button onClick={() => handleSaveEdit(post.id)}>Save</button>
                                            <button onClick={handleCancelEdit}>Cancel</button>
                                        </ButtonContainer>
                                    </PostActions>
                                </>
                            ) : (
                                <>
                                    <PostTitle>{post.title}</PostTitle>
                                    <PostMeta>Posted by: {post.username || `User ID: ${post.author_id}`}</PostMeta>
                                    <PostContent>{post.content}</PostContent>
                                    <CommentList postId={post.id} />
                                    <AddComment postId={post.id} onCommentCreated={handleCommentCreated} />
                                    <PostActions>
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
                                    </PostActions>
                                </>
                            )}
                        </PostCard>
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