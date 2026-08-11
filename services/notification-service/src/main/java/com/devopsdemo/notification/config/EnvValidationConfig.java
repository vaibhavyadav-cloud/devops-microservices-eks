package com.devopsdemo.notification.config;

import jakarta.annotation.PostConstruct;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.env.Environment;

import java.util.List;

// Fail-fast env validation, same intent as the centralized env validation
// in the Auth (Node.js) service: crash on startup with a clear error
// instead of failing confusingly later when a Kafka call or similar
// first touches a missing config value.
@Configuration
public class EnvValidationConfig {

    private static final Logger log = LoggerFactory.getLogger(EnvValidationConfig.class);

    private final Environment env;

    // Only vars that MUST be explicitly set in prod-like environments.
    // KAFKA_BOOTSTRAP_SERVERS has a localhost default for local dev, so it's
    // not in this list — but in EKS you'd typically still want it required;
    // tighten this list once we wire real per-environment config (Helm values).
    private static final List<String> REQUIRED_IN_PROD = List.of(
            "KAFKA_BOOTSTRAP_SERVERS"
    );

    public EnvValidationConfig(Environment env) {
        this.env = env;
    }

    @PostConstruct
    public void validate() {
        boolean isProd = List.of(env.getActiveProfiles()).contains("prod");

        if (!isProd) {
            log.info("env_validation_skipped reason=not_prod_profile");
            return;
        }

        List<String> missing = REQUIRED_IN_PROD.stream()
                .filter(key -> env.getProperty(key) == null || env.getProperty(key).isBlank())
                .toList();

        if (!missing.isEmpty()) {
            log.error("env_validation_failed missing={}", missing);
            throw new IllegalStateException(
                    "Missing required environment variables for prod profile: " + missing);
        }

        log.info("env_validation_passed");
    }
}
