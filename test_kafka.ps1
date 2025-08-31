# Имя топика и брокера
$topic = "orders_test"
$broker = "localhost:9092"
$api_url = "http://localhost:8081/order"

# Массив тестовых заказов
$orders = @(
    @{ order_uid = "b563feb7b2b84b6test_1"; delivery_id = "d1" },
    @{ order_uid = "b563feb7b2b84b6test_2"; delivery_id = "d2" },
    @{ order_uid = "b563feb7b2b84b6test_3"; delivery_id = "d3" }
)

# Функция для генерации JSON
function Generate-OrderJson($order) {
    return @{
        order_uid = $order.order_uid
        track_number = "WBILMTESTTRACK"
        entry = "WBIL"
        locale = "en"
        internal_signature = ""
        customer_id = "test-customer"
        delivery_service = "WBIL"
        shard_key = "9"
        sm_id = 99
        date_created = (Get-Date).ToString("yyyy-MM-ddTHH:mm:ssZ")
        oof_shard = "1"
        delivery = @{
            delivery_id = $order.delivery_id
            order_uid = $order.order_uid
            name = "Test User"
            phone = "+1234567890"
            zip = "123456"
            city = "Moscow"
            address = "Red Square, 1"
            region = "Moscow"
            email = "test@example.com"
        }
        payment = @{
            payment_id = "p_" + $order.delivery_id
            order_uid = $order.order_uid
            transaction = "trx_" + $order.delivery_id
            request_id = ""
            currency = "USD"
            provider = "wbpay"
            amount = 1500.50
            payment_dt = [int][double]::Parse((Get-Date -UFormat %s))
            bank = "TestBank"
            delivery_cost = 200.00
            goods_total = 1300.50
            custom_fee = 0
        }
        items = @(
            @{
                item_id = "i_" + $order.delivery_id
                order_uid = $order.order_uid
                chrt_id = 12345
                track_number = "WBILMTESTTRACK"
                price = 500.25
                rid = "rid_" + $order.delivery_id
                name = "T-Shirt"
                sale = 10
                size = "L"
                total_price = 450.23
                nm_id = 123456
                brand = "WB"
                status = 202
            }
        )
    } | ConvertTo-Json -Compress
}

# Отправка сообщений через bash внутри контейнера
foreach ($order in $orders) {
    $json = Generate-OrderJson $order
    Write-Host "Sending order $($order.order_uid) to Kafka..."

    # Используем bash вместо PowerShell
    $escapedJson = $json.Replace('"','\"')  # экранируем кавычки для bash
    docker exec -i kafka bash -c "echo `"$escapedJson`" | kafka-console-producer --topic $topic --bootstrap-server $broker"
}

# Ждём, чтобы Go-сервис успел обработать
Start-Sleep -Seconds 3

# Проверка через API
foreach ($order in $orders) {
    $url = "$api_url/$($order.order_uid)"
    Write-Host "Checking API for order $($order.order_uid)..."
    try {
        $response = curl $url
        Write-Host $response
    } catch {
        Write-Host "Failed to get order $($order.order_uid)"
    }
}
