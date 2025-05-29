import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import styled from 'styled-components';
import CommentList from './CommentList';
import AddComment from './AddComment';
import Chat from '../Chat/chat';
import { useAuth } from '../Chat/AuthContext';

// === Styled Components ===

const PostListContainer = styled.div`
    background-color: #f8f4fc;
    border-radius: 10px;
    box-shadow: 0 2px 8px #e1d5ee44;
    padding: 24px;
    font-family: 'Montserrat', sans-serif;
`;

const PostCard = styled.div`
    background: linear-gradient(135deg, #f8f4fc 60%, #e1d5ee 100%);
    border-radius: 16px;
    padding: 32px 28px 24px 28px;
    margin-bottom: 32px;
    box-shadow: 0 4px 24px 0 #e1d5ee88;
    border: 1.5px solid #e0e0e0;
    transition: transform 0.2s, box-shadow 0.2s;
    position: relative;
    overflow: hidden;
    &:hover {
        transform: translateY(-4px) scale(1.01);
        box-shadow: 0 8px 32px 0 #c8a2e8aa;
    }
`;

const PostTitle = styled.h2`
    color: #7c3aed;
    margin: 0 0 12px 0;
    font-size: 2.1rem;
    font-weight: 700;
    letter-spacing: 0.5px;
    background: linear-gradient(90deg, #a259ff 0%, #f8f4fc 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
`;

const PostMeta = styled.div`
    display: flex;
    align-items: center;
    gap: 18px;
    color: #b085d6;
    font-size: 1rem;
    margin-bottom: 18px;
    font-weight: 500;
    font-family: 'Montserrat', sans-serif;
`;

const PostAuthor = styled.span`
    color: #8e44ad;
    font-weight: 600;
`;

const PostDate = styled.span`
    color: #b085d6;
    font-size: 0.95rem;
    font-style: italic;
`;

const PostContent = styled.p`
    color: #3d2466;
    line-height: 1.7;
    margin-bottom: 22px;
    font-size: 1.15rem;
    font-weight: 400;
    font-family: 'Montserrat', sans-serif;
    background: #fff;
    border-radius: 8px;
    padding: 18px 16px;
    box-shadow: 0 2px 8px #e1d5ee33;
`;

const DecorativeLine = styled.div`
    width: 100%;
    height: 2px;
    background: linear-gradient(90deg, #a259ff 0%, #f8f4fc 100%);
    margin: 18px 0 22px 0;
    border-radius: 2px;
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

const ErrorText = styled.div`
    color: #d291bc;
    font-family: 'Montserrat', sans-serif;
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

const PostListWrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  font-family: 'Montserrat', sans-serif;
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

    const { isAuthenticated } = useAuth();

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
            console.error('Ошибка получения постов:', error);
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
            alert('Вы не авторизованы.');
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
            console.error('Ошибка удаления поста:', error);
            if (error.response && error.response.status === 403) {
                alert('У вас нет прав на удаление этого поста.');
            } else {
                alert('Не удалось удалить пост.');
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
            alert('Вы не авторизованы.');
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
            console.error('Ошибка обновления поста:', error);
            if (error.response && error.response.status === 403) {
                alert('У вас нет прав на редактирование этого поста.');
            } else {
                alert('Не удалось обновить пост.');
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
        return <LoadingMessage>Загрузка постов...</LoadingMessage>;
    }

    if (posts.length === 0 && !isAdmin) {
        return <LoadingMessage>Нет доступных постов.</LoadingMessage>;
    }

    const totalPages = Math.ceil(pagination.total / pagination.limit);

    return (
        <PostListWrapper>
            <PostListContainer>
                {posts.length === 0 && isAdmin ? (
                    <LoadingMessage>Нет доступных постов.</LoadingMessage>
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
                                                <button onClick={() => handleSaveEdit(post.id)}>Сохранить</button>
                                                <button onClick={handleCancelEdit}>Отменить</button>
                                            </ButtonContainer>
                                        </PostActions>
                                    </>
                                ) : (
                                    <>
                                        <PostTitle>{post.title}</PostTitle>
                                        <PostMeta>
                                            <PostAuthor>От {post.username || `ID пользователя: ${post.author_id}`}</PostAuthor>
                                        </PostMeta>
                                        <DecorativeLine />
                                        <PostContent>{post.content}</PostContent>
                                        <CommentList postId={post.id} />
                                        <AddComment postId={post.id} onCommentCreated={handleCommentCreated} />
                                        <PostActions>
                                            {(isAdmin || userId === post.author_id) && (
                                                <ButtonContainer>
                                                    <EditButton onClick={() => handleEditPost(post)}>
                                                        Редактировать
                                                    </EditButton>
                                                    <DeleteButton onClick={() => handleDeletePost(post.id)}>
                                                        Удалить
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
                                <option value={5}>5 на странице</option>
                                <option value={10}>10 на странице</option>
                                <option value={20}>20 на странице</option>
                            </PageSizeSelect>

                            <PaginationButton
                                onClick={() => handlePageChange(1)}
                                disabled={pagination.page <= 1}
                            >
                                Первая
                            </PaginationButton>

                            <PaginationButton
                                onClick={() => handlePageChange(pagination.page - 1)}
                                disabled={pagination.page <= 1}
                            >
                                Предыдущая
                            </PaginationButton>

                            <span>Страница {pagination.page} из {totalPages}</span>

                            <PaginationButton
                                onClick={() => handlePageChange(pagination.page + 1)}
                                disabled={pagination.page >= totalPages}
                            >
                                Следующая
                            </PaginationButton>

                            <PaginationButton
                                onClick={() => handlePageChange(totalPages)}
                                disabled={pagination.page >= totalPages}
                            >
                                Последняя
                            </PaginationButton>
                        </PaginationContainer>
                    </>
                )}
            </PostListContainer>
            <Chat />
        </PostListWrapper>
    );
};

export default PostList;