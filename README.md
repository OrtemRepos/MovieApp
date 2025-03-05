# Movie Example Microservices Project
Проект представляет собой пример микросервисной архитектуры для управления метаданными фильмов, рейтингами и агрегацией информации о фильмах.

## Архитектура
### Сервисы:
1. Metadata Service
    - Хранение метаданных фильмов (название, описание, режиссер)
    - Реализации: in-memory, PostgreSQL
    - API: gRPC и HTTP

2. Rating Service
    - Управление рейтингами пользователей
    - Агрегация оценок
    - Интеграция с Kafka для обработки событий
    - Реализации: in-memory, PostgreSQL
    - API: gRPC и HTTP

3. Movie Service
    - Агрегация данных из Metadata и Rating сервисов
    - Предоставление полной информации о фильме
    - API: gRPC и HTTP

Общие компоненты:
    - Service Discovery (Consul/in-memory)
    - Вспомогательные утилиты для gRPC

### Зависимости
Обязательные:
    - Go 1.19+
    - PostgreSQL 14+
    - Consul 1.13+
    - Kafka 3.3+ (для Rating Service)

Опциональные:
    - grpcurl (для тестирования gRPC)
    - Docker (для упрощения развертывания)
