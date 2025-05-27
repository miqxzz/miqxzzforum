import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import styled from 'styled-components';
import CommentList from './CommentList';
import AddComment from './AddComment';

// === Styled Components ===

const Container = styled.div`
    max-width: 700px;
    margin: 32px auto;
    padding: 32px 20px 24px 20px;
    background-color: #f9f9ff;
    border-radius: 18px;
    box-shadow: 0 0 16px rgba(162, 89, 255, 0.10);
    border: 1.5px solid rgba(162, 89, 255, 0.3);
`;

const PostItemContainer = styled.div`
    margin-bottom: 40px;
    padding: 32px 28px 24px 28px;
    border: 1.5px solid rgba(162, 89, 255, 0.25);
    border-radius: 18px;
    background-color: #fff;
    box-shadow: 0 4px 24px rgba(162, 89, 255, 0.07);
    display: flex;
    flex-direction: column;
    align-items: stretch;
    max-width: 600px;
    margin-left: auto;
    margin-right: auto;
`;

const PostTitle = styled.h3`
    color: #a259ff;
    text-align: left;
    margin-bottom: 18px;
    font-size: 2rem;
    font-weight: 800;
    letter-spacing: 0.5px;
    font-family: 'Montserrat', Arial, sans-serif;
`;

const PostContent = styled.p`
    color: #6c2eb7;
    font-size: 1.1rem;
    margin-bottom: 18px;
    text-align: left;
    font-family: 'Montserrat', Arial, sans-serif;
`;

const PostAuthor = styled.small`
    color: #a259ff;
    font-style: italic;
    margin-bottom: 18px;
    display: block;
    font-size: 1rem;
    font-family: 'Montserrat', Arial, sans-serif;
`;

const ButtonContainer = styled.div`
    display: flex;
    gap: 16px;
    margin-top: 18px;
    justify-content: flex-end;
`;

const EditButton = styled.button`
    background-color: #a259ff;
    color: #fff;
    padding: 12px 28px;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.2s, box-shadow 0.2s;
    font-family: 'Montserrat', Arial, sans-serif;
    font-size: 1rem;
    min-width: 130px;
    height: 44px;
    font-weight: 700;
    box-shadow: 0 2px 8px rgba(162, 89, 255, 0.13);
    letter-spacing: 0.5px;
    display: flex;
    align-items: center;
    justify-content: center;
    text-align: center;
    &:hover {
        background-color: #7c3aed;
        box-shadow: 0 4px 16px rgba(162, 89, 255, 0.18);
    }
`;

const DeleteButton = styled.button`
    background-color: #f44336;
    color: white;
    padding: 12px 28px;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.2s, box-shadow 0.2s;
    font-family: 'Montserrat', Arial, sans-serif;
    font-size: 1rem;
    min-width: 130px;
    height: 44px;
    font-weight: 600;
    box-shadow: 0 2px 8px rgba(244, 67, 54, 0.08);
    &:hover {
        background-color: #d32f2f;
        box-shadow: 0 4px 16px rgba(244, 67, 54, 0.13);
    }
`;

const PaginationContainer = styled.div`
    display: flex;
    justify-content: center;
    margin: 32px 0 0 0;
    gap: 10px;
    flex-wrap: wrap;
    align-items: center;
`;

const PaginationButton = styled.button`
    padding: 10px 18px;
    cursor: pointer;
    background-color: ${props => props.active ? '#a259ff' : 'transparent'};
    color: ${props => props.active ? '#fff' : '#a259ff'};
    border: 1.5px solid #a259ff;
    border-radius: 7px;
    transition: background 0.2s, color 0.2s, box-shadow 0.2s;
    font-family: 'Montserrat', Arial, sans-serif;
    font-size: 1rem;
    min-width: 110px;
    font-weight: 600;
    box-shadow: 0 2px 8px rgba(162, 89, 255, 0.06);
    &:hover:not(:disabled) {
        background: #e0c3fc;
        color: #6c2eb7;
        box-shadow: 0 4px 16px rgba(162, 89, 255, 0.13);
    }
    &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
        background: #f9f9ff;
        color: #bfa6e6;
        border: 1.5px solid #e0c3fc;
        box-shadow: none;
    }
`;

const PageSizeSelect = styled.select`
    padding: 10px 16px;
    border-radius: 7px;
    border: 2.5px solid #a259ff;
    color: #a259ff;
    background: #f9f9ff;
    font-family: 'Montserrat', Arial, sans-serif;
    font-size: 1rem;
    margin-right: 10px;
    font-weight: 600;
    outline: none;
    box-shadow: 0 2px 8px rgba(162, 89, 255, 0.08);
`;

