# Build stage
FROM node:20 AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

# Production stage
FROM nginx:alpine
WORKDIR /usr/share/nginx/html
COPY --from=builder /app/dist .
COPY --from=builder /app/public ./public
EXPOSE 5173
CMD ["nginx", "-g", "daemon off;"]
