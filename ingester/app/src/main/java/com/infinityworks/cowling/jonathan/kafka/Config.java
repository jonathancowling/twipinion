package com.infinityworks.cowling.jonathan.kafka;

import java.time.Duration;

import org.apache.kafka.clients.admin.NewTopic;
import org.apache.kafka.common.config.TopicConfig;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.config.TopicBuilder;

@Configuration
public class Config {

    @Bean
    public NewTopic tweets() {
        return TopicBuilder.name("tweets")
            .replicas(1)
            .partitions(1)
            .config(
                TopicConfig.DELETE_RETENTION_MS_CONFIG,
                String.valueOf(Duration.ofHours(1).toMillis())
            )
            .build();
    }
}
