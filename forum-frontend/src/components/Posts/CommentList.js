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
    font-weight: 400;
`;

const CommentAuthor = styled.small`
    color: #777;
    font-style: italic;
    font-weight: 300;
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
    color: #d32f2f;
    font-weight: 500;
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
            cancelTokenSource.current.cancel('Request canceled due to new request');
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
                console.error('Error fetching comments:', error);
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
            alert('You are not authenticated.');
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
            console.error('Error deleting comment:', error);
            if (error.response && error.response.status === 403) {
                alert('You are not authorized to delete this comment.');
            } else {
                alert('Failed to delete comment.');
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
                            <CommentHeader>
                                <CommentContent>{comment.content}</CommentContent>
                                {(isAdmin || userId === comment.author_id) && (
                                    <DeleteButton onClick={() => handleDeleteComment(comment.id)}>
                                        Delete
                                    </DeleteButton>
                                )}
                            </CommentHeader>
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