import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import styled from 'styled-components';

const CommentListContainer = styled.div`
    margin-top: 15px;
    padding: 10px;
    border-radius: 8px;
    background-color: #f8f8f8;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
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
    color: #333;
    margin-bottom: 5px;
`;

const CommentAuthor = styled.small`
    color: #777;
    font-style: italic;
`;

const LoadingMessage = styled.p`
    color: #555;
    font-style: italic;
`;

const ErrorMessage = styled.p`
    color: #d32f2f;
`;

const NoCommentsMessage = styled.p`
    color: #999;
    font-style: italic;
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
            cancelTokenSource.current.cancel('Request canceled due to new request');
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
                console.error('Error fetching comments:', error);
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
                cancelTokenSource.current.cancel('Component unmounted');
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

    if (loading) return <LoadingMessage>Loading comments...</LoadingMessage>;
    if (error) return <ErrorMessage>Error: {error.message} {error.response?.status && `(Status: ${error.response.status})`}</ErrorMessage>;

    const totalPages = Math.ceil(pagination.total / pagination.limit);

    return (
        <CommentListContainer>
            <h4>Comments ({pagination.total}):</h4>
            
            {comments.length ? (
                <>
                    {comments.map(comment => (
                        <CommentItem key={comment.id}>
                            <CommentContent>{comment.content}</CommentContent>
                            <CommentAuthor>By {comment.username || `User ID: ${comment.author_id}`}</CommentAuthor>
                        </CommentItem>
                    ))}
                    
                    <PaginationContainer>
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
                        
                        <PageSizeSelect 
                            value={pagination.limit} 
                            onChange={handleLimitChange}
                        >
                            <option value={3}>3 per page</option>
                            <option value={5}>5 per page</option>
                            <option value={10}>10 per page</option>
                        </PageSizeSelect>
                    </PaginationContainer>
                </>
            ) : (
                <NoCommentsMessage>No comments yet.</NoCommentsMessage>
            )}
        </CommentListContainer>
    );
};

export default CommentList;