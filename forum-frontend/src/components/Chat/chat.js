import React, { useState, useEffect, useRef } from 'react';
import styled from 'styled-components';
import { useAuth } from './AuthContext';

const ChatContainer = styled.div`
  width: 420px;
  min-height: 340px;
  background: #fff;
  border: 2px solid #a678f7;
  border-radius: 18px;
  box-shadow: 0 2px 16px #e1d5ee33;
  display: flex;
  flex-direction: column;
  margin: 0 auto;
  margin-bottom: 32px;
  font-family: 'Montserrat', sans-serif;
`;

const MessageList = styled.div`
  flex: 1;
  padding: 18px 12px 8px 12px;
  overflow-y: auto;
  background: none;
  border-radius: 16px 16px 0 0;
`;

const MessageItem = styled.div`
  margin-bottom: 12px;
  padding: 10px 18px;
  background: #faf7ff;
  border-radius: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: none;
  font-family: 'Montserrat', sans-serif;
`;

const MessageHeader = styled.div`
  display: flex;
  justify-content: space-between;
  width: 100%;
  font-size: 0.95em;
  color: #a678f7;
  margin-bottom: 2px;
`;

const MessageContent = styled.div`
  color: #3d2466;
  font-size: 1.1em;
  font-family: 'Montserrat', sans-serif;
`;

const MessageForm = styled.form`
  display: flex;
  padding: 12px 12px 16px 12px;
  border-top: 1.5px solid #f3e6fa;
  background: none;
  border-radius: 0 0 16px 16px;
`;

const MessageInput = styled.input`
  flex: 1;
  padding: 10px 15px;
  border: 2px solid #a678f7;
  border-radius: 8px;
  outline: none;
  font-size: 1em;
  font-family: 'Montserrat', sans-serif;
  color: #3d2466;
  background: #fff;
  margin-right: 10px;
  &:focus {
    border-color: #8e44ad;
  }
`;

const SendButton = styled.button`
  padding: 10px 22px;
  background: #a678f7;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1em;
  font-family: 'Montserrat', sans-serif;
  font-weight: 600;
  transition: background 0.2s;
  box-shadow: 0 2px 8px #e1d5ee44;
  &:hover {
    background: #8e44ad;
  }
  &:disabled {
    background: #e1d5ee;
    color: #fff;
    cursor: not-allowed;
  }
`;

const StatusIndicator = styled.div.attrs({
    className: 'status-indicator'
  })`
    padding: 5px 15px;
    font-size: 0.9em;
    color: ${props => props.$isConnected ? '#a678f7' : '#F44336'};
    background: none;
    border-bottom: none;
    text-align: right;
    margin-bottom: 8px;
  `;

const ChatMessage = styled.div`
  background: ${props => props.isOwn ? '#e9d6f7' : '#f3e6fa'};
  color: #3d2466;
  border-radius: 8px;
  padding: 8px 12px;
  margin-bottom: 8px;
  font-family: 'Montserrat', sans-serif;
`;

const ChatInput = styled.input`
  background: #f3e6fa;
  border: 1.5px solid #c8a2e8;
  border-radius: 5px;
  padding: 8px;
  font-family: 'Montserrat', sans-serif;
  color: #3d2466;
  transition: border 0.2s;
  &:focus {
    border: 1.5px solid #9b59b6;
    outline: none;
  }
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
        return isNaN(date.getTime()) ? 'Только что' : 
          date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      } catch {
        return 'Только что';
      }
    };
  
    // Подключение WebSocket
    useEffect(() => {
      const connectWebSocket = () => {
        const wsUrl = `ws://localhost:8081/ws?userID=${user?.id || 0}&username=${encodeURIComponent(user?.username || 'Guest')}&auth=${isAuthenticated}`;
        ws.current = new WebSocket(wsUrl);
  
        ws.current.onopen = () => {
          console.log('WebSocket подключен');
          setIsConnected(true);
        };
  
        ws.current.onclose = (e) => {
          console.log('WebSocket отключен', e.code);
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
            console.error('Ошибка обработки сообщения:', err);
          }
        };
  
        ws.current.onerror = (err) => {
          console.error('Ошибка WebSocket:', err);
        };
      };
  
      connectWebSocket();
  
      return () => {
        if (ws.current?.readyState === WebSocket.OPEN) {
          ws.current.close(1000, 'Компонент отмонтирован');
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
        console.error('Ошибка отправки сообщения:', err);
      }
    };
  
    return (
      <ChatContainer>
    
        <MessageList>
          {messages.map((msg) => (
            <MessageItem key={msg.id} $isOwn={msg.userID === user?.id}>
              <MessageHeader>
                <span>{msg.username}</span>
              </MessageHeader>
              <MessageContent>
                {msg.content}
                <span style={{ color: '#b085d6', fontSize: '0.95em', marginLeft: '10px' }}>
                  {formatDate(msg.timestamp)}
                </span>
              </MessageContent>
            </MessageItem>
          ))}
          <div ref={messagesEndRef} />
        </MessageList>
  
        {isAuthenticated ? (
          <MessageForm onSubmit={sendMessage}>
            <MessageInput
              value={newMessage}
              onChange={(e) => setNewMessage(e.target.value)}
              placeholder="Введите сообщение..."
              disabled={!isConnected}
            />
            <SendButton
              type="submit"
              disabled={!newMessage.trim() || !isConnected}
            >
                  Отправить
            </SendButton>
          </MessageForm>
        ) : (
          <MessageForm>
            <MessageInput
              placeholder="Войдите в чат"
              disabled
            />
            <SendButton disabled>Отправить</SendButton>
          </MessageForm>
        )}
      </ChatContainer>
    );
  };
  
  export default Chat;