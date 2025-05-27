CREATE TRIGGER IF NOT EXISTS cleanup_old_messages
    AFTER INSERT ON chat_messages
BEGIN
    DELETE FROM chat_messages
    WHERE timestamp < datetime('now', '-10 minutes');
END;