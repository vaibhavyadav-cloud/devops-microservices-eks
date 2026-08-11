package com.devopsdemo.notification.consumer;

import com.devopsdemo.notification.dto.OrderCreatedEvent;
import com.devopsdemo.notification.service.NotificationService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

@Component
public class OrderEventConsumer {

    private static final Logger log = LoggerFactory.getLogger(OrderEventConsumer.class);

    private final NotificationService notificationService;

    public OrderEventConsumer(NotificationService notificationService) {
        this.notificationService = notificationService;
    }

    // topics = "${...}" pulls the topic name from application.properties
    // (notification.kafka.topic.order-created) instead of hardcoding it here,
    // same env-driven config approach as the rest of the service.
    @KafkaListener(topics = "${notification.kafka.topic.order-created}")
    public void onOrderCreated(OrderCreatedEvent event) {
        log.info("order_created_event_received orderId={}", event.getOrderId());

        // NOTE: no try/catch here on purpose for now — if this throws, the
        // record is NOT acknowledged and Kafka will redeliver it (default
        // ack mode). That's the correct "at-least-once" behavior for a demo.
        // In a real system you'd add a dead-letter topic for records that
        // keep failing, so one bad message can't block the whole partition
        // forever — worth doing once we get to production-hardening.
        notificationService.send(event.getCustomerEmail(), "Your order " + event.getOrderId() + " was placed.");
    }
}
