import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import styled from 'styled-components';

const CommentListContainer = styled.div`
    margin-top: 18px;
    padding: 14px 10px 10px 10px;
    border-radius: 10px;
    background-color: #f6f2fa;
    box-shadow: 0 1px 3px rgba(162, 89, 255, 0.06);
    border: 1.2px solid rgba(162, 89, 255, 0.18);
    max-width: 600px;
    margin-left: auto;
    margin-right: auto;
`;

const CommentItem = styled.div`
    padding: 10px 8px 8px 8px;
    margin-bottom: 10px;
    border-bottom: 1px solid #e0c3fc;
    font-family: 'Montserrat', Arial, sans-serif;
    &:last-child {
        border-bottom: none;
    }
`;

const CommentContent = styled.p`
    font-size: 15px;
    color: #6c2eb7;
    margin-bottom: 4px;
    font-family: 'Montserrat', Arial, sans-serif;
`;

const CommentAuthor = styled.small`
    color: #a259ff;
    font-style: italic;
    font-family: 'Montserrat', Arial, sans-serif;
`;

const LoadingMessage = styled.p`
    color: #a259ff;
    font-style: italic;
    font-family: 'Montserrat', Arial, sans-serif;
`;

const ErrorMessage = styled.p`
    color: #d32f2f;
    text-align: center;
    font-family: 'Montserrat', Arial, sans-serif;
`;

const NoCommentsMessage = styled.p`
    color: #bfa6e6;
    font-style: italic;
    font-family: 'Montserrat', Arial, sans-serif;
`;

const PaginationContainer = styled.div`
    display: flex;
    justify-content: center;
    align-items: center;
    margin-top: 10px;
    gap: 10px;
    flex-wrap: wrap;
    width: 100%;
`;

const PaginationButton = styled.button`
    padding: 6px 12px;
    cursor: pointer;
    background-color: ${props => props.active ? '#a259ff' : 'transparent'};
    color: ${props => props.active ? '#fff' : '#a259ff'};
    border: 1.2px solid #a259ff;
    border-radius: 5px;
    transition: background 0.2s, color 0.2s;
    font-family: 'Montserrat', Arial, sans-serif;
    font-size: 13px;
    min-width: 70px;
    &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
        background: #f9f9ff;
        color: #bfa6e6;
        border: 1.2px solid #e0c3fc;
    }
`;

const PageSizeSelect = styled.select`
    padding: 6px 8px;
    border-radius: 5px;
    border: 2px solid #a259ff;
    color: #a259ff;
    background: #f9f9ff;
    font-family: 'Montserrat', Arial, sans-serif;
    font-size: 13px;
    margin-left: 8px;
    outline: none;
`;

const PaginationInfo = styled.span`
    min-width: 110px;
    text-align: center;
    font-family: 'Montserrat', Arial, sans-serif;
    font-size: 15px;
    color: #222;
`;

const CommentList = ({ postId }) => {
    const [comments, setComments] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [pagination, setPagination] = useState({
        page: 1,
        limit: 5,
        total: 0
    });
    
    const isFetching = useRef(false);
    const refreshInterval = useRef(null);
    const cancelTokenSource = useRef(null);

    const fetchComments = async () => {
        if (isFetching.current) return;
        
        isFetching.current = true;
        setLoading(true);
        setError(null);

        // Отменяем предыдущий запрос, если он существует
        if (cancelTokenSource.current) {
            cancelTokenSource.current.cancel('Запрос отменен из-за нового запроса');
        }
        
        // Создаем новый токен отмены
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
                console.error('Ошибка при получении комментариев:', error);
                setError(error);
            }
        } finally {
            isFetching.current = false;
            setLoading(false);
        }
    };

    useEffect(() => {
        // Первоначальная загрузка
        fetchComments();

        // Устанавливаем интервал обновления
        const intervalId = setInterval(fetchComments, 5000);
        refreshInterval.current = intervalId;

        // Очистка при размонтировании
        return () => {
            clearInterval(refreshInterval.current);
            if (cancelTokenSource.current) {
                cancelTokenSource.current.cancel('Компонент размонтирован');
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
                            <CommentContent>{comment.content}</CommentContent>
                            <CommentAuthor>От {comment.username || `User ID: ${comment.author_id}`}</CommentAuthor>
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
                        
                        <PaginationInfo>Страница {pagination.page} из {totalPages}</PaginationInfo>
                    </PaginationContainer>
                </>
            ) : (
                <NoCommentsMessage>Нет комментариев.</NoCommentsMessage>
            )}
        </CommentListContainer>
    );
};

export default CommentList;