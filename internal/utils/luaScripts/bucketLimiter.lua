local key = KEYS[1]

local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

local data = redis.call("HMGET", key, "tokens", "timestamp")

local tokens = tonumber(data[1])
local timestamp = tonumber(data[2])

-- первый запрос пользователя
if tokens == nil then
    tokens = capacity
    timestamp = now
end

-- сколько прошло секунд
local elapsed = now - timestamp

-- добавляем токены
local refill = elapsed * refill_rate

tokens = math.min(
    capacity,
    tokens + refill
)

local allowed = 0

-- пробуем взять токен
if tokens >= requested then
    tokens = tokens - requested
    allowed = 1
end

-- сохраняем состояние
redis.call(
    "HMSET",
    key,
    "tokens",
    tokens,
    "timestamp",
    now
)

-- TTL чтобы не хранить старых пользователей
redis.call(
    "EXPIRE",
    key,
    3600
)

return {allowed, tokens}
