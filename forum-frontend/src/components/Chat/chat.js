import React, { useState, useEffect, useRef } from 'react';
import styled from 'styled-components';
import { useAuth } from './AuthContext';

const ChatContainer = styled.div`
  position: fixed;
  bottom: 20px;
  right: 20px;
  width: 350px;
  height: 500px;
  border: 1px solid #ddd;
  border-radius: 10px;
  background: white;
  display: flex;
  flex-direction: column;
  box-shadow: 0 2px 15px rgba(0, 0, 0, 0.1);
  z-index: 1000;
`;

const MessageList = styled.div`
  flex: 1;
  padding: 15px;
  overflow-y: auto;
  background: #f9f9f9;
  border-radius: 10px 10px 0 0;
`;

const MessageItem = styled.div`
  margin-bottom: 12px;
  padding: 10px 15px;
  background: ${props => props.isOwn ? '#e3f2fd' : '#ffffff'};
  border-radius: 15px;
  align-self: ${props => props.isOwn ? 'flex-end' : 'flex-start'};
  max-width: 80%;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
`;

const MessageHeader = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
  font-size: 0.8rem;
  color: #555;
`;

const MessageContent = styled.div`
  word-wrap: break-word;
`;

const MessageForm = styled.form`
  display: flex;
  padding: 15px;
  border-top: 1px solid #eee;
  background: white;
  border-radius: 0 0 10px 10px;
`;

const MessageInput = styled.input`
  flex: 1;
  padding: 10px 15px;
  border: 1px solid #ddd;
  border-radius: 20px;
  outline: none;
  font-size: 14px;

  &:focus {
    border-color: #4d90fe;
  }
`;

const SendButton = styled.button`
  margin-left: 10px;
  padding: 0 20px;
  background: #4d90fe;
  color: white;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.2s;

  &:hover {
    background: #357ae8;
  }

  &:disabled {
    background: #cccccc;
    cursor: not-allowed;
  }
`;

const StatusIndicator = styled.div.attrs({
    className: 'status-indicator'
  })`
    padding: 5px 15px;
    font-size: 0.8rem;
    color: ${props => props.$isConnected ? '#4CAF50' : '#F44336'};
    background: #f5f5f5;
    border-bottom: 1px solid #eee;
  `;

  const Chat = () => {
    const [messages, setMessages] = useState([]);
    const [newMessage, setNewMessage] = useState('');
    const [isConnected, setIsConnected] = useState(false);
    const { user, isAuthenticated } = useAuth();
    const ws = useRef(null);
    const messagesEndRef = useRef(null);
    const lastMessageRef = useRef(null); // Для отслеживания последнего сообщения
  
    // Форматирование даты
    const formatDate = (timestamp) => {
      try {
        const date = new Date(timestamp);
        return isNaN(date.getTime()) ? 'Just now' : 
          date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      } catch {
        return 'Just now';
      }
    };
  
    // Подключение WebSocket
    useEffect(() => {
      const connectWebSocket = () => {
        const wsUrl = `ws://localhost:8081/ws?userID=${user?.id || 0}&username=${encodeURIComponent(user?.username || 'Guest')}&auth=${isAuthenticated}`;
        ws.current = new WebSocket(wsUrl);
  
        ws.current.onopen = () => {
          console.log('WebSocket connected');
          setIsConnected(true);
        };
  
        ws.current.onclose = (e) => {
          console.log('WebSocket disconnected', e.code);
          setIsConnected(false);
          if (e.code !== 1000) {
            setTimeout(connectWebSocket, 3000);
          }
        };
  
        ws.current.onmessage = (e) => {
          try {
            const data = typeof e.data === 'string' ? e.data : new TextDecoder().decode(e.data);
            const message = JSON.parse(data);
            
            // Проверяем, не является ли это сообщение тем, которое мы только что отправили
            if (lastMessageRef.current && 
                lastMessageRef.current.content === message.content && 
                lastMessageRef.current.userID === message.userID) {
              lastMessageRef.current = null;
              return;
            }
  
            setMessages(prev => {
              // Проверяем дубликаты по содержанию и времени
              const isDuplicate = prev.some(
                m => m.content === message.content && 
                     Math.abs(new Date(m.timestamp) - new Date(message.timestamp)) < 1000
              );
              
              return isDuplicate ? prev : [...prev, {
                id: message.id || Date.now(),
                userID: message.userID || 0,
                username: message.username || 'Unknown',
                content: message.content || '',
                timestamp: message.timestamp || new Date().toISOString()
              }];
            });
          } catch (err) {
            console.error('Error processing message:', err);
          }
        };
  
        ws.current.onerror = (err) => {
          console.error('WebSocket error:', err);
        };
      };
  
      connectWebSocket();
  
      return () => {
        if (ws.current?.readyState === WebSocket.OPEN) {
          ws.current.close(1000, 'Component unmounted');
        }
      };
    }, [user, isAuthenticated]);
  
    // Прокрутка к новым сообщениям
    useEffect(() => {
      messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);
  
    // Отправка сообщения
    const sendMessage = (e) => {
      e.preventDefault();
      if (!newMessage.trim() || !isConnected || !isAuthenticated || !ws.current) return;
  
      try {
        const message = {
          userID: user.id,
          username: user.username,
          content: newMessage.trim(),
          timestamp: new Date().toISOString()
        };
        
        // Сохраняем ссылку на последнее отправленное сообщение
        lastMessageRef.current = message;
        
        // Добавляем сообщение локально сразу
        setMessages(prev => [...prev, {
          ...message,
          id: Date.now()
        }]);
        
        ws.current.send(JSON.stringify(message));
        setNewMessage('');
      } catch (err) {
        console.error('Error sending message:', err);
      }
    };
  
    return (
      <ChatContainer>
        <StatusIndicator $isConnected={isConnected}>
          {isConnected ? 'Online' : 'Connecting...'}
        </StatusIndicator>
        
        <MessageList>
          {messages.map((msg) => (
            <MessageItem key={msg.id} $isOwn={msg.userID === user?.id}>
              <MessageHeader>
                <span>{msg.username}</span>
                <span>{formatDate(msg.timestamp)}</span>
              </MessageHeader>
              <MessageContent>{msg.content}</MessageContent>
            </MessageItem>
          ))}
          <div ref={messagesEndRef} />
        </MessageList>
  
        {isAuthenticated ? (
          <MessageForm onSubmit={sendMessage}>
            <MessageInput
              value={newMessage}
              onChange={(e) => setNewMessage(e.target.value)}
              placeholder="Type a message..."
              disabled={!isConnected}
            />
            <SendButton
              type="submit"
              disabled={!newMessage.trim() || !isConnected}
            >
              Send
            </SendButton>
          </MessageForm>
        ) : (
          <MessageForm>
            <MessageInput
              placeholder="Login to participate in chat"
              disabled
            />
            <SendButton disabled>Send</SendButton>
          </MessageForm>
        )}
      </ChatContainer>
    );
  };
  
  export default Chat;