package com.gozcu.backend.netmon;

import jakarta.persistence.*;
import lombok.*;
import java.time.Instant;

@Entity
@Table(name = "network_alerts")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class NetworkAlert {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false)
    private Long pid;

    @Column(nullable = false)
    private Long uid;

    @Column(nullable = false)
    private String comm;

    @Column(nullable = false)
    private String destIp;

    @Column(nullable = false)
    private Integer dstPort;

    @Column(nullable = false)
    private String proto;

    @Column(nullable = false)
    private Instant detectedAt;
}
