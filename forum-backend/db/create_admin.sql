-- Создание пользователя с ролью администратора
INSERT INTO users (username, password_hash, role)
VALUES (
    'admin77',
    -- В реальном приложении пароль должен быть хэширован
    -- Здесь используется простой хэш для демонстрации
    crypt('123456', gen_salt('bf')),
    'admin'
);

-- Предоставление прав на удаление постов и комментариев
GRANT DELETE ON posts TO admin77;
GRANT DELETE ON comments TO admin77; 