package com.devopsdemo.notification.service;

import com.devopsdemo.notification.dto.NotificationRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.util.Map;

@Service
public class NotificationService {

    // SLF4J, not System.out.println — this is what actually integrates with
    // Spring Boot's logging config (and the structured JSON logging via
    // logback-spring.xml) for the K8s log pipeline, unlike raw stdout prints.
    private static final Logger log = LoggerFactory.getLogger(NotificationService.class);

    // Single core method — both the REST controller (manual/Postman testing)
    // and the Kafka consumer (real "order.created" events) call this same
    // method, so there's one place that owns "what sending a notification
    // actually means" regardless of how the request arrived.
    public Map<String, Object> send(String recipientEmail, String message) {
        // Stand-in for real delivery (SES/SendGrid call). Swap this method's
        // body for the real provider call later — callers don't need to change.
        log.info("sending_notification recipient={}", recipientEmail);

        return Map.of(
                "status", "sent",
                "recipientEmail", recipientEmail
        );
    }

    public Map<String, Object> send(NotificationRequest request) {
        return send(request.getRecipientEmail(), request.getMessage());
    }
}