const LoadingMessage = styled.p`
    color: #a259ff;
    font-style: italic;
    text-align: center;
    font-family: 'Montserrat', Arial, sans-serif;
`;

const EditInput = styled.input`
    width: 100%;
    margin-bottom: 10px;
    font-size: 2rem;
    text-align: center;
    padding: 12px;
    border: 1.5px solid #a259ff;
    border-radius: 8px;
    font-family: 'Montserrat', Arial, sans-serif;
    font-weight: 700;
    color: #a259ff;
    background: #f9f9ff;
    outline: none;
    transition: border 0.2s;
    &:focus {
        border: 2px solid #6c2eb7;
    }
`;

const EditTextarea = styled.textarea`
    width: 100%;
    min-height: 100px;
    margin-bottom: 10px;
    padding: 12px;
    border: 1.5px solid #a259ff;
    border-radius: 8px;
    font-size: 1.1rem;
    font-family: 'Montserrat', Arial, sans-serif;
    background: #f9f9ff;
    color: #6c2eb7;
    outline: none;
    transition: border 0.2s;
    resize: vertical;
    &:focus {
        border: 2px solid #6c2eb7;
    }
`;

const EditButtonStyled = styled.button`
    background-color: #a259ff;
    color: #fff;
    padding: 12px 28px;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-family: 'Montserrat', Arial, sans-serif;
    font-size: 1rem;
    min-width: 130px;
    height: 44px;
    font-weight: 700;
    margin-right: 12px;
    transition: background 0.2s, box-shadow 0.2s;
    box-shadow: 0 2px 8px rgba(162, 89, 255, 0.13);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    text-align: center;
    &:hover {
        background-color: #7c3aed;
        box-shadow: 0 4px 16px rgba(162, 89, 255, 0.18);
    }
`;

const CancelButtonStyled = styled.button`
    background-color: #f3f3f3;
    color: #a259ff;
    padding: 12px 28px;
    border: 1.5px solid #a259ff;
    border-radius: 8px;
    cursor: pointer;
    font-family: 'Montserrat', Arial, sans-serif;
    font-size: 1rem;
    min-width: 130px;
    height: 44px;
    font-weight: 700;
    transition: background 0.2s, color 0.2s, box-shadow 0.2s;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    text-align: center;
    &:hover {
        background: #e0c3fc;
        color: #6c2eb7;
        box-shadow: 0 4px 16px rgba(162, 89, 255, 0.13);
    }
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
            console.error('Ошибка при получении постов:', error);
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
            console.error('Error deleting post:', error);
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
            console.error('Error updating post:', error);
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
        return <LoadingMessage>Нет постов.</LoadingMessage>;
    }

    const totalPages = Math.ceil(pagination.total / pagination.limit);

    return (
        <Container>
            {posts.length === 0 && isAdmin ? (
                <LoadingMessage>Нет постов.</LoadingMessage>
            ) : (
                <>
                    {posts.map(post => (
                        <PostItemContainer key={post.id}>
                            {editingPostId === post.id ? (
                                <>
                                    <EditInput 
                                        type="text" 
                                        value={editTitle} 
                                        onChange={(e) => setEditTitle(e.target.value)} 
                                    />
                                    <EditTextarea 
                                        value={editContent} 
                                        onChange={(e) => setEditContent(e.target.value)} 
                                    />
                                    <ButtonContainer>
                                        <EditButtonStyled onClick={() => handleSaveEdit(post.id)}>Сохранить</EditButtonStyled>
                                        <CancelButtonStyled onClick={handleCancelEdit}>Отменить</CancelButtonStyled>
                                    </ButtonContainer>
                                </>
                            ) : (
                                <>
                                    <PostTitle>{post.title}</PostTitle>
                                    <PostContent>{post.content}</PostContent>
                                    <PostAuthor>От: {post.username || `User ID: ${post.author_id}`}</PostAuthor>
                                    <CommentList postId={post.id} />
                                    <AddComment postId={post.id} onCommentCreated={handleCommentCreated} />
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
                                </>
                            )}
                        </PostItemContainer>
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

                        <span>Страница {pagination.page} из {totalPages}</span>
                    </PaginationContainer>
                </>
            )}
        </Container>
    );
};

export default PostList;