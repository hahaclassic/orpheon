# Используем Node.js образ
FROM node:20-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем package.json и package-lock.json
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем остальные файлы проекта
COPY . .

# # Собираем проект (если нужен build)
# RUN npm run build

# Указываем порт (обычно 5173 для Vite или 3000 для CRA)
EXPOSE 3000

# Запускаем dev-сервер
CMD ["npm", "run", "dev"]
