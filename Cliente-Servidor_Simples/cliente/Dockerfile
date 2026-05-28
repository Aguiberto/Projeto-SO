# Esse arquivo funciona configurando o container para a aplicação

FROM golang:1.22-alpine                         
WORKDIR /app    
COPY . .                                        
RUN go build -o /usr/local/bin/echo-client .    
CMD ["/usr/local/bin/echo-client"]              
