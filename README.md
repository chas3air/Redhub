# Redhub

Redhub - это новостная социальная седь для разработчиков, в которой они могут делиться всякими интересными открытиями и рассказывать о каких-то событиях в мире технологий.



Приложение предоставляет API-Gateway, к которому можно обращаться с HTTP-запросами. А он в свою очередь обрацается к бекенду приложения, который написан на микросервисах. Микросервисы приложения написаны на Go с использованием gRPC, в качестве базы данных выбрана PostgreSQL. Работает на docker а также используется docker-compose

Запуск приложения
Если у вас установлен make, то нужно просто прописать команду make build-up
Если нет то docker compose up --build