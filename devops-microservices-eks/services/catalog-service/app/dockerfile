FROM node:20-alpine AS deps
WORKDIR /app
COPY package*.json .
RUN npm ci --omit=dev

# ---- Stage 2: Final runtime image ----
FROM node:20-alpine
WORKDIR /app
ENV NODE_ENV=production

# Sirf installed node_modules copy karo pehle stage se — apna npm/build tools wapas nahi aate
COPY --from=deps /app/node_modules ./node_modules
COPY . .

USER node
EXPOSE 3000
CMD ["node", "index.js"]