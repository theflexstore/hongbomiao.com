# Base image
FROM node:14

# Set working directory
WORKDIR /app

# Add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

# Build
COPY . /app
RUN yarn build

# Start
EXPOSE 5000
CMD ["yarn", "start"]
