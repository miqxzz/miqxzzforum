import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import styled from 'styled-components';

const CommentListContainer = styled.div`
    margin-top: 15px;
    padding: 10px;
    border-radius: 8px;
    background-color: #f8f4fc;
    box-shadow: 0 2px 8px #e1d5ee44;
    font-family: 'Montserrat', sans-serif;
`;

const CommentItem = styled.div`
    padding: 8px;
    margin-bottom: 8px;
    border-bottom: 1px solid #eee;
    &:last-child {
        border-bottom: none;
    }
`;

const CommentContent = styled.p`
    font-size: 14px;
    color: #3d2466;
    margin-bottom: 5px;
    font-weight: 400;
    font-family: 'Montserrat', sans-serif;
`;

const CommentAuthor = styled.small`
    color: #b085d6;
    font-style: italic;
    font-weight: 300;
    font-family: 'Montserrat', sans-serif;
`;

const DeleteButton = styled.button`
    background: none;
    border: none;
    color: #e74c3c;
    cursor: pointer;
    padding: 5px 10px;
    font-size: 0.9rem;
    transition: color 0.2s ease;
    font-weight: 500;

    &:hover {
        color: #c0392b;
    }
`;

const CommentHeader = styled.div`
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
`;

const LoadingMessage = styled.p`
    color: #555;
    font-style: italic;
    font-weight: 400;
`;

const ErrorMessage = styled.p`
    color: #d291bc;
    font-weight: 500;
    font-family: 'Montserrat', sans-serif;
`;

const NoCommentsMessage = styled.p`
    color: #999;
    font-style: italic;
    font-weight: 400;
`;

const PaginationContainer = styled.div`
    display: flex;
    justify-content: center;
    margin-top: 15px;
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
    margin-left: 10px;
`;

const CommentContainer = styled.div`
    background: #f9f9f9;
    border-radius: 8px;
    padding: 15px;
    margin-bottom: 15px;
    border: 1px solid #e0e0e0;
`;

const CommentList = ({ postId }) => {
    const [comments, setComments] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [isAdmin, setIsAdmin] = useState(false);
    const [userId, setUserId] = useState(null);
    const [pagination, setPagination] = useState({
        page: 1,
        limit: 5,
        total: 0
    });
    
    const isFetching = useRef(false);
    const refreshInterval = useRef(null);
    const cancelTokenSource = useRef(null);

    useEffect(() => {
        const storedRole = localStorage.getItem('role');
        const storedUserId = localStorage.getItem('userID');
        setIsAdmin(storedRole === 'admin');
        setUserId(storedUserId ? parseInt(storedUserId) : null);
    }, []);

    const fetchComments = async () => {
        if (isFetching.current) return;
        
        isFetching.current = true;
        setLoading(true);
        setError(null);

        if (cancelTokenSource.current) {
            cancelTokenSource.current.cancel('Запрос отменен из-за нового запроса');
        }
        
        cancelTokenSource.current = axios.CancelToken.source();
        
        try {
            const response = await axios.get(
                `http://localhost:8081/posts/${Number(postId)}/comments`, 
                {
                    params: {
                        page: pagination.page,
                        limit: pagination.limit
                    },
                    cancelToken: cancelTokenSource.current.token
                }
            );
            
            setComments(response.data.comments || []);
            setPagination(prev => ({
                ...prev,
                total: response.data.pagination?.total || 0
            }));
        } catch (error) {
            if (!axios.isCancel(error)) {
                console.error('Ошибка получения комментариев:', error);
                setError(error);
            }
        } finally {
            isFetching.current = false;
            setLoading(false);
        }
    };

    const handleDeleteComment = async (commentId) => {
        const token = localStorage.getItem('token');

        if (!token) {
            alert('Вы не авторизованы.');
            return;
        }

        try {
            await axios.delete(`http://localhost:8081/comments/${commentId}`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            fetchComments();
        } catch (error) {
            console.error('Ошибка удаления комментария:', error);
            if (error.response && error.response.status === 403) {
                alert('У вас нет прав на удаление этого комментария.');
            } else {
                alert('Не удалось удалить комментарий.');
            }
        }
    };

    useEffect(() => {
        fetchComments();

        const intervalId = setInterval(fetchComments, 5000);
        refreshInterval.current = intervalId;

        return () => {
            clearInterval(refreshInterval.current);
            if (cancelTokenSource.current) {
                cancelTokenSource.current.cancel('Компонент отмонтирован');
            }
        };
    }, [postId, pagination.page, pagination.limit]);

    const handlePageChange = (newPage) => {
        setPagination(prev => ({
            ...prev,
            page: newPage
        }));
    };

    const handleLimitChange = (e) => {
        setPagination({
            page: 1,
            limit: parseInt(e.target.value),
            total: pagination.total
        });
    };

    if (loading) return <LoadingMessage>Загрузка комментариев...</LoadingMessage>;
    if (error) return <ErrorMessage>Ошибка: {error.message} {error.response?.status && `(Статус: ${error.response.status})`}</ErrorMessage>;

    const totalPages = Math.ceil(pagination.total / pagination.limit);

    return (
        <CommentListContainer>
            <h4>Комментарии ({pagination.total}):</h4>
            
            {comments.length ? (
                <>
                    {comments.map(comment => (
                        <CommentItem key={comment.id}>
                            <CommentHeader>
                                <CommentContent>
                                    {comment.content}
                                    <span style={{color: '#b085d6', fontSize: '12px', marginLeft: '8px'}}>
                                        {comment.created_at ? new Date(comment.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : ''}
                                    </span>
                                </CommentContent>
                                {(isAdmin || userId === comment.author_id) && (
                                    <DeleteButton onClick={() => handleDeleteComment(comment.id)}>
                                        Удалить
                                    </DeleteButton>
                                )}
                            </CommentHeader>
                            <CommentAuthor>От {comment.username || `ID пользователя: ${comment.author_id}`}</CommentAuthor>
                        </CommentItem>
                    ))}
                    
                    <PaginationContainer>
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
                        
                        <PageSizeSelect 
                            value={pagination.limit} 
                            onChange={handleLimitChange}
                        >
                            <option value={3}>3 на странице</option>
                            <option value={5}>5 на странице</option>
                            <option value={10}>10 на странице</option>
                        </PageSizeSelect>
                    </PaginationContainer>
                </>
            ) : (
                <NoCommentsMessage>Нет комментариев.</NoCommentsMessage>
            )}
        </CommentListContainer>
    );
};

export default CommentList;