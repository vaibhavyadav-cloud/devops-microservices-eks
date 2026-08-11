package com.devopsdemo.notification.dto;

// Shape of the "order.created" event published by the Order service to Kafka.
// This is the CONTRACT between Order and Notification — if Order changes
// this shape, this class needs updating too. Kept intentionally minimal
// for now (just enough to trigger a notification).
public class OrderCreatedEvent {

    private String orderId;
    private String customerEmail;

    public OrderCreatedEvent() {
        // required by Jackson for deserialization
    }

    public String getOrderId() {
        return orderId;
    }

    public void setOrderId(String orderId) {
        this.orderId = orderId;
    }

    public String getCustomerEmail() {
        return customerEmail;
    }

    public void setCustomerEmail(String customerEmail) {
        this.customerEmail = customerEmail;
    }
}
