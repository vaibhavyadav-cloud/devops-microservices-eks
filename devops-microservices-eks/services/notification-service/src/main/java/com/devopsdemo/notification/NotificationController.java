package com.devopsdemo.notification;

import com.devopsdemo.notification.dto.NotificationRequest;
import com.devopsdemo.notification.service.NotificationService;
import jakarta.validation.Valid;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
public class NotificationController {

    private final NotificationService notificationService;

    public NotificationController(NotificationService notificationService) {
        this.notificationService = notificationService;
    }

    // Note: no hand-rolled /health here anymore — Actuator's
    // /actuator/health/liveness and /actuator/health/readiness cover that
    // now (see application.properties). Controller only owns business routes.
    @PostMapping("/notify")
    public Map<String, Object> notify(@Valid @RequestBody NotificationRequest request) {
        return notificationService.send(request);
    }
}
