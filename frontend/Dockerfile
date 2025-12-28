FROM node:24.12.0-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci
RUN npm install @supabase/supabase-js

COPY . .
EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
