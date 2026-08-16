# Build Stage
FROM node:24-alpine AS build

# Enable pnpm (version pinned via the workspace root's "packageManager" field)
RUN corepack enable

WORKDIR /app

# Declare build arguments for environment variables
# These are needed at build time for Vite to embed them in the bundle
ARG VITE_PORT=5173
ARG VITE_BASE_PATH=/

# Set as environment variables for the build process
ENV VITE_PORT=${VITE_PORT}
ENV VITE_BASE_PATH=${VITE_BASE_PATH}

# Copy workspace manifests first for better layer caching
COPY package.json pnpm-workspace.yaml pnpm-lock.yaml ./
COPY apps/admin/package.json ./apps/admin/
COPY apps/consent/package.json ./apps/consent/
COPY apps/member/package.json ./apps/member/

# Install only member-portal's dependencies (and any workspace packages it depends on)
RUN pnpm install --frozen-lockfile --filter member-portal...

# Copy member-portal's source code
COPY apps/member ./apps/member

# Build the application (Vite will use the environment variables)
RUN pnpm --filter member-portal run build

# Serve Stage
FROM nginx:1.27-alpine

# Copy built static files from build stage
COPY --from=build /app/apps/member/dist /usr/share/nginx/html

# Copy custom nginx configuration
COPY apps/member/nginx.conf /etc/nginx/conf.d/default.conf

# Copy entrypoint script
COPY apps/member/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Set default environment variables (can be overridden at runtime)
ENV VITE_API_URL=""
ENV VITE_LOGS_URL=""
ENV VITE_CLIENT_ID=""
ENV VITE_BASE_URL=""
ENV VITE_SCOPE=""
ENV VITE_SIGN_IN_REDIRECT_URL=""
ENV VITE_SIGN_OUT_REDIRECT_URL=""

# Expose port 80
EXPOSE 80

# Use entrypoint script to generate config.js at runtime
ENTRYPOINT ["/entrypoint.sh"]

