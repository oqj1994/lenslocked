FROM node:latest
WORKDIR /tailwind
RUN npm init -y 
RUN  npm install tailwindcss && \
    npx tailwindcss init

CMD npx tailwindcss -c /src/tailwind.config.js -i /src/input.css -o /dst/style.css --watch --poll