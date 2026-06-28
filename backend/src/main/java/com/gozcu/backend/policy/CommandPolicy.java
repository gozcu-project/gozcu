package com.gozcu.backend.policy;

import jakarta.persistence.*;
import lombok.*;

@Entity
@Table(name = "command_policies")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class CommandPolicy {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    // Örn: "systemctl restart nginx" (tam eşleşme) veya "systemctl restart *" (wildcard)
    @Column(nullable = false, unique = true)
    private String pattern;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private RiskLevel riskLevel;
}