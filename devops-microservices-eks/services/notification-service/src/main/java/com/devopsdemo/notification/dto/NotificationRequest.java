package com.devopsdemo.notification.dto;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;

// A typed request body instead of Map<String,Object> — the other services
// (Pydantic in FastAPI, express-validator-style in Node) all validate at
// the boundary the same way; this is Spring's idiomatic version of that
// using jakarta.validation annotations + @Valid in the controller.
public class NotificationRequest {

    @NotBlank(message = "recipientEmail is required")
    @Email(message = "recipientEmail must be a valid email")
    private String recipientEmail;

    @NotBlank(message = "message is required")
    private String message;

    public String getRecipientEmail() {
        return recipientEmail;
    }

    public void setRecipientEmail(String recipientEmail) {
        this.recipientEmail = recipientEmail;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }
}
